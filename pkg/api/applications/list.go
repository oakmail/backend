package applications

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/goqu"

	"github.com/oakmail/backend/pkg/api/errors"
	"github.com/oakmail/backend/pkg/models"
)

var orderFields = map[string]struct{}{
	"id":            {},
	"date_created":  {},
	"date_modified": {},
	"owner":         {},
	"secret":        {},
	"callback":      {},
	"name":          {},
	"email":         {},
	"home_page":     {},
}

// List allows you to (who would've known!) list applications in the system
func (i *Impl) List(c *gin.Context) {
	var (
		token = c.MustGet("token").(*models.Token)
	)

	var (
		owner uint64
		err   error
	)
	if os := c.Query("owner"); os != "" {
		owner, err = strconv.ParseUint(os, 10, 64)
		if err != nil {
			errors.Abort(c, http.StatusBadRequest, errors.InvalidOwnerID)
			return
		}

		if !token.CheckOr([]string{
			"applications.owner:" + strconv.FormatUint(owner, 10) + ".read",
			"applications.list",
		}) {
			errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
			return
		}
	} else {
		if !token.Check(
			"applications.list",
		) {
			errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
			return
		}
	}

	var (
		filters = goqu.Ex{}
		verr    = []*errors.Error{}
	)
	if owner != 0 {
		filters["owner"] = owner
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if email := c.Query("email"); email != "" {
		filters["email"] = email
	}
	if homePage := c.Query("home_page"); homePage != "" {
		filters["home_page"] = homePage
	}

	var (
		dcq = c.Query("date_created")
		dcp = strings.SplitN(dcq, "~", 2)
		dcs time.Time
		dce time.Time
		dmq = c.Query("date_created")
		dmp = strings.SplitN(dmq, "~", 2)
		dms time.Time
		dme time.Time
	)
	if dcq != "" {
		if dcp[0] != "" {
			dms, err = time.Parse(time.RFC3339Nano, dcp[0])
			if err != nil {
				verr = append(verr, errors.InvalidDateCreatedStart)
			}
		}
		if dcp[1] != "" {
			dme, err = time.Parse(time.RFC3339Nano, dcp[1])
			if err != nil {
				verr = append(verr, errors.InvalidDateCreatedEnd)
			}
		}
	}
	if dmq != "" {
		if dmp[0] != "" {
			dms, err = time.Parse(time.RFC3339Nano, dmp[0])
			if err != nil {
				verr = append(verr, errors.InvalidDateModifiedStart)
			}
		}
		if dmp[1] != "" {
			dme, err = time.Parse(time.RFC3339Nano, dmp[1])
			if err != nil {
				verr = append(verr, errors.InvalidDateModifiedEnd)
			}
		}
	}
	if len(verr) != 0 {
		errors.Abort(c, http.StatusBadRequest, verr...)
		return
	}

	if dcs.IsZero() {
		if !dce.IsZero() {
			filters["date_created"] = goqu.Op{
				"lte": dce,
			}
		}
	} else {
		if dce.IsZero() {
			filters["date_created"] = goqu.Op{
				"gte": dcs,
			}
		} else {
			filters["date_created"] = goqu.Op{
				"between": goqu.RangeVal{
					Start: dcs,
					End:   dce,
				},
			}
		}
	}

	if dms.IsZero() {
		if !dme.IsZero() {
			filters["date_modified"] = goqu.Op{
				"lte": dme,
			}
		}
	} else {
		if dme.IsZero() {
			filters["date_modified"] = goqu.Op{
				"gte": dms,
			}
		} else {
			filters["date_modified"] = goqu.Op{
				"between": goqu.RangeVal{
					Start: dms,
					End:   dme,
				},
			}
		}
	}

	query := i.GQ.From("applications")
	if len(filters) != 0 {
		query = query.Where(filters)
	}

	if oq := c.Query("order"); oq != "" {
		op := strings.Split(oq, ",")
		ol := []goqu.OrderedExpression{}
		for _, field := range op {
			// Descending only then the first char is -
			ascending := field[0] != '-'

			// Trim the ordering prefix if passed
			if field[0] == '-' || field[0] == '+' {
				field = field[1:]
			}

			if _, ok := orderFields[field]; !ok {
				continue
			}

			if ascending {
				ol = append(ol, goqu.I(field).Asc())
			} else {
				ol = append(ol, goqu.I(field).Desc())
			}

		}

		if len(ol) != 0 {
			query = query.Order(ol...)
		}
	}

	if sq := c.Query("skip"); sq != "" {
		skip, err := strconv.Atoi(sq)
		if err != nil {
			errors.Abort(c, http.StatusUnprocessableEntity, errors.InvalidSkipFormat)
			return
		}

		query = query.Offset(uint(skip))
	}
	if lq := c.Query("limit"); lq != "" {
		limit, err := strconv.Atoi(lq)
		if err != nil {
			errors.Abort(c, http.StatusUnprocessableEntity, errors.InvalidLimitFormat)
			return
		}

		query = query.Limit(uint(limit))
	}

	var applications []models.Application
	if err := query.ScanStructs(applications); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, applications)
}
