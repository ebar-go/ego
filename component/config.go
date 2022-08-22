package component

import "github.com/spf13/viper"

type Config struct {
	Named
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

func NewConfig() *Config {
	return &Config{Named: componentConfig, Viper: viper.New()}
}
