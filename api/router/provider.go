package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/providersrv"
)

type ProviderRouter struct {
	r         chi.Router
	llService *providersrv.ProviderService
	uwService *providersrv.ProviderService
}

func NewProviderRouter(lambdaLabsService, unweaveService *providersrv.ProviderService) *ProviderRouter {
	return &ProviderRouter{
		llService: lambdaLabsService,
		uwService: unweaveService,
	}
}

func (p *ProviderRouter) Routes() []Route {
	var routes []Route

	_ = chi.Walk(p.r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
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

func (p *ProviderRouter) ProviderListNodeTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msgf("Executing ProviderListNodeTypes request")

	provider := types.Provider(chi.URLParam(r, "provider"))
	filterAvailable := r.URL.Query().Get("available") == "true"

	// TODO: This should be change to middleware.GetAccountIDFromContext once the route is
	//  moved to /accounts/{accountID}/providers/{provider}/node-types
	userID := middleware.GetUserIDFromContext(ctx)

	var err error
	var nodeTypes []types.NodeType

	switch provider {
	case types.LambdaLabsProvider:
		nodeTypes, err = p.llService.ListNodeTypes(ctx, userID, filterAvailable)
	case types.UnweaveProvider:
		nodeTypes, err = p.uwService.ListNodeTypes(ctx, userID, filterAvailable)
	default:
		err = &types.Error{
			Code:       http.StatusBadRequest,
			Message:    "Invalid provider",
			Suggestion: fmt.Sprintf("Valid providers are: %s, %s ", types.LambdaLabsProvider, types.UnweaveProvider),
		}
	}

	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to list node types"))
		return
	}

	res := &types.NodeTypesListResponse{NodeTypes: nodeTypes}
	render.JSON(w, r, res)
}
