package rule

import (
	"net"
	"net/http"
	"regexp"
	"time"
)

const (
	DEFAULT_ACCESS          = "block"
	DEFAULT_TYPE            = "host"
	DEFAULT_PATTERN         = ""
	DEFAULT_PASSWORD_BYPASS = true

	DEFAULT_BASIC_AUTH_REALM = "pc-proxy: Enter password to bypass block"
)

var (
	accessValues map[string]string = map[string]string{"block": "block", "allow": "allow"}
	typeValues   map[string]string = map[string]string{"host": "host", "path": "path", "url": "url"}
)

type RuleAccess string

type RuleType string

type Rule struct {
	Access         RuleAccess
	Type           RuleType
	Pattern        string
	PasswordBypass bool
}

func New() Rule {
	return Rule{
		Access:         DEFAULT_ACCESS,
		Type:           DEFAULT_TYPE,
		Pattern:        DEFAULT_PATTERN,
		PasswordBypass: DEFAULT_PASSWORD_BYPASS,
	}
}

func (r Rule) Match(req *http.Request, resp http.ResponseWriter, password string, bypassCache *map[string]map[string]time.Duration, bypassCacheDuration time.Duration) (match bool, allow bool) {
	var checkStr string

	// does this rule allow access?
	allowed := r.Access.String() == "allow"

	switch r.Type {
	case "host":
		checkStr = req.Host
	case "path":
		checkStr = req.URL.Path
	case "url":
		checkStr = req.URL.String()
	default:
		// can't match against an invalid rule
		return false, allowed
	}

	if matched, err := regexp.Match(r.Pattern, []byte(checkStr)); err != nil || !matched {
		// error occurred during match check OR not the rule we're looking for
		return false, allowed
	}

	//
	// rule is matched
	//

	if !allowed && r.PasswordBypass {
		// this request is only allowed if the bypass password has been specified
		clientIp, _, _ := net.SplitHostPort(req.RemoteAddr)

		if clientIp != "" {
			// check the bypass cache firs
			if clientCache, ok := (*bypassCache)[clientIp]; ok {
				if duration, dOk := clientCache[r.Pattern]; dOk && duration > 0 {
					// the bypass password has been specified previously, allow the request
					return true, true
				}
			}
		}

		proxyAuth := req.Header.Get("Proxy-Authorization")

		if proxyAuth != "" {
			req.Header.Set("Authorization", proxyAuth)
		}

		_, p, ok := req.BasicAuth()

		// unset the proxy authorization header
		if proxyAuth != "" {
			req.Header.Del("Authorization")
		}

		if !ok {
			// no basic auth provided, request it and exit
			// request basic auth
			resp.Header().Set("Proxy-Authenticate", "Basic realm=\""+DEFAULT_BASIC_AUTH_REALM+"\"")
			http.Error(resp, http.StatusText(http.StatusProxyAuthRequired), http.StatusProxyAuthRequired)

			return true, false
		}

		if p == password {
			// the provided password is valid, allow access
			// cache the successful bypass
			clientCache, ok := (*bypassCache)[clientIp]

			if ok {
				clientCache[r.Pattern] = bypassCacheDuration
			} else {
				// create the cache
				(*bypassCache)[clientIp] = map[string]time.Duration{
					r.Pattern: bypassCacheDuration,
				}
			}

			return true, true
		}
	}

	return true, allowed
}

func (a RuleAccess) IsValid() bool {
	_, ok := accessValues[string(a)]

	return ok
}

func (a RuleAccess) String() string {
	ret, ok := accessValues[string(a)]

	if !ok {
		return ""
	}

	return ret
}

func (t RuleType) IsValid() bool {
	_, ok := typeValues[string(t)]

	return ok
}

func (t RuleType) String() string {
	ret, ok := typeValues[string(t)]

	if !ok {
		return ""
	}

	return ret
}
