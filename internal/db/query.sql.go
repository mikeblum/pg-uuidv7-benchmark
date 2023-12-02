// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const generateSeries = `-- name: GenerateSeries :many

SELECT ts::timestamp FROM generate_series(
	date_trunc('day', now()::timestamp) - INTERVAL '1 day',
    now()::timestamp,
    INTERVAL '1 minute'
) AS ts
`

// casting resolves computation requirement
// https://github.com/sqlc-dev/sqlc/issues/1995
func (q *Queries) GenerateSeries(ctx context.Context) ([]pgtype.Timestamp, error) {
	rows, err := q.db.Query(ctx, generateSeries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Timestamp
	for rows.Next() {
		var ts pgtype.Timestamp
		if err := rows.Scan(&ts); err != nil {
			return nil, err
		}
		items = append(items, ts)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUUIDv4 = `-- name: GetUUIDv4 :one
SELECT id, created FROM uuid_v4 WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUUIDv4(ctx context.Context, id pgtype.UUID) (UuidV4, error) {
	row := q.db.QueryRow(ctx, getUUIDv4, id)
	var i UuidV4
	err := row.Scan(&i.ID, &i.Created)
	return i, err
}

const getUUIDv7 = `-- name: GetUUIDv7 :one
SELECT id, created FROM uuid_v7 WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUUIDv7(ctx context.Context, id pgtype.UUID) (UuidV7, error) {
	row := q.db.QueryRow(ctx, getUUIDv7, id)
	var i UuidV7
	err := row.Scan(&i.ID, &i.Created)
	return i, err
}

const truncateUUIDResult = `-- name: TruncateUUIDResult :exec
TRUNCATE uuid_result RESTART IDENTITY CASCADE
`

func (q *Queries) TruncateUUIDResult(ctx context.Context) error {
	_, err := q.db.Exec(ctx, truncateUUIDResult)
	return err
}

const truncateUUIDv4 = `-- name: TruncateUUIDv4 :exec
TRUNCATE uuid_v4 RESTART IDENTITY CASCADE
`

func (q *Queries) TruncateUUIDv4(ctx context.Context) error {
	_, err := q.db.Exec(ctx, truncateUUIDv4)
	return err
}

const truncateUUIDv7 = `-- name: TruncateUUIDv7 :exec
TRUNCATE uuid_v7 RESTART IDENTITY CASCADE
`

func (q *Queries) TruncateUUIDv7(ctx context.Context) error {
	_, err := q.db.Exec(ctx, truncateUUIDv7)
	return err
}
