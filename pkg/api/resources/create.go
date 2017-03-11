package resources

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/database"
	"github.com/oakmail/backend/pkg/models"
)

// Create creates a new resource
func (i *Impl) Create(c *gin.Context) {
	var input struct {
		Owner uint64                 `json:"owner"`
		Meta  map[string]interface{} `json:"meta"`
		Tags  []string               `json:"tags"`
	}
	if err := c.BindJSON(&input); err != nil {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidJSONInput)
		return
	}

	var (
		token = c.MustGet("token").(models.Token)
	)

	if !token.Check(
		"resources.owner:" + strconv.FormatUint(input.Owner, 10) + ".create",
	) {
		errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
		return
	}

	resource := models.Resource{
		DateCreated:  time.Now(),
		DateModified: time.Now(),
		Owner:        input.Owner,
		Meta:         input.Meta,
		Tags:         input.Tags,
	}
	if _, err := database.MustDataset(i.GQ.From("resources").Insert(resource).ResultingRow()).Select("id").ScanVal(&resource.ID); err != nil {
		panic(err)
	}

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
	if _, err := i.GQ.From("resources").Where(goqu.I("id").Eq(resource.ID)).Update(resource).Exec(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, resource)
}
