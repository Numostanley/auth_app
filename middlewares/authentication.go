package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Numostanley/d8er_app/utils"
)

type contextKey string

const (
	decodedTokenKey contextKey = "decodedToken"
	userKey         contextKey = "user"
)

func AuthenticationMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticator := utils.TokenAuthentication{}

		decodedToken, user, err := authenticator.Authenticate(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication error: %v", err), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), decodedTokenKey, decodedToken)
		ctx = context.WithValue(ctx, userKey, user)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
