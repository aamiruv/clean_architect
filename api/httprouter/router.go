package httprouter

import (
	"io"
	"log"
	"net/http"

	"github.com/AmirMirzayi/clean_architecture/pkg/middleware"
	"github.com/AmirMirzayi/clean_architecture/pkg/router"
)

func Register(mux *http.ServeMux, logger *log.Logger) {
	lr := func(next http.Handler) http.Handler {
		return middleware.LogRequest(next, logger)
	}

	r := router.Route{"/ping": {
		http.MethodGet:  router.NewHandler(ping),
		http.MethodPost: router.NewHandler(dong),
	},
		"/p2":  {"": router.NewHandler(final, middleware.DenyUnauthorizedClient)},
		"/log": {http.MethodPost: router.NewHandler(l, lr)},
	}

	router.BindRoutes(mux, r)
}
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func dong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dong"))
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	w.Write([]byte("OK"))
}

func l(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	w.Write(b)
}
