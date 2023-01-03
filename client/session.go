package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/types"
)

type SessionService struct {
	client *Client
}

func (s *SessionService) Create(ctx context.Context, projectID uuid.UUID, params api.SessionCreateParams) (*types.Session, error) {
	uri := fmt.Sprintf("projects/%s/sessions", projectID)
	req, err := s.client.NewAuthorizedRestRequest(Post, uri, nil, params)
	if err != nil {
		return nil, err
	}
	session := &types.Session{}
	if err = s.client.ExecuteRest(ctx, req, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) Get(ctx context.Context, projectID, sessionID uuid.UUID) (*types.Session, error) {
	uri := fmt.Sprintf("projects/%s/sessions/%s", projectID, sessionID)
	req, err := s.client.NewAuthorizedRestRequest(Get, uri, nil, nil)
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

func (s *SessionService) Terminate(ctx context.Context, projectID, sessionID uuid.UUID) error {
	uri := fmt.Sprintf("projects/%s/sessions/%s/terminate", projectID, sessionID)
	req, err := s.client.NewAuthorizedRestRequest(Put, uri, nil, nil)
	if err != nil {
		return err
	}
	res := &api.SessionTerminateResponse{}
	return s.client.ExecuteRest(ctx, req, res)
}
