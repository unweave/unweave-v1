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

func (s *SessionService) Create(ctx context.Context, projectID uuid.UUID, params api.SessionCreateRequestParams) (*types.Session, error) {
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

func (s *SessionService) Exec(ctx context.Context, cmd []string, image string, sessionID *uuid.UUID) (*types.Session, error) {
	return nil, nil
}

func (s *SessionService) Get(ctx context.Context, projectID, sessionID uuid.UUID) (*types.Session, error) {
	uri := fmt.Sprintf("projects/%s/sessions/%s", projectID, sessionID)
	req, err := s.client.NewAuthorizedRestRequest(Get, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	session := &types.Session{}
	if err = s.client.ExecuteRest(ctx, req, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) List(ctx context.Context, projectID uuid.UUID, listTerminated bool) ([]types.Session, error) {
	uri := fmt.Sprintf("projects/%s/sessions", projectID)
	query := map[string]string{
		"terminated": fmt.Sprintf("%t", listTerminated),
	}
	req, err := s.client.NewAuthorizedRestRequest(Get, uri, query, nil)
	if err != nil {
		return nil, err
	}
	res := &api.SessionsListResponse{}
	if err = s.client.ExecuteRest(ctx, req, res); err != nil {
		return nil, err
	}
	return res.Sessions, nil
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
