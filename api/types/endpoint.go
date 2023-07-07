package types

type EndpointCreate struct {
	ExecID string `json:"exec_id"`
}

type EndpointEvalAttach struct {
	EvalID string `json:"eval_id"`
}

type EndpointList struct {
	Endpoints []Endpoint `json:"endpoints"`
}

//type Endpoint struct {
//	ID        string   `json:"id"`
//	ProjectID string   `json:"project_id"`
//	Exec      Exec     `json:"exec"`
//	EvalIDs   []string `json:"eval_ids"`
//}

type Endpoint struct {
	ID           string   `json:"id"`
	ProjectID    string   `json:"project_id"`
	ExecID       string   `json:"exec_id"`
	HTTPEndpoint string   `json:"http_endpoint"`
	EvalIDs      []string `json:"eval_ids"`
}

type Eval struct {
	ID   string `json:"id"`
	Exec Exec   `json:"exec"`
}

type EvalCreate struct {
	ExecID string `json:"exec_id"`
}
