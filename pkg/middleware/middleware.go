// Package middleware contains http middlewares which acts before request get into the server handler
package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// Chain chains http middlewares onto handler.
// Last one will execute at first
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func Recovery(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Printf("recovered: %v", rec)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func MeterResponseTime(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		defer func() {
			logger.Printf("Completed in %v", time.Since(now))
		}()
		next.ServeHTTP(w, r)
	})
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

func EnforceJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			http.Error(w, "Content-Type header is not set", http.StatusBadRequest)
			return
		}
		mt, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(w, "Content-Type header is not set", http.StatusBadRequest)
			return
		}
		if mt != "application/json" {
			http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogRequest(next http.Handler, logger *log.Logger) http.Handler {
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
