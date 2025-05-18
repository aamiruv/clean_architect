package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/amirzayi/clean_architect/pkg/errs"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return fmt.Errorf("failed to encode json: %w", err)
	}
	return nil
}

func EncodeError(w http.ResponseWriter, err error) error {
	msg := fmt.Sprintf("unexpected error: %v", err)
	code := 0
	statusCode := http.StatusInternalServerError
	details := []any{}
	var appErr *errs.Error
	if errors.As(err, &appErr) {
		statusCode = appErr.Code.HttpStatus()
		msg = appErr.Error()
		code = int(appErr.Code)
		details = appErr.Details
	}
	return Encode(w, statusCode, map[string]any{
		"message": msg,
		"code":    code,
		"details": details,
	})
}
