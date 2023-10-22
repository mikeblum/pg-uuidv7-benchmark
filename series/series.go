package series

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"
	dbms "github.com/mikeblum/pg-uuidv7/internal/db"
)

const (
	attrError = "error"
)

type Series struct {
	Query      *dbms.Queries
	Logger     *slog.Logger
	timestamps []pgtype.Timestamp
}

func (s *Series) GenerateSeries() error {
	var err error
	if s.timestamps, err = s.Query.GenerateSeries(context.Background()); err != nil {
		s.Logger.Error("Error generating series", attrError, err)
		return err
	}
	return err
}

func (s *Series) InsertUUIDv4Bulk() (chan pgtype.UUID, error) {
	var err error
	batch := make([]dbms.InsertUUIDv4BulkParams, len(s.timestamps))
	out := make(chan pgtype.UUID, len(batch))
	for i, ts := range s.timestamps {
		var idv4 uuid.UUID
		if idv4, err = uuid.NewV4(); err != nil {
			return nil, err
		}
		batch[i] = dbms.InsertUUIDv4BulkParams{
			ID: pgtype.UUID{
				Bytes: [16]byte(idv4.Bytes()),
				Valid: true,
			},
			Created: ts,
		}
	}
	start := time.Now()
	results := s.Query.InsertUUIDv4Bulk(context.Background(), batch)
	results.QueryRow(func(i int, id pgtype.UUID, err error) {
		if err != nil {
			s.Logger.Error("InsertUUIDv4Bulk", "uuid", UUIDString(id), "err", err)
			return
		}
		out <- id
		s.Logger.Info("InsertUUIDv4Bulk", "uuid", UUIDString(id), "version", uuid.V4, "time", time.Since(start))
		start = time.Now()
		if i >= len(batch)-1 {
			close(out)
		}
	})
	return out, err
}

func (s *Series) GetUUIDv4(id pgtype.UUID) error {
	var err error
	start := time.Now()
	if _, err = s.Query.GetUUIDv4(context.Background(), id); err != nil {
		s.Logger.Error("GetUUIDv4", "uuid", UUIDString(id), "err", err)
		return err
	}
	s.Logger.Info("GetUUIDv4", "uuid", UUIDString(id), "version", uuid.V4, "time", time.Since(start))
	return err
}

func (s *Series) InsertUUIDv7Bulk() (chan pgtype.UUID, error) {
	var err error
	batch := make([]dbms.InsertUUIDv7BulkParams, len(s.timestamps))
	out := make(chan pgtype.UUID, len(batch))
	for i, ts := range s.timestamps {
		var idv7 uuid.UUID
		if idv7, err = uuid.NewV7(); err != nil {
			return nil, err
		}
		batch[i] = dbms.InsertUUIDv7BulkParams{
			ID: pgtype.UUID{
				Bytes: [16]byte(idv7.Bytes()),
				Valid: true,
			},
			Created: ts,
		}
	}

	results := s.Query.InsertUUIDv7Bulk(context.Background(), batch)
	start := time.Now()
	results.QueryRow(func(i int, id pgtype.UUID, err error) {
		if err != nil {
			s.Logger.Error("InsertUUIDv7Bulk", "uuid", UUIDString(id), "err", err)
			return
		}
		out <- id
		s.Logger.Info("InsertUUIDv7Bulk", "uuid", UUIDString(id), "version", uuid.V7, "time", time.Since(start))
		start = time.Now()
		if i >= len(batch)-1 {
			close(out)
		}
	})
	return out, err
}

func (s *Series) GetUUIDv7(id pgtype.UUID) error {
	var err error
	start := time.Now()
	if _, err = s.Query.GetUUIDv7(context.Background(), id); err != nil {
		s.Logger.Error("GetUUIDv7", "uuid", UUIDString(id), "err", err)
		return err
	}
	s.Logger.Info("GetUUIDv7", "uuid", UUIDString(id), "version", uuid.V7, "time", time.Since(start))
	return err
}

func UUIDString(id pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", id.Bytes[0:4], id.Bytes[4:6], id.Bytes[6:8], id.Bytes[8:10], id.Bytes[10:16])
}
