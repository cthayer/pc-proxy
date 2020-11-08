package main

import (
	"errors"
	"github.com/cthayer/pc-proxy/internal/logger"
	"go.uber.org/zap"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/spf13/cobra"

	"github.com/cthayer/pc-proxy/internal/config"
)

const (
	DEFAULT_CONFIG_FILE = "/etc/pc-proxy/config.json"
	DEFAULT_PID_FILE    = ""
	DEFAULT_LOG_LEVEL   = "info"
)

var cliRootCmd = cobra.Command{
	Use:     "pc-proxy",
	Short:   "Parental control proxy for controlling access to websites",
	Long:    "Parental control proxy for controlling access to websites",
	Example: "  pc-proxy -c /path/to/config.json",
	Args:    cobra.ExactArgs(0),
	Version: VERSION,
	Run: func(cmd *cobra.Command, args []string) {
		runProxy()
	},
}

type cliConfig struct {
	ConfigFile string
	PidFile    string
}

var k = koanf.New(".")

var cliConf cliConfig = cliConfig{
	ConfigFile: DEFAULT_CONFIG_FILE,
	PidFile:    DEFAULT_PID_FILE,
}

func init() {
	cliRootCmd.PersistentFlags().StringVarP(&cliConf.ConfigFile, "config-file", "c", DEFAULT_CONFIG_FILE, "path to JSON or HCL formatted configuration file")
	cliRootCmd.PersistentFlags().StringVarP(&cliConf.PidFile, "pid-file", "", DEFAULT_PID_FILE, "the file to write the pid to (used for initv style services")
}

func LoadConfigFile(confFile string, onChange func(conf *config.Config)) error {
	var parser koanf.Parser

	extension := filepath.Ext(confFile)

	switch extension {
	case ".json":
		parser = json.Parser()
	case ".hcl":
		parser = hcl.Parser(true)
	default:
		return errors.New("unsupported configuration file type (" + extension + ").  Must be one of: '.json' or '.hcl'")
	}

	provider := file.Provider(confFile)

	if err := k.Load(provider, parser); err != nil {
		return err
	}

	if err := loadConfig(); err != nil {
		return errors.New("error unmarshalling config: " + err.Error())
	}

	// initialize the logger
	log, err := logger.InitLogger(config.GetConfig().Logging.Level, config.GetConfig().Logging.Encoding)

	if err != nil {
		return errors.New("error initializing logger: " + err.Error())
	}

	if log == nil {
		return errors.New("logger initialized to nil")
	}

	defer log.Sync()

	log.Debug("Configuration File: " + cliConf.ConfigFile)
	log.Debug("", zap.Any("cliConf", cliConf))

	// call `onChange` function
	onChange(config.GetConfig())

	// Watch the config file and reload the config when it changes.
	// File provider always returns a nil `event`.
	if wErr := provider.Watch(func(event interface{}, err error) {
		log := logger.GetLogger()
		defer log.Sync()

		if err != nil {
			log.Error("config watch error", zap.Error(err))
			return
		}

		log.Info("config changed. Reloading ...")

		// save the old logging config before reloading the config
		oldLogLevel := config.GetConfig().Logging.Level
		oldLogEncoding := config.GetConfig().Logging.Encoding

		// reload the config
		if err := k.Load(provider, parser); err != nil {
			log.Error("error loading config file", zap.Error(err))
		}

		if err := loadConfig(); err != nil {
			log.Error("error unmarshalling config", zap.Error(err))
		}

		// update the logging config if it has changed
		if oldLogLevel != config.GetConfig().Logging.Level || oldLogEncoding != config.GetConfig().Logging.Encoding {
			// reconfigure the logger
			log.Info("Changing logging config", zap.String("oldLevel", oldLogLevel), zap.String("oldEncoding", oldLogEncoding), zap.String("newLevel", config.GetConfig().Logging.Level), zap.String("newEncoding", config.GetConfig().Logging.Encoding))

			_, err := logger.InitLogger(config.GetConfig().Logging.Level, config.GetConfig().Logging.Encoding)

			if err != nil {
				log.Error("Failed to change logging config", zap.String("oldLevel", oldLogLevel), zap.String("oldEncoding", oldLogEncoding), zap.String("newLevel", config.GetConfig().Logging.Level), zap.String("newEncoding", config.GetConfig().Logging.Encoding), zap.Error(err))
			} else {
				log.Info("Logging configuration changed")
			}
		}

		// call `onChange` function
		onChange(config.GetConfig())
	}); wErr != nil {
		return wErr
	}

	return nil
}

func loadConfig() error {
	err := k.Unmarshal("", config.GetConfig())

	return err
}
