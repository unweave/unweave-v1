package server

import (
	"context"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/builder"
	"github.com/unweave/unweave/runtime"
)

type Service struct {
	rti     runtime.Initializer
	aid     string // account ID
	cid     string // caller ID
	runtime runtime.Session
	builder builder.Builder

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

func (s *Service) InitializeBuilder(ctx context.Context, builder string) (builder.Builder, error) {
	if s.builder != nil {
		return s.builder, nil
	}
	bld, err := s.rti.InitializeBuilder(ctx, s.cid, builder)
	if err != nil {
		return nil, err
	}
	s.builder = bld
	return s.builder, nil
}

func NewCtxService(rti runtime.Initializer, accountID, callerID string) *Service {
	srv := &Service{
		rti:      rti,
		aid:      accountID,
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
