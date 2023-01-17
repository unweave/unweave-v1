package client

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/types"
)

type AccountService struct {
	client *Client
}

func (a *AccountService) PairingTokenCreate(ctx context.Context) (code string, err error) {
	uri := fmt.Sprintf("account/pair")
	req, err := a.client.NewAuthorizedRestRequest(Post, uri, nil, nil)
	if err != nil {
		return "", err
	}
	res := &types.PairingTokenCreateResponse{}
	if err = a.client.ExecuteRest(ctx, req, res); err != nil {
		return "", err
	}
	return res.Code, nil
}

func (a *AccountService) PairingTokenExchange(ctx context.Context, code string) (token, email string, err error) {
	uri := fmt.Sprintf("account/pair/%s", code)
	req, err := a.client.NewAuthorizedRestRequest(Put, uri, nil, nil)
	if err != nil {
		return "", "", err
	}
	res := &types.PairingTokenExchangeResponse{}
	if err = a.client.ExecuteRest(ctx, req, res); err != nil {
		return "", "", err
	}
	return res.Token, res.Email, nil
}
