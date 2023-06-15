package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/sshkeys"
)

type SSHKeysRouter struct {
	r       chi.Router
	service *sshkeys.Service
}

func (s *SSHKeysRouter) Routes() []Route {
	var routes []Route

	_ = chi.Walk(s.r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
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

func NewSSHKeysRouter(service *sshkeys.Service) *SSHKeysRouter {
	router := chi.NewRouter()
	sshKeysRouter := &SSHKeysRouter{
		r:       router,
		service: service,
	}
	return sshKeysRouter
}

func (s *SSHKeysRouter) SSHKeysAddHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msg("Executing SSHKeysAdd request")
	userID := middleware.GetUserIDFromContext(ctx)

	var params types.SSHKeyAddParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Failed to decode SSHKeyAdd parameters"))
		return
	}

	name, err := s.service.Add(ctx, userID, params)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to add SSH key"))
		return
	}

	res := types.SSHKeyResponse{
		Name:       name,
		PublicKey:  params.PublicKey,
		PrivateKey: "",
	}
	render.JSON(w, r, &res)
}

func (s *SSHKeysRouter) SSHKeysGenerateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msg("Executing SSHKeysGenerate request")
	userID := middleware.GetUserIDFromContext(ctx)

	var params types.SSHKeyGenerateParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Failed to decode SSHKeyGenerate parameters"))
		return
	}

	name, prv, pub, err := s.service.Generate(ctx, userID, params)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to generate SSH key"))
		return
	}

	res := types.SSHKeyResponse{
		Name:       name,
		PublicKey:  pub,
		PrivateKey: prv,
	}
	render.JSON(w, r, &res)
}

func (s *SSHKeysRouter) SSHKeysListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Ctx(ctx).Info().Msg("Executing SSHKeysList request")
	userID := middleware.GetUserIDFromContext(ctx)

	keys, err := s.service.List(ctx, userID)
	if err != nil {
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to list SSH keys"))
		return
	}

	res := types.SSHKeyListResponse{Keys: keys}
	render.JSON(w, r, res)
}
