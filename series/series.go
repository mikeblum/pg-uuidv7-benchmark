package series

import (
	"context"
	"log/slog"
	"os"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	dbms "github.com/mikeblum/pg-uuidv7/internal/db"
)

const (
	attrError = "error"
)

type Series struct {
	conn       *pgx.Conn
	queries    *dbms.Queries
	timestamps []pgtype.Timestamp
}

func NewWithConn(conn *pgx.Conn) (*Series, error) {
	s := &Series{
		conn: conn,
	}
	s.queries = dbms.New(s.conn)
	return s, s.generateSeries()
}

func (s *Series) generateSeries() error {
	var err error
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if s.timestamps, err = s.queries.GenerateSeries(context.Background()); err != nil {
		logger.Error("Error generating series", attrError, err)
		return err
	}
	return err
}

func (s *Series) InsertUUIDv4Bulk() error {
	var err error
	batch := make([]dbms.InsertUUIDv4BulkParams, len(s.timestamps))
	for i, ts := range s.timestamps {
		var idv4 uuid.UUID
		if idv4, err = uuid.NewV4(); err != nil {
			return err
		}
		batch[i] = dbms.InsertUUIDv4BulkParams{
			ID: pgtype.UUID{
				Bytes: [16]byte(idv4.Bytes()),
				Valid: true,
			},
			Created: ts,
		}
	}
	results := s.queries.InsertUUIDv4Bulk(context.Background(), batch)
	defer results.Close()
	results.Query(func(i int, row []dbms.InsertUUIDv4BulkRow, rowErr error) {
		err = rowErr
		if err != nil {
			return
		}
	})
	return err
}

func (s *Series) InsertUUIDv7Bulk() error {
	var err error
	batch := make([]dbms.InsertUUIDv7BulkParams, len(s.timestamps))
	for i, ts := range s.timestamps {
		var idv7 uuid.UUID
		if idv7, err = uuid.NewV7(); err != nil {
			return err
		}
		batch[i] = dbms.InsertUUIDv7BulkParams{
			ID: pgtype.UUID{
				Bytes: [16]byte(idv7.Bytes()),
				Valid: true,
			},
			Created: ts,
		}
	}
	results := s.queries.InsertUUIDv7Bulk(context.Background(), batch)
	defer results.Close()
	results.Query(func(i int, row []dbms.InsertUUIDv7BulkRow, rowErr error) {
		err = rowErr
		if err != nil {
			return
		}
	})
	return err
}
