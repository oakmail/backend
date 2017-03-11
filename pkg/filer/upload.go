package filer

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/models"
)

func (f *Filer) Upload(c *gin.Context) {
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

	if token.Type != models.UploadToken {
		c.String(http.StatusUnauthorized, "Invalid token type")
		return
	}

	file, _, err := f.Filesystem.Upload(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

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

		if resource.File != "" {
			// todo consider if removing it here is supposed to stay or do we cleanup with daemon
			if err := f.Filesystem.Delete(resource.File); err != nil {
				c.String(http.StatusUnauthorized, err.Error())
				return
			}
		}

		if _, err := f.GQ.From("resources").Where(goqu.I("id").Eq(resource.ID)).Update(map[string]interface{}{
			"date_modified": time.Now(),
			"upload_token":  "",
			"file":          file,
		}).Exec(); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
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
