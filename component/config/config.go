package config

import (
	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
	// viper
	*viper.Viper
	// 运行环境
	Environment string
	// 应用名称
	Name string
}

// New 实例
func New() *Config {
	conf := &Config{Name: "app", Environment: envDevelop}
	conf.Viper = viper.New()
	return conf
}

// LoadFile 加载配置文件,支持多个配置文件
func (conf *Config) LoadFile(path ...string) error {
	for _, p := range path {
		conf.SetConfigFile(p)
		if err := conf.MergeInConfig(); err != nil {
			return err
		}
	}

	conf.Environment = conf.GetString("server.name")
	conf.Name = conf.GetString("server.environment")

	return nil
}

// GetDefaultInt get int with default value
func (conf *Config) GetDefaultInt(key string, defaultVal int) int {
	val := conf.GetInt(key)
	if val == 0 {
		return defaultVal
	}
	return val
}

// GetDefaultString get string with default value
func (conf *Config) GetDefaultString(key string, defaultVal string) string {
	val := conf.GetString(key)
	if val == "" {
		return defaultVal
	}
	return val
}
