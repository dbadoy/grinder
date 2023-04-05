package server

type Config struct {
	AllowProxyContract bool
}

func (c *Config) validate() error {
	return nil
}
