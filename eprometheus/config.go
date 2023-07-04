package eprometheus

type Config struct {
	Path string
	Port string
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Path: "/metrics",
		Port: ":2112",
	}
}
