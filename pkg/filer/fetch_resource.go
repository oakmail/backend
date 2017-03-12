package filer

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/models"
)

func (f *Filer) FetchResource(c *gin.Context) {
	var (
		token = c.MustGet("token").(models.Token)
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid ID format")
		return
	}

	var resource models.Resource
	if found, err := f.GQ.From("resources").Where(goqu.I("id").Eq(id)).ScanStruct(&resource); !found || err != nil {
		if err != nil {
			panic(err)
		}

		c.String(http.StatusUnauthorized, "Resource not found")
		return
	}

	if !token.CheckOr([]string{
		"resources." + strconv.FormatUint(resource.ID, 10) + ".delete",
		"resources.owner:" + strconv.FormatUint(resource.Owner, 10) + ".delete",
	}) {
		c.String(http.StatusUnauthorized, "Insufficient token permissions.")
		return
	}

	reader, err := f.Filesystem.Fetch(resource.File)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}
