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
$ make db/sqlc
```

Start docker compose:

```bash
$ docker compose -f development/docker-compose.yaml up -d
```

Run migrations:

```bash
$ make db/migrate
```

Generate mock files for testing:

```bash
make test/generate
```

Run the `dev` target with `make`:

```bash
$ make dev
```

To attach to the postgres container:

```bash
$ docker exec -it develpment-db sh
```

You can generate a TLS certificate for local development using:

```bash
./scripts/generate-tls.sh
# Or specify your local IP to access the server from your network:
./scripts/generate-tls.sh
```

You can generate a TLS certificate for local development using:

```bash
./scripts/generate-tls.sh
# Or specify your local IP to access the server from your network:
./scripts/generate-tls.sh --local-ip <your-local-ip>
```

> **Note:**
>
> - This only works with production builds (`make build/prod`). The Templ proxy does not support HTTPS.
> - OAuth will not work because `https://localhost` is not registered as a redirect URL in GitHub or Google.
