package applications

import (
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/models"
)

// Update allows you to update an application
func (i *Impl) Update(c *gin.Context) {
	var (
		token = c.MustGet("token").(*models.Token)
	)

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

	var input struct {
		ID           uint64    `json:"id"`
		DateCreated  time.Time `json:"date_created"`
		DateModified time.Time `json:"date_modified"`
		Owner        uint64    `json:"owner"`
		Callback     string    `json:"callback"`
		Name         string    `json:"name"`
		Email        string    `json:"email"`
		HomePage     string    `json:"home_page"`
		Description  string    `json:"description"`
	}
	if err := c.BindJSON(&input); err != nil {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidJSONInput)
		return
	}

	if !token.CheckOr([]string{
		"applications." + strconv.FormatUint(application.ID, 10) + ".update",
		"applications.owner:" + strconv.FormatUint(application.Owner, 10) + ".update",
	}) {
		errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
		return
	}

	if !input.DateModified.Equal(application.DateModified) ||
		!input.DateCreated.Equal(application.DateCreated) ||
		input.ID != id || input.ID != application.ID || input.Owner != application.Owner {
		errors.Abort(c, http.StatusBadRequest, errors.OutdatedObjectUsedInPUT)
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

	if input.Callback != application.Callback {
		application.Callback = input.Callback
	}
	if input.Name != application.Name {
		application.Name = input.Name
	}
	if input.Email != application.Email {
		application.Email = input.Email
	}
	if input.HomePage != application.HomePage {
		application.HomePage = input.HomePage
	}
	if input.Description != application.Description {
		application.Description = input.Description
	}
	application.DateModified = time.Now()

	if _, err := i.GQ.From("applications").Where(goqu.I("id").Eq(id)).Update(application).Exec(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, application)
}
