// Code generated by counterfeiter. DO NOT EDIT.
package builderfakes

import (
	"context"
	"sync"

	"github.com/unweave/unweave-v1/api/types"
	"github.com/unweave/unweave-v1/builder"
)

type FakeLogDriver struct {
	GetLogsStub        func(context.Context, string) ([]types.LogEntry, error)
	getLogsMutex       sync.RWMutex
	getLogsArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	getLogsReturns struct {
		result1 []types.LogEntry
		result2 error
	}
	getLogsReturnsOnCall map[int]struct {
		result1 []types.LogEntry
		result2 error
	}
	SaveLogsStub        func(context.Context, string, []types.LogEntry) error
	saveLogsMutex       sync.RWMutex
	saveLogsArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 []types.LogEntry
	}
	saveLogsReturns struct {
		result1 error
	}
	saveLogsReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLogDriver) GetLogs(arg1 context.Context, arg2 string) ([]types.LogEntry, error) {
	fake.getLogsMutex.Lock()
	ret, specificReturn := fake.getLogsReturnsOnCall[len(fake.getLogsArgsForCall)]
	fake.getLogsArgsForCall = append(fake.getLogsArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.GetLogsStub
	fakeReturns := fake.getLogsReturns
	fake.recordInvocation("GetLogs", []interface{}{arg1, arg2})
	fake.getLogsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeLogDriver) GetLogsCallCount() int {
	fake.getLogsMutex.RLock()
	defer fake.getLogsMutex.RUnlock()
	return len(fake.getLogsArgsForCall)
}

func (fake *FakeLogDriver) GetLogsCalls(stub func(context.Context, string) ([]types.LogEntry, error)) {
	fake.getLogsMutex.Lock()
	defer fake.getLogsMutex.Unlock()
	fake.GetLogsStub = stub
}

func (fake *FakeLogDriver) GetLogsArgsForCall(i int) (context.Context, string) {
	fake.getLogsMutex.RLock()
	defer fake.getLogsMutex.RUnlock()
	argsForCall := fake.getLogsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeLogDriver) GetLogsReturns(result1 []types.LogEntry, result2 error) {
	fake.getLogsMutex.Lock()
	defer fake.getLogsMutex.Unlock()
	fake.GetLogsStub = nil
	fake.getLogsReturns = struct {
		result1 []types.LogEntry
		result2 error
	}{result1, result2}
}

func (fake *FakeLogDriver) GetLogsReturnsOnCall(i int, result1 []types.LogEntry, result2 error) {
	fake.getLogsMutex.Lock()
	defer fake.getLogsMutex.Unlock()
	fake.GetLogsStub = nil
	if fake.getLogsReturnsOnCall == nil {
		fake.getLogsReturnsOnCall = make(map[int]struct {
			result1 []types.LogEntry
			result2 error
		})
	}
	fake.getLogsReturnsOnCall[i] = struct {
		result1 []types.LogEntry
		result2 error
	}{result1, result2}
}

func (fake *FakeLogDriver) SaveLogs(arg1 context.Context, arg2 string, arg3 []types.LogEntry) error {
	var arg3Copy []types.LogEntry
	if arg3 != nil {
		arg3Copy = make([]types.LogEntry, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.saveLogsMutex.Lock()
	ret, specificReturn := fake.saveLogsReturnsOnCall[len(fake.saveLogsArgsForCall)]
	fake.saveLogsArgsForCall = append(fake.saveLogsArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 []types.LogEntry
	}{arg1, arg2, arg3Copy})
	stub := fake.SaveLogsStub
	fakeReturns := fake.saveLogsReturns
	fake.recordInvocation("SaveLogs", []interface{}{arg1, arg2, arg3Copy})
	fake.saveLogsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLogDriver) SaveLogsCallCount() int {
	fake.saveLogsMutex.RLock()
	defer fake.saveLogsMutex.RUnlock()
	return len(fake.saveLogsArgsForCall)
}

func (fake *FakeLogDriver) SaveLogsCalls(stub func(context.Context, string, []types.LogEntry) error) {
	fake.saveLogsMutex.Lock()
	defer fake.saveLogsMutex.Unlock()
	fake.SaveLogsStub = stub
}

func (fake *FakeLogDriver) SaveLogsArgsForCall(i int) (context.Context, string, []types.LogEntry) {
	fake.saveLogsMutex.RLock()
	defer fake.saveLogsMutex.RUnlock()
	argsForCall := fake.saveLogsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeLogDriver) SaveLogsReturns(result1 error) {
	fake.saveLogsMutex.Lock()
	defer fake.saveLogsMutex.Unlock()
	fake.SaveLogsStub = nil
	fake.saveLogsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeLogDriver) SaveLogsReturnsOnCall(i int, result1 error) {
	fake.saveLogsMutex.Lock()
	defer fake.saveLogsMutex.Unlock()
	fake.SaveLogsStub = nil
	if fake.saveLogsReturnsOnCall == nil {
		fake.saveLogsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.saveLogsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeLogDriver) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getLogsMutex.RLock()
	defer fake.getLogsMutex.RUnlock()
	fake.saveLogsMutex.RLock()
	defer fake.saveLogsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLogDriver) recordInvocation(key string, args []interface{}) {
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

var _ builder.LogDriver = new(FakeLogDriver)
