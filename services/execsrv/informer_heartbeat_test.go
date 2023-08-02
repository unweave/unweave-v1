package execsrv_test

import (
	"testing"
	"time"

	"github.com/unweave/unweave-v1/api/types"
	"github.com/unweave/unweave-v1/services/execsrv"
	"github.com/unweave/unweave-v1/services/execsrv/internal/execsrvfakes"
)

func TestPollingHeartbeatInformerManager(t *testing.T) {
	t.Parallel()

	driver := new(execsrvfakes.FakeDriver)

	done := make(chan struct{})
	observer := new(execsrvfakes.FakeHeartbeatObserver)
	observer.UpdateCalls(func(h execsrv.Heartbeat) {
		safeClose(done)
	})

	manager := execsrv.NewPollingHeartbeatInformerManager(driver, 1)
	manager.PollInterval = 5 * time.Millisecond

	exec := types.Exec{ID: "abc123"}

	informer := manager.Add(exec)
	informer.Register(observer)

	go informer.Watch()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("should have passed poll event to observer")
	}

	if driver.ExecGetStatusCallCount() == 0 {
		t.Error("should have called get status")
	}

	if _, id := driver.ExecGetStatusArgsForCall(0); id != "abc123" {
		t.Error("should have called with correct exec id")
	}
}

func safeClose(done chan struct{}) {
	select {
	case _, ok := <-done:
		if !ok {
			return
		}

		close(done)
	default:
		close(done)
	}
}
