package v2

import (
	"net/http"

	"github.com/amirzayi/clean_architect/api/http/handler/v2/dto"
	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/internal/service"
	"github.com/amirzayi/clean_architect/pkg/util"
	"github.com/amirzayi/rahjoo"
)

type authRouter struct {
	authService service.Auth
}

func AuthRoutes(auth service.Auth) rahjoo.Route {
	router := &authRouter{authService: auth}

	return rahjoo.NewGroupRoute("/v2/auth", rahjoo.Route{
		"/register": {
			http.MethodPost: rahjoo.NewHandler(router.register),
		},
	}) // todo: add throttle middleware
}

func (a *authRouter) register(w http.ResponseWriter, r *http.Request) {
	in, err := util.DecodeAndValidate[dto.CreateUserRequest](r)
	if err != nil {
		_ = util.Encode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err = a.authService.Register(r.Context(), domain.Auth{
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Password:    in.Password,
	})
	if err != nil {
		_ = util.Encode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
}
