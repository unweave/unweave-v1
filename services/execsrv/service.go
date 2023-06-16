package execsrv

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/infra/platform/posthog"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/volumesrv"
	"github.com/unweave/unweave/tools/random"
)

var (
	DefaultImageURI = "ubuntu:latest"
)

type ExecService struct {
	store                 Store
	driver                Driver
	volume                *volumesrv.VolumeService
	provider              types.Provider
	stateInformerFunc     StateInformerFunc
	statsInformerFunc     StatsInformerFunc
	heartbeatInformerFunc HeartbeatInformerFunc

	stateObserversFuncs     []StateObserverFunc
	statsObserversFuncs     []StatsObserverFunc
	heartbeatObserversFuncs []HeartbeatObserverFunc
}

func WithStateObserver(s *ExecService, f StateObserverFunc) *ExecService {
	s.stateObserversFuncs = append(s.stateObserversFuncs, f)
	return s
}

func WithStatsObserver(s *ExecService, f StatsObserverFunc) *ExecService {
	s.statsObserversFuncs = append(s.statsObserversFuncs, f)
	return s
}

func WithHeartbeatObserver(s *ExecService, f HeartbeatObserverFunc) *ExecService {
	s.heartbeatObserversFuncs = append(s.heartbeatObserversFuncs, f)
	return s
}

func NewService(
	store Store,
	driver Driver,
	volumeService *volumesrv.VolumeService,
	stateInformerFunc StateInformerFunc,
	statsInformerFunc StatsInformerFunc,
	heartbeatInformerFunc HeartbeatInformerFunc,
) (*ExecService, error) {
	s := &ExecService{
		store:                   store,
		driver:                  driver,
		volume:                  volumeService,
		provider:                driver.ExecProvider(),
		stateInformerFunc:       stateInformerFunc,
		statsInformerFunc:       statsInformerFunc,
		heartbeatInformerFunc:   heartbeatInformerFunc,
		stateObserversFuncs:     nil,
		statsObserversFuncs:     nil,
		heartbeatObserversFuncs: nil,
	}

	execs, err := store.List(nil, &s.provider, true)
	if err != nil {
		return nil, fmt.Errorf("failed to init StateInformer, failed list all execs: %w", err)
	}

	log.Info().
		Str(types.ProviderCtxKey, s.provider.String()).
		Msgf("Found %d existing execs", len(execs))

	for _, e := range execs {
		e := e

		ed, err := s.store.GetDriver(e.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to init StateInformer, failed get exec driver: %w", err)
		}

		if ed != s.driver.ExecDriverName() {
			continue
		}

		informer := s.stateInformerFunc(e)
		informer.Watch()

		for _, f := range s.stateObserversFuncs {
			o := f(e, informer)
			informer.Register(o)
		}
	}

	return s, nil
}

func (s *ExecService) Create(ctx context.Context, project string, creator string, params types.ExecCreateParams) (types.Exec, error) {
	image := DefaultImageURI
	if params.Image != nil && *params.Image != "" {
		image = *params.Image
	}

	spec := types.SetSpecDefaultValues(params.Spec)

	execID, err := s.driver.ExecCreate(ctx, project, image, spec, []string{params.SSHPublicKey}, nil)
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
		Status:    types.StatusPending,
		Command:   params.Command,
		Keys: []types.SSHKey{
			{
				Name:      params.SSHKeyName,
				PublicKey: &params.SSHPublicKey,
			},
		},
		Volumes:  nil,
		Spec:     spec,
		CommitID: params.CommitID,
		GitURL:   params.GitURL,
		Provider: params.Provider,
		// Set when a connection is established to the exec
		Network: types.ExecNetwork{},
		Region:  "",
	}

	if err = s.store.Create(project, exec); err != nil {
		return types.Exec{}, fmt.Errorf("failed to add exec to store: %w", err)
	}

	informer := s.stateInformerFunc(exec)
	informer.Watch()

	for _, so := range s.stateObserversFuncs {
		o := so(exec, informer)
		informer.Register(o)
	}

	posthog.Client.SendSessionCreated(exec.CreatedBy, exec.ProjectID, exec.CreatedAt)
	return exec, nil
}

func (s *ExecService) Get(ctx context.Context, id string) (types.Exec, error) {
	exec, err := s.store.Get(id)
	if err != nil {
		return types.Exec{}, err
	}
	return exec, err
}

func (s *ExecService) List(ctx context.Context, project string) ([]types.Exec, error) {
	execs, err := s.store.List(&project, &s.provider, false)
	if err != nil {
		return nil, err
	}
	return execs, nil
}

func (s *ExecService) Terminate(ctx context.Context, id string) error {
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

	if err = s.driver.ExecTerminate(ctx, exec.ID); err != nil {
		return fmt.Errorf("failed to terminate exec: %w", err)
	}

	if err = s.store.UpdateStatus(exec.ID, types.StatusTerminated); err != nil {
		return fmt.Errorf("failed to update exec status in store: %w", err)
	}

	posthog.Client.SendSessionTerminated(exec.CreatedBy, exec.ProjectID, exec.CreatedBy, exec.CreatedAt, time.Now())
	// TODO Clean up SSH keys associated with the terminated exec
	return nil
}

// Monitor starts monitoring an exec by registering observers to the stats and heartbeat
// informers.
func (s *ExecService) Monitor(ctx context.Context, execID string) error {
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
