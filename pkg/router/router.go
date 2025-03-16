// Package router provides a flexible and structured way to define and manage HTTP routes.
// It supports grouping routes by prefixes and mapping paths to handlers based on HTTP methods.
package router

import (
	"fmt"
	"net/http"

	"github.com/amirzayi/clean_architec/pkg/httpmiddleware"
)

type (
	actionHandler struct {
		handler     http.HandlerFunc
		middlewares []httpmiddleware.Middleware
	}

	// Path represents the URL path for a route e.g. /shelves/{shelf_id}/books
	Path string

	// http method e.g. http.MethodGet or empty to handle all types
	Method string

	// Route defines a mapping of URL paths to their corresponding HTTP methods and handlers.
	// It is a map where:
	// - The key is a Path (URL path).
	// - The value is another map where:
	// - The key is a Method (HTTP method).
	// - The value is an actionHandler(function to handle the request).
	Route map[Path]map[Method]actionHandler

	// Prefix represents a common prefix for a group of routes.
	// It is used to group related routes under a shared URL prefix (e.g., "/api/v1").
	Prefix string

	// GroupRoute defines a mapping of route prefixes to their corresponding Route maps.
	// It is a map where:
	// - The key is a Prefix (common URL prefix).
	// - The value is a Route (collection of paths and their handlers).
	// This allows for organizing routes into logical groups based on their prefixes.
	GroupRoute map[Prefix]Route
)

func NewGroup(group GroupRoute, middlewares ...httpmiddleware.Middleware) Route {
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
func NewHandler(handler http.HandlerFunc, middlewares ...httpmiddleware.Middleware) actionHandler {
	return actionHandler{
		handler:     handler,
		middlewares: middlewares,
	}
}

func BindRoutesToMux(mux *http.ServeMux, routes ...Route) {
	r := mergeRoutes(routes...)
	for route, handler := range r {
		for method, action := range handler {
			mux.Handle(fmt.Sprintf("%s %s", method, route), httpmiddleware.Chain(action.handler, action.middlewares...))
		}
	}
}

func mergeRoutes(routes ...Route) Route {
	r := Route{}
	for _, v := range routes {
		for k, v2 := range v {
			r[k] = v2
		}
	}
	return r
}
