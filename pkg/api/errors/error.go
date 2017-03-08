package errors

import (
	"github.com/gin-gonic/gin"
)

// Error is the basic struct for API errors.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Message
}

// Abort cancels the request and returns the error list.
func Abort(c *gin.Context, code int, errors ...*Error) {
	c.JSON(code, gin.H{
		"errors": errors,
	})
	c.Abort()
}
