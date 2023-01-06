package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/types"
)

type NodeTypesListResponse struct {
	NodeTypes []types.NodeType `json:"nodeTypes"`
}

// NodeTypesList returns a list of node types available for the user
func NodeTypesList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := getUserIDFromContext(ctx)
		provider := chi.URLParam(r, "provider")

		ctx = log.With().Stringer(ContextKeyUser, userID).Logger().WithContext(ctx)
		log.Ctx(ctx).Info().Msgf("Executing NodeTypesList request for provider %s", provider)

		rt, err := rti.FromUser(userID, types.RuntimeProvider(provider))
		if err != nil {
			render.Render(w, r, ErrHTTPError(err, "Failed to create runtime"))
			return
		}

		nodeTypes, err := rt.ListNodeTypes(ctx)
		if err != nil {
			render.Render(w, r, ErrHTTPError(err, "Failed to list node types"))
			return
		}
		res := &NodeTypesListResponse{
			NodeTypes: nodeTypes,
		}
		render.JSON(w, r, res)
	}
}
