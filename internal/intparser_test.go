package internal_test

import (
	"testing"

	"github.com/bo0mer/yamt/internal"
)

func TestParsInt(t *testing.T) {
	p := &internal.ErrParser{}
	i := p.ParseInt("42")
	if i != 42 {
		t.Errorf("expected 42, got %d\n", i)
	}
}

func TestParseUint64(t *testing.T) {
	p := &internal.ErrParser{}
	u := p.ParseUint64("42")
	if u != 42 {
		t.Errorf("expected 42, got %d\n", u)
	}
}

func TestErr(t *testing.T) {
	p := &internal.ErrParser{}
	_ = p.ParseInt("not int")
	if p.Err() == nil {
		t.Errorf("expected error, got nil\n")
	}
}
