package netstat_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/Bo0mer/yamt/metric"
	"github.com/Bo0mer/yamt/netstat"
	"github.com/Bo0mer/yamt/netstat/netstatfakes"
)

// Test that *IfStatCollector implements metric.Collector
var _ metric.Collector = (*netstat.IfStatCollector)(nil)

func TestNewIfStatCollector(t *testing.T) {
	_, err := netstat.NewIfStatCollector(netstat.DefaultIfStatReader, nil)
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	errReader := new(netstatfakes.FakeInterfaceStatReader)
	errReader.ReadStatsReturns(nil, errors.New("kaboom"))
	_, err = netstat.NewIfStatCollector(errReader, nil)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

var ifName = "eth0"

var stats = map[int][]netstat.IfStat{
	0: []netstat.IfStat{
		netstat.IfStat{
			Name:    ifName,
			TxBytes: 1000,
		},
	},
	1: []netstat.IfStat{
		netstat.IfStat{
			Name:    ifName,
			TxBytes: 2000,
		},
	},
}

func newFakedReader(t *testing.T) netstat.InterfaceStatReader {
	fakeReader := new(netstatfakes.FakeInterfaceStatReader)
	i := 0
	fakeReader.ReadStatsStub = func() ([]netstat.IfStat, error) {
		ret := stats[i]
		i++
		return ret, nil
	}
	return fakeReader
}

func TestIfStatCollectorCollect(t *testing.T) {
	want := []metric.Event{
		metric.Event{Name: "eth0 rx bytes", Value: 0.0},
		metric.Event{Name: "eth0 rx packets", Value: 0.0},
		metric.Event{Name: "eth0 rx errs", Value: 0.0},
		metric.Event{Name: "eth0 rx drop", Value: 0.0},
		metric.Event{Name: "eth0 rx fifo", Value: 0.0},
		metric.Event{Name: "eth0 rx frame", Value: 0.0},
		metric.Event{Name: "eth0 rx compressed", Value: 0.0},
		metric.Event{Name: "eth0 rx multicast", Value: 0.0},
		metric.Event{}, // tx.bytes, handled separately
		metric.Event{Name: "eth0 tx packets", Value: 0.0},
		metric.Event{Name: "eth0 tx errs", Value: 0.0},
		metric.Event{Name: "eth0 tx drop", Value: 0.0},
		metric.Event{Name: "eth0 tx fifo", Value: 0.0},
		metric.Event{Name: "eth0 tx colls", Value: 0.0},
		metric.Event{Name: "eth0 tx carrier", Value: 0.0},
		metric.Event{Name: "eth0 tx compressed", Value: 0.0},
	}

	reader := newFakedReader(t)
	c, err := netstat.NewIfStatCollector(reader, nil)
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	got, err := c.Collect()
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	if len(got) != len(want) {
		t.Errorf("expected %d events, got %d\n", len(want), len(got))
	}

	for i := range got {
		if got[i].Name == ifName+" tx bytes" {
			if f, ok := got[i].Value.(float64); !ok {
				t.Errorf("expected flaot64 value, got %T\n", got[i].Value)
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

func TestIfStatCollectorCollect_except(t *testing.T) {
	reader := newFakedReader(t)
	except := regexp.MustCompile("eth0")
	c, err := netstat.NewIfStatCollector(reader, except)
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
