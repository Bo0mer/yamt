package riemann

import (
	"github.com/bigdatadev/goryman"
	"github.com/bo0mer/yamt/metric"
)

type Option func(e *Emitter)

// Prefix sets prefix to be prepended to each event name.
func Prefix(prefix string) Option {
	return func(e *Emitter) {
		e.prefix = prefix
	}
}

// Host sets the reported host for each event. Defaults to os.Hostname.
func Host(host string) Option {
	return func(e *Emitter) {
		e.host = host
	}
}

// Attributes sets attributes to be added to each emitted event.
func Attributes(attributes map[string]string) Option {
	return func(e *Emitter) {
		e.attributes = attributes
	}
}

// Tags sets tags to be added to each emitted event.
func Tags(tags []string) Option {
	return func(e *Emitter) {
		e.tags = tags
	}
}

// Emitter sends events to Riemann.
type Emitter struct {
	c           *goryman.GorymanClient
	isConnected bool

	prefix     string
	host       string
	attributes map[string]string
	tags       []string
}

// NewEmitter returns brand new emitter.
func NewEmitter(addr string, opts ...Option) *Emitter {
	c := goryman.NewGorymanClient(addr)
	e := &Emitter{
		c:           c,
		isConnected: false,
	}

	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Emit sends the specified event to riemann.
func (e *Emitter) Emit(event metric.Event) error {
	if !e.isConnected {
		if err := e.c.Connect(); err != nil {
			return err
		}
		e.isConnected = true
	}

	err := e.c.SendEvent(&goryman.Event{
		Service:    prependPrefix(event.Name, e.prefix),
		Metric:     event.Value,
		Host:       e.host,
		Attributes: e.attributes,
		Tags:       e.tags,
		State:      "ok",
	})
	if err != nil {
		e.c.Close()
		e.isConnected = false
	}
	return err
}

func prependPrefix(service string, prefix string) string {
	if prefix == "" {
		return service
	}
	return prefix + "." + service
}
