package config

import (
	"github.com/ebar-go/egu"
	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
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


// LoadFile 加载配置文件
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
func (conf *Config) GetDefaultInt(key string, dn int) int {
	return egu.DefaultInt(conf.GetInt(key), dn)
}

// GetDefaultString get string with default value
func (conf *Config) GetDefaultString(key string, ds string) string{
	return egu.DefaultString(conf.GetString(key), ds)
}

const (
	envProduct = "product"
	envDevelop = "develop"
)

// IsProduct 是否为生产环境
func (conf *Config) IsProduct() bool {
	return envProduct == conf.Environment
}

// IsDevelop 是否为测试环境
func (conf *Config) IsDevelop() bool {
	return envDevelop == conf.Environment
}
