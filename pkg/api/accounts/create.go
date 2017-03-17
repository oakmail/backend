package accounts

import (
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/api/passwords"
	"github.com/oakmail/backend/pkg/database"
	"github.com/oakmail/backend/pkg/models"
)

// Create creates a new account in the system
func (i *Impl) Create(c *gin.Context) {
	var input struct {
		Address  string `json:"address"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&input); err != nil {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidJSONInput)
		return
	}

	verr := []*errors.Error{}
	if input.Address == "" || !govalidator.IsEmail(input.Address) {
		verr = append(verr, errors.InvalidEmailFormat)
	} else if strings.SplitN(input.Address, "@", 2)[1] != i.Config.DefaultDomain {
		verr = append(verr, errors.InvalidEmailDomain)
	} else if count, err := i.GQ.From("addresses").Where(goqu.I("id").Eq(input.Address)).Count(); count > 0 || err != nil {
		if err != nil {
			panic(err)
		}

		verr = append(verr, errors.AddressIsTaken)
	}
	if len(input.Password) != 128 {
		verr = append(verr, errors.InvalidPasswordFormat)
	}
	if len(verr) != 0 {
		errors.Abort(c, http.StatusUnprocessableEntity, verr...)
		return
	}

	// do it first to prevent race condition
	address := models.Address{
		ID:           models.RemoveDots(models.NormalizeAddress(input.Address)),
		StyledID:     input.Address,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
		Account:      0,
	}
	if _, err := database.MustDataset(i.GQ.From("addresses").Insert(address).ResultingRow()).Select("id").ScanVal(&address.ID); err != nil {
		panic(err)
	}

	account := models.Account{
		DateCreated:  time.Now(),
		DateModified: time.Now(),
		Type:         "standard",
		MainAddress:  input.Address,
		Password:     passwords.Hash(input.Password),
	}
	if _, err := database.MustDataset(i.GQ.From("accounts").Insert(account).ResultingRow()).Select("id").ScanVal(&account.ID); err != nil {
		panic(err)
	}

	address.Account = account.ID
	if _, err := i.GQ.From("addresses").Where(goqu.I("id").Eq(address.ID)).Update(address).Exec(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, struct {
		models.Account
		Address models.Address `json:"_address"`
	}{
		Account: account,
		Address: address,
	})
}
