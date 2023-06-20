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
	stateInformerFunc        StateInformerFunc
	statsInformerFunc        StatsInformerFunc
	heartbeatInformerManager *HeartbeatInformerManager

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
	heartbeatInformerManager *HeartbeatInformerManager,
) *ExecService {
	s := &ExecService{
		store:                    store,
		driver:                   driver,
		volume:                   volumeService,
		provider:                 driver.ExecProvider(),
		stateInformerFunc:        stateInformerFunc,
		statsInformerFunc:        statsInformerFunc,
		heartbeatInformerManager: heartbeatInformerManager,
		stateObserversFuncs:      nil,
		statsObserversFuncs:      nil,
		heartbeatObserversFuncs:  nil,
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

	execID, err := s.driver.ExecCreate(ctx, projectID, image, spec, volumes, []string{params.SSHPublicKey}, nil)
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
		// Set when a connection is established to the exec
		Network: types.ExecNetwork{},
		Region:  "",
	}

	if err = s.store.Create(projectID, exec); err != nil {
		return types.Exec{}, fmt.Errorf("failed to add exec to store: %w", err)
	}

	informer := s.stateInformerFunc(exec)
	informer.Watch()

	for _, so := range s.stateObserversFuncs {
		o := so(exec, informer)
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

		informer := s.stateInformerFunc(exec)
		informer.Watch()

		for _, f := range s.stateObserversFuncs {
			o := f(exec, informer)
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

	stInformer := s.statsInformerFunc(exec)
	stInformer.Watch()

	for _, so := range s.statsObserversFuncs {
		o := so(exec)
		stInformer.Register(o)
	}

	hbInformer := s.heartbeatInformerManager.Add(exec)
	hbInformer.Watch()

	for _, ho := range s.heartbeatObserversFuncs {
		o := ho(exec)
		hbInformer.Register(o)
	}
	return nil
}
