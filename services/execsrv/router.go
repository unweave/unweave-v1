//nolint:wrapcheck
package execsrv

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/types"
)

type Service interface {
	Provider() types.Provider
	Create(ctx context.Context, projectID string, creator string, params types.ExecCreateParams) (types.Exec, error)
	Get(ctx context.Context, execID string) (types.Exec, error)
	List(ctx context.Context, projectID string) ([]types.Exec, error)
	Terminate(ctx context.Context, execID string) error
	Monitor(ctx context.Context, execID string) error
	RefreshConnectionInfo(ctx context.Context, execID string) (types.Exec, error)
}

// DelegatingService is a service that routes requests to the correct provider. In most cases
// use should be using this service instead of provider specific services. This takes care
// of routing requests based on the provider and aggregating responses from multiple
// providers when needed.
type DelegatingService struct {
	store     Store
	delegates map[types.Provider]Service
}

func NewDelegatingService(store Store, services ...Service) Service {
	delegates := make(map[types.Provider]Service)

	for i := range services {
		svc := services[i]
		delegates[svc.Provider()] = svc
	}

	return &DelegatingService{
		store:     store,
		delegates: delegates,
	}
}

func (s *DelegatingService) Provider() types.Provider {
	panic("service router doesn't have a single provider")
}

func (s *DelegatingService) Create(
	ctx context.Context,
	projectID string,
	userID string,
	params types.ExecCreateParams,
) (types.Exec, error) {
	svc, err := s.service(params.Provider)
	if err != nil {
		return types.Exec{}, fmt.Errorf("establish service: %w", err)
	}

	return svc.Create(ctx, projectID, userID, params)
}

// Get returns a single session irrespective of the provider.
func (s *DelegatingService) Get(ctx context.Context, execID string) (types.Exec, error) {
	exec, err := s.store.Get(execID)
	if err != nil {
		return types.Exec{}, err
	}

	return exec, nil
}

// List returns a list of sessions for a given project irrespective of the providers.
func (s *DelegatingService) List(ctx context.Context, projectID string) ([]types.Exec, error) {
	execs, err := s.store.List(&projectID, nil, false)
	if err != nil {
		return nil, err
	}

	return execs, nil
}

// Terminate routes the exec termination request to the correct service based on the provider.
func (s *DelegatingService) Terminate(ctx context.Context, execID string) error {
	exec, err := s.store.Get(execID)
	if err != nil {
		return fmt.Errorf("failed to get exec: %w", err)
	}

	svc, err := s.service(exec.Provider)
	if err != nil {
		return fmt.Errorf("establish service: %w", err)
	}

	return svc.Terminate(ctx, execID)
}

// Monitor routes the exec monitoring request to the correct service based on the provider.
func (s *DelegatingService) Monitor(ctx context.Context, execID string) error {
	exec, err := s.store.Get(execID)
	if err != nil {
		return fmt.Errorf("failed to get exec: %w", err)
	}

	svc, err := s.service(exec.Provider)
	if err != nil {
		return fmt.Errorf("establish service: %w", err)
	}

	return svc.Monitor(ctx, execID)
}

func (s *DelegatingService) RefreshConnectionInfo(ctx context.Context, execID string) (types.Exec, error) {
	exec, err := s.store.Get(execID)
	if err != nil {
		return types.Exec{}, fmt.Errorf("failed to get exec: %w", err)
	}

	svc, err := s.service(exec.Provider)
	if err != nil {
		return types.Exec{}, fmt.Errorf("establish service: %w", err)
	}

	return svc.RefreshConnectionInfo(ctx, execID)
}

func (s *DelegatingService) service(provider types.Provider) (Service, error) {
	service, ok := s.delegates[provider]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}

	return service, nil
}
