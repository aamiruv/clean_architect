package router

import (
	"fmt"
	"net/http"

	"github.com/AmirMirzayi/clean_architecture/pkg/middleware"
)

type (
	actionHandler struct {
		handler     http.HandlerFunc
		middlewares []middleware.Middleware
	}

	// routing path eg: /shelves/{shelf_id}/books
	Path string

	// http method eg: http.MethodGet or empty to handle all types
	Method string

	// Route include
	Route map[Path]map[Method]actionHandler

	// routing prefix eg: /api/v1
	Prefix string

	// routing group allow you to set prefix and middlewares on
	// several individual routes
	GroupRoute map[Prefix]Route
)

func NewGroup(group GroupRoute, middlewares ...middleware.Middleware) Route {
	ro := Route{}
	for p, r := range group {
		for k, r2 := range r {
			ro[Path(fmt.Sprintf("%s%s", p, k))] = r2
		}
	}
	return ro
}

// NewHandler create actionHandler within given middlewares.
// Middlewares are executed from last to first.
func NewHandler(handler http.HandlerFunc, middlewares ...middleware.Middleware) actionHandler {
	return actionHandler{
		handler:     handler,
		middlewares: middlewares,
	}
}

func BindRoutesToMux(mux *http.ServeMux, routes ...Route) {
	r := MergeRoutes(routes...)
	for route, handler := range r {
		for method, action := range handler {
			mux.Handle(fmt.Sprintf("%s %s", method, route), middleware.Chain(action.handler, action.middlewares...))
		}
	}
}

func MergeRoutes(routes ...Route) Route {
	r := Route{}
	for _, v := range routes {
		for k, v2 := range v {
			r[k] = v2
		}
	}
	return r
}
