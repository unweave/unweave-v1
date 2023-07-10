package types

import "encoding/json"

type EndpointCreate struct {
	ExecID string `json:"execId"`
}

type EndpointEvalAttach struct {
	EvalID string `json:"evalId"`
}

type EndpointList struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type EvalList struct {
	Evals []Eval `json:"evals"`
}

type Endpoint struct {
	ID           string   `json:"id"`
	ProjectID    string   `json:"projectId"`
	ExecID       string   `json:"execId"`
	HTTPEndpoint string   `json:"httpEndpoint"`
	EvalIDs      []string `json:"evalIDs"`
}

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
	CheckID string
	Steps   []EndpointCheckStep
}

type EndpointCheckStep struct {
	StepID    string
	EvalID    string
	Input     json.RawMessage
	Output    json.RawMessage
	Assertion string
}
