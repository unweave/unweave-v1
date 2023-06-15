package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/services/volumesrv"
)

type VolumeRouter struct {
	r         chi.Router
	llService *volumesrv.Service
	uwService *volumesrv.Service
}

func NewVolumeRouter(llService, uwService *volumesrv.Service) *VolumeRouter {
	return &VolumeRouter{
		llService: llService,
		uwService: uwService,
	}
}

func (v *VolumeRouter) Routes() []Route {
	var routes []Route

	_ = chi.Walk(v.r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		r := Route{
			Handler: handler,
			Method:  method,
			Path:    route,
		}
		routes = append(routes, r)
		return nil
	})

	return routes
}

func (v *VolumeRouter) VolumeCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := middleware.GetProjectIDFromContext(r.Context())

	vcr := &types.VolumeCreateRequest{}
	if err := render.Bind(r, vcr); err != nil {
		render.Render(w, r, types.ErrHTTPBadRequest(err, "Failed to parse request"))
		return
	}

	vol, err := v.uwService.Create(r.Context(), projectID, vcr.Name, vcr.Size)
	if err != nil {
		err = fmt.Errorf("failed to create volume, %w", err)
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to create volume"))
	}

	render.JSON(w, r, vol)
}

func (v *VolumeRouter) VolumeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	projectID := middleware.GetProjectIDFromContext(r.Context())

	vdr := &types.VolumeDeleteRequest{}
	if err := render.Bind(r, vdr); err != nil {
		render.Render(w, r, types.ErrHTTPBadRequest(err, "Failed to parse request"))
		return
	}

	err := v.uwService.Delete(r.Context(), projectID, vdr.IDOrName)
	if err != nil {
		err = fmt.Errorf("failed to delete volume, %w", err)
		render.Render(w, r, types.ErrHTTPError(err, "Failed to delete volume"))
		return
	}

	render.Status(r, http.StatusOK)
}

func (v *VolumeRouter) VolumeGetHandler(w http.ResponseWriter, r *http.Request) {
	projectID := middleware.GetProjectIDFromContext(r.Context())
	idOrName := chi.URLParam(r, "idOrName")

	vol, err := v.uwService.Get(r.Context(), projectID, idOrName)
	if err != nil {
		err = fmt.Errorf("failed to get volume, %w", err)
		render.Render(w, r, types.ErrHTTPError(err, "Failed to get volume"))
		return
	}

	render.JSON(w, r, vol)
}

func (v *VolumeRouter) VolumeListHandler(w http.ResponseWriter, r *http.Request) {
	projectID := middleware.GetProjectIDFromContext(r.Context())

	vol, err := v.uwService.List(r.Context(), projectID)
	if err != nil {
		err = fmt.Errorf("failed to list volumes, %w", err)
		render.Render(w, r, types.ErrHTTPError(err, "Failed to list volumes"))
		return
	}

	render.JSON(w, r, vol)
}

func (v *VolumeRouter) VolumeResizeHandler(w http.ResponseWriter, r *http.Request) {
	projectID := middleware.GetProjectIDFromContext(r.Context())

	vrr := &types.VolumeResizeRequest{}
	if err := render.Bind(r, vrr); err != nil {
		render.Render(w, r, types.ErrHTTPBadRequest(err, "Failed to parse request"))
		return
	}

	err := v.uwService.Resize(r.Context(), projectID, vrr.IDOrName, vrr.Size)
	if err != nil {
		err = fmt.Errorf("failed to resize volume, %w", err)
		render.Render(w, r, types.ErrHTTPError(err, "Failed to resize volume"))
		return
	}

	render.Status(r, http.StatusOK)
}
