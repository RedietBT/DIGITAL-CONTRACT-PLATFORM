package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

//We define a custom type for the context key to avoid collisions
type contextkey string
const UserIDKey contextkey = "userID"

// AuthMiddleware is a middleware that validates JWT tokens and extracts user information
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
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	
}

// RoleMiddleware checks if the user has the required role to access the endpoint
func RoleMiddleware(requiredRole string, svc service.AuthService) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			// 1. Extract user ID from context
			userID, ok := r.Context().Value(UserIDKey).(string)

			if !ok {
				http.Error(w, "User ID not found in context", http.StatusNotFound)
			return
			}

			// 2. Call the service to get user 
			user, err := svc.GetUserByID(r.Context(),userID)
			if err !=nil{
				http.Error(w, "User not Found", http.StatusNotFound)
				return
			}

			//3. Compare userRole to requiredRole
			if user.Role != requiredRole{
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			//4. Proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}