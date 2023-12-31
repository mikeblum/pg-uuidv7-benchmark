-- name: GetUUIDv4 :one
SELECT * FROM uuid_v4 WHERE id = $1 LIMIT 1;

-- name: InsertUUIDv4Bulk :batchone

INSERT INTO uuid_v4(id, created)
VALUES($1, $2)
RETURNING id;

-- name: TruncateUUIDv4 :exec
TRUNCATE uuid_v4 RESTART IDENTITY CASCADE;

-- name: GetUUIDv7 :one
SELECT * FROM uuid_v7 WHERE id = $1 LIMIT 1;

-- name: GetUUIDv7BRIN :one
SELECT * FROM uuid_v7 WHERE id_brin = $1 LIMIT 1;

-- name: InsertUUIDv7Bulk :batchone

INSERT INTO uuid_v7(id, id_brin, created)
VALUES($1, $2, $3)
RETURNING id;

-- name: TruncateUUIDv7 :exec
TRUNCATE uuid_v7 RESTART IDENTITY CASCADE;

-- name: TruncateUUIDResult :exec
TRUNCATE uuid_result RESTART IDENTITY CASCADE;

-- name: VacuumUUIDv4 :exec
VACUUM VERBOSE ANALYZE uuid_v4;

-- name: VacuumUUIDv7 :exec
VACUUM VERBOSE ANALYZE uuid_v7;

-- name: GenerateSeries :many

-- casting resolves computation requirement
-- https://github.com/sqlc-dev/sqlc/issues/1995
SELECT ts::timestamp FROM generate_series(
	date_trunc('day', now()::timestamp) - INTERVAL '1 month',
    now()::timestamp,
    INTERVAL '1 minute'
) AS ts;

-- name: InsertUUIDResult :batchone

INSERT INTO uuid_result(id, id_idx, version, insert_duration_ns, lookup_duration_ns)
VALUES($1, $2, $3, $4, $5)
RETURNING id;
