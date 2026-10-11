package projections

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sync/atomic"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// counter is used to build a unique name for each in-memory database.
	counter atomic.Int64

	//go:embed *.sql
	schema embed.FS
)

// NewDB returns an in-memory SQLite database with the required schema.
func NewDB() (*sql.DB, error) {
	n := counter.Add(1)
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:dot%d?mode=memory&cache=shared", n))
	if err != nil {
		return nil, err
	}

	if err := CreateSchema(context.Background(), db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// MustNewDB returns an in-memory SQLite database, panicking on error.
func MustNewDB() *sql.DB {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}
	return db
}

// CreateSchema creates the schema elements required by the projection handlers.
func CreateSchema(ctx context.Context, db *sql.DB) error {
	entries, err := schema.ReadDir(".")
	if err != nil {
		return err
	}

	for _, e := range entries {
		s, err := schema.ReadFile(e.Name())
		if err != nil {
			return err
		}

		if _, err := db.ExecContext(ctx, string(s)); err != nil {
			return err
		}
	}

	return nil
}
