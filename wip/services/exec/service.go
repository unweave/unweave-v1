package exec

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/tools/random"
)

var (
	DefaultImageURI = "ubuntu:latest"
)

type Service struct {
	store                 Store
	driver                Driver
	stateInformerFunc     StateInformerFunc
	statsInformerFunc     StatsInformerFunc
	heartbeatInformerFunc HeartbeatInformerFunc

	stateObserversFuncs     []StateObserverFunc
	statsObserversFuncs     []StatsObserverFunc
	heartbeatObserversFuncs []HeartbeatObserverFunc
}

func WithStateObserver(s *Service, f StateObserverFunc) *Service {
	s.stateObserversFuncs = append(s.stateObserversFuncs, f)
	return s
}

func WithStatsObserver(s *Service, f StatsObserverFunc) *Service {
	s.statsObserversFuncs = append(s.statsObserversFuncs, f)
	return s
}

func WithHeartbeatObserver(s *Service, f HeartbeatObserverFunc) *Service {
	s.heartbeatObserversFuncs = append(s.heartbeatObserversFuncs, f)
	return s
}

func NewService(
	store Store,
	driver Driver,
	stateInformerFunc StateInformerFunc,
	statsInformerFunc StatsInformerFunc,
	heartbeatInformerFunc HeartbeatInformerFunc,
) (*Service, error) {
	s := &Service{
		store:                 store,
		driver:                driver,
		stateInformerFunc:     stateInformerFunc,
		statsInformerFunc:     statsInformerFunc,
		heartbeatInformerFunc: heartbeatInformerFunc,
	}

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

		informer := s.stateInformerFunc(e)
		informer.Watch()

		for _, f := range s.stateObserversFuncs {
			o := f(e)
			informer.Register(o)
		}
	}

	return s, nil
}

func (s *Service) Create(ctx context.Context, project string, creator string, params types.ExecCreateParams) (types.Exec, error) {
	image := DefaultImageURI
	if params.Image != nil && *params.Image != "" {
		image = *params.Image
	}

	spec := types.SetSpecDefaultValues(params.Spec)

	// TODO: currently assumes only one SSH key - need to support multiple
	execID, err := s.driver.Create(ctx, project, image, spec, []string{params.SSHPublicKey}, nil)
	if err != nil {
		return types.Exec{}, err
	}

	log.Ctx(ctx).
		Info().
		Str(types.ExecIDCtxKey, execID).
		Msgf("Created new exec with image %q", image)

	exec := types.Exec{
		ID:        execID,
		Name:      random.GenerateRandomPhrase(4, "-"),
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
		Spec:     spec,
		CommitID: params.CommitID,
		GitURL:   params.GitURL,
		Region:   "", // Set later once the exec has been successfully scheduled
		Provider: params.Provider,
	}

	if err = s.store.Create(project, exec); err != nil {
		return types.Exec{}, fmt.Errorf("failed to add exec to store: %w", err)
	}

	informer := s.stateInformerFunc(exec)
	informer.Watch()

	for _, so := range s.stateObserversFuncs {
		o := so(exec)
		informer.Register(o)
	}

	return exec, nil
}

func (s *Service) Get(ctx context.Context, id string) (types.Exec, error) {
	exec, err := s.store.Get(id)
	if err != nil {
		return types.Exec{}, err
	}
	return exec, err
}

func (s *Service) List(ctx context.Context, project string) ([]types.Exec, error) {
	execs, err := s.store.List(project)
	if err != nil {
		return nil, err
	}
	return execs, nil
}

func (s *Service) Terminate(ctx context.Context, id string) error {
	exec, err := s.store.Get(id)
	if err != nil {
		if err == ErrNotFound {
			return nil
		}
	}

	// We don't want to overwrite the status if it's already terminated, failed or errored
	if exec.Status == types.StatusTerminated ||
		exec.Status == types.StatusFailed ||
		exec.Status == types.StatusError {
		return nil
	}

	log.Ctx(ctx).
		Info().
		Str(types.ExecIDCtxKey, exec.ID).
		Msg("Terminating exec")

	if err = s.driver.Terminate(ctx, exec.ID); err != nil {
		return fmt.Errorf("failed to terminate exec: %w", err)
	}

	if err = s.store.UpdateStatus(exec.ID, types.StatusTerminated); err != nil {
		return fmt.Errorf("failed to update exec status in store: %w", err)
	}
	return nil
}

// Monitor starts monitoring an exec by registering observers to the stats and heartbeat
// informers.
func (s *Service) Monitor(ctx context.Context, execID string) error {
	exec, err := s.store.Get(execID)
	if err != nil {
		return fmt.Errorf("failed to get exec from store: %w", err)
	}

	stInformer := s.statsInformerFunc(exec)
	stInformer.Watch()

	for _, so := range s.statsObserversFuncs {
		o := so(exec)
		stInformer.Register(o)
	}

	hbInformer := s.heartbeatInformerFunc(exec)
	hbInformer.Watch()

	for _, ho := range s.heartbeatObserversFuncs {
		o := ho(exec)
		hbInformer.Register(o)
	}
	return nil
}
