package exec

import (
	"context"
	"fmt"
	"time"

	"github.com/unweave/unweave/api/types"
)

type Service struct {
	store         Store
	driver        Driver
	stateInformer StateInformer
	statsInformer StatsInformer

	stateObservers []func(exec types.Exec) StateObserver
	statsObservers []func(exec types.Exec) StatsObserver
}

func WithStateObserver(s *Service, f func(exec types.Exec) StateObserver) *Service {
	s.stateObservers = append(s.stateObservers, f)
	return s
}

func WithStatsObserver(s *Service, f func(exec types.Exec) StatsObserver) *Service {
	s.statsObservers = append(s.statsObservers, f)
	return s
}

func NewService(store Store, driver Driver) (*Service, error) {
	s := &Service{
		store:         store,
		driver:        driver,
		stateInformer: newStateInformer(store, driver),
	}

	go s.stateInformer.Watch()

	execs, err := store.ListAll()
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

		o := s.newStateObserver(e)
		s.stateInformer.Register(o)
	}

	return s, nil
}

func (s *Service) newStateObserver(exec types.Exec) StateObserver {
	return &stateObserver{exec: exec, srv: s}
}

func (s *Service) Create(ctx context.Context, project string, params types.ExecCreateParams) (types.Exec, error) {
	// TODO:
	// 	- Parse image and buildID
	//  - Parse network
	// 	- Parse volumes

	image := ""
	execID, err := s.driver.Create(ctx, project, image, params.Spec)
	if err != nil {
		return types.Exec{}, err
	}

	exec := types.Exec{
		ID:        execID,
		Name:      "",
		CreatedAt: time.Now(),
		CreatedBy: "",
		Image:     image,
		BuildID:   nil,
		Status:    types.StatusInitializing,
		Command:   params.Command,
		Keys:      nil,
		Volumes:   nil,
		Network:   types.ExecNetwork{},
		Spec:      types.HardwareSpec{},
		CommitID:  params.CommitID,
		GitURL:    params.GitURL,
		Region:    "", // Set later once
		Provider:  params.Provider,
	}
	if err = s.store.Create(project, exec); err != nil {
		return types.Exec{}, fmt.Errorf("failed to add exec to store: %w", err)
	}

	for _, so := range s.stateObservers {
		o := so(exec)
		s.stateInformer.Register(o)
	}

	for _, so := range s.statsObservers {
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
