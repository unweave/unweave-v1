package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/endpointsrv"
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

	id, err := e.endpoints.RunEndpointEvals(ctx, endpointID)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to run endpoint evals"))
		return
	}

	render.JSON(w, r, struct {
		CheckID string `json:"check_id"`
	}{CheckID: id})
}

func (e *EndpointRouter) EndpointCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := chi.URLParam(r, "project")

	var req types.EndpointCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Render(w, r, types.ErrHTTPBadRequest(err, "invalid request body"))
		return
	}

	endpoint, err := e.endpoints.EndpointExecCreate(ctx, projectID, req.ExecID)
	if err != nil {
		render.Render(w, r, types.ErrInternalServer(err, "create endpoint failed"))
		return
	}

	render.JSON(w, r, endpoint)
}

func (e *EndpointRouter) EndpointList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := chi.URLParam(r, "project")

	ends, err := e.endpoints.EndpointList(ctx, projectID)
	if err != nil {
		render.Render(w, r, types.ErrInternalServer(err, "list endpoints"))
		return
	}

	render.JSON(w, r, types.EndpointList{Endpoints: ends})
}

func (e *EndpointRouter) EndpointEvalAttach(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	endpointID := chi.URLParam(r, "endpointRef")

	var req types.EndpointEvalAttach
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Render(w, r, types.ErrHTTPBadRequest(err, "invalid request body"))
		return
	}

	if err := e.endpoints.EndpointAttachEval(ctx, endpointID, req.EvalID); err != nil {
		render.Render(w, r, types.ErrInternalServer(err, "attach eval"))
	}
}
