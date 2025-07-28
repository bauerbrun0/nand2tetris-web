# nand2tetris-web

## Development

Install the dependencies:
```bash
$ go mod tidy
$ bun install
```

Install golang migrate cli tool:
```bash
$ cd ~/Downloads
$ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz
$ mv migrate $GOPATH/bin/migrate
$ migrate -version
```

Create the `.env` file and fill in the values:
```bash
$ cp .env.example .env
```

Generate `sqlc` files:
```bash
$ make build/sqlc
```

Start docker compose:
```bash
$ docker compose -f development/docker-compose.yaml up -d
```

Run migrations:
```bash
$ make db/migrate
```

Run the `dev` target with `make`:
```bash
$ make dev
```

To attach to the postgres container:

```bash
$ docker exec -it develpment-db sh
```
