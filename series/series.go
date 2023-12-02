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

type IndexType int

const (
	attrError           = "error"
	BTREE     IndexType = 0
	BRIN      IndexType = 1
)

type Series struct {
	Query      *dbms.Queries
	Logger     *slog.Logger
	timestamps []pgtype.Timestamp
}

type UUID struct {
	ID      pgtype.UUID
	Version byte
	Insert  time.Duration
	Lookup  time.Duration
}

func (s *Series) GenerateSeries() error {
	var err error
	if s.timestamps, err = s.Query.GenerateSeries(context.Background()); err != nil {
		s.Logger.Error("Error generating series", attrError, err)
		return err
	}
	s.Logger.Info("GenerateSeries", "series #", len(s.timestamps))
	return err
}

func (s *Series) InsertUUIDv4Bulk() (chan UUID, error) {
	var err error
	batch := make([]dbms.InsertUUIDv4BulkParams, len(s.timestamps))
	out := make(chan UUID, len(batch))
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
		end := time.Since(start)
		out <- UUID{
			ID:      id,
			Version: uuid.V4,
			Insert:  end,
		}
		s.Logger.Debug("InsertUUIDv4Bulk", "uuid", UUIDString(id), "version", uuid.V4, "time", end)
		start = time.Now()
		if i >= len(batch)-1 {
			close(out)
		}
	})
	return out, err
}

func (s *Series) GetUUIDv4(id pgtype.UUID) (time.Duration, error) {
	var err error
	start := time.Now()
	if _, err = s.Query.GetUUIDv4(context.Background(), id); err != nil {
		s.Logger.Error("GetUUIDv4", "uuid", UUIDString(id), "err", err)
		return 0, err
	}
	end := time.Since(start)
	s.Logger.Debug("GetUUIDv4", "uuid", UUIDString(id), "version", uuid.V4, "time", end)
	return end, err
}

func (s *Series) InsertUUIDv7Bulk() (chan UUID, error) {
	var err error
	batch := make([]dbms.InsertUUIDv7BulkParams, len(s.timestamps))
	out := make(chan UUID, len(batch))
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
			IDBrin: pgtype.UUID{
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
		end := time.Since(start)
		out <- UUID{
			ID:      id,
			Version: uuid.V7,
			Insert:  end,
		}
		start = time.Now()
		if i >= len(batch)-1 {
			close(out)
		}
	})
	return out, err
}

func (s *Series) GetUUIDv7(id pgtype.UUID) (time.Duration, error) {
	var err error
	start := time.Now()
	if _, err = s.Query.GetUUIDv7(context.TODO(), id); err != nil {
		s.Logger.Error("GetUUIDv7", "uuid", UUIDString(id), "err", err)
		return 0, err
	}
	end := time.Since(start)
	s.Logger.Debug("GetUUIDv7", "uuid", UUIDString(id), "version", uuid.V7, "time", end)
	return end, err
}

func (s *Series) GetUUIDv7BRIN(id pgtype.UUID) (time.Duration, error) {
	var err error
	start := time.Now()
	if _, err = s.Query.GetUUIDv7BRIN(context.TODO(), id); err != nil {
		s.Logger.Error("GetUUIDv7", "uuid", UUIDString(id), "err", err)
		return 0, err
	}
	end := time.Since(start)
	s.Logger.Debug("GetUUIDv7BRIN", "uuid", UUIDString(id), "version", uuid.V7, "time", end)
	return end, err
}

// sqlc doesn't support MERGE queries:
// https://github.com/sqlc-dev/sqlc/issues/1661
func (s *Series) MergeUUIDResult(uuid UUID, index IndexType) error {
	var err error
	results := s.Query.InsertUUIDResult(context.Background(), []dbms.InsertUUIDResultParams{
		{
			ID:      uuid.ID,
			IDIdx:   IndexTypeString(index),
			Version: int16(uuid.Version),
			InsertDurationNs: pgtype.Int8{
				Int64: uuid.Insert.Nanoseconds(),
				Valid: true,
			},
			LookupDurationNs: pgtype.Int8{
				Int64: uuid.Lookup.Nanoseconds(),
				Valid: true,
			},
		},
	})
	results.QueryRow(func(i int, id pgtype.UUID, err error) {
		if err != nil {
			s.Logger.Error("MergeUUIDResult", "uuid", UUIDString(id), "err", err)
			return
		}
	})
	s.Logger.Debug("MergeUUIDResult", "uuid", UUIDString(uuid.ID), "version", uuid.Version, "insert", uuid.Insert, "lookup", uuid.Lookup)
	return err
}

func UUIDString(id pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", id.Bytes[0:4], id.Bytes[4:6], id.Bytes[6:8], id.Bytes[8:10], id.Bytes[10:16])
}

func IndexTypeString(indexType IndexType) string {
	switch indexType {
	case BTREE:
		return "BTREE"
	case BRIN:
		return "BRIN"
	default:
		return ""
	}
}
