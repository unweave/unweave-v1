//nolint:varnamelen
package router

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/unweave/unweave-v1/api/middleware"
	"github.com/unweave/unweave-v1/api/types"
	"github.com/unweave/unweave-v1/services/endpointsrv"
)

func NewEndpointRouter(endpoints endpointsrv.Service) *EndpointRouter {
	return &EndpointRouter{
		endpoints: endpoints,
	}
}

type EndpointRouter struct {
	endpoints endpointsrv.Service
}

func (e *EndpointRouter) EndpointRunCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	endpointID := chi.URLParam(r, "endpointRef")
	projectID := middleware.GetProjectIDFromContext(ctx)

	id, err := e.endpoints.RunEndpointEvals(ctx, projectID, endpointID)
	if err != nil {
		_ = render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to run endpoint evals"))

		return
	}

	render.JSON(w, r, types.EndpointCheckRun{CheckID: id})
}

func (e *EndpointRouter) EndpointCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := middleware.GetProjectIDFromContext(ctx)

	var req types.EndpointCreateParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, types.ErrHTTPBadRequest(err, "invalid request body"))

		return
	}

	endpoint, err := e.endpoints.EndpointExecCreate(ctx, projectID, req.ExecID, req.Name)
	if err != nil {
		_ = render.Render(w, r, types.ErrHTTPError(err, "create endpoint failed"))

		return
	}

	render.JSON(w, r, endpoint)
}

func (e *EndpointRouter) EndpointGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	endpointID := chi.URLParam(r, "endpointRef")
	projectID := middleware.GetProjectIDFromContext(ctx)

	endpoint, err := e.endpoints.EndpointGet(ctx, projectID, endpointID)
	if err != nil {
		_ = render.Render(w, r, types.ErrHTTPError(err, "get endpoint failed"))

		return
	}

	render.JSON(w, r, &types.EndpointGetResponse{Endpoint: endpoint})
}

func (e *EndpointRouter) EndpointList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := middleware.GetProjectIDFromContext(ctx)

	ends, err := e.endpoints.EndpointList(ctx, projectID)
	if err != nil {
		_ = render.Render(w, r, types.ErrHTTPError(err, "list endpoints"))

		return
	}

	render.JSON(w, r, types.EndpointList{Endpoints: ends})
}

func (e *EndpointRouter) EndpointEvalAttach(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	endpointID := chi.URLParam(r, "endpointRef")

	var req types.EndpointEvalAttach
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, types.ErrHTTPBadRequest(err, "invalid request body"))

		return
	}

	if err := e.endpoints.EndpointAttachEval(ctx, endpointID, req.EvalID); err != nil {
		_ = render.Render(w, r, types.ErrHTTPError(err, "attach eval"))

		return
	}
}

func (e *EndpointRouter) EndpointEvalCheckStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	checkID := chi.URLParam(r, "checkID")

	if checkID == "" {
		_ = render.Render(w, r, types.ErrHTTPBadRequest(errors.New("check id missing"), "Missing check-id from url path"))

		return
	}

	status, err := e.endpoints.EndpointCheckStatus(ctx, checkID)
	if err != nil {
		_ = render.Render(w, r, types.ErrHTTPError(err, "check status"))

		return
	}

	render.JSON(w, r, status)
}

func (e *EndpointRouter) EndpointCreateVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	endpointID := chi.URLParam(r, "endpointRef")
	projectID := middleware.GetProjectIDFromContext(ctx)

	var req types.EndpointVersionCreateParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, types.ErrHTTPBadRequest(err, "invalid request body"))

		return
	}

	version, err := e.endpoints.EndpointVersionCreate(ctx, projectID, endpointID, req.ExecID, req.Promote)
	if err != nil {
		_ = render.Render(w, r, types.ErrHTTPError(err, "create endpoint version"))

		return
	}

	render.JSON(w, r, version)
}
