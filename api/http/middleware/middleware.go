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
					http.Error(w, "request body is empty", http.StatusBadRequest)
					return
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if err = json.Unmarshal(b, &rb); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			logger.Println(rb)
			r.Body = io.NopCloser(bytes.NewBuffer(b))
			next.ServeHTTP(w, r)
		})
	}
}
