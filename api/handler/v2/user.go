package v2

import (
	"net/http"

	"github.com/amirzayi/clean_architec/api/appmiddleware"
	"github.com/amirzayi/clean_architec/internal/auth/adapter/controller"
	"github.com/amirzayi/clean_architec/pkg/httpmiddleware"
	"github.com/amirzayi/clean_architec/pkg/router"
)

func UserRoutes(logMiddleware httpmiddleware.Middleware, authUseCase controller.AuthUseCase) router.Route {
	authHandler := controller.NewAuthRestHandler(authUseCase)
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
			"/register": {
				http.MethodPost: router.NewHandler(authHandler.Register),
			},
		},
	}, appmiddleware.SuperAdminRole)
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
