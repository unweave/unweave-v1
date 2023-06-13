package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

}

func (v *VolumeRouter) VolumeDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func (v *VolumeRouter) VolumeGetHandler(w http.ResponseWriter, r *http.Request) {

}

func (v *VolumeRouter) VolumeListHandler(w http.ResponseWriter, r *http.Request) {

}

func (v *VolumeRouter) VolumeResizeHandler(w http.ResponseWriter, r *http.Request) {

}