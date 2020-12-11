// Adapter 实现读写分离的gorm适配器
// 通过实现gorm里的SQLCommon接口实现读写分离的适配
// TODO 读库可能存在多个，可以通过基于权重的副在均衡算法完成连接工作的分配
package mysql

import (
	"context"
	"database/sql"
	"github.com/zutim/egu"
	"time"
)

// ReadWriteAdapter 读写分离适配器
type ReadWriteAdapter struct {
	// Physical databases
	pdbs []*sql.DB

	// Monotonically incrementing counter on each query
	count uint64
}

// Exec
func (adapter ReadWriteAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return adapter.Master().Exec(query, args...)
}

// Prepare
func (adapter ReadWriteAdapter) Prepare(query string) (*sql.Stmt, error) {
	return adapter.Slave().Prepare(query)
}

// Query
func (adapter ReadWriteAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return adapter.Slave().Query(query, args...)
}

// QueryRow
func (adapter ReadWriteAdapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return adapter.Slave().QueryRow(query, args...)
}

// BeginTx starts a transaction with the provided context on the master.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support,
// an error will be returned.
func (adapter ReadWriteAdapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return adapter.Master().BeginTx(ctx, opts)
}

// Begin starts a transaction on the master. The isolation level is dependent on the driver.
func (adapter ReadWriteAdapter) Begin() (*sql.Tx, error) {
	return adapter.Master().Begin()
}

// Master 获取master连接
func (adapter ReadWriteAdapter) Master() *sql.DB {
	if adapter.count == 0 {
		return nil
	}

	return adapter.pdbs[0]
}

// Slave 获取slave链接
func (adapter ReadWriteAdapter) Slave() *sql.DB {
	if adapter.count <= 1 {
		return adapter.Master()
	}

	index := 1
	if adapter.count > 2 {
		index = egu.RandInt(1, int(adapter.count))
	}
	return adapter.pdbs[index]
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool for each underlying physical db.
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns then the
// new MaxIdleConns will be reduced to match the MaxOpenConns limit
// If n <= 0, no idle connections are retained.
func (adapter ReadWriteAdapter) SetMaxIdleConns(n int) {
	for i := range adapter.pdbs {
		adapter.pdbs[i].SetMaxIdleConns(n)
	}
}

// SetMaxOpenConns sets the maximum number of open connections
// to each physical database.
// If MaxIdleConns is greater than 0 and the new MaxOpenConns
// is less than MaxIdleConns, then MaxIdleConns will be reduced to match
// the new MaxOpenConns limit. If n <= 0, then there is no limit on the number
// of open connections. The default is 0 (unlimited).
func (adapter ReadWriteAdapter) SetMaxOpenConns(n int) {
	for i := range adapter.pdbs {
		adapter.pdbs[i].SetMaxOpenConns(n)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// Expired connections may be closed lazily before reuse.
// If d <= 0, connections are reused forever.
func (adapter ReadWriteAdapter) SetConnMaxLifetime(d time.Duration) {
	for i := range adapter.pdbs {
		adapter.pdbs[i].SetConnMaxLifetime(d)
	}
}

// Ping verifies if a connection to each physical database is still alive,
// establishing a connection if necessary.
func (adapter ReadWriteAdapter) Ping() error {
	return adapter.scatter(int(adapter.count), func(i int) error {
		return adapter.pdbs[i].Ping()
	})
}

// scatter 通过协程与通道实现异步执行
func (adapter ReadWriteAdapter) scatter(n int, fn func(i int) error) error {
	errors := make(chan error, n)

	var i int
	for i = 0; i < n; i++ {
		go func(i int) { errors <- fn(i) }(i)
	}

	var err, innerErr error
	for i = 0; i < cap(errors); i++ {
		if innerErr = <-errors; innerErr != nil {
			err = innerErr
		}
	}

	return err
}

// NewReadWriteAdapter 通过多个dsn打开多个实例的连接
// Open concurrently opens each underlying physical db.
// dataSourceNames must be a semi-comma separated list of DSNs with the first
// one being used as the master and the rest as slaves.
func NewReadWriteAdapter(dialectType string, dataSourceNames []string) (*ReadWriteAdapter, error) {
	length := len(dataSourceNames)
	adapter := &ReadWriteAdapter{pdbs: make([]*sql.DB, length), count: uint64(length)}

	var err error
	for index, item := range dataSourceNames {
		adapter.pdbs[index], err = sql.Open(dialectType, item)

		if err != nil {
			return nil, err
		}
	}

	return adapter, nil
}
