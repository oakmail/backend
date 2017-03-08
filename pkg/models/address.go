package models

import (
	"time"
)

type Address struct {
	ID           string    `db:"id" json:"id"`
	DateCreated  time.Time `db:"date_created" json:"date_created"`
	DateModified time.Time `db:"date_modified" json:"date_modified"`
	Account      uint64    `db:"account" json:"account"`
}
