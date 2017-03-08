package models

import (
	"database/sql/driver"
	"time"

	"github.com/oakmail/perms"
)

// Token is an implementation of both types of OAuth tokens
// implement some sorts of perms
type Token struct {
	ID           string      `db:"id" json:"id" goqu:"skipinsert"`
	DateCreated  time.Time   `db:"date_created" json:"date_created"`
	DateModified time.Time   `db:"date_modified" json:"date_modified"`
	Owner        uint64      `db:"owner" json:"owner"`
	ExpiryDate   time.Time   `db:"expiry_date" json:"expiry_date"`
	Type         TokenType   `db:"type" json:"type"`
	Perms        perms.Nodes `db:"perms" json:"perms"`
	Application  uint64      `db:"application" json:"application"`
}

// CheckAnd checks all passed perms
func (t *Token) CheckAnd(rp []string) bool {
	for _, p := range rp {
		if found, negated := t.Perms.Check(perms.MustParseNode(p)); !found || negated {
			return false
		}
	}

	return true
}

// CheckOr checks if there is at least one perm
func (t *Token) CheckOr(rp []string) bool {
	for _, p := range rp {
		if found, negated := t.Perms.Check(perms.MustParseNode(p)); found && !negated {
			return true
		}
	}

	return false
}

// Check checks if the perm is there
func (t *Token) Check(perm string) bool {
	found, negated := t.Perms.Check(perms.MustParseNode(perm))
	return found && !negated
}

// TokenType is either auth or code and defines what function a token serves.
type TokenType string

// AuthToken is for authorization, CodeToken is for 3rd party app OAuth flow.
const (
	AuthToken TokenType = "auth"
	CodeToken           = "code"
)

// Scan implements the database/sql Scanner interface
func (t *TokenType) Scan(value interface{}) error {
	*t = TokenType(string(value.([]byte)))
	return nil
}

// Value implements the database/sql Valuer interface
func (t TokenType) Value() (driver.Value, error) {
	return string(t), nil
}
