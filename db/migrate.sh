#!/bin/sh

set -e

IS_SEED=0
DB_URL=${DATABASE_URL}

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

go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./migrations postgres "$DB_URL" up
echo "Migrations complete."

if [ "$IS_SEED" -eq 1 ]; then
  echo "Seeding database..."
  goose -no-versioning -dir ./seed postgres "$DB_URL" up
fi
