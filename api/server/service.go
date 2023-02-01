package server

import (
	"github.com/google/uuid"
	"github.com/unweave/unweave/runtime"
)

func NewCtxService(rti runtime.Initializer, callerID uuid.UUID) *Service {
	srv := &Service{
		rti:     rti,
		cid:     callerID,
		Session: nil,
		SSHKey:  nil,
	}
	srv.Session = &SessionService{srv: srv}
	srv.SSHKey = &SSHKeyService{srv: srv}

	return srv
}
