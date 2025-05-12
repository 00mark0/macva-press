include .env
# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/templ:
	templ generate --watch --proxy="http://localhost:3000" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/air-verse/air@latest \
	--build.cmd "go build -o tmp/bin/main" --build.bin "./tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind:
	npx --yes tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify --watch

# Run Tailwind in the background, then start templ and server
live:
	make -j3 live/tailwind live/templ live/server 

migrategenerate:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

# PRODUCTION
# Run Tailwind once for production (no watch mode)
prod/tailwind:
	npx --yes tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify
