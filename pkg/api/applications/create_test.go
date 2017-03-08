package applications_test

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

func TestCreate(t *testing.T) {
	api := test.InitAPI()
	defer api.Cleanup()

	Convey("Given an account without privileges", t, func() {
		account := &models.Account{
			DateCreated:      time.Now(),
			DateModified:     time.Now(),
			Type:             "admin",
			MainAddress:      uniuri.NewLen(uniuri.UUIDLen) + "@oakmail.io",
			Identity:         "Test Testing",
			Password:         "",
			Subscription:     "sth",
			Blocked:          false,
			AltEmail:         uniuri.NewLen(uniuri.UUIDLen) + "@gmail.com",
			AltEmailVerified: time.Now().Truncate(time.Hour),
		}
		if _, err := database.MustDataset(api.GQ.From("accounts").Insert(account).ResultingRow()).Select("id").ScanVal(&account.ID); err != nil {
			So(err, ShouldBeNil)
		}

		Convey("Fail on insufficient privileges when trying to create an app for other user", func() {
			token := &models.Token{
				DateCreated:  time.Now(),
				DateModified: time.Now(),
				Owner:        account.ID,
				ExpiryDate:   time.Now().Add(time.Hour),
				Type:         "auth",
				Application:  10,
				Perms:        perms.MustParseNodes([]byte("some.random")),
			}
			if _, err := database.MustDataset(api.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
				So(err, ShouldBeNil)
			}

			So(
				api.Bl.Post("/applications").
					SetHeader("Authorization", "Bearer "+token.ID).
					JSON(map[string]interface{}{
						"owner": account.ID,
					}).
					Expect(t).
					StatusError().
					BodyMatchString("Insufficient token permissions.").
					Done(),
				ShouldBeNil,
			)
		})

		Convey("Fail on insufficient privileges while creating an app for yourself", func() {
			token := &models.Token{
				DateCreated:  time.Now(),
				DateModified: time.Now(),
				Owner:        account.ID,
				ExpiryDate:   time.Now().Add(time.Hour),
				Type:         "auth",
				Application:  10,
				Perms:        perms.MustParseNodes([]byte("some.random")),
			}
			if _, err := database.MustDataset(api.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
				So(err, ShouldBeNil)
			}

			So(
				api.Bl.Post("/applications").
					SetHeader("Authorization", "Bearer "+token.ID).
					JSON(map[string]interface{}{
						"owner": account.ID,
					}).
					Expect(t).
					StatusError().
					BodyMatchString("Insufficient token permissions.").
					Done(),
				ShouldBeNil,
			)
		})
	})

	Convey("Given an account with privileges", t, func() {
		account := &models.Account{
			DateCreated:      time.Now(),
			DateModified:     time.Now(),
			Type:             "admin",
			MainAddress:      uniuri.NewLen(uniuri.UUIDLen) + "@oakmail.io",
			Identity:         "Test Testing",
			Password:         "",
			Subscription:     "sth",
			Blocked:          false,
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
			Application:  1,
			Perms:        perms.MustParseNodes([]byte("*")),
		}
		if _, err := database.MustDataset(api.GQ.From("tokens").Insert(token).ResultingRow()).Select("id").ScanVal(&token.ID); err != nil {
			So(err, ShouldBeNil)
		}

		Convey("Fail on invalid JSON input", func() {
			So(
				api.Bl.Post("/applications").
					SetHeader("Authorization", "Bearer "+token.ID).
					Type("json").
					BodyString("`").
					Expect(t).
					StatusError().
					BodyMatchString("Invalid JSON input.").
					Done(),
				ShouldBeNil,
			)
		})

		Convey("Fail on invalid fields in the JSON body", func() {
			So(
				api.Bl.Post("/applications").
					SetHeader("Authorization", "Bearer "+token.ID).
					JSON(map[string]interface{}{
						"callback":    "hehe xd",
						"name":        "",
						"email":       "what is this",
						"home_page":   "hehe xd",
						"description": "",
						"owner":       account.ID,
					}).
					Expect(t).
					StatusError().
					BodyMatchString("Invalid callback format.").
					BodyMatchString("Application name is invalid.").
					BodyMatchString("Invalid e-mail format.").
					BodyMatchString("Invalid home page format.").
					Done(),
				ShouldBeNil,
			)
		})

		Convey("Succeed on application insertion", func() {
			So(
				api.Bl.Post("/applications").
					SetHeader("Authorization", "Bearer "+token.ID).
					JSON(map[string]interface{}{
						"callback":    "https://oakmail.io",
						"name":        "oakmail.io",
						"email":       "hello@oakmail.io",
						"home_page":   "https://oakmail.io",
						"description": "What is this",
						"owner":       account.ID,
					}).
					Expect(t).
					StatusOk().
					BodyMatchString("oakmail.io").
					Done(),
				ShouldBeNil,
			)
		})
	})
}
