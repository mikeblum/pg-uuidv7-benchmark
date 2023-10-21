CREATE TABLE
    IF NOT EXISTS uuid_v4 (
        id UUID PRIMARY KEY NOT NULL,
        created TIMESTAMP UNIQUE NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_uuid_v4_id ON uuid_v4 (id);

CREATE INDEX IF NOT EXISTS idx_uuid_v4_created ON uuid_v4 (created);
