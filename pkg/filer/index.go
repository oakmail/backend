package filer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index is the landing page of the API
func (f *Filer) Index(c *gin.Context) {
	c.String(http.StatusOK, "trtl filer api v0\n")
}
