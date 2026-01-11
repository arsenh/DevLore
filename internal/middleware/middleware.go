package middleware

import (
	"context"
	"net/http"

	"github.com/arsenh/DevLore/internal/auth"
)

type ctxKey string

const (
	CtxAuthError ctxKey = "authError"
	CtxUserID    ctxKey = "userID"
	CtxEmail     ctxKey = "email"
	CtxFullName  ctxKey = "fullName"
)

func JWTAuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("auth_token")
			if err != nil {
				ctx := context.WithValue(r.Context(), CtxAuthError, err)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			tokenStr := cookie.Value

			authUser, err := auth.ParseJWTToken(tokenStr)

			if err != nil {
				ctx := context.WithValue(r.Context(), CtxAuthError, err)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := context.WithValue(r.Context(), CtxAuthError, nil) // start with base context
			ctx = context.WithValue(ctx, CtxUserID, authUser.ID)     // chain correctly
			ctx = context.WithValue(ctx, CtxEmail, authUser.Email)
			ctx = context.WithValue(ctx, CtxFullName, authUser.FullName)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
