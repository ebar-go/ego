package prometheus

// code from github.com/vearne/golib/metric
import (
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

// MySQLCollector Mysql收集器
type MySQLCollector struct {
	Clients    map[string]*gorm.DB
	lock       sync.Mutex
	registered bool
}

var mySQLCollector *MySQLCollector

func init() {
	mySQLCollector = &MySQLCollector{
		Clients: make(map[string]*gorm.DB),
	}
}

// 监听MySql
func ListenMysql(client *gorm.DB, role string) {
	if client == nil {
		return
	}

	client.DB().Stats()

	mySQLCollector.lock.Lock()
	defer mySQLCollector.lock.Unlock()

	if !mySQLCollector.registered {
		mySQLCollector.register()
	}

	mySQLCollector.Clients[role] = client
}

type collector struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType

	getter getterFn
}

type getterFn func(ch chan<- prometheus.Metric, desc *prometheus.Desc, typ prometheus.ValueType)

func newCollector(
	name, help string,
	valueType prometheus.ValueType,
	labels []string,
	getter getterFn) *collector {

	desc := prometheus.NewDesc(name, help, labels, nil)
	return &collector{
		desc:      desc,
		valueType: valueType,
		getter:    getter,
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.getter(ch, c.desc, c.valueType)
}

func (mc *MySQLCollector) register() {
	prometheus.MustRegister(newCollector(
		"mysql_pool_state",
		"MySQL pool state",
		prometheus.GaugeValue,
		[]string{"role", "state"},
		mc.stat,
	))

	prometheus.MustRegister(newCollector(
		"mysql_pool_fetches_total",
		"MySQL pool fetches total",
		prometheus.CounterValue,
		[]string{"role", "state"},
		mc.fetches,
	))
	mc.registered = true
}

func (mc *MySQLCollector) stat(ch chan<- prometheus.Metric, desc *prometheus.Desc, typ prometheus.ValueType) {
	for role, client := range mc.Clients {
		st := client.DB().Stats()

		ch <- prometheus.MustNewConstMetric(
			desc,
			typ,
			float64(st.Idle),
			role, "idle",
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			typ,
			float64(st.InUse),
			role, "active",
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			typ,
			float64(st.OpenConnections),
			role, "open",
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			typ,
			float64(st.MaxOpenConnections),
			role, "poolsize",
		)
	}
}

func (mc *MySQLCollector) fetches(ch chan<- prometheus.Metric, desc *prometheus.Desc, typ prometheus.ValueType) {
	for role, client := range mc.Clients {
		st := client.DB().Stats()

		ch <- prometheus.MustNewConstMetric(
			desc,
			typ,
			float64(st.WaitCount),
			role, "wait_count",
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			typ,
			float64(st.WaitDuration),
			role, "wait_duration",
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			typ,
			float64(st.MaxIdleClosed),
			role, "max_idle_closed",
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			typ,
			float64(st.MaxLifetimeClosed),
			role, "max_life_closed",
		)
	}
}
