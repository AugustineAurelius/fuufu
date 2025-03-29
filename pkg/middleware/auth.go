package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Shield interface {
	GetSecret() string
}

func AuthMiddleware(shield Shield, next http.Handler, excludedPath ...string) http.Handler {
	slices.Sort(excludedPath)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := slices.BinarySearch(excludedPath, r.URL.Path); !ok {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Authorization header required"}`, http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(shield.GetSecret()), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
