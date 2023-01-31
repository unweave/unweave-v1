package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/providers/lambdalabs"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/tools/gonfig"
)

type Config struct {
	APIPort string    `json:"port" env:"UNWEAVE_API_PORT"`
	DB      db.Config `json:"db"`
}

// EnvInitializer is only used in development or if you're self-hosting Unweave.
type EnvInitializer struct{}

type providerConfig struct {
	LambdaLabsAPIKey string `env:"LAMBDALABS_API_KEY"`
}

func (i *EnvInitializer) Initialize(ctx context.Context, accountID uuid.UUID, provider types.RuntimeProvider, token *string) (*runtime.Runtime, error) {
	var cfg providerConfig
	gonfig.GetFromEnvVariables(&cfg)

	switch provider {
	case types.LambdaLabsProvider:
		if cfg.LambdaLabsAPIKey == "" && token == nil {
			return nil, fmt.Errorf("missing LambdaLabs API key in runtime config file")
		}
		if token != nil {
			log.Ctx(ctx).
				Info().
				Msgf("Overriding LambdaLabs API key in env config: using runtime token")
			cfg.LambdaLabsAPIKey = *token
		}

		sess, err := lambdalabs.NewSessionProvider(cfg.LambdaLabsAPIKey)
		if err != nil {
			return nil, err
		}
		return &runtime.Runtime{Session: sess}, nil

	default:
		return nil, fmt.Errorf("%q provider not supported in the env initializer", provider)
	}
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
			r.Get("/", SessionsList)

			r.Group(func(r chi.Router) {
				r.Use(withSessionCtx)
				r.Get("/{sessionID}", SessionsGet(rti))
				r.Put("/{sessionID}/terminate", SessionsTerminate(rti))
			})
		})
	})

	r.Route("/ssh-keys", func(r chi.Router) {
		r.Post("/", SSHKeyAdd)
		r.Get("/", SSHKeyList)
		r.Post("/generate", SSHKeyGenerate)
	})
	r.Get("/providers/{provider}/node-types", NodeTypesList(rti))

	log.Info().Msgf("🚀 API listening on %s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		panic(err)
	}
}
