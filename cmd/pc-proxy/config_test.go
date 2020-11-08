package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/cthayer/pc-proxy/internal/config"
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
