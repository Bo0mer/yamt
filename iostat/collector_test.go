package iostat_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/Bo0mer/yamt/iostat"
	"github.com/Bo0mer/yamt/iostat/iostatfakes"
	"github.com/Bo0mer/yamt/metric"
)

// Test that *DeviceStatCollector implements metric.Collector
var _ metric.Collector = (*iostat.DeviceStatCollector)(nil)

func TestNewDeviceStatCollector(t *testing.T) {
	_, err := iostat.NewDeviceStatCollector(iostat.DefaultDevStatReader, nil)
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	errReader := new(iostatfakes.FakeDeviceStatReader)
	errReader.ReadStatsReturns(nil, errors.New("kaboom"))
	_, err = iostat.NewDeviceStatCollector(errReader, nil)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

var devName = "sda"
var stats = map[int][]iostat.DeviceStat{
	0: []iostat.DeviceStat{
		iostat.DeviceStat{
			Name:     devName,
			Major:    1,
			Minor:    0,
			Reads:    1000,
			InFlight: 0,
		},
	},
	1: []iostat.DeviceStat{
		iostat.DeviceStat{
			Name:     devName,
			Major:    1,
			Minor:    0,
			Reads:    2000,
			InFlight: 42, // not cumulative
		},
	},
}

func newFakedReader(t *testing.T) iostat.DeviceStatReader {
	r := new(iostatfakes.FakeDeviceStatReader)
	i := 0
	r.ReadStatsStub = func() ([]iostat.DeviceStat, error) {
		ret := stats[i]
		i++
		return ret, nil
	}
	return r
}

func TestDevStatCollectorCollect(t *testing.T) {
	want := []metric.Event{
		metric.Event{}, // reads total, handled separately
		metric.Event{Name: "sda reads merged", Value: 0.0},
		metric.Event{Name: "sda reads sectors", Value: 0.0},
		metric.Event{Name: "sda reads time(ms)", Value: 0.0},
		metric.Event{Name: "sda writes total", Value: 0.0},
		metric.Event{Name: "sda writes merged", Value: 0.0},
		metric.Event{Name: "sda writes sectors", Value: 0.0},
		metric.Event{Name: "sda writes time(ms)", Value: 0.0},
		metric.Event{Name: "sda io inflight", Value: 42.0},
		metric.Event{Name: "sda io time(ms)", Value: 0.0},
		metric.Event{Name: "sda io weighted(ms)", Value: 0.0},
	}

	reader := newFakedReader(t)
	c, err := iostat.NewDeviceStatCollector(reader, nil)
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	got, err := c.Collect()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	if len(got) != len(want) {
		fmt.Printf("%#v\n", got)
		t.Errorf("expected %d events, got %d\n", len(want), len(got))
	}

	for i := range got {
		if got[i].Name == "sda reads total" {
			if f, ok := got[i].Value.(float64); !ok {
				t.Errorf("expected float64 value, got %T\n", got[i].Value)
			} else {
				if f <= 0 {
					t.Errorf("expected positive value, got %f\n", f)
				}
			}
			continue
		}
		if got[i] != want[i] {
			t.Errorf("expected %#v, got %#v\n", want[i], got[i])
		}
	}
}

func TestDevStatCollectorCollect_except(t *testing.T) {
	reader := newFakedReader(t)
	except := regexp.MustCompile("sda")
	c, err := iostat.NewDeviceStatCollector(reader, except)
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	got, err := c.Collect()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	if len(got) != 0 {
		t.Errorf("expected zero results, got %v\n", got)
	}
}
