# jwt-in-go

Basic implementation jwt in a stateful way in go.

> P.S. Yes i'm aware that jwt should be stateless (like storing them in cookies). But i thought about how to make it more secure.
> Easier way instead of creating black list store for tokens, with db it is simple to invalidate refresh token.

### Install deps

```
go mod tidy
```

### Install goose for migrations

https://github.com/pressly/goose

### Install sqlc for schema generation

https://github.com/sqlc-dev/sqlc

### Run generate schema

```
sqlc generate -f ./db/sqlc/sqlc.yaml
```

### Run migrations with goose

```
goose -dir ./db/migrations postgres "host=localhost user=postgres dbname=postgres password=postgres port=5433 sslmode=disable" up
```
