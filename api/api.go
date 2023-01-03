package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/config"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/runtime"
)

func API(cfg config.Config, rti runtime.Initializer, dbq db.Querier) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  &log.Logger,
		NoColor: true,
	}))
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

	r.Route("/sessions", func(r chi.Router) {
		r.Post("/", SessionsCreate(rti, dbq))
		r.Get("/", SessionsList(rti))
		r.Get("/{id}", SessionsGet(rti))
		r.Put("/{id}/terminate", SessionsTerminate(rti))
	})

	log.Info().Msgf("🚀 API listening on %s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		panic(err)
	}
}
