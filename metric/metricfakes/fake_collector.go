// This file was generated by counterfeiter
package metricfakes

import (
	"sync"

	"github.com/bo0mer/yamt/metric"
)

type FakeCollector struct {
	CollectStub        func() ([]metric.Event, error)
	collectMutex       sync.RWMutex
	collectArgsForCall []struct{}
	collectReturns     struct {
		result1 []metric.Event
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCollector) Collect() ([]metric.Event, error) {
	fake.collectMutex.Lock()
	fake.collectArgsForCall = append(fake.collectArgsForCall, struct{}{})
	fake.recordInvocation("Collect", []interface{}{})
	fake.collectMutex.Unlock()
	if fake.CollectStub != nil {
		return fake.CollectStub()
	} else {
		return fake.collectReturns.result1, fake.collectReturns.result2
	}
}

func (fake *FakeCollector) CollectCallCount() int {
	fake.collectMutex.RLock()
	defer fake.collectMutex.RUnlock()
	return len(fake.collectArgsForCall)
}

func (fake *FakeCollector) CollectReturns(result1 []metric.Event, result2 error) {
	fake.CollectStub = nil
	fake.collectReturns = struct {
		result1 []metric.Event
		result2 error
	}{result1, result2}
}

func (fake *FakeCollector) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.collectMutex.RLock()
	defer fake.collectMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeCollector) recordInvocation(key string, args []interface{}) {
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

var _ metric.Collector = new(FakeCollector)