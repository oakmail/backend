package models

import (
	"time"

	"regexp"
	"strings"
	"unicode"
)

// Account is pretty much an account.
type Account struct {
	ID               uint64    `db:"id" json:"id" goqu:"skipinsert"`
	DateCreated      time.Time `db:"date_created" json:"date_created"`
	DateModified     time.Time `db:"date_modified" json:"date_modified"`
	Type             string    `db:"type" json:"type"`
	MainAddress      string    `db:"main_address" json:"main_address"`
	Identity         string    `db:"identity" json:"identity"`
	Password         string    `db:"password" json:"-"`
	Subscription     string    `db:"subscription" json:"subscription"`
	Blocked          bool      `db:"blocked" json:"blocked"`
	AltEmail         string    `db:"alt_email" json:"alt_email"`
	AltEmailVerified time.Time `db:"alt_email_verified" json:"alt_email_verified"`
}

var rNotASCII = regexp.MustCompile(`[^\w\.]`)

func RemoveDots(input string) string {
	if strings.Index(input, "@") != -1 {
		parts := strings.SplitN(input, "@", 2)

		return strings.Replace(parts[0], ".", "", -1) + "@" + parts[1]
	}

	return strings.Replace(input, ".", "", -1)
}

func NormalizeUsername(input string) string {
	return rNotASCII.ReplaceAllString(
		strings.ToLowerSpecial(unicode.TurkishCase, input),
		"",
	)
}

func NormalizeAddress(input string) string {
	parts := strings.SplitN(input, "@", 2)

	return NormalizeUsername(parts[0]) + "@" + strings.ToLowerSpecial(unicode.TurkishCase, parts[1])
}
