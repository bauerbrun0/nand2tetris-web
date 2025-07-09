# run templ generation in watch mode to detect all .templ file changes and
# re-create _templ.go and _templ.txt files on change, then send reload event to browser.
# web server is running at http://localhost:3000
# live reload is at http://0.0.0.0:8080 (open this in your browser)
dev/templ:
	go tool templ generate --watch --proxy="http://localhost:3000" --proxyport 8080 --proxybind "0.0.0.0" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
dev/server:
	go tool github.com/air-verse/air \
	--build.cmd "go build -o tmp/bin/web ./cmd/web" --build.bin "tmp/bin/web" --build.delay "100" --build.args_bin "--dev" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# transpile and bundle main.ts with esbuild in watch mode
dev/esbuild:
	bunx esbuild ui/js/main.ts --bundle --outdir=ui/static/js/ --sourcemap --watch

# build css using tailwind in watch mode
dev/tailwind:
	bunx @tailwindcss/cli -i ./ui/css/main.css -o ./ui/static/css/main.css --watch --map

# watch for any js or css change in the ui/static/ folder, then reload the browser via templ proxy.
dev/sync_static:
	go run github.com/air-verse/air \
	--build.cmd "templ generate --notify-proxy --proxyport 8080 --proxybind \"0.0.0.0\"" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "ui/static" \
	--build.include_ext "js,css"

# generate _templ.go files
build/templ:
	go tool templ generate

# transpile and bundle main.ts with esbuild with sourcemap for development
build/esbuild:
	bunx esbuild ui/js/main.ts --bundle --outdir=ui/static/js/ --sourcemap

# transpile and bundle main.ts with esbuild for production
build/esbuild/prod:
	bunx esbuild ui/js/main.ts --bundle --outdir=ui/static/js/ --minify

# build css using tailwind with sourcemap for development
build/tailwind:
	bunx @tailwindcss/cli -i ./ui/css/main.css -o ./ui/static/css/main.css --map

# build css using tailwind with for production
build/tailwind/prod:
	bunx @tailwindcss/cli -i ./ui/css/main.css -o ./ui/static/css/main.css --minify

# build go web server for production
build/web:
	go build -o build/bin/web ./cmd/web

# run live dev environment
dev:
	make build/templ build/esbuild build/tailwind && bun run dev

# build for production
build/prod:
	make build/templ build/esbuild/prod build/tailwind/prod build/web

# remove generated templ files
rm/templ:
	find . -type f \( -name '*_templ.go' -o -name '*_templ.txt' \) -exec rm {} +

# remove generated files
clear:
	rm -rf build ui/static/css ui/static/js && make rm/templ
