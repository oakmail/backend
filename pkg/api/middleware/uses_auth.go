package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/models"
)

var alphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// UsesAuth loads up token and account information given a request with a token specified
func (i *Impl) UsesAuth(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidAuthorizationHeader)
		return
	}

	id := parts[1]
	if len(id) != 20 || !alphanumeric.MatchString(id) {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidAuthorizationHeader)
		return
	}

	var token models.Token
	if found, err := i.GQ.From("tokens").Where(goqu.I("id").Eq(id)).ScanStruct(&token); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusUnauthorized, errors.TokenNotFound)
		return
	}

	var account models.Account
	if found, err := i.GQ.From("accounts").Where(goqu.I("id").Eq(token.Owner)).ScanStruct(&account); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusInternalServerError, errors.DatabaseInconsistency)
		return
	}

	if token.Type != models.AuthToken {
		errors.Abort(c, http.StatusUnauthorized, errors.InvalidTokenTypeMustBeAuth)
		return
	}

	if token.ExpiryDate.Before(time.Now()) {
		errors.Abort(c, http.StatusUnauthorized, errors.AuthenticationTokenExpired)
		return
	}

	if account.Blocked {
		errors.Abort(c, http.StatusUnauthorized, errors.AccountIsBlocked)
		return
	}

	c.Set("account", account)
	c.Set("token", token)

	c.Header("X-Authenticated-As", strconv.FormatUint(account.ID, 10)+"; "+account.MainAddress)
	c.Header("X-Authenticated-Perms", token.Perms.String())
}
