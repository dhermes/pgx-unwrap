# `pgx-unwrap`

> Unwrap `database/sql` connections for `pgx` interop

In rare circumstances, a project may need to use both the standard library
driver (`github.com/jackc/pgx/v5/stdlib`) and a native `pgx` / `pgxpool`
client. In these rare cases, it may be necessary to convert connections of
one type to the other. This library provides helpers for "unwrapping"
connections for this purpose.
