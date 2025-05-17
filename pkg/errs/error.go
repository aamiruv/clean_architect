package errs

import (
	"errors"
	"fmt"
	"runtime"
)

type Error struct {
	Msg        string
	Code       ErrorCode
	StackTrace string
	Details    []any
}

func New(err error, code ErrorCode, details ...any) error {
	// hide internal errors from end users
	if code == CodeInternal {
		err = errors.New("internal error")
	}
	return &Error{
		Msg:        err.Error(),
		Code:       code,
		StackTrace: caller(),
		Details:    details,
	}
}

func NotFound(entity string) error {
	msg := fmt.Sprintf("%s not found", entity)
	return Error{
		Msg:        msg,
		Code:       CodeNotFound,
		StackTrace: caller(),
	}
}

func (e Error) Error() string {
	return e.Msg
}

// caller uses log.Lshortfile to format the caller
func caller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	return fmt.Sprintf("%s:%d", file, line)
}
