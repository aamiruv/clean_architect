package router

import "net/http"

func RegisterHttpRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ping", Ping)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
