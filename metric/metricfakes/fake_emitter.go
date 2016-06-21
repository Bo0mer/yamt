// This file was generated by counterfeiter
package metricfakes

import (
	"sync"

	"github.com/bo0mer/yamt/metric"
)

type FakeEmitter struct {
	EmitStub        func(metric.Event) error
	emitMutex       sync.RWMutex
	emitArgsForCall []struct {
		arg1 metric.Event
	}
	emitReturns struct {
		result1 error
	}
	ErrStub        func() error
	errMutex       sync.RWMutex
	errArgsForCall []struct{}
	errReturns     struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeEmitter) Emit(arg1 metric.Event) error {
	fake.emitMutex.Lock()
	fake.emitArgsForCall = append(fake.emitArgsForCall, struct {
		arg1 metric.Event
	}{arg1})
	fake.recordInvocation("Emit", []interface{}{arg1})
	fake.emitMutex.Unlock()
	if fake.EmitStub != nil {
		return fake.EmitStub(arg1)
	} else {
		return fake.emitReturns.result1
	}
}

func (fake *FakeEmitter) EmitCallCount() int {
	fake.emitMutex.RLock()
	defer fake.emitMutex.RUnlock()
	return len(fake.emitArgsForCall)
}

func (fake *FakeEmitter) EmitArgsForCall(i int) metric.Event {
	fake.emitMutex.RLock()
	defer fake.emitMutex.RUnlock()
	return fake.emitArgsForCall[i].arg1
}

func (fake *FakeEmitter) EmitReturns(result1 error) {
	fake.EmitStub = nil
	fake.emitReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEmitter) Err() error {
	fake.errMutex.Lock()
	fake.errArgsForCall = append(fake.errArgsForCall, struct{}{})
	fake.recordInvocation("Err", []interface{}{})
	fake.errMutex.Unlock()
	if fake.ErrStub != nil {
		return fake.ErrStub()
	} else {
		return fake.errReturns.result1
	}
}

func (fake *FakeEmitter) ErrCallCount() int {
	fake.errMutex.RLock()
	defer fake.errMutex.RUnlock()
	return len(fake.errArgsForCall)
}

func (fake *FakeEmitter) ErrReturns(result1 error) {
	fake.ErrStub = nil
	fake.errReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEmitter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.emitMutex.RLock()
	defer fake.emitMutex.RUnlock()
	fake.errMutex.RLock()
	defer fake.errMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeEmitter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ metric.Emitter = new(FakeEmitter)
