package handler

import (
	"net/http"

	. "github.com/AmirMirzayi/clean_architecture/api/middleware"
	"github.com/AmirMirzayi/clean_architecture/pkg/middleware"
	"github.com/AmirMirzayi/clean_architecture/pkg/router"
)

func userV2Routes(logMiddleware middleware.Middleware) router.Route {
	return router.NewGroup(router.GroupRoute{
		"/v2/users": {
			"": {
				http.MethodGet:  router.NewHandler(listUserV2),
				http.MethodPost: router.NewHandler(createUserV2),
			},
			"/{id}": {
				http.MethodGet:    router.NewHandler(getUserV2),
				http.MethodDelete: router.NewHandler(deleteUserV2, logMiddleware),
				http.MethodPut:    router.NewHandler(updateUserV2),
			},
		},
	}, SuperAdminRole)
}
func listUserV2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func createUserV2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func deleteUserV2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func updateUserV2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getUserV2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
