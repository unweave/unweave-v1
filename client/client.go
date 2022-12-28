package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/unweave/unweave-v2/api"
)

type Config struct {
	ApiUrl string `json:"apiURL"`
	Token  string `json:"token"`
}

type Client struct {
	cfg    *Config
	client *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		cfg: &cfg,
	}
}

type RestRequestType string

const (
	Get    RestRequestType = http.MethodGet
	Post   RestRequestType = http.MethodPost
	Put    RestRequestType = http.MethodPut
	Delete RestRequestType = http.MethodDelete
)

type RestRequest struct {
	Url    string
	Header http.Header
	Body   io.Reader
	Type   RestRequestType
}

func (c *Client) NewRestRequest(rtype RestRequestType, endpoint string, params map[string]string) (
	*RestRequest, error,
) {
	query := ""
	for k, v := range params {
		query += fmt.Sprintf("%s=%v&", k, v)
	}

	url := fmt.Sprintf("%s/%s?%s", c.cfg.ApiUrl, endpoint, query)
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	return &RestRequest{
		Url:    url,
		Header: header,
		Body:   &bytes.Buffer{},
		Type:   rtype,
	}, nil
}

func (c *Client) NewAuthorizedRestRequest(rtype RestRequestType, endpoint string, params map[string]string) (
	*RestRequest, error,
) {
	req, err := c.NewRestRequest(rtype, endpoint, params)
	if err != nil {
		return nil, err
	}

	if c.cfg.Token == "" {
		return nil, fmt.Errorf("no token provided")
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	return req, nil
}

func (c *Client) ExecuteRest(ctx context.Context, req *RestRequest, resp interface{}) error {
	httpReq, err := http.NewRequest(string(req.Type), req.Url, req.Body)
	if err != nil {
		return err
	}

	httpReq = httpReq.WithContext(ctx)
	httpReq.Header = req.Header
	res, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, res.Body); err != nil {
		return fmt.Errorf("status %s, fail to read response body", res.Status)
	}
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		var msg api.HTTPError
		if err = json.NewDecoder(&buf).Decode(&msg); err != nil {
			return fmt.Errorf("status %s, fail to decode response body", res.Status)
		}
		return fmt.Errorf("status %s, %s", res.Status, msg.Message)
	}

	if err = json.NewDecoder(&buf).Decode(&resp); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to decode response body")
	}
	return nil
}
