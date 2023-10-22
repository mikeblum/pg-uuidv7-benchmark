-- name: GetUUIDv4 :one
SELECT * FROM uuid_v4 WHERE id = $1 LIMIT 1;

-- name: InsertUUIDv4Bulk :batchone

INSERT INTO uuid_v4(id, created)
VALUES($1, $2)
RETURNING id;

-- name: GetUUIDv7 :one
SELECT * FROM uuid_v7 WHERE id = $1 LIMIT 1;

-- name: InsertUUIDv7Bulk :batchone

INSERT INTO uuid_v7(id, created)
VALUES($1, $2)
RETURNING id;

-- name: GenerateSeries :many

-- casting resolves computation requirement
-- https://github.com/sqlc-dev/sqlc/issues/1995
SELECT ts::timestamp FROM generate_series(
	date_trunc('day', now()::timestamp) - INTERVAL '1 year',
    now()::timestamp,
    INTERVAL '1 minute'
) AS ts;
