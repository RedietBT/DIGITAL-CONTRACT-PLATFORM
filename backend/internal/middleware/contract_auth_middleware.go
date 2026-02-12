package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ContractMiddleware accepts the secret needed to decode the JWT
func ContractMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract token from Authorizaion header
		// Format: Authorization "Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// c.Abort ensures no further handler are called
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return 
		}

		// 2. SANITIZE: Split the string "Bearer <token>" to get just the token part
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return 
		}

		tokenString := parts[1]

		// 3. VERIFY: We parse the token.
		// The second argument is a "keyFunc" where we provide our secret to verify the signature.
		token, err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error)  {
			// We ensure the token is using the expected signing method (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
			
		})

		// 4. VALIDATE: If the token is expired or the signature is wrong, we block the request.
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return 
		}

		// 5. INJECT: This is the most important part of Gin
		//We extract the "sub" (UserID) from the token claims.
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["sub"].(string)
			
			//c.Set(key, value) stores the ID in Gin's local context for THIS specific request.
			c.Set("user_id", userID)
		}

		// 6. PROCEED: If everything is fine, let the request continue to the handler.
		c.Next()
	}
}