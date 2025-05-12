#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -database "${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

echo "start the app"
exec "$@"
