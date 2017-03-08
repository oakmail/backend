package models

import (
	"time"
)

// Account is pretty much an account.
type Account struct {
	ID               uint64    `db:"id" json:"id" goqu:"skipinsert"`
	DateCreated      time.Time `db:"date_created" json:"date_created"`
	DateModified     time.Time `db:"date_modified" json:"date_modified"`
	Type             string    `db:"type" json:"type"`
	MainAddress      string    `db:"main_address" json:"main_address"`
	Identity         string    `db:"identity" json:"identity"`
	Password         string    `db:"password" json:"password"`
	Subscription     string    `db:"subscription" json:"subscription"`
	Blocked          bool      `db:"blocked" json:"blocked"`
	AltEmail         string    `db:"alt_email" json:"alt_email"`
	AltEmailVerified time.Time `db:"alt_email_verified" json:"alt_email_verified"`
}
