package v2

import (
	"net/http"

	appmiddleware "github.com/amirzayi/clean_architect/api/http/middleware"
	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/internal/service"
	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/amirzayi/rahjoo"
	"github.com/amirzayi/rahjoo/middleware"
)

type userRouter struct {
	userService service.User
}

func UserRoutes(logMiddleware middleware.Middleware, userService service.User, authManager auth.Manager) rahjoo.Route {
	user := &userRouter{userService: userService}

	return rahjoo.NewGroup(rahjoo.GroupRoute{
		"/v2/users": {
			"": {
				http.MethodGet:  rahjoo.NewHandler(user.list),
				http.MethodPost: rahjoo.NewHandler(user.create),
			},
			"/{id}": {
				http.MethodGet:    rahjoo.NewHandler(user.get),
				http.MethodDelete: rahjoo.NewHandler(user.delete, logMiddleware),
				http.MethodPut:    rahjoo.NewHandler(user.update),
			},
		},
	}, appmiddleware.MustHaveAtLeastOneRole(authManager, []domain.UserRole{domain.UserRoleAdmin}))
}

func (u *userRouter) list(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (u *userRouter) create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (u *userRouter) delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (u *userRouter) update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (u *userRouter) get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
