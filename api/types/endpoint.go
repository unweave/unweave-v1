package types

import (
	"encoding/json"
	"time"
)

type EndpointCreateParams struct {
	Name   string `json:"name"`
	ExecID string `json:"execID"`
}

type EndpointVersionCreateParams struct {
	ExecID  string `json:"execID"`
	Promote bool   `json:"promote"`
}

type EndpointEvalAttach struct {
	EvalID string `json:"evalId"`
}

type EndpointList struct {
	Endpoints []EndpointListItem `json:"endpoints"`
}

type EndpointListItem struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	ProjectID   string    `json:"projectID"`
	HTTPAddress string    `json:"httpAddress"`
	CreatedAt   time.Time `json:"createdAt"`
}

type EvalList struct {
	Evals []Eval `json:"evals"`
}

type Endpoint struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Icon        string            `json:"icon"`
	ProjectID   string            `json:"projectID"`
	HTTPAddress string            `json:"httpAddress"`
	EvalIDs     []string          `json:"evalIDs"`
	Status      EndpointStatus    `json:"status"`
	Versions    []EndpointVersion `json:"versions"`
	CreatedAt   time.Time         `json:"createdAt"`
}

type EndpointGetResponse struct {
	Endpoint Endpoint `json:"endpoint"`
}

type EndpointVersion struct {
	ID          string         `json:"id"`
	ExecID      string         `json:"execID"`
	HTTPAddress string         `json:"httpAddress"`
	Status      EndpointStatus `json:"status"`
	Primary     bool           `json:"primary"`
	CreatedAt   time.Time      `json:"createdAt"`
}

type EndpointStatus string

const (
	EndpointStatusUnknown   EndpointStatus = ""
	EndpointStatusPending   EndpointStatus = "pending"
	EndpointStatusDeploying EndpointStatus = "deploying"
	EndpointStatusDeployed  EndpointStatus = "deployed"
	EndpointStatusFailed    EndpointStatus = "failed"
)

type Eval struct {
	ID           string `json:"id"`
	ExecID       string `json:"execID"`
	HTTPEndpoint string `json:"httpEndpoint"`
}

type EvalCreate struct {
	ExecID string `json:"execID"`
}

type EndpointCheckRun struct {
	CheckID string `json:"checkID"`
}

type EndpointCheck struct {
	CheckID    string
	Steps      []EndpointCheckStep
	Status     CheckStatus
	Conclusion *CheckConclusion `json:"conclusion,omitempty"`
}

type EndpointCheckStep struct {
	StepID     string
	EvalID     string
	Input      json.RawMessage
	Output     json.RawMessage
	Assertion  string
	Status     CheckStatus
	Conclusion *CheckConclusion `json:"conclusion,omitempty"`
}
