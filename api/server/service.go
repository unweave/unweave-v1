package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/runtime"
)

type CtxService interface {
	Initialize()
}

type Service struct {
	rti     runtime.Initializer
	cid     uuid.UUID // caller ID
	runtime runtime.Session
	builder runtime.Builder

	Builder  *BuilderService
	Provider *ProviderService
	Session  *SessionService
	SSHKey   *SSHKeyService
}

// InitializeRuntime initializes the runtime a caches it in memory.
func (s *Service) InitializeRuntime(ctx context.Context, provider types.RuntimeProvider) (runtime.Session, error) {
	if s.runtime != nil {
		return s.runtime, nil
	}
	rt, err := s.rti.InitializeRuntime(ctx, s.cid, provider)
	if err != nil {
		return nil, err
	}
	s.runtime = rt.Session
	return s.runtime, nil
}

func (s *Service) InitializerBuilder(ctx context.Context) (runtime.Builder, error) {
	if s.builder != nil {
		return s.builder, nil
	}
	builder, err := s.rti.InitializeBuilder(ctx, s.cid, "docker")
	if err != nil {
		return nil, err
	}
	s.builder = builder
	return s.builder, nil
}

func NewCtxService(rti runtime.Initializer, callerID uuid.UUID) *Service {
	srv := &Service{
		rti:      rti,
		cid:      callerID,
		Provider: nil,
		Session:  nil,
		SSHKey:   nil,
	}
	srv.Builder = &BuilderService{srv: srv}
	srv.Provider = &ProviderService{srv: srv}
	srv.Session = &SessionService{srv: srv}
	srv.SSHKey = &SSHKeyService{srv: srv}

	return srv
}
