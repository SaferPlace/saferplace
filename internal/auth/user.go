package auth

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"github.com/saferplace/webserver-go/middleware"
)

var (
	ErrUserUnauthenticated = errors.New("user unauthenticated")
)

// NewUserAuthInterceptor in the future will extract the tokens sent by the users to then convert
// them into the user emails, but for now it just reads the user email and does nothing with it, but
// does not let the request through if the header is missing.
func NewUserAuthInterceptor() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			emailHeader := req.Header().Get("X-Email")
			if emailHeader == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, ErrUserUnauthenticated)
			}

			return next(ctx, req)
		})
	})
}

// NewUserAuthMiddleware
func NewUserAuthMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			emailHeader := req.Header.Get("X-Email")
			if emailHeader == "" {
				http.Error(w, ErrUserUnauthenticated.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}
