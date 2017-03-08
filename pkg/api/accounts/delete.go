package accounts

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/models"
)

// Delete deletes an account
func (i *Impl) Delete(c *gin.Context) {
	var (
		account = c.MustGet("account").(models.Account)
		token   = c.MustGet("token").(models.Token)
	)

	if ids := c.Param("id"); ids != "me" {
		id, err := strconv.ParseUint(ids, 10, 64)
		if err != nil {
			errors.Abort(c, http.StatusBadRequest, errors.InvalidIDFormat)
			return
		}

		if found, err := i.GQ.From("account").Where(goqu.I("id").Eq(id)).ScanStruct(&account); !found || err != nil {
			if err != nil {
				panic(err)
			}

			errors.Abort(c, http.StatusNotFound, errors.AccountNotFound)
			return
		}
	}

	if !token.Check(
		"accounts." + strconv.FormatUint(account.ID, 10) + ".delete",
	) {
		errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
		return
	}

	// account
	if _, err := i.GQ.From("accounts").Where(goqu.I("id").Eq(account.ID)).Delete().Exec(); err != nil {
		panic(err)
	}
	// addresses
	if _, err := i.GQ.From("addresses").Where(goqu.I("account").Eq(account.ID)).Delete().Exec(); err != nil {
		panic(err)
	}
	// tokens by app, applications
	if _, err := i.GQ.From("tokens").Where(goqu.I("application").In(
		i.GQ.From("applications").Where(goqu.I("applications.owner").Eq(account.ID)).Select("applications.id"),
	)).Delete().Exec(); err != nil {
		panic(err)
	}
	if _, err := i.GQ.From("applications").Where(goqu.I("owner").Eq(account.ID)).Delete().Exec(); err != nil {
		panic(err)
	}
	// all tokens
	if _, err := i.GQ.From("tokens").Where(goqu.I("owner").Eq(account.ID)).Delete().Exec(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, account)
}
