package webserver

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Code  int
	Cause error
}

func (e Error) Unwrap() error {
	return e.Cause
}

func (e Error) Is(err error) bool {
	return errors.Is(err, e.Cause)
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%s)", e.Cause.Error(), http.StatusText(e.Code))
}
