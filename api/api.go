package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/config"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/types"
)

func userContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx,
			types.ContextKeyUser,
			types.UserContext{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
		r.Use(userContext) // <-- fakes an authenticated user
		r.Post("/", SessionsCreate(rti, dbq))
		r.Get("/", SessionsList(rti))
		r.Get("/{id}", SessionsGet(rti))
		r.Put("/{id}/terminate", SessionsTerminate(rti, dbq))
	})

	log.Info().Msgf("🚀 API listening on %s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		panic(err)
	}
}
