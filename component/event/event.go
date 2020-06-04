package event

// Event
type Event struct {
	// name
	Name string
	// event params
	Params interface{}
}

const (
	Sync  = 0
	Async = 1
)

// Listener
type Listener struct {
	Mode    int
	Handler Handler
}

// Handler process event
type Handler func(ev Event)

const (
	BeforeHttpStart      = "http.start.before"
	AfterHttpStart       = "http.start.after"
	BeforeHttpShutdown   = "http.shutdown.before"
	AfterDatabaseConnect = "database.connect.after"
	BeforeRoute          = "route.before"
	AfterRoute           = "route.after"
)
