-- name: InsertUUIDv4Bulk :batchmany

INSERT INTO uuid_v4(id, created)
VALUES($1, $2);

-- name: InsertUUIDv7Bulk :batchmany

INSERT INTO uuid_v7(id, created)
VALUES($1, $2);

-- name: GenerateSeries :many

-- casting resolves computation requirement
-- https://github.com/sqlc-dev/sqlc/issues/1995
SELECT ts::timestamp FROM generate_series(
	date_trunc('day', now()::timestamp) - INTERVAL '1 year',
    now()::timestamp,
    INTERVAL '1 minute'
) AS ts;
