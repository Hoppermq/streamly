package config

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/zixyos/goloader/config"

	"github.com/hoppermq/streamly/conf"
)

type PlatformConfig struct {
	Platform struct {
		Service struct {
			Name    string `toml:"name"`
			Version string `toml:"version"`
		}

		HTTP struct {
			Port         int           `toml:"port"`
			ReadTimeout  time.Duration `toml:"read_timeout"`
			WriteTimeout time.Duration `toml:"write_timeout"`
		} `toml:"http"`

		Zitadel struct {
			Port              uint16 `toml:"port"`
			Domain            string `toml:"domain"`
			PatPath           string `toml:"patpath"`           // Path to PAT file (v0: file path, prod: empty if using env var)
			ServiceAccountKey string `toml:"serviceaccountkey"` // Path to service account JSON keyfile for token introspection
		} `toml:"zitadel"`

		Storage struct {
			Database struct {
				Host         string `toml:"host"`
				Port         int    `toml:"port"`
				User         string `toml:"user"`
				Password     string `toml:"password"`
				Name         string `toml:"name"`
				ReadTimeout  int    `toml:"read_timeout"`
				WriteTimeout int    `toml:"write_timeout"`
			} `toml:"database"`
		} `toml:"storage"`
	} `toml:"platform"`
}

func LoadPlatformConfig() (*PlatformConfig, error) {
	var platformConfig PlatformConfig
	err := config.Load(&platformConfig, config.WithFs(conf.FileFS))

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &platformConfig, nil
}

func (c *PlatformConfig) DatabaseDSN() string {
	hostPort := net.JoinHostPort(
		c.Platform.Storage.Database.Host,
		strconv.Itoa(c.Platform.Storage.Database.Port),
	)

	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		c.Platform.Storage.Database.User,
		c.Platform.Storage.Database.Password,
		hostPort,
		c.Platform.Storage.Database.Name,
	)
}

// ZitadelPATPath returns the PAT path (prefer env var, fallback to config).
func (c *PlatformConfig) ZitadelPATPath() string {
	return c.Platform.Zitadel.PatPath
}

// ZitadelServiceAccountKeyPath returns the service account keyfile path.
func (c *PlatformConfig) ZitadelServiceAccountKeyPath() string {
	return c.Platform.Zitadel.ServiceAccountKey
}
