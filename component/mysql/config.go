package mysql

type Config struct {
	// 最大空闲连接数
	MaxIdleConnections int

	// 最大打开连接数
	MaxOpenConnections int

	// 最长活跃时间
	MaxLifeTime int

	// dsn
	Dsn string
}

//
type ResolverItem struct {
	Sources  []string //
	Replicas []string //
	Tables   []string //
}
