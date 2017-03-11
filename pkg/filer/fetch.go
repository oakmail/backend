package filer

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/models"
)

func (f *Filer) Fetch(c *gin.Context) {
	id := c.Param("id")

	var token models.Token
	if found, err := f.GQ.From("tokens").Where(goqu.I("id").Eq(id)).ScanStruct(&token); !found || err != nil {
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
		} else {
			c.String(http.StatusUnauthorized, "Token not found")
		}
		return
	}

	if token.Type != models.FetchToken {
		c.String(http.StatusUnauthorized, "Invalid token type")
		return
	}

	var file string
	switch token.ReferenceType {
	case models.ResourceRef:
		var resource models.Resource
		if found, err := f.GQ.From("resources").Where(goqu.I("id").Eq(token.ReferenceID)).ScanStruct(resource); !found || err != nil {
			if err != nil {
				c.String(http.StatusUnauthorized, err.Error())
			} else {
				c.String(http.StatusUnauthorized, "Resource not found")
			}
			return
		}
		file = resource.File
	default:
		c.String(http.StatusBadRequest, "Invalid reference type")
		return
	}

	reader, err := f.Filesystem.Fetch(file)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := f.GQ.From("tokens").Where(goqu.I("id").Eq(token.ID)).Delete().Exec(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}
