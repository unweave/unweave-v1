package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/config"
	"github.com/unweave/unweave-v2/session/runtime"
)

func API(cfg config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
	}))

	r.Route("/session", func(r chi.Router) {

		// swagger:route POST /session/{id} session sessionCreate
		// responses:
		// 	201: sessionCreate
		r.Post("/", sessionCreateHandler)

		// swagger:route GET /session/{id} session sessionGet
		// responses:
		// 	200: sessionGet
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")

			res := &SessionGetResponse{ID: id}
			render.JSON(w, r, res)
		})

		// swagger:route GET /session/{id}/connect session sessionConnect
		//
		// Returns the SSH connection details for the session. This will return a 404 if
		// the session is not yet ready.
		//
		// responses:
		// 	200: sessionConnect
		//  404: errorResponse
		r.Get("/{id}/connect", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")

			res := &SessionConnectResponse{
				ID:     id,
				Status: runtime.StatusRunning,
				Connection: runtime.SSHConnection{
					Host:     "localhost",
					Port:     "22",
					User:     "noorvir",
					Password: "",
				},
			}
			render.JSON(w, r, res)
		})

	})

	log.Info().Msgf("🚀 API listening on %s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		panic(err)
	}
}
