package middleware

import "github.com/gin-gonic/gin"

// CORSMiddleware allows the Global Swagger UI to fetch the doc.json
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // In production, replace "*" with your specific domain
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}