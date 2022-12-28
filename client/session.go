package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/unweave/unweave-v2/api"
)

type SessionService struct {
	client *Client
}

func (s *SessionService) CreateSession(ctx context.Context, params api.SessionCreateParams) (*api.Session, error) {
	req, err := s.client.NewAuthorizedRestRequest(Post, "session", nil, params)
	if err != nil {
		return nil, err
	}
	var session *api.Session
	if err = s.client.ExecuteRest(ctx, req, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) GetSession(ctx context.Context, sessionID uuid.UUID) (*api.Session, error) {
	req, err := s.client.NewAuthorizedRestRequest(Get, "session/"+sessionID.String(), nil, nil)
	if err != nil {
		return nil, err
	}
	var session *api.Session
	if err = s.client.ExecuteRest(ctx, req, session); err != nil {
		return nil, err
	}
	return session, nil
}
