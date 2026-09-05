package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sudarshanmg/gotask/pkg/response"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				response.WriteError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				response.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			userID, err := strconv.ParseInt(claims.Subject, 10, 64)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "invalid subject claim")
				return
			}
			role := claims.Role

			ctx := context.WithValue(r.Context(), "userID", userID)
			ctx = context.WithValue(ctx, "role", role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
