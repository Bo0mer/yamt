package riemann

import (
	"fmt"
)

type Event struct {
	Name  string
	Value int
}

type Emitter struct{}

func (e *Emitter) Emit(event Event) error {
	fmt.Printf("%#v\n", event)
	return nil
}

func (e *Emitter) Err() error {
	return nil
}
