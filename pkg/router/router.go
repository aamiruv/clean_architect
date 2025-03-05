package router

import (
	"fmt"
	"net/http"

	"github.com/AmirMirzayi/clean_architecture/pkg/middleware"
)

type Route map[string]map[string]actionHandler

type actionHandler struct {
	handler     http.HandlerFunc
	middlewares []middleware.Middleware
}

func NewHandler(handler http.HandlerFunc, middlewares ...middleware.Middleware) actionHandler {
	return actionHandler{
		handler:     handler,
		middlewares: middlewares,
	}
}

func BindRoutes(mux *http.ServeMux, routes Route) {
	for route, handler := range routes {
		for method, action := range handler {
			mux.Handle(fmt.Sprintf("%s %s", method, route), middleware.Chain(action.handler, action.middlewares...))
		}
	}
}
