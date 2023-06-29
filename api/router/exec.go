package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/execsrv"
)

type ExecRouter struct {
	r       chi.Router
	store   execsrv.Store
	service execsrv.Service
}

func NewExecRouter(store execsrv.Store, services ...execsrv.Service) *ExecRouter {
	delegates := make(map[types.Provider]execsrv.Service)

	for i := range services {
		svc := services[i]
		delegates[svc.Provider()] = svc
	}

	return &ExecRouter{
		store:   store,
		service: execsrv.NewServiceRouter(store, delegates),
	}
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

func (e *ExecRouter) ExecCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msgf("Executing ExecCreate request")

	params := &types.ExecCreateParams{}
	if err := params.Bind(r); err != nil {
		err = fmt.Errorf("failed to read body: %w", err)
		render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Invalid request body"))
		return
	}

	userID := middleware.GetUserIDFromContext(ctx)
	projectID := middleware.GetProjectIDFromContext(ctx)

	exec, err := e.service.Create(ctx, projectID, userID, *params)
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

	exec, err := e.service.Get(ctx, execID)
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

	execs, err := e.service.List(ctx, projectID)
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

	err := e.service.Terminate(ctx, execID)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to terminate session"))
		return
	}

	render.Status(r, http.StatusOK)
}
