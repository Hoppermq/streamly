package config

type AuthConfig struct {
	Auth struct {
		Service struct {
			Name    string `toml:"name"`
			Version string `toml:"version"`
		} `toml:"service"`

		HTTP struct {
			Port         int `toml:"port"`
			ReadTimeout  int `toml:"read_timeout"`
			WriteTimeout int `toml:"write_timeout"`
		} `toml:"http"`

		Storage struct {
			Database struct{} `toml:"database"`
			Cache    struct{} `toml:"cache"`
		}

		Logging struct {
			Level  string `toml:"level"`
			Format string `toml:"format"`
		} `toml:"logging"`
	} `toml:"auth"`
}
