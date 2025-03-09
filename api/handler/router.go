package handler

import (
	"log"
	"net/http"

	"github.com/AmirMirzayi/clean_architecture/pkg/middleware"
	"github.com/AmirMirzayi/clean_architecture/pkg/router"
)

func Register(mux *http.ServeMux, logger *log.Logger) {
	lr := func(next http.Handler) http.Handler {
		return middleware.LogRequest(next, logger)
	}

	router.BindRoutesToMux(mux,
		userV2Routes(lr))
}
