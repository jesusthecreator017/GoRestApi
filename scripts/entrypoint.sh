#!/bin/sh
set -e

# The migrate binary expects pgx5:// scheme, but the API uses postgres://.
# Convert the DATABASE_URL for migrations.
MIGRATE_URL=$(echo "$DATABASE_URL" | sed 's|^postgres://|pgx5://|')

echo "Running migrations..."
DATABASE_URL="$MIGRATE_URL" ./migrate -direction up

echo "Starting API server..."
exec ./main
