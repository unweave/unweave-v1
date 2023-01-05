package client

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api"
)

type SSHKeyService struct {
	client *Client
}

func (s *SSHKeyService) Add(ctx context.Context, params api.SSHKeyAddParams) error {
	uri := fmt.Sprintf("ssh-keys")
	req, err := s.client.NewAuthorizedRestRequest(Post, uri, nil, params)
	if err != nil {
		return err
	}
	res := &api.SSHKeyAddResponse{}
	if err = s.client.ExecuteRest(ctx, req, res); err != nil {
		return err
	}
	return nil
}

func (s *SSHKeyService) List(ctx context.Context) ([]api.SSHKey, error) {
	uri := fmt.Sprintf("ssh-keys")
	req, err := s.client.NewAuthorizedRestRequest(Get, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	res := &api.SSHKeyListResponse{}
	if err = s.client.ExecuteRest(ctx, req, res); err != nil {
		return nil, err
	}
	return res.Keys, nil
}
