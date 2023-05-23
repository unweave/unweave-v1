package exec

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/types"
)

type Store interface {
	Create(project, id, provider string) (types.Exec, error)
	Get(project, id string) (types.Exec, error)
	List(project string) ([]types.Exec, error)
	Delete(project, id string) error
}

type Driver interface {
	Create(ctx context.Context, project string, params types.ExecCreateParams) (types.Exec, error)
	Attach(ctx context.Context, id string) error
	GetNetwork(ctx context.Context, id string) (string, error)
	Terminate(ctx context.Context, id string) error
}

type Service struct {
	store   Store
	project string
	driver  Driver
}

func NewExecService(project string, store Store) *Service {
	return &Service{
		store:   store,
		project: project,
	}
}

func (s *Service) Create(ctx context.Context, project string, params types.ExecCreateParams) (types.Exec, error) {
	exec, err := s.driver.Create(ctx, project, params)
	if err != nil {
		return types.Exec{}, err
	}

	// TODO: parse exec image

	e, err := s.store.Create(s.project, exec.ID, params.Provider.String())
	if err != nil {
		return types.Exec{}, fmt.Errorf("failed to add exec to store: %w", err)
	}

	return e, nil
}
