package middleware

import (
	"net/http"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/golang-jwt/jwt/v5/request"
)

// MustHaveAtLeastOneRole will check request for authorization header
// and prohibit access if not found any roles.
func MustHaveAtLeastOneRole(authManager auth.Manager, roles []domain.UserRole) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := request.BearerExtractor{}.ExtractToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claims, err := authManager.VerifyToken(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var hasRole bool
			for _, role := range roles {
				if claims.UserRole == string(role) {
					hasRole = true
				}
			}
			if !hasRole {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
