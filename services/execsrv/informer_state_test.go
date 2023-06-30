package execsrv_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/execsrv"
	"github.com/unweave/unweave/services/execsrv/internal/execsrvfakes"
)

func TestPollingStateInformerManager(t *testing.T) {
	// Loosely in order:
	// - Store returns StatusInitializing
	// - Observer sees StatusInitializing and returns StatusRunning
	// - Driver called and returns StatusTerminated
	t.Parallel()

	storeDone := make(chan struct{})
	driverDone := make(chan struct{})

	store := new(execsrvfakes.FakeStore)
	store.GetCalls(func(s string) (types.Exec, error) {
		assert.Equal(t, s, "abc123")
		safeClose(storeDone)

		return types.Exec{ID: "abc123", Status: types.StatusInitializing}, nil
	})

	driver := new(execsrvfakes.FakeDriver)
	driver.ExecGetStatusCalls(func(_ context.Context, s string) (types.Status, error) {
		assert.Equal(t, s, "abc123")
		select {
		case <-storeDone:
			// if the store was called first
			safeClose(driverDone)

			return types.StatusTerminated, nil
		default:
			// if the store was not yet called
			return types.StatusInitializing, nil
		}
	})

	observer := new(execsrvfakes.FakeStateObserver)
	observer.UpdateCalls(func(state execsrv.State) execsrv.State {
		if state.Status == types.StatusInitializing {
			state.Status = types.StatusRunning

			return state
		}

		return state
	})

	manager := execsrv.NewPollingStateInformerManager(store, driver)
	manager.PollInterval = 500 * time.Millisecond

	exec := types.Exec{ID: "abc123"}

	informer := manager.Add(exec)
	informer.Register(observer)

	go informer.Watch()

	chans := []chan struct{}{storeDone, driverDone, waitForCalls(observer, 3)}
	<-waitAll(chans)

	assert.GreaterOrEqual(t, observer.UpdateCallCount(), 3)

	state0 := observer.UpdateArgsForCall(0)
	assert.Equal(t, state0, execsrv.State{Status: types.StatusInitializing})

	state1 := observer.UpdateArgsForCall(1)
	assert.Equal(t, state1, execsrv.State{Status: types.StatusRunning})

	state2 := observer.UpdateArgsForCall(2)
	assert.Equal(t, state2, execsrv.State{Status: types.StatusTerminated})
}

func waitAll(chans []chan struct{}) chan struct{} {
	out := make(chan struct{})

	go func() {
		defer close(out)

		for _, ch := range chans {
			<-ch
		}
	}()

	return out
}

func waitForCalls(observer *execsrvfakes.FakeStateObserver, calls int) chan struct{} {
	out := make(chan struct{})
	go func() {
		defer close(out)

		for {
			if observer.UpdateCallCount() >= calls {
				return
			}
		}
	}()

	return out
}
