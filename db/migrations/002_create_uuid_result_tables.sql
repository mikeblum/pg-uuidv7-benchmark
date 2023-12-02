-- +goose Up
CREATE TABLE IF NOT EXISTS uuid_result (
    id UUID PRIMARY KEY NOT NULL,
    version SMALLINT NOT NULL,
    insert_duration_ns BIGINT,
    lookup_duration_ns BIGINT
);

CREATE INDEX IF NOT EXISTS idx_uuid_result_id ON uuid_result (id);

CREATE INDEX IF NOT EXISTS idx_uuid_result_version ON uuid_result (version);

-- +goose Down
DROP TABLE IF EXISTS uuid_result;

DROP INDEX IF EXISTS idx_uuid_result_id;

DROP INDEX IF EXISTS idx_uuid_result_version;
