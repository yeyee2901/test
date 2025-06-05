package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-www-form-urlencoded")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func AttachRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check in header
		reqID := c.GetHeader("X-Request-Id")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		c.Set("X-Request-Id", reqID)
		c.Next()
	}
}
