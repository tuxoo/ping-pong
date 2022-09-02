package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	defaultHttpPort           = "8888"
	defaultHttpRWTimeout      = 10 * time.Second
	defaultMaxHeaderMegabytes = 1
)

type (
	Config struct {
		ServerConfig HTTPServerConfig
		ClientConfig HTTPClientConfig
	}

	HTTPServerConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"read"`
		WriteTimeout       time.Duration `mapstructure:"write"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderMegabytes"`
	}

	HTTPClientConfig struct {
		Timeout time.Duration `mapstructure:"timeout"`
	}
)

func InitConfig(path string) (*Config, error) {
	viper.AutomaticEnv()
	preDefaults()

	if err := parseConfigFile(path); err != nil {
		return nil, err
	}

	if err := parseEnv(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshalConfig(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func preDefaults() {
	viper.SetDefault("http.server.port", defaultHttpPort)
	viper.SetDefault("http.server.max_header_megabytes", defaultMaxHeaderMegabytes)
	viper.SetDefault("http.server.timeouts.read", defaultHttpRWTimeout)
	viper.SetDefault("http.server.timeouts.write", defaultHttpRWTimeout)
}
func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0]) // folder
	viper.SetConfigName(path[1]) // config file name

	return viper.ReadInConfig()
}

func parseEnv() error {
	if err := parseServerEnv(); err != nil {
		return err
	}

	return parseClientEnv()
}

func parseServerEnv() error {
	if err := viper.BindEnv("http.server.host", "HTTP_SERVER_HOST"); err != nil {
		return err
	}

	return viper.BindEnv("http.server.port", "HTTP_SERVER_PORT")
}

func parseClientEnv() error {
	return viper.BindEnv("http.client.timeout", "HTTP_CLIENT_TIMEOUT")
}

func unmarshalConfig(cfg *Config) error {
	if err := viper.UnmarshalKey("http.server", &cfg.ServerConfig); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http.server.timeouts", &cfg.ServerConfig); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http.client", &cfg.ClientConfig); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) {
	cfg.ServerConfig.Host = viper.GetString("http.server.host")
	cfg.ServerConfig.Port = viper.GetString("http.server.port")

	cfg.ClientConfig.Timeout = viper.GetDuration("http.client.timeout")
}
