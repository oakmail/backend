package accounts

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/models"
)

// Get fetches an account
func (i *Impl) Get(c *gin.Context) {
	ids := c.Param("id")

	if ids == "me" {
		account, found := c.Get("account")
		if !found {
			errors.Abort(c, http.StatusBadRequest, errors.InvalidIDFormat)
			return
		}

		c.JSON(http.StatusOK, account)
		return
	}

	id, err := strconv.ParseUint(ids, 10, 64)
	if err != nil {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidIDFormat)
		return
	}

	var account models.Account
	if found, err := i.GQ.From("accounts").Where(goqu.I("id").Eq(id)).ScanStruct(&account); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusNotFound, errors.AccountNotFound)
		return
	}

	c.JSON(http.StatusOK, account)
	return
}
