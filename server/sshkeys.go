package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/tools"
	"github.com/unweave/unweave/tools/random"
)

// SSHKeyAdd adds an SSH key to the user's account.
//
// This does not add the key to the user's configured providers. That is done lazily
// when the user first tries to use the key.
func SSHKeyAdd(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SSHKeyAdd request")

		params := api.SSHKeyAddRequestParams{}
		if err := render.Bind(r, &params); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), api.ErrHTTPError(err, "Invalid request body"))
			return
		}

		if params.Name != nil {
			k, err := store.SSHKey.GetByName(ctx, *params.Name)
			if err == nil {
				render.Render(w, r.WithContext(ctx), &api.HTTPError{
					Code:    http.StatusNotFound,
					Message: fmt.Sprintf("SSH key already exists with name: %q", k.Name),
				})
				return
			}
			if err != nil && err != sql.ErrNoRows {
				err = fmt.Errorf("failed to get ssh key from db: %w", err)
				render.Render(w, r.WithContext(ctx), api.ErrInternalServer(err, ""))
				return
			}
		} else {
			params.Name = tools.Stringy("uw:" + random.GenerateRandomPhrase(4, "-"))
		}

		err := store.SSHKey.Add(ctx, *params.Name, params.PublicKey)
		if err != nil {
			var e *pgconn.PgError
			if errors.As(err, &e) {
				// We already check the unique constraint on the name column, so this
				// should only happen if the public key is a duplicate.
				if e.Code == pgerrcode.UniqueViolation {
					render.Render(w, r.WithContext(ctx), &api.HTTPError{
						Code:    http.StatusConflict,
						Message: "Public key already exists",
						Suggestion: "Public keys in Unweave have to be globally unique. " +
							"It could be that you added this key earlier or that you " +
							"added it to another account. If you've already added this " +
							"key to your account, remove it first.",
					})
					return
				}
			}
			err = fmt.Errorf("failed to add ssh key to db: %w", err)
			render.Render(w, r.WithContext(ctx), api.ErrInternalServer(err, ""))
			return
		}

		render.JSON(w, r, &api.SSHKeyAddResponse{Success: true})
	}
}

func SSHKeyList(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Ctx(ctx).Info().Msgf("Executing SSHKeyList request")

		keys, err := store.SSHKey.List(ctx)
		if err != nil {
			err = fmt.Errorf("failed to list ssh keys from db: %w", err)
			render.Render(w, r.WithContext(ctx), api.ErrInternalServer(err, ""))
			return
		}

		res := api.SSHKeyListResponse{
			Keys: make([]api.SSHKey, len(keys)),
		}
		for idx, key := range keys {
			res.Keys[idx] = api.SSHKey{
				Name:      key.Name,
				PublicKey: key.PublicKey,
				CreatedAt: key.CreatedAt,
			}
		}
		render.JSON(w, r, res)
	}
}
