#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

cleanup() {
  jobs -pr | xargs -r kill
}
trap cleanup EXIT INT TERM

: "${APP_ENV:=development}"
: "${PORT:=4545}"
: "${QUEUE_DIR:=/tmp/ppharma-queue}"
: "${QUEUE_TOPIC:=inventory.sync}"
: "${CRON_TICK_SECONDS:=10}"
: "${WORKER_POLL_SECONDS:=2}"
: "${WORKER_CONSUMER_ID:=worker-1}"
: "${JWT_SECRET:=dev-secret}"
: "${INTERNAL_API_KEYS:=int1:internal-secret:inventory.write|orders.item_status.write|products.write|queue.write}"

APP_ENV="$APP_ENV" PORT="$PORT" JWT_SECRET="$JWT_SECRET" INTERNAL_API_KEYS="$INTERNAL_API_KEYS" QUEUE_DIR="$QUEUE_DIR" QUEUE_TOPIC="$QUEUE_TOPIC" CRON_TICK_SECONDS="$CRON_TICK_SECONDS" WORKER_POLL_SECONDS="$WORKER_POLL_SECONDS" WORKER_CONSUMER_ID="$WORKER_CONSUMER_ID" go run ./cmd/api &
APP_ENV="$APP_ENV" QUEUE_DIR="$QUEUE_DIR" QUEUE_TOPIC="$QUEUE_TOPIC" CRON_TICK_SECONDS="$CRON_TICK_SECONDS" go run ./cmd/cron &
APP_ENV="$APP_ENV" QUEUE_DIR="$QUEUE_DIR" QUEUE_TOPIC="$QUEUE_TOPIC" WORKER_POLL_SECONDS="$WORKER_POLL_SECONDS" WORKER_CONSUMER_ID="$WORKER_CONSUMER_ID" go run ./cmd/worker &

wait
