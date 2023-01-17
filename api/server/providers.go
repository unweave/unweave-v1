package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/runtime"
)

type NodeTypesListResponse struct {
	NodeTypes []types.NodeType `json:"nodeTypes"`
}

// NodeTypesList returns a list of node types available for the user
func NodeTypesList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := GetUserIDFromContext(ctx)
		provider := types.RuntimeProvider(chi.URLParam(r, "provider"))

		ctx = log.With().Stringer(UserIDCtxKey, userID).Logger().WithContext(ctx)
		log.Ctx(ctx).Info().Msgf("Executing NodeTypesList request for provider %s", provider)

		if provider != types.LambdaLabsProvider && provider != types.UnweaveProvider {
			render.Render(w, r.WithContext(ctx), &types.HTTPError{
				Code:       http.StatusBadRequest,
				Message:    "Invalid runtime provider: " + string(provider),
				Suggestion: fmt.Sprintf("Use %q or %q as the runtime provider", types.LambdaLabsProvider, types.UnweaveProvider),
			})
			return
		}

		rt, err := rti.FromUserID(ctx, userID, provider)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to create runtime"))
			return
		}

		nodeTypes, err := rt.ListNodeTypes(ctx)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to list node types"))
			return
		}

		res := &NodeTypesListResponse{NodeTypes: nodeTypes}
		render.JSON(w, r, res)
	}
}
