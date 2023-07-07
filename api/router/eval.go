package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/evalsrv"
)

type EvalRouter struct {
	service evalsrv.Service
}

func NewEvalRouter(service evalsrv.Service) *EvalRouter {
	return &EvalRouter{service: service}
}

func (e *EvalRouter) EvalCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := chi.URLParam(r, "project")

	var req types.EvalCreate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Render(w, r, types.ErrHTTPBadRequest(err, "decode request"))
		return
	}

	defer r.Body.Close()

	eval, err := e.service.EvalCreate(ctx, projectID, req.ExecID)
	if err != nil {
		render.Render(w, r, types.ErrInternalServer(err, "create eval"))
		return
	}

	render.JSON(w, r, eval)
}
