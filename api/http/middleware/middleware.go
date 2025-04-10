package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5/request"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/pkg/auth"
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

func MeterResponseTime(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			defer func() {
				logger.Printf("Completed in %v", time.Since(now))
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func LogRequestBody(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var rb any
			b, err := io.ReadAll(r.Body)
			if err != nil {
				if errors.Is(err, io.EOF) {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if err = json.Unmarshal(b, &rb); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			logger.Print(rb)
			r.Body = io.NopCloser(bytes.NewBuffer(b))
			next.ServeHTTP(w, r)
		})
	}
}
