package ego

import (
	"github.com/ebar-go/ego/utils/runtime"
)

// Aggregator define engine with name.
type Aggregator struct {
	runners []runtime.Runnable
	name    string

	done chan struct{}
}

// NewAggregator creates a new named aggregator.
func NewAggregator(name string) *Aggregator {
	return &Aggregator{
		name:    name,
		runners: make([]runtime.Runnable, 0),
		done:    make(chan struct{}),
	}
}

// Aggregate aggregates some runner
func (engine *Aggregator) Aggregate(runner ...runtime.Runnable) *Aggregator {
	engine.runners = append(engine.runners, runner...)
	return engine
}

// Run runs the engine with blocking mode.
func (engine *Aggregator) Run(stopCh <-chan struct{}) {
	for _, runner := range engine.runners {
		go runner.Run(engine.done)
	}

	runtime.WaitClose(stopCh, func() {
		close(engine.done)
	})
}
