package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/unweave/unweave/types"
)

type SessionStore interface {
	Add(ctx context.Context, session types.Session) (sessionID uuid.UUID, err error)
	Get(ctx context.Context, sessionID uuid.UUID) (session types.Session, err error)
	List(ctx context.Context) (sessions []types.Session, err error)
	SetTerminated(ctx context.Context, sessionID uuid.UUID) (err error)
}

type SSHKeyStore interface {
	Add(ctx context.Context, name string, publicKey string) (err error)
	Get(ctx context.Context, keyID uuid.UUID) (key SSHKey, err error)
	GetByName(ctx context.Context, name string) (key SSHKey, err error)
	GetByPublicKey(ctx context.Context, publicKey string) (key SSHKey, err error)
	List(ctx context.Context) (keys []SSHKey, err error)
}

type Store struct {
	Session SessionStore
	SSHKey  SSHKeyStore
}
