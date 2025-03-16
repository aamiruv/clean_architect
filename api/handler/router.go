package handler

import (
	"log"
	"net/http"

	v2 "github.com/amirzayi/clean_architec/api/handler/v2"
	"github.com/amirzayi/clean_architec/internal/auth/adapter/controller"
	"github.com/amirzayi/clean_architec/pkg/httpmiddleware"
	"github.com/amirzayi/clean_architec/pkg/router"
)

func Register(mux *http.ServeMux, logger *log.Logger, authUseCase controller.AuthUseCase) {
	lr := func(next http.Handler) http.Handler {
		return httpmiddleware.LogRequestBody(next, logger)
	}

	router.BindRoutesToMux(mux,
		v2.UserRoutes(lr, authUseCase))
}
