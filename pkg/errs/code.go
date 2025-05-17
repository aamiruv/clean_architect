package errs

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type ErrorCode int

const (
	CodeUnauthorized ErrorCode = iota + 1
	CodeForbiddenAccess
	CodeNotFound
	CodeExisted
	CodeInvalidArgument
	CodeInternal
)

func (e ErrorCode) HttpStatus() int {
	switch e {
	case CodeUnauthorized:
		return http.StatusUnauthorized

	case CodeForbiddenAccess:
		return http.StatusForbidden

	case CodeNotFound:
		return http.StatusNotFound

	case CodeExisted:
		return http.StatusConflict

	case CodeInvalidArgument:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

func (e ErrorCode) GRPCStatus() codes.Code {
	switch e {
	case CodeUnauthorized:
		return codes.Unauthenticated

	case CodeForbiddenAccess:
		return codes.PermissionDenied

	case CodeNotFound:
		return codes.NotFound

	case CodeExisted:
		return codes.AlreadyExists

	case CodeInvalidArgument:
		return codes.InvalidArgument

	default:
		return codes.Internal
	}
}
