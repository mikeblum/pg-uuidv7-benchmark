-- +goose Up
CREATE TABLE
    IF NOT EXISTS uuid_v4 (
        id UUID PRIMARY KEY NOT NULL,
        created TIMESTAMP UNIQUE NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_uuid_v4_id ON uuid_v4 (id);

CREATE INDEX IF NOT EXISTS idx_uuid_v4_created ON uuid_v4 (created);

CREATE TABLE
    IF NOT EXISTS uuid_v7 (
        id UUID PRIMARY KEY NOT NULL,
        id_brin UUID NOT NULL,
        created TIMESTAMP UNIQUE NOT NULL
    );

-- B-tree index by default
CREATE INDEX IF NOT EXISTS idx_uuid_v7_id ON uuid_v7 (id);

-- BRIN index
CREATE INDEX IF NOT EXISTS idx_uuid_v7_id_brin ON uuid_v7 USING BRIN (id_brin);

CREATE INDEX IF NOT EXISTS idx_uuid_v7_created ON uuid_v7 (created);

-- +goose Down
DROP TABLE IF EXISTS uuid_v4;

DROP INDEX IF EXISTS idx_uuid_v4_id;
DROP INDEX IF EXISTS idx_uuid_v4_created;

DROP TABLE IF EXISTS uuid_v7;

DROP INDEX IF EXISTS idx_uuid_v7_id;
DROP INDEX IF EXISTS idx_uuid_v7_created;
