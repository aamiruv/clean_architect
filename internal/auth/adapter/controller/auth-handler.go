package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/amirzayi/clean_architec/api/handler/v2/dto"
	"github.com/amirzayi/clean_architec/internal/auth/domain"
	"github.com/amirzayi/clean_architec/pkg/util"
)

type AuthRestHandler struct {
	useCase AuthUseCase
}

func NewAuthRestHandler(authUseCase AuthUseCase) *AuthRestHandler {
	return &AuthRestHandler{useCase: authUseCase}
}

func (h *AuthRestHandler) Register(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	in, err := util.Decode[dto.CreateUserRequest](r)
	if err != nil {
		_ = util.Encode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	err = validate.Struct(in)
	if err != nil {
		_ = util.Encode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	err = h.useCase.Register(r.Context(), domain.Auth{
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
