package netstat

import (
	"regexp"
	"testing"
	"time"
)

func TestInterval(t *testing.T) {
	d := time.Second * 42
	r := NewReporter(nil, Interval(d))
	if r.interval != d {
		t.Errorf("expected interval %d, got %d\n", d, r.interval)
	}
}

func TestExcept(t *testing.T) {
	except := regexp.MustCompile("lo")
	r := NewReporter(nil, Except(except))
	if r.except != except {
		t.Errorf("expected %v, got %v\n", except, r.except)
	}
}
