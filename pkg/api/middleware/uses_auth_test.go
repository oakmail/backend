package middleware_test

import (
	"testing"
	"time"

	"github.com/dchest/uniuri"
	"github.com/oakmail/perms"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/oakmail/backend/pkg/api/test"
	"github.com/oakmail/backend/pkg/database"
	"github.com/oakmail/backend/pkg/models"
)

func TestUsesAuth(t *testing.T) {
	api := test.InitAPI()
	defer api.Cleanup()

	Convey("UsesAuth middleware should", t, func() {
		Convey("Pass through with no Authorization passed", func() {
			So(
				api.Bl.Get("/").
					Expect(t).
					BodyMatchString("trtl backend api v0").
					StatusOk().
					Done(),
				ShouldBeNil,
			)
		})

		Convey("Fail on invalid Authorization header's structure", func() {
			// Without Bearer
			So(
				api.Bl.Get("/").
					SetHeader("Authorization", "wew lad").
					Expect(t).
					BodyMatchString("Invalid Authorization header.").
					StatusError().
					Done(),
				ShouldBeNil,
			)

			// Without token
			So(
				api.Bl.Get("/").
					SetHeader("Authorization", "Bearer").
					Expect(t).
					BodyMatchString("Invalid Authorization header.").
					StatusError().
					Done(),
				ShouldBeNil,
			)

			// Invalid token length
			So(
				api.Bl.Get("/").
					SetHeader("Authorization", "Bearer wew").
					Expect(t).
					BodyMatchString("Invalid Authorization header.").
					StatusError().
					Done(),
				ShouldBeNil,
			)

			// Non-alphanumeric token
			So(
				api.Bl.Get("/").
					SetHeader("Authorization", "Bearer 0123456789_abcdefghi").
					Expect(t).
					BodyMatchString("Invalid Authorization header.").
					StatusError().
					Done(),
				ShouldBeNil,
			)
		})

		Convey("Fail on non-existing token", func() {
			So(
				api.Bl.Get("/").
					SetHeader("Authorization", "Bearer 00000000000000000000").
					Expect(t).
					BodyMatchString("Token not found.").
					StatusError().
					Done(),
				ShouldBeNil,
			)
		})

		Convey("Fail on a database inconsistency", func() {
			token := &models.Token{
				DateCreated:  time.Now(),
				DateModified: time.Now(),
				Owner:        10,
				ExpiryDate:   time.Now().Add(time.Hour),
				Type:         "auth",
				Application:  10,
				Perms:        perms.Nodes{},
			}
			if _, err := database.MustDataset(api.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
				So(err, ShouldBeNil)
			}

			// The account is missing, which is the inconsistency.
			So(
				api.Bl.Get("/").
					SetHeader("Authorization", "Bearer "+token.ID).
					Expect(t).
					BodyMatchString("Database inconsistency.").
					StatusError().
					Done(),
				ShouldBeNil,
			)
		})

		Convey("Fail on a blocked account", func() {
			account := &models.Account{
				DateCreated:      time.Now(),
				DateModified:     time.Now(),
				Type:             "standard",
				MainAddress:      uniuri.New() + "@oakmail.io",
				Identity:         "Test Testing",
				Password:         "",
				Subscription:     "free",
				Blocked:          true,
				AltEmail:         uniuri.NewLen(uniuri.UUIDLen) + "@gmail.com",
				AltEmailVerified: time.Now().Truncate(time.Hour),
			}
			if _, err := database.MustDataset(api.GQ.From("accounts").Insert(account).ResultingRow()).Select("id").ScanVal(&account.ID); err != nil {
				So(err, ShouldBeNil)
			}

			token := &models.Token{
				DateCreated:  time.Now(),
				DateModified: time.Now(),
				Owner:        account.ID,
				ExpiryDate:   time.Now().Add(time.Hour),
				Type:         "auth",
				Application:  10,
				Perms:        perms.Nodes{},
			}
			if _, err := database.MustDataset(api.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
				So(err, ShouldBeNil)
			}

			So(
				api.Bl.Get("/").
					SetHeader("Authorization", "Bearer "+token.ID).
					Expect(t).
					BodyMatchString(`Account is blocked.`).
					StatusError().
					Done(),
				ShouldBeNil,
			)
		})

		Convey("Given a working account", func() {
			account := &models.Account{
				DateCreated:      time.Now(),
				DateModified:     time.Now(),
				Type:             "standard",
				MainAddress:      uniuri.New() + "@oakmail.io",
				Identity:         "Test Testing",
				Password:         "",
				Subscription:     "free",
				Blocked:          false,
				AltEmail:         uniuri.NewLen(uniuri.UUIDLen) + "@gmail.com",
				AltEmailVerified: time.Now().Truncate(time.Hour),
			}
			if _, err := database.MustDataset(api.GQ.From("accounts").Insert(account).ResultingRow()).Select("id").ScanVal(&account.ID); err != nil {
				So(err, ShouldBeNil)
			}

			Convey("Fail on a token that is not auth", func() {
				token := &models.Token{
					DateCreated:  time.Now(),
					DateModified: time.Now(),
					Owner:        account.ID,
					ExpiryDate:   time.Now().Add(time.Hour),
					Type:         "code",
					Application:  10,
					Perms:        perms.Nodes{},
				}
				if _, err := database.MustDataset(api.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
					So(err, ShouldBeNil)
				}

				So(
					api.Bl.Get("/").
						SetHeader("Authorization", "Bearer "+token.ID).
						Expect(t).
						BodyMatchString(`Invalid token type, must be \\"auth\\".`).
						StatusError().
						Done(),
					ShouldBeNil,
				)
			})

			Convey("Fail on an expired token", func() {
				token := &models.Token{
					DateCreated:  time.Now(),
					DateModified: time.Now(),
					Owner:        account.ID,
					ExpiryDate:   time.Now().Truncate(time.Hour),
					Type:         "auth",
					Application:  10,
					Perms:        perms.Nodes{},
				}
				if _, err := database.MustDataset(api.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
					So(err, ShouldBeNil)
				}

				So(
					api.Bl.Get("/").
						SetHeader("Authorization", "Bearer "+token.ID).
						Expect(t).
						BodyMatchString(`Authentication token expired.`).
						StatusError().
						Done(),
					ShouldBeNil,
				)
			})

			Convey("Succeed given a proper token", func() {
				token := &models.Token{
					DateCreated:  time.Now(),
					DateModified: time.Now(),
					Owner:        account.ID,
					ExpiryDate:   time.Now().Add(time.Hour),
					Type:         "auth",
					Application:  10,
					Perms:        perms.Nodes{},
				}
				if _, err := database.MustDataset(api.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
					So(err, ShouldBeNil)
				}

				So(
					api.Bl.Get("/").
						SetHeader("Authorization", "Bearer "+token.ID).
						Expect(t).
						BodyMatchString(`trtl backend api v0`).
						StatusOk().
						Done(),
					ShouldBeNil,
				)
			})
		})
	})
}
