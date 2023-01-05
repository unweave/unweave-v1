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
