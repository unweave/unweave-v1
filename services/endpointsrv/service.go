//nolint:noctx,godox
package endpointsrv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/services/evalsrv"
	"github.com/unweave/unweave/services/execsrv"
	"go.jetpack.io/typeid"
)

type Service interface {
	EndpointExecCreate(ctx context.Context, projectID, execID string) (types.Endpoint, error)
	EndpointAttachEval(ctx context.Context, endpointID, evalID string) error
	EndpointList(ctx context.Context, projectID string) ([]types.Endpoint, error)
	EndpointGet(ctx context.Context, id string) (types.Endpoint, error)
	RunEndpointEvals(ctx context.Context, endpointID string) (string, error)
}

type EndpointService struct {
	store Store
	evals evalsrv.Service
	execs execsrv.Service
}

type Store interface {
	EndpointCreate(ctx context.Context, arg db.EndpointCreateParams) error
	EndpointGet(ctx context.Context, id string) (db.UnweaveEndpoint, error)
	EndpointsForProject(ctx context.Context, id string) ([]db.UnweaveEndpoint, error)
	EndpointEval(ctx context.Context, endpointID string) ([]db.UnweaveEndpointEval, error)
	EndpointEvalAttach(ctx context.Context, arg db.EndpointEvalAttachParams) error
}

func NewEndpointService(
	store Store,
	evals evalsrv.Service,
	execs execsrv.Service,
) *EndpointService {
	return &EndpointService{
		store: store,
		evals: evals,
		execs: execs,
	}
}

func (e *EndpointService) EndpointExecCreate(ctx context.Context, projectID, execID string) (types.Endpoint, error) {
	// TODO; we probably want to change the attachment of a service onto an "endpoint"
	// and make it's hostname specific to the endpoint, not the exec
	endpointID := typeid.Must(typeid.New("end")).String()

	exec, err := e.execs.Get(ctx, execID)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("get exec: %w", err)
	}

	if exec.Network.HTTPService == nil {
		return types.Endpoint{}, &types.Error{
			Code:       400,
			Message:    "Endpoints can only be created on execs that expose a service",
			Suggestion: "Try creating an exec with a port exposed",
		}
	}

	dbe := db.EndpointCreateParams{
		ID:        endpointID,
		ExecID:    execID,
		ProjectID: projectID,
	}

	if err := e.store.EndpointCreate(ctx, dbe); err != nil {
		return types.Endpoint{}, fmt.Errorf("save endpoint: %w", err)
	}

	return endpoint(endpointID, projectID, exec, []string{}), nil
}

func (e *EndpointService) EndpointAttachEval(ctx context.Context, endpointID, evalID string) error {
	if err := e.store.EndpointEvalAttach(ctx, db.EndpointEvalAttachParams{
		EndpointID: endpointID,
		EvalID:     evalID,
	}); err != nil {
		return fmt.Errorf("attach eval: %w", err)
	}

	return nil
}

func (e *EndpointService) EndpointList(ctx context.Context, projectID string) ([]types.Endpoint, error) {
	dbEnds, err := e.store.EndpointsForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get endpoint: %w", err)
	}

	projectExecs, err := e.execs.List(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list exec: %w", err)
	}

	execByID := make(map[string]types.Exec)

	for _, pe := range projectExecs {
		pe := pe
		execByID[pe.ID] = pe
	}

	out := []types.Endpoint{}

	for _, end := range dbEnds {
		exec, ok := execByID[end.ExecID]

		if !ok {
			continue
		}

		// TODO: fix query in loop
		endEvals, err := e.store.EndpointEval(ctx, end.ID)
		if err != nil {
			return nil, fmt.Errorf("get endpoint evals: %w", err)
		}

		ids := make([]string, len(endEvals))
		for i, eval := range endEvals {
			ids[i] = eval.EvalID
		}

		out = append(out, endpoint(end.ID, end.ProjectID, exec, ids))
	}

	return out, nil
}

func (e *EndpointService) EndpointGet(ctx context.Context, id string) (types.Endpoint, error) {
	end, err := e.store.EndpointGet(ctx, id)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("get endpoint: %w", err)
	}

	exec, err := e.execs.Get(ctx, end.ExecID)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("get exec: %w", err)
	}

	endEvals, err := e.store.EndpointEval(ctx, end.ID)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("get endpoint evals: %w", err)
	}

	ids := make([]string, len(endEvals))
	for i, eval := range endEvals {
		ids[i] = eval.EvalID
	}

	return endpoint(end.ID, end.ProjectID, exec, ids), nil
}

func (e *EndpointService) RunEndpointEvals(ctx context.Context, endpointID string) (string, error) {
	checkID := typeid.Must(typeid.New("check")).String()

	endpoint, err := e.EndpointGet(ctx, endpointID)
	if err != nil {
		return "", fmt.Errorf("get endpoint: %w", err)
	}

	evals, err := e.evals.Evals(ctx, endpoint.EvalIDs)
	if err != nil {
		return "", fmt.Errorf("get evals: %w", err)
	}

	if err := verifyCanRunChecks(endpoint, evals); err != nil {
		return "", fmt.Errorf("verify checks: %w", err)
	}

	checks, err := buildEndpointCheckSteps(endpoint, evals)
	if err != nil {
		return "", fmt.Errorf("build endpoint checks: %w", err)
	}

	go func() {
		for _, check := range checks {
			check.callEndpoint()
			check.assertResponse()

			if check.err != nil {
				log.Error().Err(check.err).Msg("check failed")
			}

			// TODO: store check progress / data somewhere
			// make it pollable
			log.Info().
				Str("check_id", checkID).
				Str("input", string(check.input)).
				Str("response", string(check.endpointResponse)).
				Str("assertion", check.assertion).
				Send()
		}
	}()

	return checkID, nil
}

func verifyCanRunChecks(endpoint types.Endpoint, evals []types.Eval) error {
	if endpoint.HTTPEndpoint == "" {
		return errors.New("endpoint must be exposed and have hostname attached")
	}

	if !allHaveHTTPServiceHostname(evals...) {
		return errors.New("evals must have http service exposed")
	}

	if len(evals) == 0 {
		return errors.New("no evals on endpoint")
	}

	return nil
}

func buildEndpointCheckSteps(endpoint types.Endpoint, evals []types.Eval) ([]checkEndpointStep, error) {
	var checks []checkEndpointStep

	for _, eval := range evals {
		datasetPath := "https://" + eval.HTTPEndpoint + "/dataset"
		assertPath := "https://" + eval.HTTPEndpoint + "/assert"

		d, err := fetchDataset(datasetPath)
		if err != nil {
			return nil, fmt.Errorf("fetch dataset: %w", err)
		}

		for _, item := range d.Data {
			checks = append(checks, checkEndpointStep{
				input:      item.Input,
				assertPath: assertPath,
				endpoint:   "https://" + endpoint.HTTPEndpoint,
			})
		}
	}

	return checks, nil
}

type checkEndpointStep struct {
	err error

	assertPath       string
	endpoint         string
	input            json.RawMessage
	endpointResponse json.RawMessage
	assertion        string
}

func (c *checkEndpointStep) callEndpoint() {
	buf := bytes.NewBuffer(c.input)

	req, err := http.NewRequest(http.MethodPost, c.endpoint, buf)
	if err != nil {
		c.err = fmt.Errorf("build endpoint request: %w", err)

		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.err = fmt.Errorf("call endpoint: %w", err)

		return
	}

	defer resp.Body.Close()

	c.endpointResponse, c.err = io.ReadAll(resp.Body)
}

func (c *checkEndpointStep) assertResponse() {
	buf := bytes.NewBuffer(c.endpointResponse)

	req, err := http.NewRequest(http.MethodPost, c.assertPath, buf)
	if err != nil {
		c.err = fmt.Errorf("build assert request: %w", err)

		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.err = fmt.Errorf("call assert: %w", err)

		return
	}

	defer resp.Body.Close()

	var response struct {
		Result string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.err = fmt.Errorf("decode assert response: %w", err)

		return
	}

	c.assertion = response.Result
}

type dataset struct {
	Data []datasetItem `json:"data"`
}

type datasetItem struct {
	Input json.RawMessage `json:"input"`
}

func endpoint(endpointID, projectID string, exec types.Exec, ids []string) types.Endpoint {
	endpoint := types.Endpoint{
		ID:        endpointID,
		ProjectID: projectID,
		ExecID:    exec.ID,
		EvalIDs:   ids,
	}

	if exec.Network.HTTPService != nil && exec.Network.HTTPService.Hostname != "" {
		endpoint.HTTPEndpoint = exec.Network.HTTPService.Hostname
	}

	return endpoint
}

func fetchDataset(datasetPath string) (dataset, error) {
	//nolint:gosec
	resp, err := http.Get(datasetPath)
	if err != nil {
		return dataset{}, fmt.Errorf("get dataset: %w", err)
	}

	defer resp.Body.Close()

	var data dataset

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return dataset{}, fmt.Errorf("decode dataset: %w", err)
	}

	return data, nil
}

func allHaveHTTPServiceHostname(evals ...types.Eval) bool {
	for _, e := range evals {
		if e.HTTPEndpoint == "" {
			return false
		}
	}

	return true
}
