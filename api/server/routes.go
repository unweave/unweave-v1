package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	middleware2 "github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/router"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/runtime"
)

type Config struct {
	APIPort string    `json:"port" env:"UNWEAVE_API_PORT"`
	DB      db.Config `json:"db"`
}

func API(cfg Config, rti runtime.Initializer, execRouter *router.ExecRouter) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  &log.Logger,
		NoColor: true,
	}))
	// Initialize contextual logger
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := log.With().Logger().WithContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
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

	r.Use(middleware2.WithAccountCtx) // fakes an authenticated user
	r.Route("/projects/{owner}/{project}", func(r chi.Router) {
		r.Use(middleware2.WithProjectCtx)

		r.Route("/builds", func(r chi.Router) {
			r.Post("/", BuildsCreate(rti))
			r.Get("/{buildID}", BuildsGet(rti))
		})

		r.Route("/sessions", func(r chi.Router) {
			r.Post("/", execRouter.ExecCreateHandler)
		})

		r.Route("/sessions", func(r chi.Router) {

			r.Get("/", execRouter.ExecListHandler)

			r.Route("/{exec}", func(r chi.Router) {
				r.Use(middleware2.WithExecCtx)
				r.Get("/", execRouter.ExecGetHandler)
				r.Put("/terminate", execRouter.ExecTerminateHandler)
			})
		})

		r.Route("/volumes", func(r chi.Router) {
			r.Post("/", VolumeCreate(rti))
		})
	})

	r.Route("/ssh-keys/{owner}", func(r chi.Router) {
		r.Post("/", SSHKeyAdd(rti))
		r.Get("/", SSHKeyList(rti))
		r.Post("/generate", SSHKeyGenerate(rti))
	})
	r.Get("/providers/{provider}/node-types", NodeTypesList(rti))

	ctx := context.Background()
	ctx = log.With().Logger().WithContext(ctx)

	log.Info().Msgf("🚀 API listening on %s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		panic(err)
	}
}
