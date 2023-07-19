package endpointsrv

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"go.jetpack.io/typeid"
)

type endpointChecker struct {
	checkID string
	loaded  bool
	checks  []checkEndpointStep
}

func newEndpointChecker(ctx context.Context, checkID string) (*endpointChecker, error) {
	checker := &endpointChecker{
		checkID: checkID,
	}

	return checker, nil
}

func (c *endpointChecker) Run(ctx context.Context) error {
	if !c.loaded {
		panic("check steps must be loaded before calling run")
	}

	go func() {
		for _, check := range c.checks {
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

	return nil
}

var ErrEndpointUnavailable = fmt.Errorf("endpoint unavailable")
var ErrInvalidManifest = fmt.Errorf("manifest is invalid")
var ErrInvalidDataset = fmt.Errorf("dataset is invalid")

func fetchManifest(ctx context.Context, endpoint types.Endpoint, eval types.Eval) (types.EvalManifest, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+eval.HTTPEndpoint+"/", nil)
	if err != nil {
		return types.EvalManifest{}, fmt.Errorf("create request: %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return types.EvalManifest{}, fmt.Errorf("fetch manifest: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 {
		return types.EvalManifest{}, fmt.Errorf("fetch manifest: %w", ErrEndpointUnavailable)
	}

	if response.StatusCode >= 400 {
		return resolveManifestURLs(types.DefaultEvalManifest(), endpoint, eval)
	}

	if response.StatusCode != http.StatusOK {
		return types.EvalManifest{}, fmt.Errorf(
			"manifest endpoint status code must be 200, was %d: %w",
			response.StatusCode,
			ErrInvalidManifest,
		)
	}

	var manifest types.EvalManifest
	err = json.NewDecoder(response.Body).Decode(&manifest)
	if err != nil {
		return types.EvalManifest{}, fmt.Errorf("decoding manifest: %w", err)
	}

	return resolveManifestURLs(manifest, endpoint, eval)
}

// TODO: Security improvement needed to avoid calling loopback or internal endpoints
func resolveManifestURLs(manifest types.EvalManifest, endpoint types.Endpoint, eval types.Eval) (types.EvalManifest, error) {
	if manifest.RunURL == "" {
		manifest.RunURL = "https://" + endpoint.HTTPAddress + "/"
	}

	var err error
	manifest.RunURL, err = absoluteURL(manifest.RunURL, "https://"+eval.HTTPEndpoint)
	if err != nil {
		return types.EvalManifest{}, fmt.Errorf("absolute run url: %w", err)
	}

	manifest.DatasetURL, err = absoluteURL(manifest.DatasetURL, "https://"+eval.HTTPEndpoint+"/dataset")
	if err != nil {
		return types.EvalManifest{}, fmt.Errorf("absolute dataset url: %w", err)
	}

	manifest.AssertURL, err = absoluteURL(manifest.AssertURL, "https://"+eval.HTTPEndpoint+"/assert")
	if err != nil {
		return types.EvalManifest{}, fmt.Errorf("absolute assert url: %w", err)
	}

	return manifest, nil
}

func absoluteURL(ref string, base string) (string, error) {
	if base == "" {
		return "", fmt.Errorf("no base url provided")
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parse base url: %w", err)
	}

	if baseURL.Scheme == "" {
		return "", fmt.Errorf("base url must have a scheme")
	}

	refURL, err := url.Parse(ref)
	if err != nil {
		return "", fmt.Errorf("parse ref url: %w", err)
	}

	result := baseURL.ResolveReference(refURL)

	return result.String(), nil
}

func (c *endpointChecker) CreateCheckSteps(ctx context.Context, store Store, endpoint types.Endpoint, evals []types.Eval) error {
	var checks []checkEndpointStep

	for _, eval := range evals {
		manifest, err := fetchManifest(ctx, endpoint, eval)
		if err != nil {
			return fmt.Errorf("fetch manifest: %w", err)
		}

		d, err := fetchDataset(ctx, manifest.DatasetURL)
		if err != nil {
			return fmt.Errorf("fetch dataset: %w", err)
		}

		for _, item := range d.Data {
			stepID := typeid.Must(typeid.New("step")).String()

			checks = append(checks, checkEndpointStep{
				checkID:    c.checkID,
				stepID:     stepID,
				input:      item.Input,
				runPath:    manifest.RunURL,
				assertPath: manifest.AssertURL,
				endpoint:   "https://" + endpoint.HTTPAddress + "/",
				store:      store,
			})

			if err := store.EndpointCheckStepCreate(ctx, db.EndpointCheckStepCreateParams{
				ID:      stepID,
				CheckID: c.checkID,
				EvalID:  eval.ID,
				Input:   sql.NullString{String: string(item.Input), Valid: true},
			}); err != nil {
				return fmt.Errorf("create check step: %w", err)
			}
		}
	}

	c.checks = checks
	c.loaded = true
	return nil
}

type checkEndpointStep struct {
	err error

	checkID          string
	stepID           string
	runPath          string
	assertPath       string
	endpoint         string
	input            json.RawMessage
	endpointResponse json.RawMessage
	assertion        string
	store            Store
}

func (c *checkEndpointStep) callEndpoint(ctx context.Context) {
	buf := bytes.NewBuffer(c.input)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.runPath, buf)
	if err != nil {
		c.err = fmt.Errorf("build endpoint request: %w", err)

		return
	}
	req.Header.Set("X-Unweave-Target-Endpoint-URL", c.endpoint)

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.assertPath, buf)
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

	if resp.StatusCode != http.StatusOK {
		c.err = fmt.Errorf("assert status code: %d", resp.StatusCode)

		return
	}

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

func fetchDataset(ctx context.Context, datasetPath string) (dataset, error) {
	//nolint:gosec
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, datasetPath, nil)
	if err != nil {
		return dataset{}, fmt.Errorf("build dataset request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return dataset{}, fmt.Errorf("get dataset: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dataset{}, fmt.Errorf("get dataset: status %d %w", resp.StatusCode, ErrEndpointUnavailable)
	}

	var data dataset

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return dataset{}, fmt.Errorf("decode dataset - %w: %w", ErrInvalidDataset, err)
	}

	return data, nil
}
