package index_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/oakmail/backend/pkg/api/test"
)

func TestIndex(t *testing.T) {
	api := test.InitAPI()
	defer api.Cleanup()

	Convey("Index method should work correctly", t, func() {
		So(api.Bl.Get("/").
			Expect(t).
			StatusOk().
			BodyMatchString("trtl backend api v0").
			Done(), ShouldBeNil)
	})
}
