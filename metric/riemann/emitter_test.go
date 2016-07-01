package riemann

import "testing"

func TestPrefix(t *testing.T) {
	prefix := "woho"
	e := NewEmitter("", Prefix(prefix))
	if e.prefix != prefix {
		t.Errorf("expected prefix %q, got %q\n", prefix, e.prefix)
	}
}

func TestHost(t *testing.T) {
	host := "local"
	e := NewEmitter("", Host(host))
	if e.host != host {
		t.Errorf("expected host %q, got %q\n", host, e.host)
	}
}

func TestAttributes(t *testing.T) {
	attr := map[string]string{"attr1": "value1", "x": "y"}
	e := NewEmitter("", Attributes(attr))
	if len(attr) != len(e.attributes) {
		t.Errorf("expected attributes %v, got %v\n", attr, e.attributes)
	}
	for k, v := range attr {
		if actualValue, ok := e.attributes[k]; !ok || v != actualValue {
			t.Errorf("expected attributes %v, got %v\n", attr, e.attributes)
		}
	}
}

func TestTags(t *testing.T) {
	tags := []string{"hash", "tag"}
	e := NewEmitter("", Tags(tags))
	if len(tags) != len(e.tags) {
		t.Errorf("expected tags %v, got %v\n", tags, e.tags)
	}
	for i, tag := range tags {
		if tag != e.tags[i] {
			t.Errorf("expected tags %v, got %v\n", tags, e.tags)
		}
	}
}

func TestPrependPrefix(t *testing.T) {
	srv, prefix := "service", "some"
	got := prependPrefix(srv, prefix)
	if got != "some.service" {
		t.Errorf("expected 'some.service', got: %q\n", got)
	}
	prefix = ""
	got = prependPrefix(srv, prefix)
	if got != "service" {
		t.Errorf("expected 'service', got: %q\n", got)
	}
}
