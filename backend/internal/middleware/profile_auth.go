package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ProfileAuthMiddleware accepts the secret needed to decode the JWT
func ProfileAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract token from Authorization header
		// Format expected: "Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// c.Abort ensures no furter handlers are called.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// 2. SANITIZE: Spilt the string "Bareer <token>" to get just the token part.
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else if strings.Contains(authHeader, " ") {
			// If it has a space but isn't "Bearer", it's definitely wrong
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		
		// 3. VERIFY: We parse the token. 
        // The second argument is a "KeyFunc" where we provide our secret to verify the signature.
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

		// 5. INJECT: This is the most important part of Gin.
        // We extract the "sub" (UserID) from the token claims.
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            userID, ok := claims["sub"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
				c.Abort()
				return
			}
            
            // c.Set(key, value) stores the ID in Gin's local context for THIS specific request.
            c.Set("userID", userID) 
        }

		// 6. PROCEED: If everything is fine, let the request continue to the handler.
        c.Next()
	}
}