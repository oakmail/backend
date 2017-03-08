package base

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oakmail/goqu"
	"github.com/oakmail/logrus"

	"github.com/oakmail/backend/pkg/config"
	"github.com/oakmail/backend/pkg/filesystem"
	"github.com/oakmail/backend/pkg/queue"
)

// API contains all the dependencies of the API server
type API struct {
	Config     *config.Config
	Log        *logrus.Logger
	DB         *sqlx.DB
	Gin        *gin.Engine
	GQ         *goqu.Database
	Filesystem filesystem.Filesystem
	Queue      queue.Queue
}

// Start binds the API and starts listening.
func (a *API) Start() error {
	return a.Gin.Run(a.Config.API.Address)
}
