package main

import (
	"context"
	"database/sql"
	"embed"
	"os"

	"log/slog"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mikeblum/pg-uuidv7/series"
	"github.com/pressly/goose/v3"
)

const (
	attrError = "error"

	cmdPgx            = "pgx"
	envDatabaseURL    = "DATABASE_URL"
	migrationsDir     = "db/migrations"
	migrationsDialect = "postgres"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var err error
	var db *sql.DB
	if db, err = sql.Open(cmdPgx, os.Getenv(envDatabaseURL)); err != nil {
		logger.Error("Error connecting to Postgres 🐘", attrError, err)
		os.Exit(1)
	}
	if err = bootstrap(db); err != nil {
		logger.Error("Error bootstrapping schema for Postgres 🐘", attrError, err)
		os.Exit(1)
	}
	var conn *pgx.Conn
	if conn, err = pgx.Connect(context.Background(), os.Getenv(envDatabaseURL)); err != nil {
		logger.Error("Error connecting to Postgres 🐘", attrError, err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	var s *series.Series
	if s, err = series.NewWithConn(conn); err != nil {
		logger.Error("Error generating series", attrError, err)
		os.Exit(1)
	}
	if err = s.InsertUUIDv4Bulk(); err != nil {
		logger.Error("Error generating UUIDv4 series", attrError, err)
		os.Exit(1)
	}
	if err = s.InsertUUIDv7Bulk(); err != nil {
		logger.Error("Error generating UUIDv7 series", attrError, err)
		os.Exit(1)
	}
}

func bootstrap(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	var err error
	if err = goose.SetDialect(migrationsDialect); err != nil {
		return err
	}
	if err = goose.Down(db, migrationsDir); err != nil {
		return err
	}
	if err = goose.Up(db, migrationsDir); err != nil {
		return err
	}
	return err
}
