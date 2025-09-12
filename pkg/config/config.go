package cfg

type StreamlyConfig struct{
	Service struct {
		Name string `toml:"name"`
		Version string `toml:"version"`
	} `toml:"service"`

	Storage struct{
		Clickhouse struct {} `toml:"clickhouse"`
		Database struct {} `toml:"database"`
	} `toml:"storage"`

	Redis struct{
		Addr string `toml:"addr"`
		Password string `toml:"password"`
	} `toml:"redis"`

	Logging struct {
		Level string `toml:"level"`
		Format string `toml:"format"`
	} `toml:"logging"`
}
