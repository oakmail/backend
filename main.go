package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/koding/multiconfig"
	"github.com/oakmail/goqu"
	"github.com/oakmail/logrus"

	"github.com/oakmail/backend/pkg/api"
	"github.com/oakmail/backend/pkg/config"
	"github.com/oakmail/backend/pkg/database"
	"github.com/oakmail/backend/pkg/filesystem"
	"github.com/oakmail/backend/pkg/queue"
)

func main() {
	m := multiconfig.NewWithPath(os.Getenv("config"))
	cfg := &config.Config{}
	m.MustLoad(cfg)

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	logger := logrus.StandardLogger()
	logger.Level = logLevel

	logger.WithFields(logrus.Fields{
		"module": "init",
	}).Info("Starting up the application")

	// Initialize the database
	var (
		db *sqlx.DB
		gq *goqu.Database
	)
	switch *cfg.Database {
	case config.Postgres:
		db, gq, err = database.NewPostgres(cfg.Postgres)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"module": "init",
				"error":  err.Error(),
				"cstr":   cfg.Postgres.ConnectionString,
			}).Fatal("Unable to connect to the Postgres server")
			return
		}

		logger.WithFields(logrus.Fields{
			"module": "init",
			"cstr":   cfg.Postgres.ConnectionString,
		}).Info("Connected to a Postgres server")
	case config.SQLite:
		db, gq, err = database.NewSQLite(cfg.SQLite)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"module": "init",
				"error":  err.Error(),
				"cstr":   cfg.SQLite.ConnectionString,
			}).Fatal("Unable to load the SQLite database")
			return
		}

		logger.WithFields(logrus.Fields{
			"module": "init",
			"cstr":   cfg.SQLite.ConnectionString,
		}).Info("Loaded a SQLite database")
	}

	// Initialize the filesystem
	var fs filesystem.Filesystem
	switch *cfg.Filesystem {
	case config.Flat:
		fs, err = filesystem.NewFlat(cfg.Flat)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"module": "init",
				"error":  err.Error(),
				"path":   cfg.Flat.Path,
			}).Fatal("Unable to initialize the flat filesystem")
			return
		}

		logger.WithFields(logrus.Fields{
			"module": "init",
			"path":   cfg.Flat.Path,
		}).Info("Initialized the flat filesystem")
	case config.Seaweed:
		fs, err = filesystem.NewSeaweed(cfg.Seaweed)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"module": "init",
				"error":  err.Error(),
				"url":    cfg.Seaweed.MasterURL,
			}).Fatal("Unable to initialize the seaweed filesystem")
			return
		}

		logger.WithFields(logrus.Fields{
			"module": "init",
			"url":    cfg.Seaweed.MasterURL,
		}).Info("Initialized the Seaweed filesystem")
	}

	// Initialize the queue
	var qu queue.Queue
	switch *cfg.Queue {
	case config.Memory:
		qu, err = queue.NewMemory()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"module": "init",
				"error":  err.Error(),
			}).Fatal("Unable to initialize the memory queue")
			return
		}

		logger.WithFields(logrus.Fields{
			"module": "init",
		}).Info("Initialized the memory queue")
	case config.NSQ:
		qu, err = queue.NewNSQ(cfg.NSQ)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"module":   "init",
				"error":    err.Error(),
				"nsqds":    cfg.NSQ.NSQdAddresses,
				"lookupds": cfg.NSQ.LookupdAddresses,
			}).Fatal("Unable to initialize the NSQ connection")
			return
		}

		logger.WithFields(logrus.Fields{
			"module":   "init",
			"nsqds":    cfg.NSQ.NSQdAddresses,
			"lookupds": cfg.NSQ.LookupdAddresses,
		}).Info("Initialized the NSQ queue")
	}

	api := api.NewAPI(cfg, logger, db, gq, fs, qu)
	go func() {
		logger.WithFields(logrus.Fields{
			"module": "init",
			"bind":   cfg.API.Address,
		}).Info("Starting the API server")

		if err := api.Start(); err != nil {
			logger.WithFields(logrus.Fields{
				"module": "init",
				"error":  err.Error(),
			}).Fatal("Failed to start the API server")
		}
	}()

	x := make(chan struct{})
	<-x
}
