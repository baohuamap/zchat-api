package gorm

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBAdapter interface represent adapter connect to DB
type DBAdapter interface {
	Connect(dsn string) error
	Close()
	Begin() DBAdapter
	Rollback()
	Commit()
	IsCommitted() bool
	Gormer() *gorm.DB
	DB() *sql.DB
}

type adapter struct {
	gormer      *gorm.DB
	isCommitted bool
}

// NewDB returns a new instance of DB.
func NewDB() DBAdapter {
	return &adapter{}
}

func (db *adapter) Connect(dsn string) error {
	gormer, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
	})

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	db.gormer = gormer

	return nil
}

// Gormer returns an instance of gorm.DB.
func (db *adapter) Gormer() *gorm.DB {
	return db.gormer
}

// DB returns an instance of sql.DB.
func (db *adapter) DB() *sql.DB {
	database, _ := db.gormer.DB()
	return database
}

// Close closes DB connection.
func (db *adapter) Close() {
	_ = db.DB().Close()
}

// Begin starts a DB transaction.
func (db *adapter) Begin() DBAdapter {
	tx := db.gormer.Begin()

	return &adapter{
		gormer:      tx,
		isCommitted: false,
	}
}

// Rollback rollbacks useless DB transaction committed.
func (db *adapter) Rollback() {
	if !db.isCommitted {
		db.gormer.Rollback()
	}
}

// Commit commits a DB transaction.
func (db *adapter) Commit() {
	if !db.isCommitted {
		db.gormer.Commit()
		db.isCommitted = true
	}
}

func (db *adapter) IsCommitted() bool {
	return db.isCommitted
}
