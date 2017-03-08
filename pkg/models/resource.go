package models

import (
	"time"
)

type Resource struct {
	ID           uint64                 `db:"id" json:"id" goqu:"skipinsert"`
	DateCreated  time.Time              `db:"date_created" json:"date_created"`
	DateModified time.Time              `db:"date_modified" json:"date_modified"`
	Owner        uint64                 `db:"owner" json:"owner"`
	Meta         map[string]interface{} `db:"meta" json:"meta"`
	Tags         []string               `db:"tags" json:"tags"`
	Body         []byte                 `db:"body" json:"body"`
}
