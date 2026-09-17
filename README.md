# rss-gator

rss-gator is a command-line RSS feed aggregator written in Go. It uses PostgreSQL to store users, feeds and posts.

## Requirements

You need Go and PostgreSQL to run rss-gator.

- [Go](https://go.dev/dl/)
- PostgreSQL

If you want to use the included `startPostgres.sh` script, you also need either [Podman](https://podman.io/) or [Docker](https://www.docker.com/).

## Installation

Install the `gator` CLI with `go install`:

```bash
go install github.com/nikbonk/rss-gator@latest
```

This installs the `gator` binary into your Go bin directory. Make sure that directory is in your `PATH`.

You can then run:

```bash
gator
```

`go run .` is mainly for development. For normal use, use the `gator` binary.

## Setup

The repo includes a script to start PostgreSQL in either Podman or Docker.

```bash
./startPostgres.sh
```

The script starts PostgreSQL and initializes the database for a fresh container.

The database connection used by rss-gator is:

```text
postgres://postgres:postgres@localhost:5432/gator?sslmode=disable
```

## Config

rss-gator uses `~/.gatorconfig.json` for its configuration.

Create the file with:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": "your_username"
}
```

The `current_user_name` value is changed by the `register` and `login` commands.

## Commands

Register a user:

```bash
gator register alice
```

Log in:

```bash
gator login alice
```

List users:

```bash
gator users
```

Add a feed:

```bash
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
```

List feeds:

```bash
gator feeds
```

Follow a feed:

```bash
gator follow "https://news.ycombinator.com/rss"
```

See the feeds you follow:

```bash
gator following
```

Unfollow a feed:

```bash
gator unfollow "https://news.ycombinator.com/rss"
```

Start the aggregator:

```bash
gator agg 1m
```

Browse posts:

```bash
gator browse 10
```

## Development

When working on the source code, you can use:

```bash
go run . <command> [arguments]
```

For example:

```bash
go run . users
go run . feeds
go run . browse 10
```

### Database migrations

The database migrations are stored in `sql/schema` and are managed with goose.

First install goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Then from the repo root:

```bash
export GOOSE_MIGRATION_DIR="$(pwd)/sql/schema"
goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

The `init` directory contains the schema used when starting a fresh PostgreSQL container, so you do not need to run the migrations just to get a new development database running with `startPostgres.sh`.

### Building

You can also build the binary yourself:

```bash
go build -o gator
```

After `go build` or `go install`, you can run the compiled binary without needing the Go toolchain.

## Repository

https://github.com/nikbonk/rss-gator
