package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/types"
	execsrv "github.com/unweave/unweave/wip/services/exec"
)

type ExecRouter struct {
	store execsrv.Store

	llService        *execsrv.Service
	unweaveService   *execsrv.Service
	conductorService *execsrv.Service
}

func NewExecRouter(store execsrv.Store, lambdaLabsDriver, unweaveDriver execsrv.Driver) *ExecRouter {
	lls, err := execsrv.NewService(store, lambdaLabsDriver)
	if err != nil {
		panic(fmt.Errorf("failed to create lambda labs service: %w", err))
	}

	return &ExecRouter{
		store:            store,
		llService:        lls,
		unweaveService:   nil,
		conductorService: nil,
	}
}

func (e *ExecRouter) route(provider types.Provider) *execsrv.Service {
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

	exec, err := e.route(scr.Provider).Create(ctx, projectID, *scr)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to create session"))
		return
	}
	render.JSON(w, r, exec)
}
