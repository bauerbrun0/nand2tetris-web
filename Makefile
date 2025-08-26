###################
### development ###
###################


# run templ generation in watch mode to detect all .templ file changes and
# re-create _templ.go and _templ.txt files on change, then send reload event to browser.
# web server is running at http://localhost:3000
# Keep the --proxy in sync with the PORT environment variable.
# live reload is at http://0.0.0.0:8080 (or http://localhost:8080) (open this in your browser)
dev/templ:
	go tool templ generate --watch --proxy="http://localhost:3000" --proxyport 8080 --proxybind "0.0.0.0" --open-browser=false -v

# run air to detect any go or yaml (translation files)
# file changes to re-build and re-run the server.
dev/server:
	go tool github.com/air-verse/air \
	--build.cmd "go build -o tmp/bin/web ./cmd/web" --build.bin "tmp/bin/web" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,yaml" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# watch for any change in the ui/static/ folder, then reload the browser via templ proxy.
dev/sync_static:
	go run github.com/air-verse/air \
	--build.cmd "go tool templ generate --notify-proxy --proxyport 8080 --proxybind \"0.0.0.0\"" \
	--build.full_bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "ui/static" \
	--build.include_ext "js,css"

# run svelte-check in watch mode
dev/svelte-check:
	bunx svelte-check --watch

# transpile and bundle main.ts with esbuild in watch mode
dev/esbuild:
	bunx esbuild ui/js/entries/main.ts --bundle --outdir=ui/static/js/ --sourcemap --watch

# compile and bundle svelte files in watch mode
dev/esbuild/svelte:
	bun scripts/esbuild/svelte/build-watch.js

# build css using tailwind in watch mode
dev/tailwind:
	bunx @tailwindcss/cli -i ./ui/css/main.css -o ./ui/static/css/main.css --watch --map

# run live dev environment
dev:
	bun run dev


##################
### production ###
##################


# generate _templ.go files
build/templ:
	go tool templ generate

# run svelte-check before building
build/svelte-check:
	bunx svelte-check --fail-on-warnings

# transpile and bundle main.ts with esbuild for production
build/esbuild:
	bunx esbuild ui/js/entries/main.ts --bundle --outdir=ui/static/js/ --minify

# compile and bundle svelte files for production
build/esbuild/svelte:
	bun scripts/esbuild/svelte/build-prod.js

# build css using tailwind with for production
build/tailwind:
	bunx @tailwindcss/cli -i ./ui/css/main.css -o ./ui/static/css/main.css --minify

# build go web server for production
build/web:
	go build -o build/bin/web ./cmd/web

# build for production
build/prod:
	make build/templ build/svelte-check build/esbuild build/esbuild/svelte build/tailwind db/sqlc build/web


##########
### db ###
##########


# generate sqlc files
db/sqlc:
	go tool sqlc generate

# migrate db
db/migrate:
	migrate -path=./db/migrations -database="postgres://nand2tetris_web_migration:password@localhost/nand2tetris_web?sslmode=disable" up

db/migrate/down:
	migrate -path=./db/migrations -database="postgres://nand2tetris_web_migration:password@localhost/nand2tetris_web?sslmode=disable" down


##############################
### linting and formatting ###
##############################


# check linting
lint/check:
	bunx eslint .

# check formatting
format/check:
	./scripts/check-go-formatting.sh && \
	go tool templ fmt -fail . && \
	bunx prettier . --check

# fix formatting
format/write:
	go fmt ./... && \
	go tool templ fmt . && \
	bunx prettier . --write


###############
### testing ###
###############


# generate mock files
test/generate:
	go tool mockery

# run all tests
test/all:
	go test ./...


###################
### cleaning up ###
###################


# remove generated templ files
rm/templ:
	find . -type f \( -name '*_templ.go' -o -name '*_templ.txt' \) -exec rm {} +

# remove sqlc generated db access files
rm/sqlc:
	rm internal/models/db.go internal/models/models.go && find . -type f \( -name '*.sql.go' \) -exec rm {} +

# remove generated files
rm:
	rm -rf build ui/static/css ui/static/js && make rm/templ

# remove all generated files, including slower-to-generate ones like sqlc output
rm/all:
	rm -rf build ui/static/css ui/static/js && make rm/templ rm/sqlc
