package logger

import (
	"encoding/json"
	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger(logLevel string, encoding string) (*zap.Logger, error) {
	additionalEncoderConfig := func(encoding string) string {
		if encoding != "console" {
			return `,
				"timeKey": "time",
				"timeEncoder": "iso8601",
				"callerKey": "caller",
				"callerEncoder": "short"`
		}

		return ""
	}

	rawJSON := []byte(`{
	  "level": "` + logLevel + `",
	  "encoding": "` + encoding + `",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"` + additionalEncoderConfig(encoding) + `
	  }
	}`)

	var cfg zap.Config
	var err error

	if err = json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, err
	}

	logger, err = cfg.Build()

	return logger, err
}

func GetLogger() *zap.Logger {
	return logger
}
