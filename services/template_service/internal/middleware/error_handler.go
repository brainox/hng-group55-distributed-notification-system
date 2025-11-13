package middleware

import (
	"net/http"

	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/pkg/logger"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			logger.Log.Error("request error",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
			)

			response.Error(c, http.StatusInternalServerError, err, "Internal server error")
		}
	}
}
