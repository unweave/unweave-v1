package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/providersrv"
)

type ProviderRouter struct {
	r         chi.Router
	delegates map[types.Provider]*providersrv.ProviderService
	supported []string
}

func NewProviderRouter(services ...*providersrv.ProviderService) *ProviderRouter {
	supported := make([]string, 0, len(services))
	delegates := make(map[types.Provider]*providersrv.ProviderService)

	for _, svc := range services {
		supported = append(supported, svc.Provider().String())
		delegates[svc.Provider()] = svc
	}

	return &ProviderRouter{
		delegates: delegates,
		supported: supported,
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

func (p *ProviderRouter) ProviderListNodeTypesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msgf("Executing ProviderListNodeTypes request")

	provider := types.Provider(chi.URLParam(r, "provider"))
	filterAvailable := r.URL.Query().Get("available") == "true"

	// TODO: This should be change to middleware.GetAccountIDFromContext once the route is
	//  moved to /accounts/{accountID}/providers/{provider}/node-types
	userID := middleware.GetUserIDFromContext(ctx)

	service, ok := p.delegates[provider]
	if !ok {
		err := &types.Error{
			Code:       http.StatusBadRequest,
			Message:    "Invalid provider",
			Suggestion: fmt.Sprintf("Valid providers are: %s ", strings.Join(p.supported, ", ")),
		}
		_ = render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to list node types"))

		return
	}

	nodeTypes, err := service.ListNodeTypes(ctx, userID, filterAvailable)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to list node types"))
		return
	}

	res := &types.NodeTypesListResponse{NodeTypes: nodeTypes}
	render.JSON(w, r, res)
}
