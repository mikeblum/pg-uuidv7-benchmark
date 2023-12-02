package main

import (
	"context"
	"database/sql"
	"embed"
	"os"

	"log/slog"

	"github.com/jackc/pgx/v5"
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

	query := dbms.New(conn)

	if err = truncate(query); err != nil {
		logger.Error("Error truncating db", attrError, err)
		os.Exit(1)
	}

	// UUIDv4
	var v4 *series.Series = &series.Series{
		Query:  query,
		Logger: logger,
	}
	logger.Info("Generating UUIDv4 series...")
	if err = v4.GenerateSeries(); err != nil {
		logger.Error("Error generating UUIDv4 series", attrError, err)
		os.Exit(1)
	}
	var out chan series.UUID
	if out, err = v4.InsertUUIDv4Bulk(); err != nil {
		logger.Error("Error generating UUIDv4 series", attrError, err)
		os.Exit(1)
	}
	logger.Info("Vacuuming UUIDv4s...")
	// sqlc doesn't appear to support `VACUUUM`
	var vacuumResult *sql.Rows
	if vacuumResult, err = db.Query("VACUUM ( analyze true, index_cleanup true, verbose true ) uuid_v4;"); err != nil {
		logger.Error("Error vacuuming uuid_v4", attrError, err)
		os.Exit(1)
	}
	for vacuumResult.Next() {
		var analysis *string
		if err = vacuumResult.Scan(&analysis); err != nil {
			logger.Error("Error analyzing vacuum", attrError, err)
			os.Exit(1)
		}
		logger.Info(*analysis)
	}
	logger.Info("Fetching UUIDv4s...")
	for uuid := range out {
		logger.Debug("UUIDv4", "id", series.UUIDString(uuid.ID))
		if uuid.Lookup, err = v4.GetUUIDv4(uuid.ID); err != nil {
			logger.Error("Error getting UUIDv4", "id", series.UUIDString(uuid.ID), attrError, err)
			continue
		}
		if err = v4.MergeUUIDResult(uuid, series.BTREE); err != nil {
			logger.Error("Error merging UUIDv4 result", attrError, err)
			continue
		}
	}

	// UUIDv7
	var v7 *series.Series = &series.Series{
		Query:  query,
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
	logger.Info("Vacuuming UUIDv7s...")
	// sqlc doesn't appear to support `VACUUUM`
	if vacuumResult, err = db.Query("VACUUM ( analyze true, index_cleanup true, verbose true ) uuid_v7;"); err != nil {
		logger.Error("Error vacuuming uuid_v7", attrError, err)
		os.Exit(1)
	}
	for vacuumResult.Next() {
		var analysis *string
		if err = vacuumResult.Scan(&analysis); err != nil {
			logger.Error("Error analyzing vacuum", attrError, err)
			os.Exit(1)
		}
		logger.Info(*analysis)
	}
	logger.Info("Fetching UUIDv7s...")
	for uuid := range out {
		logger.Debug("UUIDv7", "id", series.UUIDString(uuid.ID))
		if uuid.Lookup, err = v7.GetUUIDv7(uuid.ID); err != nil {
			logger.Error("Error getting UUIDv7 BTree", "id", series.UUIDString(uuid.ID), attrError, err)
			continue
		}
		if err = v7.MergeUUIDResult(uuid, series.BTREE); err != nil {
			logger.Error("Error merging UUIDv7 result", attrError, err)
			continue
		}
		logger.Debug("UUIDv7", "id", series.UUIDString(uuid.ID))
		if uuid.Lookup, err = v7.GetUUIDv7BRIN(uuid.ID); err != nil {
			logger.Error("Error getting UUIDv7 BRIN", "id", series.UUIDString(uuid.ID), attrError, err)
			continue
		}
		if err = v7.MergeUUIDResult(uuid, series.BRIN); err != nil {
			logger.Error("Error merging UUIDv7 result", attrError, err)
			continue
		}
	}
}

func bootstrap(db *sql.DB) error {
	goose.SetVerbose(true)
	goose.SetBaseFS(embedMigrations)
	var err error
	if err = goose.SetDialect(migrationsDialect); err != nil {
		return err
	}
	if err = goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return err
}

func truncate(query *dbms.Queries) error {
	var err error
	if err = query.TruncateUUIDv4(context.Background()); err != nil {
		return err
	}
	if err = query.TruncateUUIDv7(context.Background()); err != nil {
		return err
	}
	if err = query.TruncateUUIDResult(context.Background()); err != nil {
		return err
	}
	return err
}
