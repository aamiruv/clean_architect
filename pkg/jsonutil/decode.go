package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/amirzayi/clean_architect/pkg/errs"
	"github.com/go-playground/validator/v10"
)

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, errs.New(fmt.Errorf("failed to decode json: %w", err), errs.CodeInvalidArgument)
	}
	return v, nil
}

var validate = validator.New()

func DecodeAndValidate[T any](r *http.Request) (T, error) {
	v, err := Decode[T](r)
	if err != nil {
		return v, err
	}

	if err = validate.Struct(v); err != nil {
		if validationError, ok := err.(validator.ValidationErrors); ok {
			errFields := make(map[string][]string)
			for _, err := range validationError {
				fieldName := strings.ToLower(err.Field())
				errField := ""
				switch err.Tag() {
				case "required":
					errField = fmt.Sprintf("the %s is required.", fieldName)

				case "email":
					errField = fmt.Sprintf("the %s must be a valid email address.", fieldName)

				case "min":
					errField = fmt.Sprintf("the %s must be at least %s characters.", fieldName, err.Param())

				case "max":
					errField = fmt.Sprintf("the %s may not be greater than %s characters.", fieldName, err.Param())

				default:
					errField = fmt.Sprintf("the %s is invalid.", fieldName)
				}
				errFields[fieldName] = append(errFields[fieldName], errField)
			}
			return v, errs.New(errors.New("given body is not valid"), errs.CodeInvalidArgument, errFields)
		}
		return v, errs.New(fmt.Errorf("failed to validate json: %w", err), errs.CodeInvalidArgument)
	}
	return v, nil
}
