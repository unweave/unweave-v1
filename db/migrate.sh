#!/bin/sh

set -e

IS_SEED=0
DB_URL=${DATABASE_URL}
BASEDIR=$(dirname "$0")

if [ $# -eq 1 ]; then
  if [ "$1" = "seed" ]; then
    IS_SEED=1
  else
    echo "Usage: $0 [seed]"
    exit 1
  fi
elif [ $# -gt 1 ]; then
  echo "Invalid number of arguments. Max one arg allowed: [seed]."
fi

psql "$DB_URL" -f "$BASEDIR/init-scripts/1_initial_schema.sql"
goose -dir ./migrations postgres "$DB_URL" up
echo "Migrations complete."

if [ "$IS_SEED" -eq 1 ]; then
  echo "Seeding database..."
  goose -no-versioning -dir ./seed postgres "$DB_URL" up
fi
