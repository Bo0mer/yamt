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

func TestEventBuilder(t *testing.T) {
	buildEvent := eventBuilder("eth0", time.Second*5)
	e := buildEvent(20, "tx.bytes")
	if e.Name != "eth0.tx.bytes" {
		t.Errorf("expected event name %q, got %q\n", "eth0.tx.bytes", e.Name)
	}
	if e.Value != 20/5 {
		t.Errorf("expected event value %d, got %d\n", 20/5, e.Value)
	}
}
