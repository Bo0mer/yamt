package metric_test

import (
	"testing"
	"time"

	"github.com/Bo0mer/yamt/metric"
	"github.com/Bo0mer/yamt/metric/metricfakes"
)

func TestReporter(t *testing.T) {
	emitter := new(metricfakes.FakeEmitter)
	want1 := metric.Event{
		Name:  "c1",
		Value: 42.0,
	}
	want2 := metric.Event{
		Name:  "c2",
		Value: -42.0,
	}
	c1, c2 := new(metricfakes.FakeCollector), new(metricfakes.FakeCollector)

	c1.CollectReturns([]metric.Event{want1}, nil)
	c2.CollectReturns([]metric.Event{want2}, nil)

	collectors := []metric.Collector{c1, c2}
	interval := time.Millisecond * 20
	r := metric.NewReporter(emitter, collectors,
		metric.Interval(interval))

	r.Start()
	defer r.Close()

	timeout := time.After(interval * 3)
	for {
		select {
		case <-timeout:
			t.Error("expected two calls to emitter, none received")
			return
		default:
			if emitter.EmitCallCount() >= 2 {
				if got1 := emitter.EmitArgsForCall(0); got1 != want1 {
					t.Errorf("expected call to emitter with %v, got %v\n", want1, got1)
				}
				if got2 := emitter.EmitArgsForCall(1); got2 != want2 {
					t.Errorf("expected call to emitter with %v, got %v\n", want2, got2)
				}
				return
			}
			time.Sleep(time.Millisecond * 5)
		}
	}
}
