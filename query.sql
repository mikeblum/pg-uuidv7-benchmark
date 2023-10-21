-- name: GenerateUUIDv4 :exec

INSERT INTO uuid_v4(id, created)
-- gen_random_uuid generates a UUIDv4
SELECT gen_random_uuid(), ts FROM generate_series(
	date_trunc('day', now()::timestamp) - INTERVAL '1 year',
    now()::timestamp,
    INTERVAL '1 minute'
   ) as ts ON CONFLICT(created) DO NOTHING;
