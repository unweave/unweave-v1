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
	EvalListForProject(ctx context.Context, projectID string) ([]types.Eval, error)
	EvalCreate(ctx context.Context, projectID, execID string) (types.Eval, error)
}

type Store interface {
	EvalGet(ctx context.Context, id string) (db.EvalGetRow, error)
	EvalList(ctx context.Context, ids []string) ([]db.EvalListRow, error)
	EvalCreate(ctx context.Context, arg db.EvalCreateParams) error
	EvalListForProject(ctx context.Context, projectID string) ([]db.EvalListForProjectRow, error)
}

type EndpointDriver interface {
	EndpointCreate(
		ctx context.Context,
		project,
		endpointID,
		subdomain string) (string, error)

	EndpointVersionCreate(
		ctx context.Context,
		project,
		endpointID,
		versionID,
		execID string,
		internalPort int32) (string, error)
}

func NewEvalService(store Store, execService execsrv.Service, driver EndpointDriver) *EvalService {
	return &EvalService{
		store:       store,
		execService: execService,
		driver:      driver,
	}
}

type EvalService struct {
	store       Store
	execService execsrv.Service
	driver      EndpointDriver
}

func (e *EvalService) EvalCreate(ctx context.Context, projectID, execID string) (types.Eval, error) {
	evalID := typeid.Must(typeid.New("eval")).String()

	exec, err := e.execService.Get(ctx, execID)
	if err != nil {
		return types.Eval{}, fmt.Errorf("get exec: %w", err)
	}

	if exec.Provider != types.UnweaveProvider {
		return types.Eval{}, &types.Error{
			Code:       400,
			Message:    fmt.Sprintf("Cannot create eval for provider %q", exec.Provider),
			Suggestion: "Only unweave provider is supported for evals",
		}
	}

	validExecNetwork := exec.Network.HTTPService != nil &&
		exec.Network.HTTPService.InternalPort != 0

	if !validExecNetwork {
		return types.Eval{}, &types.Error{
			Code:       400,
			Message:    "Cannot create eval for exec with no port",
			Suggestion: "Create an exec exposing a port",
		}
	}

	addr, err := e.driver.EndpointVersionCreate(
		ctx,
		projectID,
		evalID,
		evalID,
		execID,
		exec.Network.HTTPService.InternalPort,
	)
	if err != nil {
		return types.Eval{}, fmt.Errorf("failed to create endpoint for eval: %w", err)
	}

	if err := e.store.EvalCreate(ctx, db.EvalCreateParams{
		ID:          evalID,
		ExecID:      exec.ID,
		HttpAddress: addr,
		ProjectID:   projectID,
	}); err != nil {
		return types.Eval{}, fmt.Errorf("create eval: %w", err)
	}

	return types.Eval{
		ID:           evalID,
		ExecID:       exec.ID,
		HTTPEndpoint: addr,
	}, nil
}

func (e *EvalService) Evals(ctx context.Context, ids []string) ([]types.Eval, error) {
	dbe, err := e.store.EvalList(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get evals: %w", err)
	}

	out := make([]types.Eval, len(dbe))

	for idx, eval := range dbe {
		out[idx] = types.Eval{
			ID:           eval.ID,
			ExecID:       eval.ExecID,
			HTTPEndpoint: eval.HttpAddress,
		}
	}

	return out, nil
}

func (e *EvalService) EvalListForProject(ctx context.Context, projectID string) ([]types.Eval, error) {
	rows, err := e.store.EvalListForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list evals: %w", err)
	}

	out := []types.Eval{}

	for _, row := range rows {
		out = append(out, types.Eval{
			ID:           row.ID,
			ExecID:       row.ExecID,
			HTTPEndpoint: row.HttpAddress,
		})
	}

	return out, nil
}
