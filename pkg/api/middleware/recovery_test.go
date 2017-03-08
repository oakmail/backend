package middleware_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/appleboy/gofight.v1"

	"github.com/oakmail/backend/pkg/api/middleware"
	"github.com/oakmail/backend/pkg/api/test"
)

// TestPanicInHandler assert that panic has been recovered.
func TestRecovery(t *testing.T) {
	api := test.InitAPI()
	defer api.Cleanup()

	mw := middleware.Impl{API: api.API}

	Convey("Panics should be recovered", t, func() {
		Convey("When they are simple panics", func() {
			router := gin.New()
			router.Use(mw.Recovery)

			fight := gofight.New()

			router.GET("/recovery", func(_ *gin.Context) {
				panic("Oupps, Houston, we have a problem")
			})

			fight.GET("/recovery").Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				So(r.Code, ShouldEqual, 500)
			})
		})

		Convey("When they are database errors", func() {
			router := gin.New()
			router.Use(mw.Recovery)

			fight := gofight.New()

			router.GET("/recovery", func(_ *gin.Context) {
				panic(sqlite3.ErrTooBig)
			})

			fight.GET("/recovery").Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				So(r.Code, ShouldEqual, 500)
				So(r.Body.String(), ShouldContainSubstring, "Database error.")
			})
		})
	})
}
