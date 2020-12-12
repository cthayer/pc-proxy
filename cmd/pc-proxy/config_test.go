package main

import (
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/cthayer/pc-proxy/internal/config"
	"github.com/cthayer/pc-proxy/internal/proxy"
)

func TestLoadConfigFile(t *testing.T) {
	configFile := filepath.Join("..", "..", "test", "config")
	configFileExts := []string{".json", ".hcl"}

	for _, ext := range configFileExts {
		confFile := configFile + ext

		if err := LoadConfigFile(confFile, func(conf *config.Config) {}); err != nil {
			t.Errorf("Failed to load "+strings.ToUpper(ext)+" configuration file. (%s)  %v", confFile, err)
		}

		conf := config.GetConfig()

		if conf.TLS.Enabled {
			t.Errorf("expected TLS.Enabled to be false, got %v", conf.TLS.Enabled)
		}

		if conf.TLS.Cert != "" {
			t.Errorf("expected TLS.Cert to be empty, got %v", conf.TLS.Cert)
		}

		if conf.TLS.Key != "" {
			t.Errorf("expected TLS.Key to be empty, got %v", conf.TLS.Key)
		}

		if len(conf.Rules) != 1 {
			t.Errorf("expected 1 rule, got %v", len(conf.Rules))
		}

		access, aOk := conf.Rules[0]["access"]

		if !aOk {
			t.Errorf("config rule missing 'access' property")
		}

		if access != "block" {
			t.Errorf("expected %v, got %v", "block", access)
		}

		pattern, pOk := conf.Rules[0]["pattern"]

		if !pOk {
			t.Errorf("config rule missing 'pattern' property")
		}

		if pattern != "" {
			t.Errorf("expected %v, got %v", "", pattern)
		}

		passwordBypass, pbOk := conf.Rules[0]["passwordBypass"]

		if !pbOk {
			t.Errorf("config rule missing 'passwordBypass' property")
		}

		if !passwordBypass.(bool) {
			t.Errorf("expected %v, got %v", true, passwordBypass)
		}

		ty, tOk := conf.Rules[0]["type"]

		if !tOk {
			t.Errorf("config rule missing 'type' property")
		}

		if ty != "host" {
			t.Errorf("expected %v, got %v", "host", ty)
		}
	}
}

func TestLoadConfigFileAllowRulesPatternEscape(t *testing.T) {
	configFile := filepath.Join("..", "..", "test", "config-allow-rules-pattern-escape")
	configFileExts := []string{".hcl"}

	for _, ext := range configFileExts {
		confFile := configFile + ext

		if err := LoadConfigFile(confFile, func(conf *config.Config) {}); err != nil {
			t.Errorf("Failed to load "+strings.ToUpper(ext)+" configuration file. (%s)  %v", confFile, err)
		}

		conf := config.GetConfig()

		if conf.TLS.Enabled {
			t.Errorf("expected TLS.Enabled to be false, got %v", conf.TLS.Enabled)
		}

		if conf.TLS.Cert != "" {
			t.Errorf("expected TLS.Cert to be empty, got %v", conf.TLS.Cert)
		}

		if conf.TLS.Key != "" {
			t.Errorf("expected TLS.Key to be empty, got %v", conf.TLS.Key)
		}

		if len(conf.Rules) != 17 {
			t.Errorf("expected 17 rules, got %v", len(conf.Rules))
		}

		access, aOk := conf.Rules[0]["access"]

		if !aOk {
			t.Errorf("config rule missing 'access' property")
		}

		if access != "allow" {
			t.Errorf("expected %v, got %v", "allow", access)
		}

		pattern, pOk := conf.Rules[0]["pattern"]

		if !pOk {
			t.Errorf("config rule missing 'pattern' property")
		}

		if pattern != "zoom\\.us" {
			t.Errorf("expected %v, got %v", "zoom\\.us", pattern)
		}

		passwordBypass, pbOk := conf.Rules[0]["passwordBypass"]

		if !pbOk {
			t.Errorf("config rule missing 'passwordBypass' property")
		}

		if !passwordBypass.(bool) {
			t.Errorf("expected %v, got %v", true, passwordBypass)
		}

		ty, tOk := conf.Rules[0]["type"]

		if !tOk {
			t.Errorf("config rule missing 'type' property")
		}

		if ty != "host" {
			t.Errorf("expected %v, got %v", "host", ty)
		}
	}
}

func TestLoadConfigFileWithProxy(t *testing.T) {
	pxy := proxy.New()
	patternReplaceRegex := regexp.MustCompile(`[^-\w\d.]`)

	configFile := filepath.Join("..", "..", "test", "config-allow-rules-pattern-escape")
	configFileExt := ".hcl"
	confFile := configFile + configFileExt

	if err := LoadConfigFile(confFile, pxy.LoadConfig); err != nil {
		t.Errorf("Failed to load "+strings.ToUpper(configFileExt)+" configuration file. (%s)  %v", confFile, err)
	}
	
	conf := config.GetConfig()

	if len(pxy.Rules) != len(conf.Rules) {
		t.Errorf("expected %v rules, got %v", len(conf.Rules), len(pxy.Rules))
	}
	
	for i, r := range pxy.Rules {
		if string(r.Access) != conf.Rules[i]["access"] {
			t.Errorf("expected %v, got %v", conf.Rules[i]["access"], r.Access)
		}

		if r.Pattern != conf.Rules[i]["pattern"] {
			t.Errorf("expected %v, got %v", conf.Rules[i]["pattern"], r.Pattern)
		}

		if r.PasswordBypass != conf.Rules[i]["passwordBypass"] {
			t.Errorf("expected %v, got %v", conf.Rules[i]["passwordBypass"], r.PasswordBypass)
		}

		if string(r.Type) != conf.Rules[i]["type"] {
			t.Errorf("expected %v, got %v", conf.Rules[i]["type"], r.Type)
		}

		target := "https://" + patternReplaceRegex.ReplaceAllString(r.Pattern, "") + "/foo/bar"

		req := httptest.NewRequest("GET", target, nil)
		bypassCache := map[string]map[string]time.Duration{}

		match, allowed := r.Match(req, httptest.NewRecorder(), "test", &bypassCache, time.Minute)

		if !match {
			t.Errorf("expected r %v to match %v", r.Pattern, target)
		}

		if conf.Rules[i]["access"] == "allow" && !allowed {
			t.Errorf("expected allowed to be true, got %v", allowed)
		}

		if conf.Rules[i]["access"] == "block" && allowed {
			t.Errorf("expected allowed to be false, got %v", allowed)
		}
	}
}
