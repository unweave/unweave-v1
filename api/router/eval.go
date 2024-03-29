//nolint:varnamelen
package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/unweave/unweave-v1/api/middleware"
	"github.com/unweave/unweave-v1/api/types"
	"github.com/unweave/unweave-v1/services/evalsrv"
)

type EvalRouter struct {
	service evalsrv.Service
}

func NewEvalRouter(service evalsrv.Service) *EvalRouter {
	return &EvalRouter{service: service}
}

func (e *EvalRouter) EvalCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := middleware.GetProjectIDFromContext(ctx)

	var req types.EvalCreate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, types.ErrHTTPBadRequest(err, "decode request"))

		return
	}

	defer r.Body.Close()

	eval, err := e.service.EvalCreate(ctx, projectID, req.ExecID)
	if err != nil {
		_ = render.Render(w, r, types.ErrHTTPError(err, "create eval"))

		return
	}

	render.JSON(w, r, eval)
}

func (e *EvalRouter) EvalList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := middleware.GetProjectIDFromContext(ctx)

	evals, err := e.service.EvalListForProject(ctx, projectID)
	if err != nil {
		_ = render.Render(w, r, types.ErrHTTPError(err, "create eval"))

		return
	}

	render.JSON(w, r, types.EvalList{Evals: evals})
}
