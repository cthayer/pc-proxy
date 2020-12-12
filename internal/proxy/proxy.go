package proxy

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/smartystreets/cproxy/v2"
	"go.uber.org/zap"

	"github.com/cthayer/pc-proxy/internal/config"
	"github.com/cthayer/pc-proxy/internal/logger"
	"github.com/cthayer/pc-proxy/internal/rule"
)

const (
	BYPASS_PASSWD_ENV_NAME   = "BYPASS_PASSWORD"
	BYPASS_PASSWD_CACHE_TIME = time.Millisecond * 72000000 // 20 hours
	TLS_MIN_VERSION          = tls.VersionTLS12
	HTTP_SERVER_STOP_TIMEOUT = 300
)

type Proxy struct {
	Rules               []rule.Rule
	password            string
	handler             http.Handler
	logger              *zap.Logger
	tlsConf             config.TLSConfig
	listenConf          config.ListenConfig
	httpSrv             *http.Server
	tlsSrv              *http.Server
	waitGroup           sync.WaitGroup
	netListener         net.Listener
	tlsNetListener      net.Listener
	passwordBypassCache map[string]map[string]time.Duration
}

func New() *Proxy {
	p := Proxy{
		Rules:               []rule.Rule{},
		password:            "",
		handler:             nil,
		logger:              logger.GetLogger(),
		tlsConf:             config.GetConfig().TLS,
		listenConf:          config.GetConfig().Listen,
		httpSrv:             nil,
		tlsSrv:              nil,
		waitGroup:           sync.WaitGroup{},
		netListener:         nil,
		tlsNetListener:      nil,
		passwordBypassCache: map[string]map[string]time.Duration{},
	}

	// start the passwordBypassCache manager
	go p.managePasswordBypassCache()

	return &p
}

func (p *Proxy) Start() error {
	var err error = nil

	p.handler = cproxy.New(cproxy.Options.Filter(p))

	p.httpSrv = &http.Server{Addr: p.listenConf.Host + ":" + strconv.Itoa(p.listenConf.Port)}

	// setup the http server
	p.httpSrv.Handler = p.handler

	p.netListener, err = net.Listen("tcp", p.httpSrv.Addr)

	if err != nil {
		return err
	}

	p.logger.Info("HTTP server listening", zap.String("listen address", p.httpSrv.Addr))

	// setup the https server
	if p.tlsConf.Enabled {
		p.tlsSrv = &http.Server{Addr: p.listenConf.Host + ":" + strconv.Itoa(p.listenConf.TlsPort)}

		if err = p.setupTls(); err != nil {
			return err
		}

		p.tlsSrv.Handler = p.handler

		p.tlsNetListener, err = net.Listen("tcp", p.tlsSrv.Addr)

		if err != nil {
			return err
		}

		p.logger.Info("HTTPS server listening", zap.String("listen address", p.tlsSrv.Addr))

		// start the HTTPS proxy async
		p.waitGroup.Add(1)
		go func() {
			defer p.waitGroup.Done()

			// this will block until the server is stopped or an error occurs
			err := p.tlsSrv.ServeTLS(p.tlsNetListener, "", "")

			// if the server exits for any reason other than being stopped, log the error
			if err != http.ErrServerClosed {
				p.logger.Error("Error serving HTTPS requests", zap.Error(err))
			}
		}()
	}

	// start the HTTP proxy async
	p.waitGroup.Add(1)
	go func() {
		defer p.waitGroup.Done()

		// this will block until the server is stopped or an error occurs
		err := p.httpSrv.Serve(p.netListener)

		// if the server exits for any reason other than being stopped, log the error
		if err != http.ErrServerClosed {
			p.logger.Error("Error serving HTTP requests", zap.Error(err))
		}
	}()

	return err
}

func (p *Proxy) Stop() []error {
	var errs []error

	// stop the HTTP server async
	go func() {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*HTTP_SERVER_STOP_TIMEOUT)

		if err := p.httpSrv.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}

		cancel()

		p.logger.Debug("HTTP server shutdown")
	}()

	if p.tlsSrv != nil {
		// stop the HTTPS server async
		go func() {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*HTTP_SERVER_STOP_TIMEOUT)

			if err := p.tlsSrv.Shutdown(ctx); err != nil {
				errs = append(errs, err)
			}

			cancel()

			p.logger.Debug("HTTPS server shutdown")
		}()
	}

	// wait for shutdown to finish
	p.waitGroup.Wait()

	return errs
}

func (p *Proxy) LoadConfig(conf *config.Config) {
	// get new logger
	p.logger = logger.GetLogger()

	p.password = os.Getenv(BYPASS_PASSWD_ENV_NAME)
	p.tlsConf = conf.TLS
	p.listenConf = conf.Listen // this config will not update without a restart of the service

	p.updateRules(conf.Rules)

	if p.tlsSrv != nil {
		_ = p.setupTls()
	}
}

func (p *Proxy) IsAuthorized(resp http.ResponseWriter, req *http.Request) bool {
	p.logger.Debug("request received", zap.Any("headers", req.Header), zap.String("client address", req.RemoteAddr))

	// check the rules to see if this request is allowed
	for _, r := range p.Rules {
		if match, allow := r.Match(req, resp, p.password, &p.passwordBypassCache, BYPASS_PASSWD_CACHE_TIME); match {
			p.logger.Debug("processed request", zap.String("url", req.URL.String()), zap.Any("rule", r), zap.Bool("match", match), zap.Bool("allow", allow), zap.Any("respHeaders", resp.Header().Get("Proxy-Authenticate")))

			return allow
		}
	}

	p.logger.Debug("no matching rules.  allowing access", zap.String("url", req.URL.String()))

	// by default we allow access
	return true
}

func (p *Proxy) updateRules(rules []map[string]interface{}) {
	var newRules []rule.Rule

	for _, v := range rules {
		a, aOk := v["access"].(string)
		t, tOk := v["type"].(string)
		pat, pOk := v["pattern"].(string)
		b, bOk := v["passwordBypass"].(bool)

		r := rule.New()

		if aOk && rule.RuleAccess(a).IsValid() {
			r.Access = rule.RuleAccess(a)
		}

		if tOk && rule.RuleType(t).IsValid() {
			r.Type = rule.RuleType(t)
		}

		if pOk {
			r.Pattern = pat
		}

		if bOk {
			r.PasswordBypass = b
		}

		if r.Pattern != "" {
			newRules = append(newRules, r)
		}
	}

	p.Rules = newRules

	p.logger.Info("new rules loaded", zap.Any("rules", p.Rules))
}

func (p *Proxy) setupTls() error {
	keyPair, err := tls.LoadX509KeyPair(p.tlsConf.Cert, p.tlsConf.Key)

	if err != nil {
		return err
	}

	if p.tlsSrv.TLSConfig == nil {
		p.tlsSrv.TLSConfig = &tls.Config{}
		p.tlsSrv.TLSConfig.PreferServerCipherSuites = true
		p.tlsSrv.TLSConfig.MinVersion = TLS_MIN_VERSION
	}

	if p.tlsSrv.TLSConfig.Certificates == nil {
		p.tlsSrv.TLSConfig.Certificates = make([]tls.Certificate, 1)
	}

	p.tlsSrv.TLSConfig.Certificates[0] = keyPair

	// load ciphers
	var cipherIds []uint16
	cipherNames := strings.Split(p.tlsConf.Ciphers, ":")

	for _, cipher := range tls.CipherSuites() {
		for _, cn := range cipherNames {
			if cn == cipher.Name {
				cipherIds = append(cipherIds, cipher.ID)
			}
		}
	}

	// allow insecure ciphers to be specified, but log warnings when using them
	for _, cipher := range tls.InsecureCipherSuites() {
		for _, cn := range cipherNames {
			if cn == cipher.Name {
				p.logger.Warn("Allowing insecure TLS cipher", zap.String("cipher", cipher.Name))
				cipherIds = append(cipherIds, cipher.ID)
			}
		}
	}

	p.tlsSrv.TLSConfig.CipherSuites = cipherIds

	return nil
}

func (p *Proxy) managePasswordBypassCache() {
	// this function is run in a background go thread
	for {
		<-time.After(time.Minute)

		p.logger.Debug("updating password bypass cache", zap.Any("passwordBypassCache", p.passwordBypassCache))

		for ip, cache := range p.passwordBypassCache {
			for pattern, duration := range cache {
				p.passwordBypassCache[ip][pattern] = duration - time.Minute

				if p.passwordBypassCache[ip][pattern] <= 0 {
					delete(p.passwordBypassCache[ip], pattern)
				}
			}

			if len(p.passwordBypassCache[ip]) < 1 {
				delete(p.passwordBypassCache, ip)
			}
		}

		p.logger.Debug("finished updating password bypass cache", zap.Any("passwordBypassCache", p.passwordBypassCache))
	}
}
