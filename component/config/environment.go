package config

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
