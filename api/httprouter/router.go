package httprouter

import "net/http"

func New() *http.ServeMux {
	handler := http.NewServeMux()
	handler.HandleFunc("/ping", ping)

	return handler
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
