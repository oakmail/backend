package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/jmoiron/sqlx"
	"github.com/oakmail/goqu"
	"github.com/oakmail/logrus"

	"github.com/oakmail/backend/pkg/api/accounts"
	"github.com/oakmail/backend/pkg/api/applications"
	"github.com/oakmail/backend/pkg/api/base"
	"github.com/oakmail/backend/pkg/api/index"
	"github.com/oakmail/backend/pkg/api/middleware"
	"github.com/oakmail/backend/pkg/api/oauth"
	"github.com/oakmail/backend/pkg/config"
	"github.com/oakmail/backend/pkg/filesystem"
	"github.com/oakmail/backend/pkg/queue"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// NewAPI sets up a new API module.
func NewAPI(
	cfg *config.Config,
	log *logrus.Logger,
	db *sqlx.DB,
	gq *goqu.Database,
	fs filesystem.Filesystem,
	qu queue.Queue,
) *base.API {
	r := gin.New()

	r.RedirectTrailingSlash = false
	r.Use(middleware.Logger(log, time.RFC3339Nano, true))
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	a := &base.API{
		Config:     cfg,
		Log:        log,
		DB:         db,
		Gin:        r,
		GQ:         gq,
		Filesystem: fs,
		Queue:      qu,
	}

	m := middleware.Impl{API: a}

	r.Use(m.Recovery)
	r.Use(m.UsesAuth)

	accounts := accounts.Impl{API: a}
	r.POST("/accounts", accounts.Create)
	r.DELETE("/accounts/:id", m.RequiresAuth, accounts.Delete)
	r.GET("/accounts/:id", m.RequiresAuth, accounts.Get)

	applications := applications.Impl{API: a}
	r.POST("/applications", m.RequiresAuth, applications.Create)
	r.GET("/applications", m.RequiresAuth, applications.List)
	r.GET("/applications/:id", applications.Get)
	r.PUT("/applications/:id", m.RequiresAuth, applications.Update)
	r.DELETE("/applications/:id", m.RequiresAuth, applications.Delete)

	index := index.Impl{API: a}
	r.GET("", index.Index)

	oauth := oauth.Impl{API: a}
	r.POST("/oauth", oauth.OAuth)

	return a
}
