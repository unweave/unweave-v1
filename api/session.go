package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/unweave/unweave-v2/runtime"
	"github.com/unweave/unweave-v2/types"
)

func sessionsGet(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		res := &types.Session{ID: id}
		render.JSON(w, r, res)
	}
}

func sessionsList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := []*types.Session{
			{ID: "1"},
		}
		render.JSON(w, r, res)
	}
}

func sessionsTerminate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		res := &types.Session{ID: id}
		render.JSON(w, r, res)
	}
}
