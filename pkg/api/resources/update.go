package resources

import (
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/database"
	"github.com/oakmail/backend/pkg/models"
)

// Update allows you to update an resource
func (i *Impl) Update(c *gin.Context) {
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

	var input struct {
		ID           uint64                 `json:"id"`
		DateCreated  time.Time              `json:"date_created"`
		DateModified time.Time              `json:"date_modified"`
		Owner        uint64                 `json:"owner"`
		Meta         map[string]interface{} `json:"meta"`
		Tags         []string               `json:"tags"`
		Upload       bool                   `json:"_upload"`
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
	resource.DateModified = time.Now()

	if input.Upload {
		// todo remove existing upload token

		upload := models.Token{
			DateCreated:   time.Now(),
			DateModified:  time.Now(),
			Owner:         input.Owner,
			ExpiryDate:    time.Now().Add(time.Hour),
			Type:          models.UploadToken,
			ReferenceType: "resource",
			ReferenceID:   resource.ID,
		}
		if _, err := database.MustDataset(i.GQ.From("tokens").Insert(upload).ResultingRow()).Select("id").ScanVal(&upload.ID); err != nil {
			panic(err)
		}

		resource.UploadToken = upload.ID
	}

	if _, err := i.GQ.From("resources").Where(goqu.I("id").Eq(id)).Update(resource).Exec(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, resource)
}
