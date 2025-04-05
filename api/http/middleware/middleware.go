package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

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

// DenyUnauthorizedClient will check request for authorization header
// and returns error if not found
func DenyUnauthorizedClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fixme: validate authorization token
		if r.Header.Get("Authorization") == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
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
