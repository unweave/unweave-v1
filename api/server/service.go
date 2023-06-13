package server

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/builder"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/vault"
)

type Service struct {
	rti     runtime.Initializer
	aid     string // account ID
	cid     string // caller ID
	builder builder.Builder
	vault   vault.Vault

	Builder *BuilderService
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
		rti:     rti,
		aid:     accountID,
		cid:     callerID,
		vault:   vlt,
		builder: nil,
		Builder: nil,
	}
	srv.Builder = &BuilderService{srv: srv}

	return srv
}
