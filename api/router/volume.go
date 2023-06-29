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
	r       chi.Router
	service volumesrv.Service
}

func NewVolumeRouter(store volumesrv.Store, services ...volumesrv.Service) *VolumeRouter {
	delegates := make(map[types.Provider]volumesrv.Service)

	for i := range services {
		svc := services[i]
		delegates[svc.Provider()] = svc
	}

	return &VolumeRouter{
		service: volumesrv.NewServiceRouter(store, delegates),
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

	accountID := middleware.GetAccountIDFromContext(ctx)
	projectID := middleware.GetProjectIDFromContext(ctx)

	vcr := &types.VolumeCreateRequest{}
	if err := render.Bind(r, vcr); err != nil {
		render.Render(w, r, types.ErrHTTPBadRequest(err, "Failed to parse request"))
		return
	}

	vol, err := v.service.Create(ctx, accountID, projectID, vcr.Provider, vcr.Name, vcr.Size)
	if err != nil {
		err = fmt.Errorf("failed to create volume: %w", err)
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to create volume"))
		return
	}

	render.JSON(w, r, vol)
}

func (v *VolumeRouter) VolumeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idOrName := chi.URLParam(r, "volumeRef")
	projectID := middleware.GetProjectIDFromContext(ctx)

	err := v.service.Delete(ctx, projectID, idOrName)
	if err != nil {
		err = fmt.Errorf("failed to delete volume, %w", err)
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to delete volume"))
		return
	}

	render.Status(r, http.StatusOK)
}

func (v *VolumeRouter) VolumeGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectID := middleware.GetProjectIDFromContext(ctx)
	idOrName := chi.URLParam(r, "volumeRef")

	vol, err := v.service.Get(ctx, projectID, idOrName)
	if err != nil {
		err = fmt.Errorf("failed to get volume, %w", err)
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to get volume"))
		return
	}

	render.JSON(w, r, vol)
}

func (v *VolumeRouter) VolumeListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := middleware.GetProjectIDFromContext(ctx)

	vols, err := v.service.List(ctx, projectID)
	if err != nil {
		err = fmt.Errorf("failed to list volumes, %w", err)
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to list volumes"))
		return
	}

	render.JSON(w, r, types.VolumesListResponse{Volumes: vols})
}

func (v *VolumeRouter) VolumeResizeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := middleware.GetProjectIDFromContext(ctx)

	vrr := &types.VolumeResizeRequest{}
	if err := render.Bind(r, vrr); err != nil {
		render.Render(w, r, types.ErrHTTPBadRequest(err, "Failed to parse request"))
		return
	}

	err := v.service.Resize(ctx, projectID, vrr.IDOrName, vrr.Size)
	if err != nil {
		err = fmt.Errorf("failed to resize volume, %w", err)
		render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to resize volume"))
		return
	}

	render.Status(r, http.StatusOK)
}
