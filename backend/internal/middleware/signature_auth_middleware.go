package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// SignatureMiddleware follows the same pattern as ContractMiddleware for consistency
func SignatureMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// 2. SANITIZE: Handle the "Bearer " prefix exactly like ContractMiddleware
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else if strings.Contains(authHeader, " ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// 3. VERIFY: Parse using the HMAC signing method
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		// 4. VALIDATE: Block if expired or invalid
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 5. INJECT: Extract the "sub" and convert to UUID
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			sub, ok := claims["sub"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
				c.Abort()
				return
			}

			// Convert string to UUID to match the Contract service pattern
			userID, err := uuid.Parse(sub)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format"})
				c.Abort()
				return
			}

			// Store as uuid.UUID in context
			c.Set("userID", userID)
		}

		// 6. PROCEED
		c.Next()
	}
}