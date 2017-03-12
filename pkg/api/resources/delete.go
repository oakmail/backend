package resources

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/models"
)

// Delete allows you to delete resources
func (i *Impl) Delete(c *gin.Context) {
	var (
		token = c.MustGet("token").(models.Token)
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidIDFormat)
		return
	}

	var resource models.Resource
	if found, err := i.GQ.From("resources").Where(goqu.I("id").Eq(id)).ScanStruct(&resource); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusNotFound, errors.ResourceNotFound)
		return
	}

	if !token.CheckOr([]string{
		"resources." + strconv.FormatUint(resource.ID, 10) + ".delete",
		"resources.owner:" + strconv.FormatUint(resource.Owner, 10) + ".delete",
	}) {
		errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
		return
	}

	if _, err := i.GQ.From("resources").Where(goqu.I("id").Eq(id)).Delete().Exec(); err != nil {
		panic(err)
	}
	if _, err := i.GQ.From("tokens").Where(goqu.I("resource").Eq(id)).Delete().Exec(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, resource)
}
