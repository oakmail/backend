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
	"github.com/oakmail/backend/pkg/api/resources"
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

	/*
		addresses := addresses.Impl{API: a}
		r.POST("/addresses", m.RequiresAuth, addresses.Create)
		r.GET("/addresses/:id", m.RequiresAuth, addresses.Get)
		r.GET("/addresses", m.RequiresAuth, addresses.List)
		r.GET("/accounts/:id/addresses", m.RequiresAuth, addresses.ListByAccount)
		r.PUT("/addresses/:id", m.RequiresAuth, addresses.Update)
		r.DELETE("/addresses/:id", m.RequiresAuth, addresses.Delete)
	*/

	applications := applications.Impl{API: a}
	r.POST("/applications", m.RequiresAuth, applications.Create)
	r.GET("/applications", m.RequiresAuth, applications.List)
	r.GET("/applications/:id", applications.Get)
	r.PUT("/applications/:id", m.RequiresAuth, applications.Update)
	r.DELETE("/applications/:id", m.RequiresAuth, applications.Delete)

	index := index.Impl{API: a}
	r.GET("", index.Index)

	/*
		keys := keys.Impl{API: a}
		r.POST("/keys", m.RequiresAuth, keys.Create)
		r.GET("/keys/:id", keys.Get)
		r.GET("/keys", keys.List)
		r.PUT("/keys/:id", m.RequiresAuth, keys.Update)
		r.DELETE("/keys/:id", m.RequiresAuth, keys.Delete)
	*/

	oauth := oauth.Impl{API: a}
	r.POST("/oauth", oauth.OAuth)

	resources := resources.Impl{API: a}
	r.POST("/resources", m.RequiresAuth, resources.Create)
	r.GET("/resources/:id", m.RequiresAuth, resources.Get)
	//r.GET("/resources", m.RequiresAuth, resources.List)
	//r.GET("/accounts/:id/resources", m.RequiresAuth, resources.ListByAccount)
	r.PUT("/resources/:id", m.RequiresAuth, resources.Update)
	r.DELETE("/resources/:id", m.RequiresAuth, resources.Delete)

	/*
		tokens := tokens.Impl{API: a}
		r.POST("/tokens", m.RequiresAuth, tokens.Create)
		r.GET("/tokens/:id", m.RequiresAuth, tokens.Get)
		r.GET("/tokens", m.RequiresAuth, tokens.List)
		r.GET("/accounts/:id/tokens", m.RequiresAuth, tokens.ListByAccount)
		r.PUT("/tokens/:id", m.RequiresAuth, tokens.Update)
		r.DELETE("/tokens/:id", m.RequiresAuth, tokens.Delete)
	*/

	return a
}
