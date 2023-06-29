// Code generated by counterfeiter. DO NOT EDIT.
package execsrvfakes

import (
	"sync"
	"time"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/execsrv"
)

type FakeStore struct {
	CreateStub        func(string, types.Exec) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 string
		arg2 types.Exec
	}
	createReturns struct {
		result1 error
	}
	createReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStub        func(string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 string
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	GetStub        func(string) (types.Exec, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 string
	}
	getReturns struct {
		result1 types.Exec
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 types.Exec
		result2 error
	}
	GetDriverStub        func(string) (string, error)
	getDriverMutex       sync.RWMutex
	getDriverArgsForCall []struct {
		arg1 string
	}
	getDriverReturns struct {
		result1 string
		result2 error
	}
	getDriverReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	ListStub        func(*string, *types.Provider, bool) ([]types.Exec, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		arg1 *string
		arg2 *types.Provider
		arg3 bool
	}
	listReturns struct {
		result1 []types.Exec
		result2 error
	}
	listReturnsOnCall map[int]struct {
		result1 []types.Exec
		result2 error
	}
	UpdateStub        func(string, types.Exec) error
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		arg1 string
		arg2 types.Exec
	}
	updateReturns struct {
		result1 error
	}
	updateReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateConnectionInfoStub        func(string, types.ConnectionInfo) error
	updateConnectionInfoMutex       sync.RWMutex
	updateConnectionInfoArgsForCall []struct {
		arg1 string
		arg2 types.ConnectionInfo
	}
	updateConnectionInfoReturns struct {
		result1 error
	}
	updateConnectionInfoReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateStatusStub        func(string, types.Status, time.Time, time.Time) error
	updateStatusMutex       sync.RWMutex
	updateStatusArgsForCall []struct {
		arg1 string
		arg2 types.Status
		arg3 time.Time
		arg4 time.Time
	}
	updateStatusReturns struct {
		result1 error
	}
	updateStatusReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStore) Create(arg1 string, arg2 types.Exec) error {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 string
		arg2 types.Exec
	}{arg1, arg2})
	stub := fake.CreateStub
	fakeReturns := fake.createReturns
	fake.recordInvocation("Create", []interface{}{arg1, arg2})
	fake.createMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStore) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeStore) CreateCalls(stub func(string, types.Exec) error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = stub
}

func (fake *FakeStore) CreateArgsForCall(i int) (string, types.Exec) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	argsForCall := fake.createArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeStore) CreateReturns(result1 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) CreateReturnsOnCall(i int, result1 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) Delete(arg1 string) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.DeleteStub
	fakeReturns := fake.deleteReturns
	fake.recordInvocation("Delete", []interface{}{arg1})
	fake.deleteMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStore) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeStore) DeleteCalls(stub func(string) error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeStore) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStore) DeleteReturns(result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) DeleteReturnsOnCall(i int, result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) Get(arg1 string) (types.Exec, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetStub
	fakeReturns := fake.getReturns
	fake.recordInvocation("Get", []interface{}{arg1})
	fake.getMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStore) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeStore) GetCalls(stub func(string) (types.Exec, error)) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *FakeStore) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStore) GetReturns(result1 types.Exec, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 types.Exec
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) GetReturnsOnCall(i int, result1 types.Exec, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 types.Exec
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 types.Exec
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) GetDriver(arg1 string) (string, error) {
	fake.getDriverMutex.Lock()
	ret, specificReturn := fake.getDriverReturnsOnCall[len(fake.getDriverArgsForCall)]
	fake.getDriverArgsForCall = append(fake.getDriverArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetDriverStub
	fakeReturns := fake.getDriverReturns
	fake.recordInvocation("GetDriver", []interface{}{arg1})
	fake.getDriverMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStore) GetDriverCallCount() int {
	fake.getDriverMutex.RLock()
	defer fake.getDriverMutex.RUnlock()
	return len(fake.getDriverArgsForCall)
}

func (fake *FakeStore) GetDriverCalls(stub func(string) (string, error)) {
	fake.getDriverMutex.Lock()
	defer fake.getDriverMutex.Unlock()
	fake.GetDriverStub = stub
}

func (fake *FakeStore) GetDriverArgsForCall(i int) string {
	fake.getDriverMutex.RLock()
	defer fake.getDriverMutex.RUnlock()
	argsForCall := fake.getDriverArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStore) GetDriverReturns(result1 string, result2 error) {
	fake.getDriverMutex.Lock()
	defer fake.getDriverMutex.Unlock()
	fake.GetDriverStub = nil
	fake.getDriverReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) GetDriverReturnsOnCall(i int, result1 string, result2 error) {
	fake.getDriverMutex.Lock()
	defer fake.getDriverMutex.Unlock()
	fake.GetDriverStub = nil
	if fake.getDriverReturnsOnCall == nil {
		fake.getDriverReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getDriverReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) List(arg1 *string, arg2 *types.Provider, arg3 bool) ([]types.Exec, error) {
	fake.listMutex.Lock()
	ret, specificReturn := fake.listReturnsOnCall[len(fake.listArgsForCall)]
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		arg1 *string
		arg2 *types.Provider
		arg3 bool
	}{arg1, arg2, arg3})
	stub := fake.ListStub
	fakeReturns := fake.listReturns
	fake.recordInvocation("List", []interface{}{arg1, arg2, arg3})
	fake.listMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStore) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeStore) ListCalls(stub func(*string, *types.Provider, bool) ([]types.Exec, error)) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = stub
}

func (fake *FakeStore) ListArgsForCall(i int) (*string, *types.Provider, bool) {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	argsForCall := fake.listArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeStore) ListReturns(result1 []types.Exec, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 []types.Exec
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) ListReturnsOnCall(i int, result1 []types.Exec, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	if fake.listReturnsOnCall == nil {
		fake.listReturnsOnCall = make(map[int]struct {
			result1 []types.Exec
			result2 error
		})
	}
	fake.listReturnsOnCall[i] = struct {
		result1 []types.Exec
		result2 error
	}{result1, result2}
}

func (fake *FakeStore) Update(arg1 string, arg2 types.Exec) error {
	fake.updateMutex.Lock()
	ret, specificReturn := fake.updateReturnsOnCall[len(fake.updateArgsForCall)]
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		arg1 string
		arg2 types.Exec
	}{arg1, arg2})
	stub := fake.UpdateStub
	fakeReturns := fake.updateReturns
	fake.recordInvocation("Update", []interface{}{arg1, arg2})
	fake.updateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStore) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *FakeStore) UpdateCalls(stub func(string, types.Exec) error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = stub
}

func (fake *FakeStore) UpdateArgsForCall(i int) (string, types.Exec) {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	argsForCall := fake.updateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeStore) UpdateReturns(result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) UpdateReturnsOnCall(i int, result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	if fake.updateReturnsOnCall == nil {
		fake.updateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) UpdateConnectionInfo(arg1 string, arg2 types.ConnectionInfo) error {
	fake.updateConnectionInfoMutex.Lock()
	ret, specificReturn := fake.updateConnectionInfoReturnsOnCall[len(fake.updateConnectionInfoArgsForCall)]
	fake.updateConnectionInfoArgsForCall = append(fake.updateConnectionInfoArgsForCall, struct {
		arg1 string
		arg2 types.ConnectionInfo
	}{arg1, arg2})
	stub := fake.UpdateConnectionInfoStub
	fakeReturns := fake.updateConnectionInfoReturns
	fake.recordInvocation("UpdateConnectionInfo", []interface{}{arg1, arg2})
	fake.updateConnectionInfoMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStore) UpdateConnectionInfoCallCount() int {
	fake.updateConnectionInfoMutex.RLock()
	defer fake.updateConnectionInfoMutex.RUnlock()
	return len(fake.updateConnectionInfoArgsForCall)
}

func (fake *FakeStore) UpdateConnectionInfoCalls(stub func(string, types.ConnectionInfo) error) {
	fake.updateConnectionInfoMutex.Lock()
	defer fake.updateConnectionInfoMutex.Unlock()
	fake.UpdateConnectionInfoStub = stub
}

func (fake *FakeStore) UpdateConnectionInfoArgsForCall(i int) (string, types.ConnectionInfo) {
	fake.updateConnectionInfoMutex.RLock()
	defer fake.updateConnectionInfoMutex.RUnlock()
	argsForCall := fake.updateConnectionInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeStore) UpdateConnectionInfoReturns(result1 error) {
	fake.updateConnectionInfoMutex.Lock()
	defer fake.updateConnectionInfoMutex.Unlock()
	fake.UpdateConnectionInfoStub = nil
	fake.updateConnectionInfoReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) UpdateConnectionInfoReturnsOnCall(i int, result1 error) {
	fake.updateConnectionInfoMutex.Lock()
	defer fake.updateConnectionInfoMutex.Unlock()
	fake.UpdateConnectionInfoStub = nil
	if fake.updateConnectionInfoReturnsOnCall == nil {
		fake.updateConnectionInfoReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateConnectionInfoReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) UpdateStatus(arg1 string, arg2 types.Status, arg3 time.Time, arg4 time.Time) error {
	fake.updateStatusMutex.Lock()
	ret, specificReturn := fake.updateStatusReturnsOnCall[len(fake.updateStatusArgsForCall)]
	fake.updateStatusArgsForCall = append(fake.updateStatusArgsForCall, struct {
		arg1 string
		arg2 types.Status
		arg3 time.Time
		arg4 time.Time
	}{arg1, arg2, arg3, arg4})
	stub := fake.UpdateStatusStub
	fakeReturns := fake.updateStatusReturns
	fake.recordInvocation("UpdateStatus", []interface{}{arg1, arg2, arg3, arg4})
	fake.updateStatusMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStore) UpdateStatusCallCount() int {
	fake.updateStatusMutex.RLock()
	defer fake.updateStatusMutex.RUnlock()
	return len(fake.updateStatusArgsForCall)
}

func (fake *FakeStore) UpdateStatusCalls(stub func(string, types.Status, time.Time, time.Time) error) {
	fake.updateStatusMutex.Lock()
	defer fake.updateStatusMutex.Unlock()
	fake.UpdateStatusStub = stub
}

func (fake *FakeStore) UpdateStatusArgsForCall(i int) (string, types.Status, time.Time, time.Time) {
	fake.updateStatusMutex.RLock()
	defer fake.updateStatusMutex.RUnlock()
	argsForCall := fake.updateStatusArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeStore) UpdateStatusReturns(result1 error) {
	fake.updateStatusMutex.Lock()
	defer fake.updateStatusMutex.Unlock()
	fake.UpdateStatusStub = nil
	fake.updateStatusReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) UpdateStatusReturnsOnCall(i int, result1 error) {
	fake.updateStatusMutex.Lock()
	defer fake.updateStatusMutex.Unlock()
	fake.UpdateStatusStub = nil
	if fake.updateStatusReturnsOnCall == nil {
		fake.updateStatusReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateStatusReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.getDriverMutex.RLock()
	defer fake.getDriverMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	fake.updateConnectionInfoMutex.RLock()
	defer fake.updateConnectionInfoMutex.RUnlock()
	fake.updateStatusMutex.RLock()
	defer fake.updateStatusMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeStore) recordInvocation(key string, args []interface{}) {
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

var _ execsrv.Store = new(FakeStore)
