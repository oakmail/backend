package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oakmail/backend/pkg/api/errors"
)

// RequiresAuth ensures that user is logged in
func (i *Impl) RequiresAuth(c *gin.Context) {
	if _, ok := c.Get("token"); !ok {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidAuthorizationHeader)
		return
	}
}
