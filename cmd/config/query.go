package config

import (
	"fmt"
	"time"

	"github.com/zixyos/goloader/config"

	"github.com/hoppermq/streamly/conf"
)

type QueryConfig struct {
	Query struct {
		Service struct {
			Name    string `toml:"name"`
			Version string `toml:"version"`
		} `toml:"service"`

		HTTP struct {
			Port         int           `toml:"port"`
			ReadTimeout  time.Duration `toml:"read_timeout"`
			WriteTimeout time.Duration `toml:"write_timeout"`
		} `toml:"http"`

		Storage struct {
			Clickhouse struct {
				Address      string        `toml:"address"`
				Port         string        `toml:"port"`
				UserName     string        `toml:"username"`
				Password     string        `toml:"password"`
				Database     string        `toml:"database"`
				ReadTimeout  time.Duration `toml:"read_timeout"`
				WriteTimeout time.Duration `toml:"write_timeout"`
			} `toml:"clickhouse"`
			Database struct{} `toml:"database"`
		} `toml:"storage"`

		Redis struct {
			Addr     string `toml:"addr"`
			Password string `toml:"password"`
		} `toml:"redis"`

		Logging struct {
			Level  string `toml:"level"`
			Format string `toml:"format"`
		} `toml:"logging"`
	} `toml:"query"`
}

func LoadQueryConfig() (*QueryConfig, error) {
	var queryConfig QueryConfig
	err := config.Load(&queryConfig, config.WithFs(conf.FileFS))

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &queryConfig, nil
}
