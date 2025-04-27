package v2

import (
	"net/http"

	"github.com/amirzayi/clean_architect/api/http/handler/v2/dto"
	appmiddleware "github.com/amirzayi/clean_architect/api/http/middleware"
	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/internal/service"
	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/amirzayi/clean_architect/pkg/jsonutil"
	"github.com/amirzayi/rahjoo"
	"github.com/amirzayi/rahjoo/middleware"
	"github.com/google/uuid"
)

type userRouter struct {
	userService service.User
}

func UserRoutes(logMiddleware middleware.Middleware, userService service.User, authManager auth.Manager) rahjoo.Route {
	user := &userRouter{userService: userService}
	return rahjoo.NewGroupRoute("/v2/users", rahjoo.Route{
		"": {
			http.MethodGet:  rahjoo.NewHandler(user.list),
			http.MethodPost: rahjoo.NewHandler(user.create),
		},
		"/{id}": {
			http.MethodGet:    rahjoo.NewHandler(user.get),
			http.MethodDelete: rahjoo.NewHandler(user.delete),
			http.MethodPut:    rahjoo.NewHandler(user.update),
		},
	}.SetMiddleware(
		appmiddleware.MustHaveAtLeastOneRole(authManager, []domain.UserRole{domain.UserRoleAdmin}),
	),
	)
}

func (u *userRouter) list(w http.ResponseWriter, r *http.Request) {
	users, err := u.userService.List(r.Context())
	if err != nil {
		return
	}
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserDomainToDTO(user))
	}

	jsonutil.Encode(w, http.StatusOK, map[string]any{"data": userResponses})
}

func (u *userRouter) create(w http.ResponseWriter, r *http.Request) {
	req, err := jsonutil.DecodeAndValidate[dto.CreateUserRequest](r)
	if err != nil {
		jsonutil.Encode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	user, err := u.userService.Create(r.Context(), req.ToDomain())
	if err != nil {
		jsonutil.Encode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonutil.Encode(w, http.StatusCreated, user)
}

func (u *userRouter) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		jsonutil.Encode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	err = u.userService.Delete(r.Context(), uid)
	if err != nil {
		jsonutil.Encode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (u *userRouter) update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		jsonutil.Encode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	req, err := jsonutil.DecodeAndValidate[dto.UpdateUserRequest](r)
	if err != nil {
		jsonutil.Encode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	user := req.ToDomain()
	user.ID = uid
	err = u.userService.Update(r.Context(), user)
	if err != nil {
		jsonutil.Encode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (u *userRouter) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		jsonutil.Encode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	user, err := u.userService.GetByID(r.Context(), uid)
	if err != nil {
		jsonutil.Encode(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	jsonutil.Encode(w, http.StatusOK, user)
}
