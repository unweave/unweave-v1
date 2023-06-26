package execsrv

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/volumesrv"
	"github.com/unweave/unweave/tools/random"
)

var (
	DefaultImageURI = "ubuntu:latest"
)

type ExecService struct {
	store                    Store
	driver                   Driver
	volume                   *volumesrv.VolumeService
	provider                 types.Provider
	stateInformerManager     StateInformerManger
	statsInformerManager     StatsInformerManger
	heartbeatInformerManager HeartbeatInformerManger

	stateObserverFactories     []StateObserverFactory
	statsObserverFactories     []StatsObserverFactory
	heartbeatObserverFactories []HeartbeatObserverFactory
}

func WithStateObserver(s *ExecService, f StateObserverFactory) *ExecService {
	s.stateObserverFactories = append(s.stateObserverFactories, f)
	return s
}

func WithStatsObserver(s *ExecService, f StatsObserverFactory) *ExecService {
	s.statsObserverFactories = append(s.statsObserverFactories, f)
	return s
}

func WithHeartbeatObserver(s *ExecService, f HeartbeatObserverFactory) *ExecService {
	s.heartbeatObserverFactories = append(s.heartbeatObserverFactories, f)
	return s
}

func NewService(
	store Store,
	driver Driver,
	volumeService *volumesrv.VolumeService,
	stateInformerManager StateInformerManger,
	statsInformerManager StatsInformerManger,
	heartbeatInformerManager HeartbeatInformerManger,
) *ExecService {
	s := &ExecService{
		store:                      store,
		driver:                     driver,
		volume:                     volumeService,
		provider:                   driver.ExecProvider(),
		stateInformerManager:       stateInformerManager,
		statsInformerManager:       statsInformerManager,
		heartbeatInformerManager:   heartbeatInformerManager,
		stateObserverFactories:     nil,
		statsObserverFactories:     nil,
		heartbeatObserverFactories: nil,
	}

	return s
}

func (s *ExecService) Create(ctx context.Context, projectID string, creator string, params types.ExecCreateParams) (types.Exec, error) {
	image := DefaultImageURI
	if params.Image != nil && *params.Image != "" {
		image = *params.Image
	}

	volumes, err := s.parseVolumes(ctx, projectID, params.Volumes)
	if err != nil {
		return types.Exec{}, fmt.Errorf("volume verification failed: %w", err)
	}

	spec := types.SetSpecDefaultValues(params.Spec)

	network := types.ExecNetwork{}

	if params.InternalPort != 0 {
		network.HTTPService = &types.HTTPService{
			InternalPort: params.InternalPort,
		}
	}

	execID, err := s.driver.ExecCreate(ctx, projectID, image, spec, network, volumes, []string{params.SSHPublicKey}, nil)
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
		Volumes:  volumes,
		Spec:     spec,
		CommitID: params.CommitID,
		GitURL:   params.GitURL,
		Provider: params.Provider,
		// Full network details are filled in when
		// the Exec transitions to the Running state.
		Network: network,
		Region:  "",
	}

	if err = s.store.Create(projectID, exec); err != nil {
		return types.Exec{}, fmt.Errorf("failed to add exec to store: %w", err)
	}

	informer := s.stateInformerManager.Add(exec)
	informer.Watch()

	for _, factory := range s.stateObserverFactories {
		o := factory.New(exec)
		informer.Register(o)
	}

	return exec, nil
}

func (s *ExecService) Get(ctx context.Context, id string) (types.Exec, error) {
	exec, err := s.store.Get(id)
	if err != nil {
		return types.Exec{}, err
	}
	return exec, err
}

func (s *ExecService) Init() error {
	execs, err := s.store.List(nil, &s.provider, true)
	if err != nil {
		return fmt.Errorf("failed to init StateInformer, failed list all execs: %w", err)
	}

	log.Info().
		Str(types.ProviderCtxKey, s.provider.String()).
		Msgf("Found %d existing execs", len(execs))

	for _, exec := range execs {
		driver, err := s.store.GetDriver(exec.ID)
		if err != nil {
			return fmt.Errorf("failed to init StateInformer, failed get exec driver: %w", err)
		}

		if driver != s.driver.ExecDriverName() {
			continue
		}

		informer := s.stateInformerManager.Add(exec)
		informer.Watch()

		for _, factory := range s.stateObserverFactories {
			o := factory.New(exec)
			informer.Register(o)
		}

		if exec.Status == types.StatusRunning {
			if err = s.Monitor(context.Background(), exec.ID); err != nil {
				log.Error().
					Err(err).
					Str(types.ExecIDCtxKey, exec.ID).
					Msg("Failed to monitor exec")
			}
		}
	}
	return nil
}

func (s *ExecService) List(ctx context.Context, project string) ([]types.Exec, error) {
	execs, err := s.store.List(&project, &s.provider, false)
	if err != nil {
		return nil, err
	}
	return execs, nil
}

func (s *ExecService) parseVolumes(ctx context.Context, projectID string, volumes []types.VolumeAttachParams) ([]types.ExecVolume, error) {
	vols := make([]types.ExecVolume, len(volumes))

	for idx, v := range volumes {
		vol, err := s.volume.Get(ctx, projectID, v.VolumeRef)
		if err != nil {
			return nil, fmt.Errorf("failed to get volume %q: %w", v.VolumeRef, err)
		}

		vols[idx] = types.ExecVolume{
			VolumeID:  vol.ID,
			MountPath: v.MountPath,
		}
	}
	return vols, nil
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

	err = s.store.Delete(exec.ID)
	if err != nil {
		return fmt.Errorf("failed to delete shared volumes in store: %w", err)
	}

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

	stInformer := s.statsInformerManager.Add(exec)
	stInformer.Watch()

	for _, factory := range s.statsObserverFactories {
		o := factory.New(exec)
		stInformer.Register(o)
	}

	hbInformer := s.heartbeatInformerManager.Add(exec)
	hbInformer.Watch()

	for _, factory := range s.heartbeatObserverFactories {
		o := factory.New(exec)
		hbInformer.Register(o)
	}
	return nil
}
