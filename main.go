package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/unweave/unweave-v2/config"
	"github.com/unweave/unweave-v2/session"
)

func main() {

	// init db
	// init config
	// init and

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

	r.Post("/session", func(w http.ResponseWriter, r *http.Request) {

		// parse request body for json to get runtime type

		cfg := config.SessionConfig{}
		sess := session.NewSession(cfg)

	})

	r.Post("/session/{id}/stop", func(w http.ResponseWriter, r *http.Request) {

	})
}
