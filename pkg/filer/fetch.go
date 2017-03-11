package filer

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oakmail/goqu"
	"github.com/oakmail/logrus"

	"github.com/oakmail/backend/pkg/models"
)

func (f *Filer) Fetch(c *gin.Context) {
	id := c.Param("id")

	var token models.Token
	if found, err := f.GQ.From("tokens").Where(goqu.I("id").Eq(id)).ScanStruct(&token); !found || err != nil {
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		} else {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Token not found"))
		}
		return
	}

	if token.Type != models.FetchToken {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token type"))
		return
	}

	var file string
	switch token.ReferenceType {
	case models.ResourceRef:
		var resource models.Resource
		if found, err := f.GQ.From("resources").Where(goqu.I("id").Eq(token.ResourceID)).ScanStruct(resource); !found || err != nil {
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
			} else {
				c.AbortWithError(http.StatusUnauthorized, errors.New("Resource not found"))
			}
			return
		}
		file = resource.File
	default:
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid reference type"))
		return
	}

	reader, err := f.Filesystem.Fetch(file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := io.Copy(c.Writer, reader); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err := f.GQ.From("tokens").Where(goqu.I("id").Eq(token.ID)).Delete(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
