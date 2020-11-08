package main

import (
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cthayer/pc-proxy/internal/logger"
	"github.com/cthayer/pc-proxy/internal/proxy"
)

const (
	PROXY_START_TIMEOUT = 120
)

var (
	// set at build time using -ldflags=" -X 'main.VERSION=$version'"
	//   - edit build/versions.json to set the version when building
	VERSION = "dev"
)

func main() {
	cmdErr := cliRootCmd.Execute()

	if cmdErr != nil {
		os.Exit(1)
	}
}

func runProxy() {
	// setup the Proxy
	pxy := proxy.New()

	// load configuration
	if err := LoadConfigFile(cliConf.ConfigFile, pxy.LoadConfig); err != nil {
		_, _ = os.Stderr.WriteString("Failed to load configuration\n")
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	log := logger.GetLogger()
	defer log.Sync()

	log.Info("Configuration Loaded")

	if cliConf.PidFile != "" {
		// write the main PID to the PidFile (for initV services)
		err := writePidFile()

		if err != nil {
			// log the error, but otherwise ignore it
			log.Error("Error writing PidFile", zap.Error(err), zap.String("pidFile", cliConf.PidFile))
		}
	}

	// start the proxy
	startErr := pxy.Start()

	if startErr != nil {
		log.Error("Server failed to start", zap.Error(startErr))
		log.Info("Exiting")
		return
	}

	// setup OS signal handler
	done := setupSignalHandler()

	// wait until signaled to exit (SIGINT or SIGTERM)
	<-done

	log.Info("Shutting down")

	// stop the server
	stopErrs := pxy.Stop()

	if len(stopErrs) > 0 {
		log.Error("Proxy failed to stop", zap.Any("errors", stopErrs))
		log.Info("Exiting")
		return
	}

	log.Info("Shutdown complete")
}

func setupSignalHandler() chan bool {
	log := logger.GetLogger()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// we don't need to handle SIGHUP for config reload because we're watching
	// the config file for changes
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer close(done)

		for {
			sig := <-sigs

			log.Debug("Got signal: " + sig.String())

			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				done <- true
			}
		}
	}()

	return done
}

func writePidFile() error {
	return ioutil.WriteFile(cliConf.PidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
}
