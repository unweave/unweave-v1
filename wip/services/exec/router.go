package exec

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/types"
)

type Service interface {
	Create(ctx context.Context, projectID string, creator string, params types.ExecCreateParams) (types.Exec, error)
	Get(ctx context.Context, execID string) (types.Exec, error)
	List(ctx context.Context, projectID string) ([]types.Exec, error)
	Terminate(ctx context.Context, execID string) error
	Monitor(ctx context.Context, execID string) error
}

// ServiceRouter is a service that routes requests to the correct provider. In most cases
// use should be using this service instead of provider specific services. This takes care
// of routing requests based on the provider and aggregating responses from multiple
// providers when needed.
type ServiceRouter struct {
	store          Store
	llService      *ProviderService
	unweaveService *ProviderService
}

func NewServiceRouter(lambdaLabsService, unweaveService *ProviderService) Service {
	return &ServiceRouter{
		llService:      lambdaLabsService,
		unweaveService: unweaveService,
	}
}

func (s *ServiceRouter) Create(ctx context.Context, project string, creator string, params types.ExecCreateParams) (types.Exec, error) {
	switch params.Provider {
	case types.LambdaLabsProvider:
		return s.llService.Create(ctx, project, creator, params)
	case types.UnweaveProvider:
		return s.unweaveService.Create(ctx, project, creator, params)
	default:
		return types.Exec{}, fmt.Errorf("unknown provider: %s", params.Provider)
	}
}

// Get returns a single session irrespective of the provider.
func (s *ServiceRouter) Get(ctx context.Context, execID string) (types.Exec, error) {
	// TODO: we probably want to clean this up somewhat and use a global store instead
	exec, err := s.llService.Get(ctx, execID)
	if err == nil {
		return exec, nil
	}
	return s.unweaveService.Get(ctx, execID)
}

// List returns a list of sessions for a given project irrespective of the providers.
func (s *ServiceRouter) List(ctx context.Context, projectID string) ([]types.Exec, error) {
	llExecs, err := s.llService.List(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list LambdaLabs execs: %w", err)
	}

	uwExecs, err := s.unweaveService.List(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list Unweave execs: %w", err)
	}

	execs := append(llExecs, uwExecs...)
	return execs, nil
}

// Terminate routes the exec termination request to the correct service based on the provider.
func (s *ServiceRouter) Terminate(ctx context.Context, execID string) error {
	exec, err := s.store.Get(execID)
	if err != nil {
		return fmt.Errorf("failed to get exec: %w", err)
	}

	switch exec.Provider {
	case types.LambdaLabsProvider:
		return s.llService.Terminate(ctx, execID)
	case types.UnweaveProvider:
		return s.unweaveService.Terminate(ctx, execID)
	default:
		return fmt.Errorf("unknown provider: %s", exec.Provider)
	}
}

// Monitor routes the exec monitoring request to the correct service based on the provider.
func (s *ServiceRouter) Monitor(ctx context.Context, execID string) error {
	exec, err := s.store.Get(execID)
	if err != nil {
		return fmt.Errorf("failed to get exec: %w", err)
	}

	switch exec.Provider {
	case types.LambdaLabsProvider:
		return s.llService.Monitor(ctx, execID)
	case types.UnweaveProvider:
		return s.unweaveService.Monitor(ctx, execID)
	default:
		return fmt.Errorf("unknown provider: %s", exec.Provider)
	}
}
