# pg-uuid7

## Overview

With the upcoming acceptance of time-based UUIDs (v7) I wanted to validate using v4 (random) vs v7 (time-based) as primary keys under Postgres 15. 

In theory we should see faster insert and lookup times for UUIDv7 primary keys compared to random UUIDv4 keys which exibit poor index locality:

[CYBERTEC - Unexpected downsides of UUID keys in PostgreSQL](https://www.cybertec-postgresql.com/en/unexpected-downsides-of-uuid-keys-in-postgresql/)

Normally I'd prefer the default [google/uuid](https://pkg.go.dev/github.com/google/uuid) package but it doesn't support `v7` UUIDs just yet.

### Sources

[UUIDv4 - RFC 4122](https://www.rfc-editor.org/rfc/rfc4122)

[UUIDv7 Proposal](https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-04)

[Timescale - generate_series](https://www.timescale.com/blog/how-to-create-lots-of-sample-time-series-data-with-postgresql-generate_series/)

## Setup

- golang 1.21.x
- postgres 15.x
- gnu make for build tooling
- goose for schema migrations
- sqlc for SQL DDL

`SELECT version();`

> PostgreSQL 15.4, compiled by Visual C++ build 1914, 64-bit

### Bootstrap SQL DDL

```
make build
```

### Bootstrap Postgres Schema

Execute as the `postgres` user":

```
CREATE USER app WITH PASSWORD '????';

GRANT CONNECT ON DATABASE postgres TO app;

CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION app;

GRANT USAGE ON SCHEMA app TO app;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA app TO app;
```

### Run

```
export DATABASE_URL="postgres://app:????@localhost:5432/postgres?search_path=app"
make run
```
