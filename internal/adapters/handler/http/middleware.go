package http

import (
	"code-typing-auth-service/internal/core/ports"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type errorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status,omitempty"`
	Error     string    `json:"error,omitempty"`
	Message   string    `json:"message,omitempty"`
	Path      string    `json:"path,omitempty"`
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			var responseStatus int
			if errors.Is(err, ports.BadRequestError) {
				responseStatus = http.StatusBadRequest
			} else if errors.Is(err, ports.UnauthorizedError) {
				responseStatus = http.StatusUnauthorized
			} else if errors.Is(err, ports.ForbiddenError) {
				responseStatus = http.StatusForbidden
			} else if errors.Is(err, ports.NotFoundError) {
				responseStatus = http.StatusNotFound
			} else {
				responseStatus = http.StatusInternalServerError
			}
			c.JSON(responseStatus, errorResponse{
				Timestamp: time.Now(),
				Status:    responseStatus,
				Error:     http.StatusText(responseStatus),
				Message:   err.Err.Error(),
				Path:      c.Request.URL.Path,
			})
			c.Abort()
			return
		}
	}
}
