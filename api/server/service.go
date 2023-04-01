package server

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/builder"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/vault"
)

type Service struct {
	rti     runtime.Initializer
	aid     string // account ID
	cid     string // caller ID
	runtime *runtime.Runtime
	builder builder.Builder
	vault   vault.Vault

	Builder  *BuilderService
	Provider *ProviderService
	Exec     *ExecService
	SSHKey   *SSHKeyService
}

// InitializeRuntime initializes the runtime a caches it in memory.
func (s *Service) InitializeRuntime(ctx context.Context, provider types.Provider) (*runtime.Runtime, error) {
	if s.runtime != nil {
		return s.runtime, nil
	}
	rt, err := s.rti.InitializeRuntime(ctx, s.cid, provider)
	if err != nil {
		return nil, err
	}
	s.runtime = rt
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
	vlt, err := rti.InitializeVault(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to initialize vault: %v", err))
	}

	srv := &Service{
		rti:      rti,
		aid:      accountID,
		cid:      callerID,
		vault:    vlt,
		runtime:  nil,
		builder:  nil,
		Builder:  nil,
		Provider: nil,
		Exec:     nil,
		SSHKey:   nil,
	}
	srv.Builder = &BuilderService{srv: srv}
	srv.Provider = &ProviderService{srv: srv}
	srv.Exec = &ExecService{srv: srv}
	srv.SSHKey = &SSHKeyService{srv: srv}

	return srv
}
