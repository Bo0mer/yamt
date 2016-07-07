package iostat

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bo0mer/yamt/internal"
	"github.com/bo0mer/yamt/metric"
)

type state map[string]DeviceStat

type DeviceStatCollector struct {
	reader   DeviceStatReader
	except   *regexp.Regexp
	last     state
	lastTime time.Time
}

func NewDeviceStatCollector(r DeviceStatReader, except *regexp.Regexp) (*DeviceStatCollector, error) {
	c := &DeviceStatCollector{
		reader: r,
		except: except,
	}
	if err := c.init(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *DeviceStatCollector) Collect() ([]metric.Event, error) {
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

func (c *DeviceStatCollector) init() error {
	state, err := c.getState()
	if err != nil {
		return err
	}
	c.last = state
	c.lastTime = time.Now()
	return nil
}

func (c *DeviceStatCollector) getState() (state, error) {
	state := make(map[string]DeviceStat)
	stats, err := c.reader.ReadStats()
	if err != nil {
		return nil, fmt.Errorf("collector: error reading stats: %v", err)
	}
	for _, stat := range stats {
		state[stat.Name] = stat
	}
	return state, nil
}

func (c *DeviceStatCollector) buildEvents(actual, last DeviceStat, interval float64) []metric.Event {
	events := make([]metric.Event, 0)
	event := eventBuilder(actual.Name)
	rate := internal.RateComputer(interval)

	events = append(events, event("reads total", rate(actual.Reads, last.Reads)))
	events = append(events, event("reads merged", rate(actual.ReadsMerged, last.ReadsMerged)))
	events = append(events, event("reads sectors", rate(actual.ReadsSectors, last.ReadsSectors)))
	events = append(events, event("reads time(ms)", rate(actual.ReadsTimeMs, last.ReadsTimeMs)))

	events = append(events, event("writes total", rate(actual.Writes, last.Writes)))
	events = append(events, event("writes merged", rate(actual.WritesMerged, last.WritesMerged)))
	events = append(events, event("writes sectors", rate(actual.WritesSectors, last.WritesSectors)))
	events = append(events, event("writes time(ms)", rate(actual.WritesTimeMs, last.WritesTimeMs)))

	events = append(events, event("io inflight", float64(actual.InFlight)))
	events = append(events, event("io time(ms)", rate(actual.IOTimeMs, last.IOTimeMs)))
	events = append(events, event("io weighted(ms)", rate(actual.WeightedIOTimeMS, last.WeightedIOTimeMS)))

	return events
}

func eventBuilder(devName string) func(string, float64) metric.Event {
	return func(name string, value float64) metric.Event {
		return metric.Event{
			Name:  devName + " " + name,
			Value: value,
		}
	}
}
