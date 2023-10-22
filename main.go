package main

import (
	"context"
	"database/sql"
	"embed"
	"os"

	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"
	dbms "github.com/mikeblum/pg-uuidv7/internal/db"
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

	// UUIDv4
	var v4 *series.Series = &series.Series{
		Query:  dbms.New(conn),
		Logger: logger,
	}
	logger.Info("Generating UUIDv4 series...")
	if err = v4.GenerateSeries(); err != nil {
		logger.Error("Error generating UUIDv4 series", attrError, err)
		os.Exit(1)
	}
	var out chan pgtype.UUID
	if out, err = v4.InsertUUIDv4Bulk(); err != nil {
		logger.Error("Error generating UUIDv4 series", attrError, err)
		os.Exit(1)
	}
	logger.Info("Waiting for UUIDv4s...")
	logger.Info("Fetching UUIDv4s...")
	for id := range out {
		logger.Info("UUIDv4", "id", series.UUIDString(id))
		if err = v4.GetUUIDv4(id); err != nil {
			logger.Error("Error getting UUIDv4", "id", series.UUIDString(id), attrError, err)
		}
	}

	// UUIDv7
	var v7 *series.Series = &series.Series{
		Query:  dbms.New(conn),
		Logger: logger,
	}
	logger.Info("Generating UUIDv7 series...")
	if err = v7.GenerateSeries(); err != nil {
		logger.Error("Error generating series", attrError, err)
		os.Exit(1)
	}
	if out, err = v7.InsertUUIDv7Bulk(); err != nil {
		logger.Error("Error generating UUIDv7 series", attrError, err)
		os.Exit(1)
	}
	logger.Info("Waiting for UUIDv7s...")
	logger.Info("Fetching UUIDv7s...")
	for id := range out {
		logger.Info("UUIDv7", "id", series.UUIDString(id))
		if err = v7.GetUUIDv7(id); err != nil {
			logger.Error("Error getting UUIDv7", "id", series.UUIDString(id), attrError, err)
		}
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
