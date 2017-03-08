package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // pq adapter for sql
	"github.com/oakmail/goqu"
	_ "github.com/oakmail/goqu/adapters/postgres" // pq adapter for goqu

	"github.com/oakmail/backend/pkg/config"
)

// NewPostgres connects to the database and returns a query generator
func NewPostgres(cfg config.PostgresConfig) (*sqlx.DB, *goqu.Database, error) {
	db, err := sqlx.Connect("postgres", cfg.ConnectionString)
	if err != nil {
		return nil, nil, err
	}

	gq := goqu.New("postgres", db.DB)

	return db, gq, nil
}
