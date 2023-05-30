package exec

import (
	"context"
	"fmt"
	"time"

	"github.com/unweave/unweave/api/types"
)

var (
	DefaultImageURI = "ubuntu:latest"
)

type Service struct {
	store         Store
	driver        Driver
	stateInformer StateInformer
	statsInformer StatsInformer

	stateObserversFuncs []StateObserverFunc
	statsObserversFuncs []StatsObserverFunc
}

func WithStateObserver(s *Service, f StateObserverFunc) *Service {
	s.stateObserversFuncs = append(s.stateObserversFuncs, f)
	return s
}

func WithStatsObserver(s *Service, f StatsObserverFunc) *Service {
	s.statsObserversFuncs = append(s.statsObserversFuncs, f)
	return s
}

func NewService(
	store Store,
	driver Driver,
	stateInformer StateInformer,
	statsInformer StatsInformer,
) (*Service, error) {
	s := &Service{
		store:         store,
		driver:        driver,
		stateInformer: stateInformer,
		statsInformer: statsInformer,
	}

	go s.stateInformer.Watch()

	execs, err := store.ListByProvider(driver.Provider(), true)
	if err != nil {
		return nil, fmt.Errorf("failed to init StateInformer, failed list all execs: %w", err)
	}

	for _, e := range execs {
		e := e

		ed, err := s.store.GetDriver(e.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to init StateInformer, failed get exec driver: %w", err)
		}

		if ed != s.driver.DriverName() {
			continue
		}

		for _, f := range s.stateObserversFuncs {
			o := f(e)
			s.stateInformer.Register(o)
		}
	}

	return s, nil
}

func (s *Service) Create(ctx context.Context, project string, creator string, params types.ExecCreateParams) (types.Exec, error) {
	image := DefaultImageURI
	if params.Image != nil && *params.Image != "" {
		image = *params.Image
	}

	// TODO: currently assumes only one SSH key - need to support multiple
	execID, err := s.driver.Create(ctx, project, image, params.Spec, []string{params.SSHPublicKey})
	if err != nil {
		return types.Exec{}, err
	}

	exec := types.Exec{
		ID:        execID,
		CreatedAt: time.Now(),
		CreatedBy: creator,
		Image:     image,
		BuildID:   nil,
		Status:    types.StatusInitializing,
		Command:   params.Command,
		Keys: []types.SSHKey{
			{
				Name:      params.SSHKeyName,
				PublicKey: &params.SSHPublicKey,
			},
		},
		Volumes:  nil,
		Network:  types.ExecNetwork{},
		Spec:     types.HardwareSpec{},
		CommitID: params.CommitID,
		GitURL:   params.GitURL,
		Region:   "", // Set later once the exec has been successfully scheduled
		Provider: params.Provider,
	}
	if err = s.store.Create(project, exec); err != nil {
		return types.Exec{}, fmt.Errorf("failed to add exec to store: %w", err)
	}

	for _, so := range s.stateObserversFuncs {
		o := so(exec)
		s.stateInformer.Register(o)
	}

	for _, so := range s.statsObserversFuncs {
		o := so(exec)
		s.statsInformer.Register(o)
	}

	return exec, nil
}

func (s *Service) Get(ctx context.Context, id string) (types.Exec, error) {
	return types.Exec{}, nil
}

func (s *Service) List(ctx context.Context, project string) ([]types.Exec, error) {
	return nil, nil
}

func (s *Service) Terminate(ctx context.Context, id string) error {
	return nil
}
