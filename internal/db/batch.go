// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: batch.go

package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBatchAlreadyClosed = errors.New("batch already closed")
)

const insertUUIDResult = `-- name: InsertUUIDResult :batchone

INSERT INTO uuid_result(id, id_idx, version, insert_duration_ns, lookup_duration_ns)
VALUES($1, $2, $3, $4, $5)
RETURNING id
`

type InsertUUIDResultBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type InsertUUIDResultParams struct {
	ID               pgtype.UUID
	IDIdx            string
	Version          int16
	InsertDurationNs pgtype.Int8
	LookupDurationNs pgtype.Int8
}

func (q *Queries) InsertUUIDResult(ctx context.Context, arg []InsertUUIDResultParams) *InsertUUIDResultBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.ID,
			a.IDIdx,
			a.Version,
			a.InsertDurationNs,
			a.LookupDurationNs,
		}
		batch.Queue(insertUUIDResult, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &InsertUUIDResultBatchResults{br, len(arg), false}
}

func (b *InsertUUIDResultBatchResults) QueryRow(f func(int, pgtype.UUID, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var id pgtype.UUID
		if b.closed {
			if f != nil {
				f(t, id, ErrBatchAlreadyClosed)
			}
			continue
		}
		row := b.br.QueryRow()
		err := row.Scan(&id)
		if f != nil {
			f(t, id, err)
		}
	}
}

func (b *InsertUUIDResultBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}

const insertUUIDv4Bulk = `-- name: InsertUUIDv4Bulk :batchone

INSERT INTO uuid_v4(id, created)
VALUES($1, $2)
RETURNING id
`

type InsertUUIDv4BulkBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type InsertUUIDv4BulkParams struct {
	ID      pgtype.UUID
	Created pgtype.Timestamp
}

func (q *Queries) InsertUUIDv4Bulk(ctx context.Context, arg []InsertUUIDv4BulkParams) *InsertUUIDv4BulkBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.ID,
			a.Created,
		}
		batch.Queue(insertUUIDv4Bulk, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &InsertUUIDv4BulkBatchResults{br, len(arg), false}
}

func (b *InsertUUIDv4BulkBatchResults) QueryRow(f func(int, pgtype.UUID, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var id pgtype.UUID
		if b.closed {
			if f != nil {
				f(t, id, ErrBatchAlreadyClosed)
			}
			continue
		}
		row := b.br.QueryRow()
		err := row.Scan(&id)
		if f != nil {
			f(t, id, err)
		}
	}
}

func (b *InsertUUIDv4BulkBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}

const insertUUIDv7Bulk = `-- name: InsertUUIDv7Bulk :batchone

INSERT INTO uuid_v7(id, id_brin, created)
VALUES($1, $2, $3)
RETURNING id
`

type InsertUUIDv7BulkBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type InsertUUIDv7BulkParams struct {
	ID      pgtype.UUID
	IDBrin  pgtype.UUID
	Created pgtype.Timestamp
}

func (q *Queries) InsertUUIDv7Bulk(ctx context.Context, arg []InsertUUIDv7BulkParams) *InsertUUIDv7BulkBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.ID,
			a.IDBrin,
			a.Created,
		}
		batch.Queue(insertUUIDv7Bulk, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &InsertUUIDv7BulkBatchResults{br, len(arg), false}
}

func (b *InsertUUIDv7BulkBatchResults) QueryRow(f func(int, pgtype.UUID, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var id pgtype.UUID
		if b.closed {
			if f != nil {
				f(t, id, ErrBatchAlreadyClosed)
			}
			continue
		}
		row := b.br.QueryRow()
		err := row.Scan(&id)
		if f != nil {
			f(t, id, err)
		}
	}
}

func (b *InsertUUIDv7BulkBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}
