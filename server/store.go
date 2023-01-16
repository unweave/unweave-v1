package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/unweave/unweave/api"
)

type SessionStore interface {
	Add(ctx context.Context, session api.Session) (sessionID uuid.UUID, err error)
	Get(ctx context.Context, sessionID uuid.UUID) (session api.Session, err error)
	List(ctx context.Context) (sessions []api.Session, err error)
	SetTerminated(ctx context.Context, sessionID uuid.UUID) (err error)
}

type SSHKeyStore interface {
	Add(ctx context.Context, name string, publicKey string) (err error)
	Get(ctx context.Context, keyID uuid.UUID) (key api.SSHKey, err error)
	GetByName(ctx context.Context, name string) (key api.SSHKey, err error)
	GetByPublicKey(ctx context.Context, publicKey string) (key api.SSHKey, err error)
	List(ctx context.Context) (keys []api.SSHKey, err error)
}

type Store struct {
	Session SessionStore
	SSHKey  SSHKeyStore
}
