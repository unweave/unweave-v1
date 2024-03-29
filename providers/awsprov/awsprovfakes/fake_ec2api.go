// Code generated by counterfeiter. DO NOT EDIT.
package awsprovfakes

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/unweave/unweave-v1/providers/awsprov"
)

type FakeEc2API struct {
	CreateTagsStub        func(context.Context, *ec2.CreateTagsInput, ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
	createTagsMutex       sync.RWMutex
	createTagsArgsForCall []struct {
		arg1 context.Context
		arg2 *ec2.CreateTagsInput
		arg3 []func(*ec2.Options)
	}
	createTagsReturns struct {
		result1 *ec2.CreateTagsOutput
		result2 error
	}
	createTagsReturnsOnCall map[int]struct {
		result1 *ec2.CreateTagsOutput
		result2 error
	}
	CreateVolumeStub        func(context.Context, *ec2.CreateVolumeInput, ...func(*ec2.Options)) (*ec2.CreateVolumeOutput, error)
	createVolumeMutex       sync.RWMutex
	createVolumeArgsForCall []struct {
		arg1 context.Context
		arg2 *ec2.CreateVolumeInput
		arg3 []func(*ec2.Options)
	}
	createVolumeReturns struct {
		result1 *ec2.CreateVolumeOutput
		result2 error
	}
	createVolumeReturnsOnCall map[int]struct {
		result1 *ec2.CreateVolumeOutput
		result2 error
	}
	DeleteVolumeStub        func(context.Context, *ec2.DeleteVolumeInput, ...func(*ec2.Options)) (*ec2.DeleteVolumeOutput, error)
	deleteVolumeMutex       sync.RWMutex
	deleteVolumeArgsForCall []struct {
		arg1 context.Context
		arg2 *ec2.DeleteVolumeInput
		arg3 []func(*ec2.Options)
	}
	deleteVolumeReturns struct {
		result1 *ec2.DeleteVolumeOutput
		result2 error
	}
	deleteVolumeReturnsOnCall map[int]struct {
		result1 *ec2.DeleteVolumeOutput
		result2 error
	}
	DescribeInstanceStatusStub        func(context.Context, *ec2.DescribeInstanceStatusInput, ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error)
	describeInstanceStatusMutex       sync.RWMutex
	describeInstanceStatusArgsForCall []struct {
		arg1 context.Context
		arg2 *ec2.DescribeInstanceStatusInput
		arg3 []func(*ec2.Options)
	}
	describeInstanceStatusReturns struct {
		result1 *ec2.DescribeInstanceStatusOutput
		result2 error
	}
	describeInstanceStatusReturnsOnCall map[int]struct {
		result1 *ec2.DescribeInstanceStatusOutput
		result2 error
	}
	DescribeInstanceTypesStub        func(context.Context, *ec2.DescribeInstanceTypesInput, ...func(*ec2.Options)) (*ec2.DescribeInstanceTypesOutput, error)
	describeInstanceTypesMutex       sync.RWMutex
	describeInstanceTypesArgsForCall []struct {
		arg1 context.Context
		arg2 *ec2.DescribeInstanceTypesInput
		arg3 []func(*ec2.Options)
	}
	describeInstanceTypesReturns struct {
		result1 *ec2.DescribeInstanceTypesOutput
		result2 error
	}
	describeInstanceTypesReturnsOnCall map[int]struct {
		result1 *ec2.DescribeInstanceTypesOutput
		result2 error
	}
	DescribeInstancesStub        func(context.Context, *ec2.DescribeInstancesInput, ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	describeInstancesMutex       sync.RWMutex
	describeInstancesArgsForCall []struct {
		arg1 context.Context
		arg2 *ec2.DescribeInstancesInput
		arg3 []func(*ec2.Options)
	}
	describeInstancesReturns struct {
		result1 *ec2.DescribeInstancesOutput
		result2 error
	}
	describeInstancesReturnsOnCall map[int]struct {
		result1 *ec2.DescribeInstancesOutput
		result2 error
	}
	ModifyVolumeStub        func(context.Context, *ec2.ModifyVolumeInput, ...func(*ec2.Options)) (*ec2.ModifyVolumeOutput, error)
	modifyVolumeMutex       sync.RWMutex
	modifyVolumeArgsForCall []struct {
		arg1 context.Context
		arg2 *ec2.ModifyVolumeInput
		arg3 []func(*ec2.Options)
	}
	modifyVolumeReturns struct {
		result1 *ec2.ModifyVolumeOutput
		result2 error
	}
	modifyVolumeReturnsOnCall map[int]struct {
		result1 *ec2.ModifyVolumeOutput
		result2 error
	}
	RunInstancesStub        func(context.Context, *ec2.RunInstancesInput, ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
	runInstancesMutex       sync.RWMutex
	runInstancesArgsForCall []struct {
		arg1 context.Context
		arg2 *ec2.RunInstancesInput
		arg3 []func(*ec2.Options)
	}
	runInstancesReturns struct {
		result1 *ec2.RunInstancesOutput
		result2 error
	}
	runInstancesReturnsOnCall map[int]struct {
		result1 *ec2.RunInstancesOutput
		result2 error
	}
	TerminateInstancesStub        func(context.Context, *ec2.TerminateInstancesInput, ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
	terminateInstancesMutex       sync.RWMutex
	terminateInstancesArgsForCall []struct {
		arg1 context.Context
		arg2 *ec2.TerminateInstancesInput
		arg3 []func(*ec2.Options)
	}
	terminateInstancesReturns struct {
		result1 *ec2.TerminateInstancesOutput
		result2 error
	}
	terminateInstancesReturnsOnCall map[int]struct {
		result1 *ec2.TerminateInstancesOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeEc2API) CreateTags(arg1 context.Context, arg2 *ec2.CreateTagsInput, arg3 ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error) {
	fake.createTagsMutex.Lock()
	ret, specificReturn := fake.createTagsReturnsOnCall[len(fake.createTagsArgsForCall)]
	fake.createTagsArgsForCall = append(fake.createTagsArgsForCall, struct {
		arg1 context.Context
		arg2 *ec2.CreateTagsInput
		arg3 []func(*ec2.Options)
	}{arg1, arg2, arg3})
	stub := fake.CreateTagsStub
	fakeReturns := fake.createTagsReturns
	fake.recordInvocation("CreateTags", []interface{}{arg1, arg2, arg3})
	fake.createTagsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeEc2API) CreateTagsCallCount() int {
	fake.createTagsMutex.RLock()
	defer fake.createTagsMutex.RUnlock()
	return len(fake.createTagsArgsForCall)
}

func (fake *FakeEc2API) CreateTagsCalls(stub func(context.Context, *ec2.CreateTagsInput, ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)) {
	fake.createTagsMutex.Lock()
	defer fake.createTagsMutex.Unlock()
	fake.CreateTagsStub = stub
}

func (fake *FakeEc2API) CreateTagsArgsForCall(i int) (context.Context, *ec2.CreateTagsInput, []func(*ec2.Options)) {
	fake.createTagsMutex.RLock()
	defer fake.createTagsMutex.RUnlock()
	argsForCall := fake.createTagsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeEc2API) CreateTagsReturns(result1 *ec2.CreateTagsOutput, result2 error) {
	fake.createTagsMutex.Lock()
	defer fake.createTagsMutex.Unlock()
	fake.CreateTagsStub = nil
	fake.createTagsReturns = struct {
		result1 *ec2.CreateTagsOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) CreateTagsReturnsOnCall(i int, result1 *ec2.CreateTagsOutput, result2 error) {
	fake.createTagsMutex.Lock()
	defer fake.createTagsMutex.Unlock()
	fake.CreateTagsStub = nil
	if fake.createTagsReturnsOnCall == nil {
		fake.createTagsReturnsOnCall = make(map[int]struct {
			result1 *ec2.CreateTagsOutput
			result2 error
		})
	}
	fake.createTagsReturnsOnCall[i] = struct {
		result1 *ec2.CreateTagsOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) CreateVolume(arg1 context.Context, arg2 *ec2.CreateVolumeInput, arg3 ...func(*ec2.Options)) (*ec2.CreateVolumeOutput, error) {
	fake.createVolumeMutex.Lock()
	ret, specificReturn := fake.createVolumeReturnsOnCall[len(fake.createVolumeArgsForCall)]
	fake.createVolumeArgsForCall = append(fake.createVolumeArgsForCall, struct {
		arg1 context.Context
		arg2 *ec2.CreateVolumeInput
		arg3 []func(*ec2.Options)
	}{arg1, arg2, arg3})
	stub := fake.CreateVolumeStub
	fakeReturns := fake.createVolumeReturns
	fake.recordInvocation("CreateVolume", []interface{}{arg1, arg2, arg3})
	fake.createVolumeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeEc2API) CreateVolumeCallCount() int {
	fake.createVolumeMutex.RLock()
	defer fake.createVolumeMutex.RUnlock()
	return len(fake.createVolumeArgsForCall)
}

func (fake *FakeEc2API) CreateVolumeCalls(stub func(context.Context, *ec2.CreateVolumeInput, ...func(*ec2.Options)) (*ec2.CreateVolumeOutput, error)) {
	fake.createVolumeMutex.Lock()
	defer fake.createVolumeMutex.Unlock()
	fake.CreateVolumeStub = stub
}

func (fake *FakeEc2API) CreateVolumeArgsForCall(i int) (context.Context, *ec2.CreateVolumeInput, []func(*ec2.Options)) {
	fake.createVolumeMutex.RLock()
	defer fake.createVolumeMutex.RUnlock()
	argsForCall := fake.createVolumeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeEc2API) CreateVolumeReturns(result1 *ec2.CreateVolumeOutput, result2 error) {
	fake.createVolumeMutex.Lock()
	defer fake.createVolumeMutex.Unlock()
	fake.CreateVolumeStub = nil
	fake.createVolumeReturns = struct {
		result1 *ec2.CreateVolumeOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) CreateVolumeReturnsOnCall(i int, result1 *ec2.CreateVolumeOutput, result2 error) {
	fake.createVolumeMutex.Lock()
	defer fake.createVolumeMutex.Unlock()
	fake.CreateVolumeStub = nil
	if fake.createVolumeReturnsOnCall == nil {
		fake.createVolumeReturnsOnCall = make(map[int]struct {
			result1 *ec2.CreateVolumeOutput
			result2 error
		})
	}
	fake.createVolumeReturnsOnCall[i] = struct {
		result1 *ec2.CreateVolumeOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) DeleteVolume(arg1 context.Context, arg2 *ec2.DeleteVolumeInput, arg3 ...func(*ec2.Options)) (*ec2.DeleteVolumeOutput, error) {
	fake.deleteVolumeMutex.Lock()
	ret, specificReturn := fake.deleteVolumeReturnsOnCall[len(fake.deleteVolumeArgsForCall)]
	fake.deleteVolumeArgsForCall = append(fake.deleteVolumeArgsForCall, struct {
		arg1 context.Context
		arg2 *ec2.DeleteVolumeInput
		arg3 []func(*ec2.Options)
	}{arg1, arg2, arg3})
	stub := fake.DeleteVolumeStub
	fakeReturns := fake.deleteVolumeReturns
	fake.recordInvocation("DeleteVolume", []interface{}{arg1, arg2, arg3})
	fake.deleteVolumeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeEc2API) DeleteVolumeCallCount() int {
	fake.deleteVolumeMutex.RLock()
	defer fake.deleteVolumeMutex.RUnlock()
	return len(fake.deleteVolumeArgsForCall)
}

func (fake *FakeEc2API) DeleteVolumeCalls(stub func(context.Context, *ec2.DeleteVolumeInput, ...func(*ec2.Options)) (*ec2.DeleteVolumeOutput, error)) {
	fake.deleteVolumeMutex.Lock()
	defer fake.deleteVolumeMutex.Unlock()
	fake.DeleteVolumeStub = stub
}

func (fake *FakeEc2API) DeleteVolumeArgsForCall(i int) (context.Context, *ec2.DeleteVolumeInput, []func(*ec2.Options)) {
	fake.deleteVolumeMutex.RLock()
	defer fake.deleteVolumeMutex.RUnlock()
	argsForCall := fake.deleteVolumeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeEc2API) DeleteVolumeReturns(result1 *ec2.DeleteVolumeOutput, result2 error) {
	fake.deleteVolumeMutex.Lock()
	defer fake.deleteVolumeMutex.Unlock()
	fake.DeleteVolumeStub = nil
	fake.deleteVolumeReturns = struct {
		result1 *ec2.DeleteVolumeOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) DeleteVolumeReturnsOnCall(i int, result1 *ec2.DeleteVolumeOutput, result2 error) {
	fake.deleteVolumeMutex.Lock()
	defer fake.deleteVolumeMutex.Unlock()
	fake.DeleteVolumeStub = nil
	if fake.deleteVolumeReturnsOnCall == nil {
		fake.deleteVolumeReturnsOnCall = make(map[int]struct {
			result1 *ec2.DeleteVolumeOutput
			result2 error
		})
	}
	fake.deleteVolumeReturnsOnCall[i] = struct {
		result1 *ec2.DeleteVolumeOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) DescribeInstanceStatus(arg1 context.Context, arg2 *ec2.DescribeInstanceStatusInput, arg3 ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error) {
	fake.describeInstanceStatusMutex.Lock()
	ret, specificReturn := fake.describeInstanceStatusReturnsOnCall[len(fake.describeInstanceStatusArgsForCall)]
	fake.describeInstanceStatusArgsForCall = append(fake.describeInstanceStatusArgsForCall, struct {
		arg1 context.Context
		arg2 *ec2.DescribeInstanceStatusInput
		arg3 []func(*ec2.Options)
	}{arg1, arg2, arg3})
	stub := fake.DescribeInstanceStatusStub
	fakeReturns := fake.describeInstanceStatusReturns
	fake.recordInvocation("DescribeInstanceStatus", []interface{}{arg1, arg2, arg3})
	fake.describeInstanceStatusMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeEc2API) DescribeInstanceStatusCallCount() int {
	fake.describeInstanceStatusMutex.RLock()
	defer fake.describeInstanceStatusMutex.RUnlock()
	return len(fake.describeInstanceStatusArgsForCall)
}

func (fake *FakeEc2API) DescribeInstanceStatusCalls(stub func(context.Context, *ec2.DescribeInstanceStatusInput, ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error)) {
	fake.describeInstanceStatusMutex.Lock()
	defer fake.describeInstanceStatusMutex.Unlock()
	fake.DescribeInstanceStatusStub = stub
}

func (fake *FakeEc2API) DescribeInstanceStatusArgsForCall(i int) (context.Context, *ec2.DescribeInstanceStatusInput, []func(*ec2.Options)) {
	fake.describeInstanceStatusMutex.RLock()
	defer fake.describeInstanceStatusMutex.RUnlock()
	argsForCall := fake.describeInstanceStatusArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeEc2API) DescribeInstanceStatusReturns(result1 *ec2.DescribeInstanceStatusOutput, result2 error) {
	fake.describeInstanceStatusMutex.Lock()
	defer fake.describeInstanceStatusMutex.Unlock()
	fake.DescribeInstanceStatusStub = nil
	fake.describeInstanceStatusReturns = struct {
		result1 *ec2.DescribeInstanceStatusOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) DescribeInstanceStatusReturnsOnCall(i int, result1 *ec2.DescribeInstanceStatusOutput, result2 error) {
	fake.describeInstanceStatusMutex.Lock()
	defer fake.describeInstanceStatusMutex.Unlock()
	fake.DescribeInstanceStatusStub = nil
	if fake.describeInstanceStatusReturnsOnCall == nil {
		fake.describeInstanceStatusReturnsOnCall = make(map[int]struct {
			result1 *ec2.DescribeInstanceStatusOutput
			result2 error
		})
	}
	fake.describeInstanceStatusReturnsOnCall[i] = struct {
		result1 *ec2.DescribeInstanceStatusOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) DescribeInstanceTypes(arg1 context.Context, arg2 *ec2.DescribeInstanceTypesInput, arg3 ...func(*ec2.Options)) (*ec2.DescribeInstanceTypesOutput, error) {
	fake.describeInstanceTypesMutex.Lock()
	ret, specificReturn := fake.describeInstanceTypesReturnsOnCall[len(fake.describeInstanceTypesArgsForCall)]
	fake.describeInstanceTypesArgsForCall = append(fake.describeInstanceTypesArgsForCall, struct {
		arg1 context.Context
		arg2 *ec2.DescribeInstanceTypesInput
		arg3 []func(*ec2.Options)
	}{arg1, arg2, arg3})
	stub := fake.DescribeInstanceTypesStub
	fakeReturns := fake.describeInstanceTypesReturns
	fake.recordInvocation("DescribeInstanceTypes", []interface{}{arg1, arg2, arg3})
	fake.describeInstanceTypesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeEc2API) DescribeInstanceTypesCallCount() int {
	fake.describeInstanceTypesMutex.RLock()
	defer fake.describeInstanceTypesMutex.RUnlock()
	return len(fake.describeInstanceTypesArgsForCall)
}

func (fake *FakeEc2API) DescribeInstanceTypesCalls(stub func(context.Context, *ec2.DescribeInstanceTypesInput, ...func(*ec2.Options)) (*ec2.DescribeInstanceTypesOutput, error)) {
	fake.describeInstanceTypesMutex.Lock()
	defer fake.describeInstanceTypesMutex.Unlock()
	fake.DescribeInstanceTypesStub = stub
}

func (fake *FakeEc2API) DescribeInstanceTypesArgsForCall(i int) (context.Context, *ec2.DescribeInstanceTypesInput, []func(*ec2.Options)) {
	fake.describeInstanceTypesMutex.RLock()
	defer fake.describeInstanceTypesMutex.RUnlock()
	argsForCall := fake.describeInstanceTypesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeEc2API) DescribeInstanceTypesReturns(result1 *ec2.DescribeInstanceTypesOutput, result2 error) {
	fake.describeInstanceTypesMutex.Lock()
	defer fake.describeInstanceTypesMutex.Unlock()
	fake.DescribeInstanceTypesStub = nil
	fake.describeInstanceTypesReturns = struct {
		result1 *ec2.DescribeInstanceTypesOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) DescribeInstanceTypesReturnsOnCall(i int, result1 *ec2.DescribeInstanceTypesOutput, result2 error) {
	fake.describeInstanceTypesMutex.Lock()
	defer fake.describeInstanceTypesMutex.Unlock()
	fake.DescribeInstanceTypesStub = nil
	if fake.describeInstanceTypesReturnsOnCall == nil {
		fake.describeInstanceTypesReturnsOnCall = make(map[int]struct {
			result1 *ec2.DescribeInstanceTypesOutput
			result2 error
		})
	}
	fake.describeInstanceTypesReturnsOnCall[i] = struct {
		result1 *ec2.DescribeInstanceTypesOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) DescribeInstances(arg1 context.Context, arg2 *ec2.DescribeInstancesInput, arg3 ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	fake.describeInstancesMutex.Lock()
	ret, specificReturn := fake.describeInstancesReturnsOnCall[len(fake.describeInstancesArgsForCall)]
	fake.describeInstancesArgsForCall = append(fake.describeInstancesArgsForCall, struct {
		arg1 context.Context
		arg2 *ec2.DescribeInstancesInput
		arg3 []func(*ec2.Options)
	}{arg1, arg2, arg3})
	stub := fake.DescribeInstancesStub
	fakeReturns := fake.describeInstancesReturns
	fake.recordInvocation("DescribeInstances", []interface{}{arg1, arg2, arg3})
	fake.describeInstancesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeEc2API) DescribeInstancesCallCount() int {
	fake.describeInstancesMutex.RLock()
	defer fake.describeInstancesMutex.RUnlock()
	return len(fake.describeInstancesArgsForCall)
}

func (fake *FakeEc2API) DescribeInstancesCalls(stub func(context.Context, *ec2.DescribeInstancesInput, ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)) {
	fake.describeInstancesMutex.Lock()
	defer fake.describeInstancesMutex.Unlock()
	fake.DescribeInstancesStub = stub
}

func (fake *FakeEc2API) DescribeInstancesArgsForCall(i int) (context.Context, *ec2.DescribeInstancesInput, []func(*ec2.Options)) {
	fake.describeInstancesMutex.RLock()
	defer fake.describeInstancesMutex.RUnlock()
	argsForCall := fake.describeInstancesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeEc2API) DescribeInstancesReturns(result1 *ec2.DescribeInstancesOutput, result2 error) {
	fake.describeInstancesMutex.Lock()
	defer fake.describeInstancesMutex.Unlock()
	fake.DescribeInstancesStub = nil
	fake.describeInstancesReturns = struct {
		result1 *ec2.DescribeInstancesOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) DescribeInstancesReturnsOnCall(i int, result1 *ec2.DescribeInstancesOutput, result2 error) {
	fake.describeInstancesMutex.Lock()
	defer fake.describeInstancesMutex.Unlock()
	fake.DescribeInstancesStub = nil
	if fake.describeInstancesReturnsOnCall == nil {
		fake.describeInstancesReturnsOnCall = make(map[int]struct {
			result1 *ec2.DescribeInstancesOutput
			result2 error
		})
	}
	fake.describeInstancesReturnsOnCall[i] = struct {
		result1 *ec2.DescribeInstancesOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) ModifyVolume(arg1 context.Context, arg2 *ec2.ModifyVolumeInput, arg3 ...func(*ec2.Options)) (*ec2.ModifyVolumeOutput, error) {
	fake.modifyVolumeMutex.Lock()
	ret, specificReturn := fake.modifyVolumeReturnsOnCall[len(fake.modifyVolumeArgsForCall)]
	fake.modifyVolumeArgsForCall = append(fake.modifyVolumeArgsForCall, struct {
		arg1 context.Context
		arg2 *ec2.ModifyVolumeInput
		arg3 []func(*ec2.Options)
	}{arg1, arg2, arg3})
	stub := fake.ModifyVolumeStub
	fakeReturns := fake.modifyVolumeReturns
	fake.recordInvocation("ModifyVolume", []interface{}{arg1, arg2, arg3})
	fake.modifyVolumeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeEc2API) ModifyVolumeCallCount() int {
	fake.modifyVolumeMutex.RLock()
	defer fake.modifyVolumeMutex.RUnlock()
	return len(fake.modifyVolumeArgsForCall)
}

func (fake *FakeEc2API) ModifyVolumeCalls(stub func(context.Context, *ec2.ModifyVolumeInput, ...func(*ec2.Options)) (*ec2.ModifyVolumeOutput, error)) {
	fake.modifyVolumeMutex.Lock()
	defer fake.modifyVolumeMutex.Unlock()
	fake.ModifyVolumeStub = stub
}

func (fake *FakeEc2API) ModifyVolumeArgsForCall(i int) (context.Context, *ec2.ModifyVolumeInput, []func(*ec2.Options)) {
	fake.modifyVolumeMutex.RLock()
	defer fake.modifyVolumeMutex.RUnlock()
	argsForCall := fake.modifyVolumeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeEc2API) ModifyVolumeReturns(result1 *ec2.ModifyVolumeOutput, result2 error) {
	fake.modifyVolumeMutex.Lock()
	defer fake.modifyVolumeMutex.Unlock()
	fake.ModifyVolumeStub = nil
	fake.modifyVolumeReturns = struct {
		result1 *ec2.ModifyVolumeOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) ModifyVolumeReturnsOnCall(i int, result1 *ec2.ModifyVolumeOutput, result2 error) {
	fake.modifyVolumeMutex.Lock()
	defer fake.modifyVolumeMutex.Unlock()
	fake.ModifyVolumeStub = nil
	if fake.modifyVolumeReturnsOnCall == nil {
		fake.modifyVolumeReturnsOnCall = make(map[int]struct {
			result1 *ec2.ModifyVolumeOutput
			result2 error
		})
	}
	fake.modifyVolumeReturnsOnCall[i] = struct {
		result1 *ec2.ModifyVolumeOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) RunInstances(arg1 context.Context, arg2 *ec2.RunInstancesInput, arg3 ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error) {
	fake.runInstancesMutex.Lock()
	ret, specificReturn := fake.runInstancesReturnsOnCall[len(fake.runInstancesArgsForCall)]
	fake.runInstancesArgsForCall = append(fake.runInstancesArgsForCall, struct {
		arg1 context.Context
		arg2 *ec2.RunInstancesInput
		arg3 []func(*ec2.Options)
	}{arg1, arg2, arg3})
	stub := fake.RunInstancesStub
	fakeReturns := fake.runInstancesReturns
	fake.recordInvocation("RunInstances", []interface{}{arg1, arg2, arg3})
	fake.runInstancesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeEc2API) RunInstancesCallCount() int {
	fake.runInstancesMutex.RLock()
	defer fake.runInstancesMutex.RUnlock()
	return len(fake.runInstancesArgsForCall)
}

func (fake *FakeEc2API) RunInstancesCalls(stub func(context.Context, *ec2.RunInstancesInput, ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)) {
	fake.runInstancesMutex.Lock()
	defer fake.runInstancesMutex.Unlock()
	fake.RunInstancesStub = stub
}

func (fake *FakeEc2API) RunInstancesArgsForCall(i int) (context.Context, *ec2.RunInstancesInput, []func(*ec2.Options)) {
	fake.runInstancesMutex.RLock()
	defer fake.runInstancesMutex.RUnlock()
	argsForCall := fake.runInstancesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeEc2API) RunInstancesReturns(result1 *ec2.RunInstancesOutput, result2 error) {
	fake.runInstancesMutex.Lock()
	defer fake.runInstancesMutex.Unlock()
	fake.RunInstancesStub = nil
	fake.runInstancesReturns = struct {
		result1 *ec2.RunInstancesOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) RunInstancesReturnsOnCall(i int, result1 *ec2.RunInstancesOutput, result2 error) {
	fake.runInstancesMutex.Lock()
	defer fake.runInstancesMutex.Unlock()
	fake.RunInstancesStub = nil
	if fake.runInstancesReturnsOnCall == nil {
		fake.runInstancesReturnsOnCall = make(map[int]struct {
			result1 *ec2.RunInstancesOutput
			result2 error
		})
	}
	fake.runInstancesReturnsOnCall[i] = struct {
		result1 *ec2.RunInstancesOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) TerminateInstances(arg1 context.Context, arg2 *ec2.TerminateInstancesInput, arg3 ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error) {
	fake.terminateInstancesMutex.Lock()
	ret, specificReturn := fake.terminateInstancesReturnsOnCall[len(fake.terminateInstancesArgsForCall)]
	fake.terminateInstancesArgsForCall = append(fake.terminateInstancesArgsForCall, struct {
		arg1 context.Context
		arg2 *ec2.TerminateInstancesInput
		arg3 []func(*ec2.Options)
	}{arg1, arg2, arg3})
	stub := fake.TerminateInstancesStub
	fakeReturns := fake.terminateInstancesReturns
	fake.recordInvocation("TerminateInstances", []interface{}{arg1, arg2, arg3})
	fake.terminateInstancesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeEc2API) TerminateInstancesCallCount() int {
	fake.terminateInstancesMutex.RLock()
	defer fake.terminateInstancesMutex.RUnlock()
	return len(fake.terminateInstancesArgsForCall)
}

func (fake *FakeEc2API) TerminateInstancesCalls(stub func(context.Context, *ec2.TerminateInstancesInput, ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)) {
	fake.terminateInstancesMutex.Lock()
	defer fake.terminateInstancesMutex.Unlock()
	fake.TerminateInstancesStub = stub
}

func (fake *FakeEc2API) TerminateInstancesArgsForCall(i int) (context.Context, *ec2.TerminateInstancesInput, []func(*ec2.Options)) {
	fake.terminateInstancesMutex.RLock()
	defer fake.terminateInstancesMutex.RUnlock()
	argsForCall := fake.terminateInstancesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeEc2API) TerminateInstancesReturns(result1 *ec2.TerminateInstancesOutput, result2 error) {
	fake.terminateInstancesMutex.Lock()
	defer fake.terminateInstancesMutex.Unlock()
	fake.TerminateInstancesStub = nil
	fake.terminateInstancesReturns = struct {
		result1 *ec2.TerminateInstancesOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) TerminateInstancesReturnsOnCall(i int, result1 *ec2.TerminateInstancesOutput, result2 error) {
	fake.terminateInstancesMutex.Lock()
	defer fake.terminateInstancesMutex.Unlock()
	fake.TerminateInstancesStub = nil
	if fake.terminateInstancesReturnsOnCall == nil {
		fake.terminateInstancesReturnsOnCall = make(map[int]struct {
			result1 *ec2.TerminateInstancesOutput
			result2 error
		})
	}
	fake.terminateInstancesReturnsOnCall[i] = struct {
		result1 *ec2.TerminateInstancesOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeEc2API) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createTagsMutex.RLock()
	defer fake.createTagsMutex.RUnlock()
	fake.createVolumeMutex.RLock()
	defer fake.createVolumeMutex.RUnlock()
	fake.deleteVolumeMutex.RLock()
	defer fake.deleteVolumeMutex.RUnlock()
	fake.describeInstanceStatusMutex.RLock()
	defer fake.describeInstanceStatusMutex.RUnlock()
	fake.describeInstanceTypesMutex.RLock()
	defer fake.describeInstanceTypesMutex.RUnlock()
	fake.describeInstancesMutex.RLock()
	defer fake.describeInstancesMutex.RUnlock()
	fake.modifyVolumeMutex.RLock()
	defer fake.modifyVolumeMutex.RUnlock()
	fake.runInstancesMutex.RLock()
	defer fake.runInstancesMutex.RUnlock()
	fake.terminateInstancesMutex.RLock()
	defer fake.terminateInstancesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeEc2API) recordInvocation(key string, args []interface{}) {
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

var _ awsprov.Ec2API = new(FakeEc2API)
