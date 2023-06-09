package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/types"
	execsrv "github.com/unweave/unweave/wip/services/exec"
)

type ExecRouter struct {
	r     chi.Router
	store execsrv.Store

	llService        *execsrv.Service
	unweaveService   *execsrv.Service
	conductorService *execsrv.Service
}

func (e *ExecRouter) Routes() []Route {
	var routes []Route

	_ = chi.Walk(e.r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		r := Route{
			Handler: handler,
			Method:  method,
			Path:    route,
		}
		routes = append(routes, r)
		return nil
	})

	return routes
}

func NewExecRouter(store execsrv.Store, lambdaLabsService, unweaveService *execsrv.Service) *ExecRouter {
	return &ExecRouter{
		store:            store,
		llService:        lambdaLabsService,
		unweaveService:   unweaveService,
		conductorService: nil,
	}
}

func (e *ExecRouter) service(provider types.Provider) *execsrv.Service {
	switch provider {
	case types.LambdaLabsProvider:
		return e.llService
	case types.UnweaveProvider:
		return e.unweaveService
	default:
		// This is unreachable. Using panic for now until we have the conductor
		// implemented for AWS, GCP etc
		panic(fmt.Errorf("unknown provider: %s", provider))
	}
}

func (e *ExecRouter) ExecCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msgf("Executing ExecCreate request")

	scr := &types.ExecCreateParams{}
	if err := scr.Bind(r); err != nil {
		err = fmt.Errorf("failed to read body: %w", err)
		render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Invalid request body"))
		return
	}

	userID := middleware.GetUserIDFromContext(ctx)
	projectID := middleware.GetProjectIDFromContext(ctx)

	exec, err := e.service(scr.Provider).Create(ctx, projectID, userID, *scr)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to create session"))
		return
	}
	render.JSON(w, r, exec)
}

func (e *ExecRouter) ExecGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msgf("Executing ExecGet request")

	execID := chi.URLParam(r, "exec")
	if execID == "" {
		err := fmt.Errorf("missing execID")
		render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Invalid request"))
		return
	}

	exec, err := e.store.Get(execID)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to get session"))
		return
	}
	render.JSON(w, r, exec)
}

func (e *ExecRouter) ExecListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msgf("Executing ExecList request")

	listTerminated := r.URL.Query().Get("terminated") == "true"
	projectID := middleware.GetProjectIDFromContext(ctx)

	execs, err := e.store.List(projectID)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to list sessions"))
		return
	}

	if listTerminated {
		render.JSON(w, r, types.ExecsListResponse{Execs: execs})
		return
	}

	var res []types.Exec

	for _, exec := range execs {
		if exec.Status == types.StatusTerminated ||
			exec.Status == types.StatusError ||
			exec.Status == types.StatusFailed {
			continue
		}
		res = append(res, exec)
	}

	render.JSON(w, r, types.ExecsListResponse{Execs: res})
}

func (e *ExecRouter) ExecTerminateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msgf("Executing ExecTerminate request")

	execID := chi.URLParam(r, "exec")
	if execID == "" {
		err := fmt.Errorf("missing execID")
		render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Invalid request"))
		return
	}

	exec, err := e.store.Get(execID)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to get session"))
		return
	}

	err = e.service(exec.Provider).Terminate(ctx, execID)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to terminate session"))
		return
	}

	render.Status(r, http.StatusOK)
}
