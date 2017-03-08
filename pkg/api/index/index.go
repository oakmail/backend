package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index is the landing page of the API
func (i *Impl) Index(c *gin.Context) {
	c.String(http.StatusOK, "trtl backend api v0\n")
}
