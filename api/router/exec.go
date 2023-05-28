package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/providers/lambdalabs"
	execsrv "github.com/unweave/unweave/wip/services/exec"
)

type ExecRouter struct {
	store execsrv.Store

	llService        *execsrv.Service
	unweaveService   *execsrv.Service
	conductorService *execsrv.Service
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
		panic("not implemented") // this is only available in the hosted version
	default:
		// This is unreachable update this. Using panic for now until we have the
		// conductor implemented for AWS, GCP etc
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

	projectID := middleware.GetProjectIDFromContext(ctx)

	exec, err := e.service(scr.Provider).Create(ctx, projectID, *scr)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to create session"))
		return
	}
	render.JSON(w, r, exec)
}

func RegisterExecRoutes(r chi.Router, store execsrv.Store) {
	lls, err := execsrv.NewService(store, lambdalabs.ExecDriver{})
	if err != nil {
		panic(err)
	}

	lls = execsrv.WithStateObserver(lls, func(e types.Exec) execsrv.StateObserver { return execsrv.NewStateObserver(e, lls) })
	lls = execsrv.WithStatsObserver(lls, func(e types.Exec) execsrv.StatsObserver { return execsrv.NewTerminateIdleObserver(e, lls, 2*time.Hour) })

	execRouter := NewExecRouter(store, lls, nil)

	r.Post("/", execRouter.ExecCreateHandler)
}
