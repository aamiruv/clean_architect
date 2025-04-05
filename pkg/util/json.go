package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("failed to decode json: %w", err)
	}
	return v, nil
}

var validate = validator.New()

func DecodeAndValidate[T any](r *http.Request) (T, error) {
	v, err := Decode[T](r)
	if err = validate.Struct(v); err != nil {
		return v, fmt.Errorf("failed to validate json: %w", err)
	}
	return v, nil
}
