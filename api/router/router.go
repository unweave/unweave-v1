package router

import (
	"net/http"
)

type Router interface {
	Routes() []Route
}

type Route struct {
	Handler http.Handler
	Method  string
	Path    string
}
