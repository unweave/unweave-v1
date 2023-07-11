//nolint:noctx,godox
package endpointsrv

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/services/evalsrv"
	"github.com/unweave/unweave/services/execsrv"
	"go.jetpack.io/typeid"
)

type Driver interface {
	EndpointDriverName() string
	EndpointProvider() types.Provider
	EndpointCreate(ctx context.Context, project, endpointID, execID string, internalPort int32) (string, error)
}

type Service interface {
	EndpointExecCreate(ctx context.Context, projectID, execID string) (types.Endpoint, error)
	EndpointAttachEval(ctx context.Context, endpointID, evalID string) error
	EndpointList(ctx context.Context, projectID string) ([]types.Endpoint, error)
	EndpointGet(ctx context.Context, id string) (types.Endpoint, error)
	RunEndpointEvals(ctx context.Context, endpointID string) (string, error)
	EndpointCheckStatus(ctx context.Context, checkID string) (types.EndpointCheck, error)
}

type EndpointService struct {
	store  Store
	evals  evalsrv.Service
	execs  execsrv.Service
	driver Driver
}

type Store interface {
	EndpointCreate(ctx context.Context, arg db.EndpointCreateParams) error
	EndpointGet(ctx context.Context, id string) (db.UnweaveEndpoint, error)
	EndpointsForProject(ctx context.Context, id string) ([]db.UnweaveEndpoint, error)
	EndpointEval(ctx context.Context, endpointID string) ([]db.UnweaveEndpointEval, error)
	EndpointEvalAttach(ctx context.Context, arg db.EndpointEvalAttachParams) error
	EndpointCheckCreate(ctx context.Context, arg db.EndpointCheckCreateParams) error
	EndpointCheckStepCreate(ctx context.Context, arg db.EndpointCheckStepCreateParams) error
	EndpointCheckStepUpdate(ctx context.Context, arg db.EndpointCheckStepUpdateParams) error
	EndpointCheckSteps(ctx context.Context, checkID string) ([]db.UnweaveEndpointCheckStep, error)
	EndpointCheck(ctx context.Context, checkID string) (db.UnweaveEndpointCheck, error)
}

func NewEndpointService(
	store Store,
	evals evalsrv.Service,
	execs execsrv.Service,
	driver Driver,
) *EndpointService {
	return &EndpointService{
		store:  store,
		evals:  evals,
		execs:  execs,
		driver: driver,
	}
}

func (e *EndpointService) EndpointExecCreate(ctx context.Context, projectID, execID string) (types.Endpoint, error) {
	endpointID := typeid.Must(typeid.New("end")).String()

	exec, err := e.execs.Get(ctx, execID)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("get exec: %w", err)
	}

	if exec.Provider != e.driver.EndpointProvider() {
		return types.Endpoint{}, &types.Error{
			Message:    "Cannot create endpoint",
			Suggestion: fmt.Sprintf("Endpoints can only be created with the %s provider", e.driver.EndpointProvider().String()),
			Provider:   e.driver.EndpointProvider(),
		}
	}

	httpSerivice := exec.Network.HTTPService

	if httpSerivice == nil {
		return types.Endpoint{}, &types.Error{
			Code:       400,
			Message:    "Endpoints can only be created on execs that expose a service",
			Suggestion: "Try creating an exec with a port exposed",
		}
	}

	httpAddr, err := e.driver.EndpointCreate(
		ctx,
		projectID,
		endpointID,
		execID,
		httpSerivice.InternalPort,
	)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("drvier create endpoint: %w", err)
	}

	dbe := db.EndpointCreateParams{
		ID:          endpointID,
		ExecID:      execID,
		ProjectID:   projectID,
		HttpAddress: httpAddr,
		CreatedAt:   time.Now(),
	}

	if err := e.store.EndpointCreate(ctx, dbe); err != nil {
		return types.Endpoint{}, fmt.Errorf("save endpoint: %w", err)
	}

	return endpoint(endpointID, projectID, dbe.CreatedAt, httpAddr, execID, []string{}), nil
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

	out := []types.Endpoint{}

	for _, end := range dbEnds {
		// TODO: fix query in loop
		endEvals, err := e.store.EndpointEval(ctx, end.ID)
		if err != nil {
			return nil, fmt.Errorf("get endpoint evals: %w", err)
		}

		ids := make([]string, len(endEvals))
		for i, eval := range endEvals {
			ids[i] = eval.EvalID
		}

		out = append(out, endpoint(end.ID, end.ProjectID, end.CreatedAt, end.HttpAddress, end.ExecID, ids))
	}

	return out, nil
}

func (e *EndpointService) EndpointGet(ctx context.Context, id string) (types.Endpoint, error) {
	end, err := e.store.EndpointGet(ctx, id)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("get endpoint: %w", err)
	}

	endEvals, err := e.store.EndpointEval(ctx, end.ID)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("get endpoint evals: %w", err)
	}

	ids := make([]string, len(endEvals))
	for i, eval := range endEvals {
		ids[i] = eval.EvalID
	}

	return endpoint(end.ID, end.ProjectID, end.CreatedAt, end.HttpAddress, end.ExecID, ids), nil
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

	if err := e.store.EndpointCheckCreate(ctx, db.EndpointCheckCreateParams{
		ID:         checkID,
		EndpointID: endpoint.ID,
		ProjectID:  endpoint.ProjectID,
	}); err != nil {
		return "", fmt.Errorf("create eval check: %w", err)
	}

	checks, err := buildEndpointCheckSteps(ctx, checkID, e.store, endpoint, evals)
	if err != nil {
		return "", fmt.Errorf("build endpoint checks: %w", err)
	}

	go func() {
		ctx := context.Background()

		for _, check := range checks {
			check.callEndpoint(ctx)
			check.assertResponse(ctx)

			if check.err != nil {
				log.Error().Err(check.err).Msg("check failed")
			}

			log.Debug().
				Str("check_id", check.checkID).
				Str("step_id", check.stepID).
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

func buildEndpointCheckSteps(ctx context.Context, checkID string, store Store, endpoint types.Endpoint, evals []types.Eval) ([]checkEndpointStep, error) {
	var checks []checkEndpointStep

	for _, eval := range evals {
		datasetPath := "https://" + eval.HTTPEndpoint + "/dataset"
		assertPath := "https://" + eval.HTTPEndpoint + "/assert"

		d, err := fetchDataset(datasetPath)
		if err != nil {
			return nil, fmt.Errorf("fetch dataset: %w", err)
		}

		for _, item := range d.Data {
			stepID := typeid.Must(typeid.New("step")).String()

			checks = append(checks, checkEndpointStep{
				checkID:    checkID,
				stepID:     stepID,
				input:      item.Input,
				assertPath: assertPath,
				endpoint:   "https://" + endpoint.HTTPEndpoint,
				store:      store,
			})

			if err := store.EndpointCheckStepCreate(ctx, db.EndpointCheckStepCreateParams{
				ID:      stepID,
				CheckID: checkID,
				EvalID:  eval.ID,
				Input:   sql.NullString{String: string(item.Input), Valid: true},
			}); err != nil {
				return nil, fmt.Errorf("create check step: %w", err)
			}
		}
	}

	return checks, nil
}

type checkEndpointStep struct {
	err error

	checkID          string
	stepID           string
	assertPath       string
	endpoint         string
	input            json.RawMessage
	endpointResponse json.RawMessage
	assertion        string
	store            Store
}

func (c *checkEndpointStep) callEndpoint(ctx context.Context) {
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

	if err := c.store.EndpointCheckStepUpdate(ctx, db.EndpointCheckStepUpdateParams{
		ID:     sql.NullString{String: c.stepID, Valid: true},
		Output: sql.NullString{String: string(c.endpointResponse), Valid: true},
	}); err != nil {
		c.err = err
	}
}

func (c *checkEndpointStep) assertResponse(ctx context.Context) {
	buf := &bytes.Buffer{}

	if err := json.
		NewEncoder(buf).
		Encode(datasetItemEndpointResponse{
			EndpointResponse: c.endpointResponse,
		}); err != nil {
		c.err = fmt.Errorf("build assert body: %w", err)

		return
	}

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
	if err := c.store.EndpointCheckStepUpdate(ctx, db.EndpointCheckStepUpdateParams{
		ID:        sql.NullString{String: c.stepID, Valid: true},
		Assertion: sql.NullString{String: c.assertion, Valid: true},
	}); err != nil {
		c.err = err
	}
}

type dataset struct {
	Data []datasetItem `json:"data"`
}

type datasetItem struct {
	Input json.RawMessage `json:"input"`
}

type datasetItemEndpointResponse struct {
	EndpointResponse json.RawMessage `json:"endpointResponse"`
}

func endpoint(endpointID, projectID string, createdAt time.Time, httpAddress, execID string, ids []string) types.Endpoint {
	endpoint := types.Endpoint{
		ID:        endpointID,
		ProjectID: projectID,
		ExecID:    execID,
		EvalIDs:   ids,
		CreatedAt: createdAt,
		Status:    types.EndpointStatusDeployed,
	}

	endpoint.HTTPEndpoint = httpAddress

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

func (e *EndpointService) EndpointCheckStatus(ctx context.Context, checkID string) (types.EndpointCheck, error) {
	steps, err := e.store.EndpointCheckSteps(ctx, checkID)
	if err != nil {
		return types.EndpointCheck{}, fmt.Errorf("get steps: %w", err)
	}

	out := make([]types.EndpointCheckStep, len(steps))
	for idx, step := range steps {
		out[idx] = types.EndpointCheckStep{
			StepID:    step.ID,
			EvalID:    step.EvalID,
			Input:     []byte(step.Input.String),
			Output:    []byte(step.Output.String),
			Assertion: step.Assertion.String,
		}
	}

	return types.EndpointCheck{
		CheckID: checkID,
		Steps:   out,
	}, nil
}
