package server

import (
	"github.com/google/uuid"
	"github.com/unweave/unweave/runtime"
)

type Service struct {
	rti      runtime.Initializer
	cid      uuid.UUID
	Provider *ProviderService
	Session  *SessionService
	SSHKey   *SSHKeyService
}

func NewCtxService(rti runtime.Initializer, callerID uuid.UUID) *Service {
	srv := &Service{
		rti:      rti,
		cid:      callerID,
		Provider: nil,
		Session:  nil,
		SSHKey:   nil,
	}
	srv.Provider = &ProviderService{srv: srv}
	srv.Session = &SessionService{srv: srv}
	srv.SSHKey = &SSHKeyService{srv: srv}

	return srv
}
