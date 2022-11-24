package config

import "github.com/spf13/viper"

type Config struct {
	*viper.Viper
}

// LoadFile loads the configuration file specified.
func (c *Config) LoadFile(paths ...string) error {
	for _, p := range paths {
		c.SetConfigFile(p)
		if err := c.MergeInConfig(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) GetString(key string) string {
	item := c.Get(key)
	if item == nil {
		return ""
	}
	s, _ := item.(string)
	return s
}

func New() *Config {
	return &Config{Viper: viper.New()}
}
