#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

DATABASE_URL="${DATABASE_URL:-postgres://taskbridge:taskbridge@localhost:5432/taskbridge?sslmode=disable}"
MIGRATIONS_PATH="${MIGRATIONS_PATH:-${ROOT_DIR}/migrations}"

if ! command -v migrate >/dev/null 2>&1; then
  echo "golang-migrate is not installed."
  echo "Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
  exit 1
fi

cmd="${1:-up}"

case "${cmd}" in
  up)
    migrate -path "${MIGRATIONS_PATH}" -database "${DATABASE_URL}" up
    ;;
  down)
    migrate -path "${MIGRATIONS_PATH}" -database "${DATABASE_URL}" down 1
    ;;
  version)
    migrate -path "${MIGRATIONS_PATH}" -database "${DATABASE_URL}" version
    ;;
  force)
    if [[ $# -lt 2 ]]; then
      echo "usage: $0 force <version>"
      exit 1
    fi
    migrate -path "${MIGRATIONS_PATH}" -database "${DATABASE_URL}" force "$2"
    ;;
  create)
    if [[ $# -lt 2 ]]; then
      echo "usage: $0 create <name>"
      exit 1
    fi
    migrate create -ext sql -dir "${MIGRATIONS_PATH}" -seq "$2"
    ;;
  *)
    echo "usage: $0 {up|down|version|force|create}"
    exit 1
    ;;
esac
