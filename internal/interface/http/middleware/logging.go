package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"ppharma/backend/internal/domain/common"
)

func RequestLogger(log common.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info(c.Request.Context(), "request completed",
			common.Field{Key: "path", Value: c.Request.URL.Path},
			common.Field{Key: "method", Value: c.Request.Method},
			common.Field{Key: "status", Value: c.Writer.Status()},
			common.Field{Key: "duration_ms", Value: time.Since(start).Milliseconds()},
		)
	}
}
