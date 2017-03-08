package test

import (
	"fmt"
	"io/ioutil"
	stdLogger "log"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oakmail/goqu"
	"github.com/oakmail/logrus/hooks/test"
	"gopkg.in/h2non/baloo.v1"

	"github.com/oakmail/backend/pkg/api"
	"github.com/oakmail/backend/pkg/api/base"
	"github.com/oakmail/backend/pkg/config"
	"github.com/oakmail/backend/pkg/database"
	"github.com/oakmail/backend/pkg/filesystem"
	"github.com/oakmail/backend/pkg/queue"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// API is a superset of the basic API with additional testing stuff
type API struct {
	*base.API

	Bl *baloo.Client
	TS *httptest.Server

	Config  *config.Config
	LogHook *test.Hook

	DatabaseName   string
	FilesystemName string
}

// InitAPI initializes a new test API given env params
func InitAPI() *API {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	filedir := path.Dir(filename)

	cfg := &config.Config{
		LogLevel: "debug",
		API: config.APIConfig{
			Enabled: true,
			Address: "0.0.0.0:2137",
		},
		Mailer: config.MailerConfig{
			Enabled: true,
			Address: "0.0.0.0:2138",
		},
		Worker: config.WorkerConfig{
			Enabled: true,
		},
	}

	log, hook := test.NewNullLogger()

	var (
		gq     *goqu.Database
		db     *sqlx.DB
		dbName string
	)

	dbs := os.Getenv("DATABASE")
	if dbs == "sqlite" || dbs == "" {
		schema, err := ioutil.ReadFile(
			filepath.Join(filedir, "../../../schema/sqlite.sql"),
		)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		dbName, err = ioutil.TempDir("", "wapi")
		if err != nil {
			fmt.Println(err)
			return nil
		}

		dbPath := filepath.Join(dbName, "sqlite.db")

		var dbt = config.SQLite
		cfg.Database = &dbt
		cfg.SQLite = config.SQLiteConfig{
			ConnectionString: dbPath,
		}

		db, gq, err = database.NewSQLite(cfg.SQLite)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if _, err := db.Exec(string(schema)); err != nil {
			fmt.Println(err)
			return nil
		}
	} else if dbs == "postgres" {
		schema, err := ioutil.ReadFile(
			filepath.Join(filedir, "../../../schema/postgres.sql"),
		)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		uri := os.Getenv("POSTGRES_URI")
		dbName = "wapi_" + uniuri.New()

		var dbt = config.Postgres
		cfg.Database = &dbt
		cfg.Postgres = config.PostgresConfig{
			ConnectionString: uri + "/" + dbName,
		}

		// First do some preparation
		sdb, err := sqlx.Connect("postgres", uri)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if _, err := sdb.Exec("CREATE DATABASE " + dbName); err != nil {
			fmt.Println(err)
			return nil
		}
		if _, err := sdb.Exec("USE " + dbName + "; " + string(schema)); err != nil {
			fmt.Println(err)
			return nil
		}
		if err := sdb.Close(); err != nil {
			fmt.Println(err)
			return nil
		}

		// Then start the session
		db, gq, err = database.NewPostgres(cfg.Postgres)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	var (
		fs     filesystem.Filesystem
		fss    = os.Getenv("FILESYSTEM")
		fsName string
	)
	if fss == "flat" || fss == "" {
		var err error
		fsName, err = ioutil.TempDir("", "wapi")
		if err != nil {
			fmt.Println(err)
			return nil
		}

		var fst = config.Flat
		cfg.Filesystem = &fst
		cfg.Flat = config.FlatConfig{
			Path: fsName,
		}

		fs, err = filesystem.NewFlat(cfg.Flat)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	} else if fss == "seaweed" {
		wurl := os.Getenv("SEAWEED_URL")

		var fst = config.Seaweed
		cfg.Filesystem = &fst
		cfg.Seaweed = config.SeaweedConfig{
			MasterURL: wurl,
		}

		var err error
		fs, err = filesystem.NewSeaweed(cfg.Seaweed)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	var (
		qu  queue.Queue
		qus = os.Getenv("QUEUE")
	)
	if qus == "memory" || qus == "" {
		var qut = config.Memory
		cfg.Queue = &qut

		var err error
		qu, err = queue.NewMemory()
		if err != nil {
			fmt.Println(err)
			return nil
		}
	} else if qus == "nsq" {
		var qut = config.NSQ
		cfg.Queue = &qut
		cfg.NSQ = config.NSQConfig{
			NSQdAddresses:    strings.Split(os.Getenv("NSQD_ADDRESSES"), ","),
			LookupdAddresses: strings.Split(os.Getenv("LOOKUPD_ADDRESSES"), ","),
		}

		var err error
		qu, err = queue.NewNSQ(cfg.NSQ)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	if os.Getenv("DEBUG") == "true" {
		// Turn on logging in goqu
		gq.Logger(
			stdLogger.New(os.Stderr, "", stdLogger.LstdFlags),
		)

		log.Out = os.Stderr
	}

	api := api.NewAPI(cfg, log, db, gq, fs, qu)

	ts := httptest.NewServer(api.Gin)
	bl := baloo.New(ts.URL)

	if os.Getenv("DEBUG") == "true" {
		// Response tracing
		bl = bl.Use(BalooResponse)
	}

	return &API{
		API: api,

		Bl: bl,
		TS: ts,

		Config:  cfg,
		LogHook: hook,

		DatabaseName:   dbName,
		FilesystemName: fsName,
	}
}

// Cleanup removes all temporary stuff after tests
func (t *API) Cleanup() error {
	// Stop the server
	t.TS.Close()

	if *t.Config.Database == config.SQLite {
		if err := t.DB.Close(); err != nil {
			fmt.Println(err)
			return err
		}

		if err := t.GQ.Db.Close(); err != nil {
			fmt.Println(err)
			return err
		}

		if err := os.RemoveAll(t.DatabaseName); err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	} else if *t.Config.Database == config.Postgres {
		if _, err := t.DB.Exec("DROP DATABASE " + t.DatabaseName); err != nil {
			fmt.Println(err)
			return err
		}

		if err := t.DB.Close(); err != nil {
			fmt.Println(err)
			return err
		}

		if err := t.GQ.Db.Close(); err != nil {
			fmt.Println(err)
			return err
		}
	}

	if *t.Config.Filesystem == config.Flat {
		if err := os.RemoveAll(t.FilesystemName); err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
