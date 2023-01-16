package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/runtime"
)

type Config struct {
	APIPort string `json:"port" env:"UNWEAVE_API_PORT"`
}

func Server(cfg Config, rti runtime.Initializer, store *Store) {
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
		r.Post("/", SessionsCreate(rti, store))
		r.Get("/", SessionsList(rti, store))

		r.Group(func(r chi.Router) {
			r.Use(withSessionCtx(store))
			r.Get("/{sessionID}", SessionsGet(rti))
			r.Put("/{sessionID}/terminate", SessionsTerminate(rti, store))
		})
	})

	r.Route("/ssh-keys", func(r chi.Router) {
		r.Post("/", SSHKeyAdd(store))
		r.Get("/", SSHKeyList(store))
	})
	r.Get("/providers/{provider}/node-types", NodeTypesList(rti))

	log.Info().Msgf("🚀 API listening on %s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		panic(err)
	}
}
