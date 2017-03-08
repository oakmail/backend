package errors_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/appleboy/gofight.v1"

	"github.com/oakmail/backend/pkg/api/errors"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestErrors(t *testing.T) {
	Convey("Abort should properly stop the execution flow", t, func() {
		r := gin.New()
		f := gofight.New()

		r.GET("/", func(c *gin.Context) {
			errors.Abort(c, http.StatusInternalServerError, errors.DatabaseError)
		})

		f.GET("/").Run(r, func(r gofight.HTTPResponse, req gofight.HTTPRequest) {
			So(r.Code, ShouldEqual, 500)
			So(r.Body.String(), ShouldContainSubstring, "Database error.")
		})
	})

	Convey("Error should do what expected", t, func() {
		So(errors.DatabaseError.Error(), ShouldEqual, "Database error.")
	})
}
