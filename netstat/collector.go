package netstat

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Bo0mer/yamt/internal"
	"github.com/Bo0mer/yamt/metric"
)

type state map[string]IfStat

// IfStatCollector computes metrics for network interfaces.
type IfStatCollector struct {
	reader   InterfaceStatReader
	except   *regexp.Regexp
	last     state
	lastTime time.Time
}

// NewIfStatCollector returns brand new interface stats collector.
func NewIfStatCollector(reader InterfaceStatReader, except *regexp.Regexp) (*IfStatCollector, error) {
	c := &IfStatCollector{
		reader: reader,
		except: except,
	}
	if err := c.init(); err != nil {
		return nil, err
	}
	return c, nil
}

// Collect collects stats and creates events for network interfaces.
func (c *IfStatCollector) Collect() ([]metric.Event, error) {
	actual, err := c.getState()
	if err != nil {
		return nil, err
	}

	actualTime := time.Now()
	interval := actualTime.Sub(c.lastTime).Seconds()

	events := make([]metric.Event, 0)

	for _, stat := range actual {
		if c.except != nil && c.except.MatchString(stat.Name) {
			continue
		}
		last, ok := c.last[stat.Name]
		if !ok {
			continue
		}

		events = append(events, c.buildEvents(stat, last, interval)...)
	}

	c.last = actual
	c.lastTime = actualTime

	return events, nil
}

// init loads the initial state of the collector.
func (c *IfStatCollector) init() error {
	state, err := c.getState()
	if err != nil {
		return err
	}
	c.last = state
	c.lastTime = time.Now()
	return nil
}

// getState reads current state for all network interfaces.
func (c *IfStatCollector) getState() (state, error) {
	state := make(map[string]IfStat)
	stats, err := c.reader.ReadStats()
	if err != nil {
		return nil, fmt.Errorf("collector: error reading stats: %v", err)
	}
	for _, stat := range stats {
		state[stat.Name] = stat
	}
	return state, nil
}

// buildEvents build all events for a single network interface.
func (c *IfStatCollector) buildEvents(actual, last IfStat, interval float64) []metric.Event {
	events := make([]metric.Event, 0)
	event := eventBuilder(actual.Name)
	rate := internal.RateComputer(interval)

	events = append(events, event("rx bytes", rate(actual.RxBytes, last.RxBytes)))
	events = append(events, event("rx packets", rate(actual.RxPackets, last.RxPackets)))
	events = append(events, event("rx errs", rate(actual.RxErrs, last.RxErrs)))
	events = append(events, event("rx drop", rate(actual.RxDrop, last.RxDrop)))
	events = append(events, event("rx fifo", rate(actual.RxFIFO, last.RxFIFO)))
	events = append(events, event("rx frame", rate(actual.RxFrame, last.RxFrame)))
	events = append(events, event("rx compressed", rate(actual.RxCompressed, last.RxCompressed)))
	events = append(events, event("rx multicast", rate(actual.RxMulticast, last.RxMulticast)))

	events = append(events, event("tx bytes", rate(actual.TxBytes, last.TxBytes)))
	events = append(events, event("tx packets", rate(actual.TxPackets, last.TxPackets)))
	events = append(events, event("tx errs", rate(actual.TxErrs, last.TxErrs)))
	events = append(events, event("tx drop", rate(actual.TxDrop, last.TxDrop)))
	events = append(events, event("tx fifo", rate(actual.TxFIFO, last.TxFIFO)))
	events = append(events, event("tx colls", rate(actual.TxColls, last.TxColls)))
	events = append(events, event("tx carrier", rate(actual.TxCarrier, last.TxCarrier)))
	events = append(events, event("tx compressed", rate(actual.TxCompressed, last.TxCompressed)))

	return events
}

func eventBuilder(ifName string) func(string, float64) metric.Event {
	return func(name string, value float64) metric.Event {
		return metric.Event{
			Name:  ifName + " " + name,
			Value: value,
		}
	}
}
