package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token != os.Getenv("BEARER_TOKEN") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
