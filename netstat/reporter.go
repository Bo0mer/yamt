package netstat

import (
	"log"
	"regexp"
	"time"

	"github.com/bo0mer/yamt/internal"
	"github.com/bo0mer/yamt/metric/riemann"
)

type state map[string]IfStat

// Option configures reporter.
type Option func(r *Reporter)

// Interval configures the interval between reports.
func Interval(d time.Duration) Option {
	return func(r *Reporter) {
		r.interval = d
	}
}

// Except configures which interfaces to skip.
func Except(re *regexp.Regexp) Option {
	return func(r *Reporter) {
		r.except = re
	}
}

// Reporter reads network interface statistics and emits them via specified
// emitter.
type Reporter struct {
	emitter  *riemann.Emitter
	interval time.Duration
	except   *regexp.Regexp
	last     state
	stop     chan struct{}
}

// NewReporter creates brand new reporter.
func NewReporter(emitter *riemann.Emitter, opts ...Option) *Reporter {
	r := &Reporter{
		emitter:  emitter,
		interval: time.Second * 1,
		last:     make(map[string]IfStat),
		stop:     make(chan struct{}),
	}

	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Starts starts to read and emit metrics. It does so in its own goroutine.
func (r *Reporter) Start() {
	go r.readAndReport()
}

// readAndReport reads and reports statistics for all interfaces, except
// one that are matched against except regexp.
func (r *Reporter) readAndReport() {
	ticker := time.NewTicker(r.interval)
	for {
		select {
		case <-ticker.C:
			stats, err := ReadIfStats()
			if err != nil {
				log.Printf("reporter: error reading stats: %v\n", err)
				continue
			}
			if err := r.report(stats); err != nil {
				log.Printf("reporter: error reporting stats: %v\n", err)
			}
		case <-r.stop:
			ticker.Stop()
			return
		}
	}
}

func (r *Reporter) report(stats []IfStat) error {
	rate := internal.RateComputer(r.interval.Seconds())
	for _, stat := range stats {
		if r.except != nil && r.except.MatchString(stat.Name) {
			continue
		}
		last, ok := r.last[stat.Name]
		r.last[stat.Name] = stat
		if !ok {
			continue
		}
		event := func(name string, value float64) riemann.Event {
			return riemann.Event{
				Name:  stat.Name + "." + name,
				Value: value,
			}
		}
		r.emitter.Emit(event("rx.bytes", rate(stat.RxBytes, last.RxBytes)))
		r.emitter.Emit(event("rx.packets", rate(stat.RxPackets, last.RxPackets)))
		r.emitter.Emit(event("rx.errs", rate(stat.RxErrs, last.RxErrs)))
		r.emitter.Emit(event("rx.drop", rate(stat.RxDrop, last.RxDrop)))
		r.emitter.Emit(event("rx.fifo", rate(stat.RxFIFO, last.RxFIFO)))
		r.emitter.Emit(event("rx.frame", rate(stat.RxFrame, last.RxFrame)))
		r.emitter.Emit(event("rx.compressed", rate(stat.RxCompressed, last.RxCompressed)))
		r.emitter.Emit(event("rx.multicast", rate(stat.RxMulticast, last.RxMulticast)))

		r.emitter.Emit(event("tx.bytes", rate(stat.TxBytes, last.TxBytes)))
		r.emitter.Emit(event("tx.packets", rate(stat.TxPackets, last.TxPackets)))
		r.emitter.Emit(event("tx.errs", rate(stat.TxErrs, last.TxErrs)))
		r.emitter.Emit(event("tx.drop", rate(stat.TxDrop, last.TxDrop)))
		r.emitter.Emit(event("tx.fifo", rate(stat.TxFIFO, last.TxFIFO)))
		r.emitter.Emit(event("tx.colls", rate(stat.TxColls, last.TxColls)))
		r.emitter.Emit(event("tx.carrier", rate(stat.TxCarrier, last.TxCarrier)))
		r.emitter.Emit(event("tx.compressed", rate(stat.TxCompressed, last.TxCompressed)))

		if err := r.emitter.Err(); err != nil {
			return err
		}
	}
	return nil
}

// Close releases all resources allocated by the reporter.
func (r *Reporter) Close() {
	close(r.stop)
}
