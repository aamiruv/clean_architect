package middleware

import "net/http"

// SuperAdminRole will check request for authorization header
// and returns error if not found
func SuperAdminRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// todo: validate authorization token and role
		next.ServeHTTP(w, r)
	})
}
