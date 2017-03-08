package database

import (
	"database/sql"
	"sync"
	"time"

	"github.com/dchest/uniuri"
	"github.com/jmoiron/sqlx"
	sqlite "github.com/mattn/go-sqlite3"
	"github.com/oakmail/goqu"
	_ "github.com/oakmail/goqu/adapters/sqlite3" // sqlite adapter for goqu

	"github.com/oakmail/backend/pkg/config"
)

var (
	epoch uint64 = 1487595908000
	index uint64
	mutex sync.Mutex
)

func nextID() uint64 {
	mutex.Lock()
	defer mutex.Unlock()

	timestamp := uint64(time.Now().UTC().UnixNano()) / uint64(time.Millisecond)

	// First 41 bits = timestamp with custom epoch
	result := (timestamp - epoch) << 23

	// Then the next 31 bits are the sharding id, ignore for SQLite
	result = result | ((0 % 8192) << 10)

	// Then the last 10 bits are the counter
	index = (index + 1) % 1024
	result = result | index

	return result
}

func uuid() string {
	return uniuri.NewLen(uniuri.UUIDLen)
}

func init() {
	sql.Register("sqlite3_custom", &sqlite.SQLiteDriver{
		ConnectHook: func(conn *sqlite.SQLiteConn) error {
			if err := conn.RegisterFunc("next_id", nextID, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("uuid", uuid, true); err != nil {
				return err
			}
			return nil
		},
	})
}

// NewSQLite sets up the database connection and returns a query generator
func NewSQLite(cfg config.SQLiteConfig) (*sqlx.DB, *goqu.Database, error) {
	db, err := sqlx.Connect("sqlite3_custom", cfg.ConnectionString)
	if err != nil {
		return nil, nil, err
	}

	gq := goqu.New("sqlite3", db.DB)

	return db, gq, nil
}
