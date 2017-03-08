package models

import (
	"time"
)

// Application is an entity used for OAuth token creation
type Application struct {
	ID           uint64    `db:"id" json:"id" goqu:"skipinsert"`
	DateCreated  time.Time `db:"date_created" json:"date_created"`
	DateModified time.Time `db:"date_modified" json:"date_modified"`
	Owner        uint64    `db:"owner" json:"owner"`
	Secret       string    `db:"secret" json:"secret"`
	Callback     string    `db:"callback" json:"callback"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	HomePage     string    `db:"home_page" json:"home_page"`
	Description  string    `db:"description" json:"description"`
}
