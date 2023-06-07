package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

func SetAccountIDInContext(ctx context.Context, aid string) context.Context {
	return context.WithValue(ctx, types.AccountIDCtxKey, aid)
}

func GetAccountIDFromContext(ctx context.Context) string {
	uid, ok := ctx.Value(types.AccountIDCtxKey).(string)
	if !ok || uid == "" {
		// This should never happen at runtime.
		log.Error().Msg("account not found in context")
		panic("account not found in context")
	}
	return uid
}

func SetUserIDInContext(ctx context.Context, aid string) context.Context {
	return context.WithValue(ctx, types.UserIDCtxKey, aid)
}

func GetUserIDFromContext(ctx context.Context) string {
	uid, ok := ctx.Value(types.UserIDCtxKey).(string)
	if !ok || uid == "" {
		// This should never happen at runtime.
		log.Error().Msg("account not found in context")
		panic("account not found in context")
	}
	return uid
}

func SetProjectIDInContext(ctx context.Context, projectID string) context.Context {
	return context.WithValue(ctx, types.ProjectIDCtxKey, projectID)
}

func GetProjectIDFromContext(ctx context.Context) string {
	projectID, ok := ctx.Value(types.ProjectIDCtxKey).(string)
	if !ok || projectID == "" {
		// This should never happen at runtime.
		log.Error().Msg("project not found in context")
		panic("project not found in context")
	}
	return projectID
}

func SetExecIDInContext(ctx context.Context, execID string) context.Context {
	return context.WithValue(ctx, types.ExecIDCtxKey, execID)
}

func GetExecIDFromContext(ctx context.Context) string {
	execID, ok := ctx.Value(types.ExecIDCtxKey).(string)
	if !ok || execID == "" {
		// This should never happen at runtime.
		log.Error().Msg("exec not found in context")
		panic("exec not found in context")
	}
	return execID
}

// WithAccountCtx is a helper middleware that fakes an authenticated account. It should only
// be user for development or when self-hosting.
func WithAccountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := "uid_1234"
		ctx = SetUserIDInContext(ctx, userID)
		ctx = SetAccountIDInContext(ctx, userID)
		ctx = log.With().Str(types.UserIDCtxKey, userID).Logger().WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithProjectCtx is a helper middleware that parsed the project id from the url and
// verifies it exists in the db.
func WithProjectCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID := chi.URLParam(r, "project")

		_, err := db.Q.ProjectGet(ctx, projectID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r.WithContext(ctx), &types.Error{
					Code:       http.StatusNotFound,
					Message:    "Project not found",
					Suggestion: "Make sure the project id is valid",
				})
				return
			}

			err = fmt.Errorf("failed to fetch project from db %q: %w", projectID, err)
			render.Render(w, r.WithContext(ctx),
				types.ErrInternalServer(err, "Failed to fetch project"))
			return
		}

		ctx = context.WithValue(ctx, types.ProjectIDCtxKey, projectID)
		ctx = log.With().Str(types.ProjectIDCtxKey, projectID).Logger().WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithExecCtx is a helper middleware that parsed the session id from the url and
// verifies it exists in the db.
func WithExecCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ref := chi.URLParam(r, "exec")

		exec, err := db.Q.ExecGet(ctx, ref)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r.WithContext(ctx), &types.Error{
					Code:       http.StatusNotFound,
					Message:    "Session not found",
					Suggestion: "Make sure the session id is valid",
				})
				return
			}

			err = fmt.Errorf("failed to fetch exec from db %q: %w", ref, err)
			render.Render(w, r.WithContext(ctx), types.ErrInternalServer(err, "Failed to fetch session"))
			return
		}

		ctx = context.WithValue(ctx, types.ExecIDCtxKey, exec)
		ctx = log.With().Str(types.ExecIDCtxKey, exec.ID).Logger().WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
