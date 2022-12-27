package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/config"
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
		r.Post("/{id}", func(w http.ResponseWriter, r *http.Request) {
			scr := SessionCreateRequest{}
			if err := render.Bind(r, &scr); err != nil {
				log.Warn().Err(err).Msg("failed to read body")
				render.Render(w, r, ErrBadRequest("Invalid request body: "+err.Error()))
				return
			}

			fmt.Println(scr.Runtime)
			w.WriteHeader(http.StatusOK)
		})

	})

	log.Info().Msgf("🚀 API listening on %s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		panic(err)
	}
}
