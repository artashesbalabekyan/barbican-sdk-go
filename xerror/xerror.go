package xerror

import "net/http"

var (
	ErrKeyExists   = NewError(http.StatusBadRequest, "key already exists")
	ErrKeyNotFound = NewError(http.StatusNotFound, "key does not exist")
)

type Error struct {
	code    int
	message string
}

func NewError(code int, msg string) Error {
	return Error{
		code:    code,
		message: msg,
	}
}

// Status returns the HTTP status code of the error.
func (e Error) Status() int { return e.code }

func (e Error) Error() string { return e.message }
