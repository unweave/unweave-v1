package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/unweave/unweave/api"
)

type Config struct {
	ApiURL string `json:"apiURL"`
	Token  string `json:"token"`
}

type Client struct {
	cfg    *Config
	client *http.Client

	Session *SessionService
	SSHKey  *SSHKeyService
}

func NewClient(cfg Config) *Client {
	c := &Client{
		cfg:    &cfg,
		client: &http.Client{},
	}
	c.Session = &SessionService{client: c}
	c.SSHKey = &SSHKeyService{client: c}
	return c
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

func (c *Client) NewRestRequest(rtype RestRequestType, endpoint string, params map[string]string, body interface{}) (
	*RestRequest, error,
) {
	query := ""
	for k, v := range params {
		query += fmt.Sprintf("%s=%v&", k, v)
	}

	url := fmt.Sprintf("%s/%s?%s", c.cfg.ApiURL, endpoint, query)
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	buf := &bytes.Buffer{}
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	return &RestRequest{
		Url:    url,
		Header: header,
		Body:   buf,
		Type:   rtype,
	}, nil
}

func (c *Client) NewAuthorizedRestRequest(rtype RestRequestType, endpoint string, query map[string]string, body interface{}) (
	*RestRequest, error,
) {
	req, err := c.NewRestRequest(rtype, endpoint, query, body)
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
		var errResp api.HTTPError
		if err = json.NewDecoder(&buf).Decode(&errResp); err != nil {
			return fmt.Errorf("status %s, fail to decode response body", res.Status)
		}
		return &api.HTTPError{
			Code:       errResp.Code,
			Message:    errResp.Message,
			Suggestion: errResp.Suggestion,
			Provider:   errResp.Provider,
		}
	}

	if err = json.NewDecoder(&buf).Decode(&resp); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to decode response body")
	}
	return nil
}
