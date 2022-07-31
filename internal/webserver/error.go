package webserver

import (
	"errors"
	"fmt"
	"net/http"
)

// Error is an internal Webserver error
type Error struct {
	Code  int
	Cause error
}

func (e Error) Unwrap() error {
	return e.Cause
}

// Is implements the errors.Is error comparision.
func (e Error) Is(err error) bool {
	return errors.Is(err, e.Cause)
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%s)", e.Cause.Error(), http.StatusText(e.Code))
}
