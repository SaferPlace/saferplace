package service

import (
	"net/http"

	"connectrpc.com/connect"
)

// Service is webserver registered function to create a new service, aliased for convenience
type Service func(...connect.Interceptor) (string, http.Handler)
