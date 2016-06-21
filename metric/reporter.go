package metric

import (
	"log"
	"time"
)

//go:generate counterfeiter . Collector
//go:generate counterfeiter . Emitter

// Event repesents generic metric event.
type Event struct {
	Name string
	// could be int, float32 or float64
	Value interface{}
}

// Collector collects metric events.
type Collector interface {
	Collect() ([]Event, error)
}

// Emitter emits specified events.
type Emitter interface {
	// Emit should try to emit the specified event. If error occurrs, it
	// should be persisted and no new emits should happen until the error is
	// reset by call to Err.
	Emit(Event) error
	// Err reads and resets the lat error occurred during emit.
	Err() error
}

type Option func(*Reporter)

func Interval(d time.Duration) Option {
	return func(r *Reporter) {
		r.interval = d
	}
}

// Reporter periodically collects and emits metrics.
type Reporter struct {
	emitter    Emitter
	collectors []Collector

	interval time.Duration
	stop     chan struct{}
}

// NewReporter returns brand new reporter.
func NewReporter(e Emitter, collectors []Collector, opts ...Option) *Reporter {
	r := &Reporter{
		emitter:    e,
		collectors: collectors,
		interval:   time.Second,
		stop:       make(chan struct{}),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Start starts collecting and emitting metric events.
// It does so in its own goroutine. See Close for stopping.
func (r *Reporter) Start() {
	go r.start()
}

func (r *Reporter) start() {
	t := time.NewTicker(r.interval)
	for {
		select {
		case <-t.C:
			for _, c := range r.collectors {
				r.collectAndEmit(c)
			}
		case <-r.stop:
			return
		}
	}
}

func (r *Reporter) collectAndEmit(c Collector) {
	events, err := c.Collect()
	if err != nil {
		log.Printf("reporter: error collecting metrics: %v\n", err)
		return
	}
	for _, event := range events {
		if err := r.emitter.Emit(event); err != nil {
			log.Printf("reporter: error emitting metric: %v\n", err)
			continue
		}
	}
}

// Close releases all resources allocated by the reporter.
func (r *Reporter) Close() {
	close(r.stop)
}
