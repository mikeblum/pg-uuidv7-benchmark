package main

import (
	"context"
	"os"

	"log/slog"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mikeblum/pg-uuidv7/internal/db"
)

const (
	attrError   = "error"
	attrVersion = "version"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var idv4 uuid.UUID
	var err error
	if idv4, err = uuid.NewV4(); err != nil {
		logger.Error("Error generating V4 UUID", attrError, err)
	}
	logger.Info(idv4.String(), attrVersion, idv4.Version())
	var idv7 uuid.UUID
	if idv7, err = uuid.NewV7(); err != nil {
		logger.Error("Error generating V7 UUID", attrError, err)
	}
	logger.Info(idv7.String(), attrVersion, idv7.Version())

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("Error connecting to Postgres 🐘", attrError, err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	q := db.New(conn)
	if err = q.GenerateUUIDv4(context.Background()); err != nil {
		logger.Error("Error generating UUIDv4", attrError, err)
		os.Exit(1)
	}
}
