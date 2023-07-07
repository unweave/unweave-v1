package evalsrv

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/services/execsrv"
	"go.jetpack.io/typeid"
)

type Service interface {
	Evals(ctx context.Context, ids []string) ([]types.Eval, error)
	EvalCreate(ctx context.Context, projectID, execID string) (types.Eval, error)
}

type Store interface {
	EvalGet(ctx context.Context, id string) (db.EvalGetRow, error)
	EvalList(ctx context.Context, ids []string) ([]db.EvalListRow, error)
	EvalCreate(ctx context.Context, arg db.EvalCreateParams) error
}

func NewEvalService(store Store, execService execsrv.Service) *EvalService {
	return &EvalService{
		store:       store,
		execService: execService,
	}
}

type EvalService struct {
	store       Store
	execService execsrv.Service
}

func (e *EvalService) EvalCreate(ctx context.Context, projectID, execID string) (types.Eval, error) {
	evalID := typeid.Must(typeid.New("eval")).String()

	exec, err := e.execService.Get(ctx, execID)
	if err != nil {
		return types.Eval{}, fmt.Errorf("get exec: %w", err)
	}

	if err := e.store.EvalCreate(ctx, db.EvalCreateParams{
		ID:        evalID,
		ExecID:    exec.ID,
		ProjectID: projectID,
	}); err != nil {
		return types.Eval{}, fmt.Errorf("create eval: %w", err)
	}

	return types.Eval{
		ID:   evalID,
		Exec: exec,
	}, nil
}

func (e *EvalService) Evals(ctx context.Context, ids []string) ([]types.Eval, error) {
	dbe, err := e.store.EvalList(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get evals: %w", err)
	}

	var out []types.Eval

	for _, eval := range dbe {
		// TODO: fix query in loop
		exec, err := e.execService.Get(ctx, eval.ExecID)
		if err != nil {
			return nil, fmt.Errorf("get exec: %w", err)
		}

		out = append(out, types.Eval{
			ID:   eval.ID,
			Exec: exec,
		})
	}

	return out, nil
}
