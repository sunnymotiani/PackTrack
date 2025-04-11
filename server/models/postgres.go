package models

import (
	"database/sql"
	"fmt"
	"io/fs"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.User, cfg.Port, cfg.Password, cfg.Database, cfg.SSLMode)
}

func Open(cfg PostgresConfig) (*sql.DB, error) {
	// Open the database connection
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		return nil, fmt.Errorf("open : %w", err)
	}
	// Set connection pool parameters
	db.SetMaxIdleConns(10)           // Max number of idle connections
	db.SetMaxOpenConns(20)           // Max number of open connections
	db.SetConnMaxLifetime(time.Hour) // Max lifetime of a connection

	return db, nil
}
func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate :%w", err)
	}
	return nil
}

func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	if dir == "" {
		dir = "."
	}
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
