-- +goose Up
CREATE TABLE IF NOT EXISTS uuid_result (
    id UUID NOT NULL,
    id_idx TEXT NOT NULL,
    version SMALLINT NOT NULL,
    insert_duration_ns BIGINT,
    lookup_duration_ns BIGINT,
    PRIMARY KEY(id, id_idx)
);

CREATE INDEX IF NOT EXISTS idx_uuid_result_id ON uuid_result (id);

CREATE INDEX IF NOT EXISTS idx_uuid_result_version ON uuid_result (version);

CREATE INDEX IF NOT EXISTS idx_uuid_result_idx ON uuid_result (id_idx);

-- +goose Down
DROP TABLE IF EXISTS uuid_result;

DROP INDEX IF EXISTS idx_uuid_result_id;

DROP INDEX IF EXISTS idx_uuid_result_version;

DROP INDEX IF EXISTS idx_uuid_result_idx;
