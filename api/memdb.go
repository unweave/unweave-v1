package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/unweave/unweave/types"
)

type sessionDB struct {
	mutex    *sync.Mutex
	sessions map[uuid.UUID]types.Session
}

func (db *sessionDB) Add(ctx context.Context, session types.Session) (sessionID uuid.UUID, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	db.sessions[sessionID] = session
	return sessionID, nil
}

func (db *sessionDB) Get(ctx context.Context, sessionID uuid.UUID) (session types.Session, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	session, ok := db.sessions[sessionID]
	if !ok {
		return types.Session{}, fmt.Errorf("session not found")
	}
	return session, nil
}

func (db *sessionDB) List(ctx context.Context) (sessions []types.Session, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for _, session := range db.sessions {
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (db *sessionDB) SetTerminated(ctx context.Context, sessionID uuid.UUID) (err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	delete(db.sessions, sessionID)
	return nil
}

type sshKeyDB struct {
	mutex   *sync.Mutex
	sshKeys map[uuid.UUID]SSHKey
}

func (db *sshKeyDB) Add(ctx context.Context, name string, publicKey string) (err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	id := uuid.New()
	db.sshKeys[id] = SSHKey{
		Name:      name,
		PublicKey: publicKey,
		CreatedAt: time.Now(),
	}
	return nil

}

func (db *sshKeyDB) Get(ctx context.Context, keyID uuid.UUID) (key SSHKey, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	key, ok := db.sshKeys[keyID]
	if !ok {
		return SSHKey{}, fmt.Errorf("key not found")
	}
	return key, nil
}

func (db *sshKeyDB) GetByName(ctx context.Context, name string) (key SSHKey, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for _, key := range db.sshKeys {
		key := key
		if key.Name == name {
			return key, nil
		}
	}
	return SSHKey{}, fmt.Errorf("key not found")
}

func (db *sshKeyDB) GetByPublicKey(ctx context.Context, publicKey string) (key SSHKey, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for _, key := range db.sshKeys {
		key := key
		if key.PublicKey == publicKey {
			return key, nil
		}
	}
	return SSHKey{}, fmt.Errorf("key not found")

}

func (db *sshKeyDB) List(ctx context.Context) (keys []SSHKey, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for _, key := range db.sshKeys {
		key := key
		keys = append(keys, key)
	}
	return keys, nil
}

func NewMemDB() *Store {
	mutex := &sync.Mutex{}
	return &Store{
		Session: &sessionDB{mutex: mutex},
		SSHKey:  &sshKeyDB{mutex: mutex},
	}
}
