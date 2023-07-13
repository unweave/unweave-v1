package types

import (
	"encoding/json"
	"time"
)

type EndpointCreate struct {
	Name   string `json:"name"`
	ExecID string `json:"execId"`
}

type EndpointVersionCreate struct {
	ExecID  string `json:"execId"`
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
	ProjectID   string    `json:"projectId"`
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
	ProjectID   string            `json:"projectId"`
	HTTPAddress string            `json:"httpAddress"`
	EvalIDs     []string          `json:"evalIDs"`
	Status      EndpointStatus    `json:"status"`
	Versions    []EndpointVersion `json:"versions"`
	CreatedAt   time.Time         `json:"createdAt"`
}

type EndpointVersion struct {
	ID          string         `json:"id"`
	ExecID      string         `json:"execId"`
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
	ExecID       string `json:"execId"`
	HTTPEndpoint string `json:"httpEndpoint"`
}

type EvalCreate struct {
	ExecID string `json:"execId"`
}

type EndpointCheckRun struct {
	CheckID string `json:"checkId"`
}

type EndpointCheck struct {
	CheckID    string
	Steps      []EndpointCheckStep
	Status     CheckStatus
	Conclusion *CheckConclusion `json:",omitempty"`
}

type EndpointCheckStep struct {
	StepID     string
	EvalID     string
	Input      json.RawMessage
	Output     json.RawMessage
	Assertion  string
	Status     CheckStatus
	Conclusion *CheckConclusion `json:",omitempty"`
}
