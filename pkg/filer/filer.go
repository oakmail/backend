package filer

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/jmoiron/sqlx"
	"github.com/oakmail/goqu"
	"github.com/oakmail/logrus"

	"github.com/oakmail/backend/pkg/config"
	"github.com/oakmail/backend/pkg/filesystem"
	"github.com/oakmail/backend/pkg/queue"
)

type Filer struct {
	Config     *config.Config
	Log        *logrus.Logger
	DB         *sqlx.DB
	Gin        *gin.Engine
	GQ         *goqu.Database
	Filesystem filesystem.Filesystem
	Queue      queue.Queue
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func NewFiler(
	cfg *config.Config,
	log *logrus.Logger,
	db *sqlx.DB,
	gq *goqu.Database,
	fs filesystem.Filesystem,
	qu queue.Queue,
) *Filer {
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

	f := &Filer{
		Config:     cfg,
		Log:        log,
		DB:         db,
		Gin:        r,
		GQ:         gq,
		Filesystem: fs,
		Queue:      qu,
	}

	r.POST("/", f.Upload)
	r.GET("/:id", f.Fetch)

	return f
}

// Start binds the Filer API and starts listening.
func (f *Filer) Start() error {
	return a.Gin.Run(a.Config.Filer.Address)
}
