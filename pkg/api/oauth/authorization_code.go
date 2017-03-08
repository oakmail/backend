package oauth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/database"
	"github.com/oakmail/backend/pkg/models"
)

func (i *Impl) authorizationCode(c *gin.Context, codeID string, appID uint64, appSecret string) {
	var application models.Application
	if found, err := i.GQ.From("applications").Where(goqu.I("id").Eq(appID)).ScanStruct(&application); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusBadRequest, errors.ApplicationNotFound)
		return
	}

	if application.Secret != appSecret {
		errors.Abort(c, http.StatusBadRequest, errors.InvalidApplicationSecret)
		return
	}

	var code models.Token
	if found, err := i.GQ.From("tokens").Where(goqu.I("id").Eq(codeID)).ScanStruct(&code); !found || err != nil {
		if err != nil {
			panic(err)
		}

		errors.Abort(c, http.StatusBadRequest, errors.AuthorizationCodeNotFound)
		return
	}

	if code.Type != models.CodeToken || code.Application != application.ID {
		errors.Abort(c, http.StatusBadRequest, errors.AuthorizationCodeNotFound)
		return
	}

	// delete code first to prevent race cond
	if _, err := i.GQ.From("tokens").Where(goqu.I("id").Eq(code.ID)).Delete().Exec(); err != nil {
		panic(err)
	}

	// derive an auth token from the code
	token := models.Token{
		DateCreated:  time.Now(),
		DateModified: time.Now(),
		Owner:        code.Owner,
		ExpiryDate:   code.ExpiryDate,
		Type:         models.AuthToken,
		Perms:        code.Perms,
		Application:  code.Application,
	}
	if _, err := database.MustDataset(i.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, token)
	return
}
