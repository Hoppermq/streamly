package config

import (
	"fmt"
	"time"

	"github.com/hoppermq/streamly/conf"
	"github.com/zixyos/goloader/config"
)

type AuthConfig struct {
	Auth struct {
		Service struct {
			Name    string `toml:"name"`
			Version string `toml:"version"`
		} `toml:"service"`

		HTTP struct {
			Port         int           `toml:"port"`
			ReadTimeout  time.Duration `toml:"read_timeout"`
			WriteTimeout time.Duration `toml:"write_timeout"`
		} `toml:"http"`

		// TODO: Add TLS Configuration for storage conn.
		Storage struct {
			Database struct {
				Host     string `toml:"host"`
				Port     int    `toml:"port"`
				User     string `toml:"user"`
				Password string `toml:"password"`
				Database string `toml:"database"`
			} `toml:"database"`

			Cache struct {
				Host     string `toml:"host"`
				Port     int    `toml:"port"`
				User     string `toml:"user"`
				Password string `toml:"password"`
				Database string `toml:"database"`
			} `toml:"cache"`
		}

		Transport struct {
			Host     string `toml:"host"`
			Port     int    `toml:"port"`
			User     string `toml:"user"`
			Password string `toml:"password"`
		} `toml:"transport"`

		Logging struct {
			Level  string `toml:"level"`
			Format string `toml:"format"`
		} `toml:"logging"`
	} `toml:"auth"`
}

func LoadAuthConfig() (*AuthConfig, error) {
	var authConfig AuthConfig
	err := config.Load(&authConfig, config.WithFs(conf.FileFS))
	if err != nil {
		return nil, fmt.Errorf("load auth config: %w", err)
	}

	return &authConfig, nil
}
