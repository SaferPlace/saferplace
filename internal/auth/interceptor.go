package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"safer.place/realtime/internal/database"
)

// NewAuthInterceptor checks each request for valid session.
func NewAuthInterceptor(db database.Database) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			session := extractSession(req.Header().Values("Cookie"))
			if session == "" {
				return nil,
					connect.NewError(connect.CodeUnauthenticated, errors.New("no valid token"))
			}

			if err := db.IsValidSession(ctx, session); err != nil {
				return nil,
					connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid session: %w", err))
			}

			return next(ctx, req)
		})
	})
}

func extractSession(cookies []string) string {
	for _, cookie := range cookies {
		if strings.HasPrefix(cookie, `Authorization="Bearer `) {
			session := strings.TrimPrefix(cookie, `Authorization="Bearer `)
			session = strings.TrimSuffix(session, `"`)
			return session
		}
	}

	return ""
}
