package resources

import (
	"bytes"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/models"
)

// Update allows you to update an resource
func (i *Impl) Update(c *gin.Context) {
	var (
		token = c.MustGet("token").(*models.Token)
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

	var input struct {
		ID           uint64                 `json:"id"`
		DateCreated  time.Time              `json:"date_created"`
		DateModified time.Time              `json:"date_modified"`
		Owner        uint64                 `json:"owner"`
		Meta         map[string]interface{} `json:"meta"`
		Tags         []string               `json:"tags"`
		Body         []byte                 `json:"body"`
	}
	if err := c.BindJSON(&input); err != nil {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidJSONInput)
		return
	}

	if !token.CheckOr([]string{
		"resources." + strconv.FormatUint(resource.ID, 10) + ".update",
		"resources.owner:" + strconv.FormatUint(resource.Owner, 10) + ".update",
	}) {
		errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
		return
	}

	if !input.DateModified.Equal(resource.DateModified) ||
		!input.DateCreated.Equal(resource.DateCreated) ||
		input.ID != id || input.ID != resource.ID || input.Owner != resource.Owner {
		errors.Abort(c, http.StatusBadRequest, errors.OutdatedObjectUsedInPUT)
		return
	}

	if !reflect.DeepEqual(input.Meta, resource.Meta) {
		resource.Meta = input.Meta
	}
	if !reflect.DeepEqual(input.Tags, resource.Tags) {
		resource.Tags = input.Tags
	}
	if bytes.Compare(input.Body, resource.Body) != 0 {
		resource.Body = input.Body
	}
	resource.DateModified = time.Now()

	if _, err := i.GQ.From("resources").Where(goqu.I("id").Eq(id)).Update(resource).Exec(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, resource)
}
