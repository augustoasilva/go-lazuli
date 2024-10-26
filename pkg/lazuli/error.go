package lazuli

import (
	"fmt"
	"io"
	"net/http"
)

type Error struct {
	Code    int
	Message string
	Details string
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"status: %d, error: %s, details: %s",
		e.Code,
		e.Message,
		e.Details,
	)
}

func newError(code int, message string, details string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func newErrorFromResponse(res *http.Response, message string) *Error {
	defer res.Body.Close()
	details := "could not get detail"
	if res.Body != nil {
		b, err := io.ReadAll(res.Body)
		if err == nil {
			details = string(b)
		}
	}
	return newError(
		res.StatusCode,
		message,
		details,
	)
}
