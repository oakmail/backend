package applications

import (
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/database"
	"github.com/oakmail/backend/pkg/models"
)

// Create creates a new application for the OAuth flow
func (i *Impl) Create(c *gin.Context) {
	var input struct {
		Owner       uint64 `json:"owner"`
		Callback    string `json:"callback"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		HomePage    string `json:"home_page"`
		Description string `json:"description"`
	}
	if err := c.BindJSON(&input); err != nil {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidJSONInput)
		return
	}

	var (
		token = c.MustGet("token").(models.Token)
	)

	if !token.Check(
		"applications.owner:" + strconv.FormatUint(input.Owner, 10) + ".create",
	) {
		errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
		return
	}

	verr := []*errors.Error{}
	if input.Callback == "" || !govalidator.IsURL(input.Callback) {
		verr = append(verr, errors.InvalidCallbackFormat)
	}
	if input.Name == "" {
		verr = append(verr, errors.ApplicationNameIsInvalid)
	}
	if input.Email == "" || !govalidator.IsEmail(input.Email) {
		verr = append(verr, errors.InvalidEmailFormat)
	}
	if input.HomePage == "" || !govalidator.IsURL(input.HomePage) {
		verr = append(verr, errors.InvalidHomePageFormat)
	}
	if len(verr) != 0 {
		errors.Abort(c, http.StatusUnprocessableEntity, verr...)
		return
	}

	application := models.Application{
		DateCreated:  time.Now(),
		DateModified: time.Now(),
		Owner:        input.Owner,
		Secret:       uniuri.NewLen(uniuri.UUIDLen),
		Callback:     input.Callback,
		Name:         input.Name,
		Email:        input.Email,
		HomePage:     input.HomePage,
		Description:  input.Description,
	}
	if _, err := database.MustDataset(i.GQ.From("applications").Insert(application).ResultingRow()).Select("id").ScanVal(&application.ID); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, application)
}
