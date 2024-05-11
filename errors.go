package core

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code int `json:"code"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error %d", e.Code)
}

func NewError(code int) *Error {
	return &Error{Code: code}
}

type ErrorHandler struct {
	Handle func(c *CaesarCtx, err error)
}

func RetrieveErrorCode(err error) int {
	if e, ok := err.(*Error); ok {
		return e.Code
	}

	return http.StatusInternalServerError
}
