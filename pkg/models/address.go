package models

import (
	"time"
)

// Address is the mapping between email addresses and accounts
type Address struct {
	ID           string    `db:"id" json:"id"`
	StyledID     string    `db:"styled_id" json:"styled_id"`
	DateCreated  time.Time `db:"date_created" json:"date_created"`
	DateModified time.Time `db:"date_modified" json:"date_modified"`
	Account      uint64    `db:"account" json:"account"`
	PublicKey    uint64    `db:"public_key" json:"public_key"`
}
