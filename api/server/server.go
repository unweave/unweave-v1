package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/runtime"
)

type Config struct {
	APIPort string    `json:"port" env:"UNWEAVE_API_PORT"`
	DB      db.Config `json:"db"`
}

func HandleRestart(ctx context.Context, rti runtime.Initializer) error {
	// Re-watch all sessions
	sessions, err := db.Q.SessionGetAllActive(ctx)
	if err != nil {
		return err
	}

	log.Ctx(ctx).Info().Msgf("Restarting watching %d sessions", len(sessions))

	for _, s := range sessions {
		sess := s
		go func() {
			c := context.Background()
			c = log.With().
				Stringer(UserIDCtxKey, sess.CreatedBy).
				Stringer(ProjectCtxKey, sess.ProjectID).
				Stringer(SessionCtxKey, sess.ID).
				Logger().WithContext(c)

			srv := NewCtxService(rti, sess.CreatedBy)
			if e := srv.Session.Watch(c, sess.ID); e != nil {
				log.Ctx(ctx).Error().Err(e).Msgf("Failed to watch session")
			}
		}()
	}
	return nil
}

func API(cfg Config, rti runtime.Initializer) {
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

	r.Use(withUserCtx) // fakes an authenticated user
	r.Route("/projects/{projectID}", func(r chi.Router) {
		r.Use(withProjectCtx)

		r.Route("/sessions", func(r chi.Router) {
			r.Post("/", SessionsCreate(rti))
			r.Get("/", SessionsList(rti))

			r.Group(func(r chi.Router) {
				r.Use(withSessionCtx)
				r.Get("/{sessionID}", SessionsGet(rti))
				r.Put("/{sessionID}/terminate", SessionsTerminate(rti))
			})
		})
	})

	r.Route("/ssh-keys", func(r chi.Router) {
		r.Post("/", SSHKeyAdd(rti))
		r.Get("/", SSHKeyList(rti))
		r.Post("/generate", SSHKeyGenerate(rti))
	})
	r.Get("/providers/{provider}/node-types", NodeTypesList(rti))

	ctx := context.Background()
	ctx = log.With().Logger().WithContext(ctx)
	if err := HandleRestart(ctx, rti); err != nil {
		panic(err)
	}

	log.Info().Msgf("🚀 API listening on %s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		panic(err)
	}
}
