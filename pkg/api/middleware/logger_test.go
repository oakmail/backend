package middleware_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oakmail/logrus/hooks/test"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/appleboy/gofight.v1"

	"github.com/oakmail/backend/pkg/api/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestLogger(t *testing.T) {
	Convey("Logger should work correctly", t, func() {
		r := gin.New()
		f := gofight.New()

		log, hook := test.NewNullLogger()

		r.Use(middleware.Logger(log, time.RFC3339, true))

		r.GET("/succeed", func(c *gin.Context) {
			c.String(http.StatusOK, "hello there %s", "folks")
		})

		r.GET("/fail", func(c *gin.Context) {
			c.AbortWithError(http.StatusInternalServerError, errors.New("wot"))
		})

		f.GET("/succeed").Run(r, func(r gofight.HTTPResponse, req gofight.HTTPRequest) {
			So(r.Code, ShouldEqual, 200)
			So(r.Body.String(), ShouldContainSubstring, "hello there folks")
			So(hook.LastEntry().Message, ShouldContainSubstring, "Completed a request")
		})

		f.GET("/fail").Run(r, func(r gofight.HTTPResponse, req gofight.HTTPRequest) {
			So(r.Code, ShouldEqual, 500)
			So(hook.LastEntry().Message, ShouldContainSubstring, "Request failed")
		})
	})
}
