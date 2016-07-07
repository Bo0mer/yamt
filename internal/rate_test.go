package internal_test

import (
	"testing"

	"github.com/Bo0mer/yamt/internal"
)

var cases = []struct {
	last     uint64
	actual   uint64
	interval float64
	want     float64
}{
	{
		last:     5,
		actual:   10,
		interval: 1.0,
		want:     5.0,
	},
	{
		last:     510,
		actual:   500,
		interval: 5.0,
		want:     2.0,
	},
}

func TestComputeRate(t *testing.T) {
	for _, c := range cases {
		got := internal.ComputeRate(c.actual, c.last, c.interval)
		if got != c.want {
			t.Errorf("want %f, got %f\n", c.want, got)
		}
	}
}

func TestRateComputer(t *testing.T) {
	for _, c := range cases {
		rc := internal.RateComputer(c.interval)
		got := rc(c.actual, c.last)
		if got != c.want {
			t.Errorf("want %f, got %f\n", c.want, got)
		}
	}
}
