package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

//We define a custom type for the context key to avoid collisions
type contextkey string
const UserIDkey contextkey = "userID"

func AuthMiddleware(jwtSecret string ) func (http.Handler) http.Handler{
	return func (next http.Handler) http.Handler{
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
			// 1. Get the Authorization header (Format: "Bearer <token>")
			authHeader := r.Header.Get("Authorization")
			if authHeader == ""{
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// 2. Parse and validate the JWT token
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error){
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid{
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// 3. Extract user ID from token claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			}

			userID := claims["sub"].(string)

			// 4. Add user ID to request context and move to next handler
			ctx := context.WithValue(r.Context(), UserIDkey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	
}