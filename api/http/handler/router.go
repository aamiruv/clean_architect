package handler

import (
	"log"
	"net/http"

	v2 "github.com/amirzayi/clean_architect/api/http/handler/v2"
	"github.com/amirzayi/clean_architect/api/http/middleware"
	"github.com/amirzayi/clean_architect/internal/service"
	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/amirzayi/rahjoo"
)

func Register(mux *http.ServeMux, logger *log.Logger, services *service.Services, authManager auth.Manager) {
	rahjoo.BindRoutesToMux(mux,
		v2.UserRoutes(middleware.LogRequestBody(logger), services.User, authManager),
		v2.AuthRoutes(services.Auth),
	)
}
