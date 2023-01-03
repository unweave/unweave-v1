package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/types"
)

type SessionService struct {
	client *Client
}

func (s *SessionService) Create(ctx context.Context, params api.SessionCreateParams) (*types.Session, error) {
	req, err := s.client.NewAuthorizedRestRequest(Post, "sessions", nil, params)
	if err != nil {
		return nil, err
	}
	session := &types.Session{}
	if err = s.client.ExecuteRest(ctx, req, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) Get(ctx context.Context, sessionID uuid.UUID) (*types.Session, error) {
	req, err := s.client.NewAuthorizedRestRequest(Get, "sessions/"+sessionID.String(), nil, nil)
	if err != nil {
		return nil, err
	}
	var session *types.Session
	if err = s.client.ExecuteRest(ctx, req, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) Exec(ctx context.Context, cmd []string, image string, sessionID *uuid.UUID) (*types.Session, error) {

	return nil, nil
}

func (s *SessionService) Terminate(ctx context.Context, sessionID uuid.UUID) error {
	req, err := s.client.NewAuthorizedRestRequest(Put, "sessions/"+sessionID.String()+"/terminate", nil, nil)
	if err != nil {
		return err
	}
	res := &api.SessionTerminateResponse{}
	return s.client.ExecuteRest(ctx, req, res)
}
