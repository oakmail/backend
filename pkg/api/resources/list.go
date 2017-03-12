package resources

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
}

// List allows you to list resources
func (i *Impl) List(c *gin.Context) {
	var (
		token = c.MustGet("token").(models.Token)
	)

	var (
		owner uint64
		err   error
	)
	if ps := c.Param("id"); ps != "" {
		owner, err = strconv.ParseUint(ps, 10, 64)
		if err != nil {
			errors.Abort(c, http.StatusBadRequest, errors.InvalidOwnerID)
			return
		}
	} else if qs := c.Param("owner"); qs != "" {
		owner, err = strconv.ParseUint(qs, 10, 64)
		if err != nil {
			errors.Abort(c, http.StatusBadRequest, errors.InvalidOwnerID)
			return
		}
	}

	if owner != 0 {
		if !token.CheckOr([]string{
			"resources.owner:" + strconv.FormatUint(owner, 10) + ".read",
			"resources.list",
		}) {
			errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
			return
		}
	} else {
		if !token.Check(
			"resources.list",
		) {
			errors.Abort(c, http.StatusUnauthorized, errors.InsufficientTokenPermissions)
			return
		}
	}

	var (
		filters = []goqu.Expression{}
		verr    = []*errors.Error{}
	)
	if owner != 0 {
		filters = append(filters, goqu.I("owner").Eq(owner))
	}
	if file := c.Query("file"); file != "" {
		filters = append(filters, goqu.I("file").Eq(file))
	}

	if ts := c.Query("tags"); ts != "" {
		tl := strings.Split(ts, ",")

		if i.GQ.Dialect == "sqlite3" {
			// sqlite3
			for _, tag := range tl {
				filters = append(filters,
					goqu.I("tags").Like("%\""+tag+"\"%"),
				)
			}
		} else {
			// postgres
			for _, tag := range tl {
				filters = append(filters, goqu.L("ANY(tags) = ?", tag))
			}
		}
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
			filters = append(filters, goqu.I("date_created").Lte(dce))
		}
	} else {
		if dce.IsZero() {
			filters = append(filters, goqu.I("date_created").Gte(dcs))
		} else {
			filters = append(filters, goqu.I("date_created").Between(goqu.RangeVal{
				Start: dcs,
				End:   dce,
			}))
		}
	}

	if dms.IsZero() {
		if !dme.IsZero() {
			filters = append(filters, goqu.I("date_modified").Lte(dme))
		}
	} else {
		if dme.IsZero() {
			filters = append(filters, goqu.I("date_modified").Gte(dms))
		} else {
			filters = append(filters, goqu.I("date_modified").Between(goqu.RangeVal{
				Start: dms,
				End:   dme,
			}))
		}
	}

	query := i.GQ.From("resources")

	if len(filters) != 0 {
		query = query.Where(filters...)
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

	var resources []models.Resource
	if err := query.ScanStructs(&resources); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, resources)
}
