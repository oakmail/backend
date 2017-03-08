package applications

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/models"
)

// Get allows you to fetch an application
func (i *Impl) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidIDFormat)
		return
	}

	var application models.Application
	if found, err := i.GQ.From("applications").Where(goqu.I("id").Eq(id)).ScanStruct(&application); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusNotFound, errors.ApplicationNotFound)
		return
	}

	c.JSON(http.StatusOK, application)
}
