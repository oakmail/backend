package middleware_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/oakmail/backend/pkg/api/test"
)

func TestRequiresAuth(t *testing.T) {
	api := test.InitAPI()
	defer api.Cleanup()

	Convey("RequiresAuth middleware should", t, func() {
		Convey("Fail with no Authorization passed", func() {
			So(
				api.Bl.Post("/applications").
					Expect(t).
					BodyMatchString("Invalid Authorization header.").
					StatusError().
					Done(),
				ShouldBeNil,
			)
		})
	})
}
