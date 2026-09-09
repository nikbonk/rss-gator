# Requirements

We need to grab some tools to get up and running:

1. goose (used for database migrations)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

2. sqlc (used to generate our database queries into Go code)

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

3. Postgres driver (we just need it to connect to the database thats all)

```bash
go get github.com/lib/pq
```

# Setup Postgres

Just run the script `./startPostgres.sh` wait a second and then:

```bash
psql -h localhost -p 5432 -U postgres
```

password is `postgres`

## Using goose

For ease of use set the `GOOSE_MIGRATION_DIR` environment variable to the path of the `migrations` directory:

```bash
export GOOSE_MIGRATION_DIR="$(pwd)/sql/schema"
```

To run goose migrations, use the following command to up the database schema:

```bash
goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

or to down the database schema:

```bash
goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" down
```
