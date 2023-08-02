package execsrv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unweave/unweave-v1/api/types"
	"github.com/unweave/unweave-v1/services/execsrv"
	"github.com/unweave/unweave-v1/services/execsrv/internal/execsrvfakes"
)

func TestPostgresStore_UpdateConnectionInfo(t *testing.T) {
	t.Parallel()

	t.Run("should update connection info", func(t *testing.T) {
		t.Parallel()
		querier := new(execsrvfakes.FakeQuerier)
		pgStore := execsrv.NewPostgresStoreDB(querier)

		execID := "abc123"
		update := types.ConnectionInfo{
			Host: "hello world",
			User: "my-user",
			Port: 2222,
		}

		assert.NoError(t, pgStore.UpdateConnectionInfo(execID, update))

		_, params := querier.ExecUpdateConnectionInfoArgsForCall(0)
		want := `{"host": "hello world", "user": "my-user", "port": 2222}`
		got := params.ConnectionInfo
		assert.JSONEq(t, want, string(got))
	})
}
