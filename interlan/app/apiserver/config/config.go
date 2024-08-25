package config

// Config ...
type Config struct {
	PORT        string `toml:"bind_addr"`
	DatabaseURL string `toml:"database_url"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}
