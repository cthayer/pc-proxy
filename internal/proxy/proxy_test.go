package proxy

import (
	"github.com/cthayer/pc-proxy/internal/logger"
	"testing"

	"github.com/cthayer/pc-proxy/internal/config"
)

func TestNew(t *testing.T) {
	_ = New()
}

func TestProxy_Start_Stop(t *testing.T) {
	conf := config.GetConfig()
	logger.InitLogger("info", "console")

	pxy := New()

	pxy.LoadConfig(conf)

	err := pxy.Start()

	if err != nil {
		t.Errorf("Error starting proxy: %v", err)
		return
	}

	errs := pxy.Stop()

	if len(errs) > 0 {
		t.Errorf("Errors stopping proxy: %v", errs)
	}
}

func TestProxy_LoadConfig(t *testing.T) {

}
