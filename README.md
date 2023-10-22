# pg-uuid7

### Bootstrap Postgres Schema

Execute as the `postgres` user":

```
CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION app;

GRANT USAGE ON SCHEMA app TO app;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA app TO app;
```
