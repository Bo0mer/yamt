package netstat

import (
	"fmt"
	"log"
	"regexp"
	"time"

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
	for _, stat := range stats {
		if r.except != nil && r.except.MatchString(stat.Name) {
			continue
		}
		last, ok := r.last[stat.Name]
		r.last[stat.Name] = stat
		if !ok {
			continue
		}
		event := eventBuilder(stat.Name, r.interval)
		r.emitter.Emit(event(stat.RxBytes-last.RxBytes, "rx.bytes"))
		r.emitter.Emit(event(stat.RxPackets-last.RxPackets, "rx.packets"))
		r.emitter.Emit(event(stat.RxErrs-last.RxErrs, "rx.errs"))
		r.emitter.Emit(event(stat.RxDrop-last.RxDrop, "rx.drop"))
		r.emitter.Emit(event(stat.RxFIFO-last.RxFIFO, "rx.fifo"))
		r.emitter.Emit(event(stat.RxFrame-last.RxFrame, "rx.frame"))
		r.emitter.Emit(event(stat.RxCompressed-last.RxCompressed, "rx.compressed"))
		r.emitter.Emit(event(stat.RxMulticast-last.RxMulticast, "rx.multicast"))

		r.emitter.Emit(event(stat.TxBytes-last.TxBytes, "tx.bytes"))
		r.emitter.Emit(event(stat.TxPackets-last.TxPackets, "tx.packets"))
		r.emitter.Emit(event(stat.TxErrs-last.TxErrs, "tx.errs"))
		r.emitter.Emit(event(stat.TxDrop-last.TxDrop, "tx.drop"))
		r.emitter.Emit(event(stat.TxFIFO-last.TxFIFO, "tx.fifo"))
		r.emitter.Emit(event(stat.TxColls-last.TxColls, "tx.colls"))
		r.emitter.Emit(event(stat.TxCarrier-last.TxCarrier, "tx.carrier"))
		r.emitter.Emit(event(stat.TxCompressed-last.TxCompressed, "tx.compressed"))
	}
	return r.emitter.Err()
}

// Close releases all resources allocated by the reporter.
func (r *Reporter) Close() {
	close(r.stop)
}

func eventBuilder(ifName string, interval time.Duration) func(uint64, string) riemann.Event {
	return func(value uint64, metricName string) riemann.Event {
		return riemann.Event{
			Name:  fmt.Sprintf("%s.%s", ifName, metricName),
			Value: int(value) / int(interval/time.Second),
		}
	}
}
