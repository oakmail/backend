package oauth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oakmail/backend/pkg/api/errors"
)

// OAuth handles the /oauth token creation calls.
func (i *Impl) OAuth(c *gin.Context) {
	var input struct {
		GrantType   string    `json:"grant_type"`
		Code        string    `json:"code"`
		Application uint64    `json:"application"`
		Secret      string    `json:"secret"`
		Address     string    `json:"address"`
		Password    string    `json:"password"`
		ExpiryDate  time.Time `json:"expiry_date"`
	}
	if err := c.BindJSON(&input); err != nil {
		i.Log.Print(err)
		errors.Abort(c, http.StatusBadRequest, errors.InvalidJSONInput)
		return
	}

	if input.GrantType == "authorization_code" {
		i.authorizationCode(
			c,
			input.Code,
			input.Application,
			input.Secret,
		)
		return
	} else if input.GrantType == "password_grant" {
		i.passwordGrant(
			c,
			input.Address,
			input.Password,
			input.Application,
			input.ExpiryDate,
		)
		return
	}

	errors.Abort(c, http.StatusBadRequest, errors.InvalidGrantType)
	return
}
