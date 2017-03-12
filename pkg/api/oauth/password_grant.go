package oauth

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"
	"github.com/oakmail/perms"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/api/passwords"
	"github.com/oakmail/backend/pkg/database"
	"github.com/oakmail/backend/pkg/models"
)

const standardScope = `accounts.[uid].delete
accounts.[uid].update
applications.owner:[uid].create
applications.owner:[uid].delete
applications.owner:[uid].read
applications.owner:[uid].update
resources.owner:[uid].create
resources.owner:[uid].delete
resources.owner:[uid].read
resources.owner:[uid].update
tokens.owner:[uid].create
tokens.owner:[uid].delete
tokens.owner:[uid].read`

func (i *Impl) passwordGrant(c *gin.Context, addr, password string, appID uint64, expiryDate time.Time) {
	var application models.Application
	if found, err := i.GQ.From("applications").Where(goqu.I("id").Eq(appID)).ScanStruct(&application); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusBadRequest, errors.ApplicationNotFound)
		return
	}

	// now < expiry < 2 weeks later for pwd grants
	if expiryDate.Before(time.Now()) || expiryDate.After(time.Now().Add(time.Hour*24*14)) {
		errors.Abort(c, http.StatusUnprocessableEntity, errors.InvalidExpiryDate)
		return
	}

	var address models.Address
	if found, err := i.GQ.From("addresses").Where(goqu.I("id").Eq(addr)).ScanStruct(&address); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusBadRequest, errors.AccountNotFound)
		return
	}

	var account models.Account
	if found, err := i.GQ.From("accounts").Where(goqu.I("id").Eq(address.Account)).ScanStruct(&account); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusBadRequest, errors.AccountNotFound)
		return
	}

	if !passwords.Verify(account.Password, password) {
		errors.Abort(c, http.StatusUnauthorized, errors.InvalidPassword)
		return
	}

	// derive pwd grant tokens based on acc type
	var permissions perms.Nodes
	if account.Type == "admin" {
		permissions = perms.Nodes{perms.MustParseNode("*")} // TODO: might not match everything, check manually later
	} else if account.Type == "standard" {
		permissions = perms.MustParseNodes(
			[]byte(strings.Replace(standardScope, "[uid]", strconv.FormatUint(account.ID, 10), -1)),
		)
	}

	token := &models.Token{
		DateCreated:  time.Now(),
		DateModified: time.Now(),
		Owner:        account.ID,
		ExpiryDate:   expiryDate,
		Type:         models.AuthToken,
		Perms:        permissions,
		Application:  application.ID,
	}
	if _, err := database.MustDataset(i.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, token)
	return
}
