package config

const (
	DEFAULT_LOGGING_LEVEL    = "info"
	DEFAULT_LOGGING_ENCODING = "console"

	DEFAULT_LISTEN_HOST     = "0.0.0.0"
	DEFAULT_LISTEN_PORT     = 80
	DEFAULT_LISTEN_TLS_PORT = 443

	DEFAULT_TLS_CIPHERS = "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
)

type Config struct {
	Rules   []map[string]interface{}
	TLS     TLSConfig
	Logging LoggingConfig
	Listen  ListenConfig
}

type TLSConfig struct {
	Enabled bool
	Cert    string
	Key     string
	Ciphers string
}

type LoggingConfig struct {
	Level    string
	Encoding string
}

type ListenConfig struct {
	Host    string
	Port    int
	TlsPort int
}

var conf Config = Config{
	Rules: nil,
	TLS: TLSConfig{
		Ciphers: DEFAULT_TLS_CIPHERS,
	},
	Logging: LoggingConfig{
		Level:    DEFAULT_LOGGING_LEVEL,
		Encoding: DEFAULT_LOGGING_ENCODING,
	},
	Listen: ListenConfig{
		Host:    DEFAULT_LISTEN_HOST,
		Port:    DEFAULT_LISTEN_PORT,
		TlsPort: DEFAULT_LISTEN_TLS_PORT,
	},
}

func GetConfig() *Config {
	return &conf
}
