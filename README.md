# PPharma Backend

Go backend for a pharma-style platform with three runnable services in one repository:

- `api`: Gin HTTP API
- `cron`: scheduled publisher service
- `worker`: queue consumer service

The project is organized around domain-first interfaces, with concrete third-party adapters isolated under `support-pkg`.

## Current Status

This repository is a backend foundation and scaffold, not a fully connected production system yet.

- HTTP routes, middleware, versioned routing, Swagger, Docker, config loading, queue, and service entrypoints are implemented.
- Order flows currently use an in-memory seeded repository for demo behavior.
- MongoDB wrappers and index helpers exist, but the main API is not yet fully wired to Mongo for all modules.
- Product, payment, consultation, and auth handlers are currently simple stubs/placeholders.

## Project Structure

```text
.
├── cmd/
│   ├── api/
│   ├── cron/
│   └── worker/
├── internal/
│   ├── app/
│   ├── config/
│   ├── domain/
│   ├── http/
│   │   ├── docs/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes/
│   │       └── v1/
│   ├── repository/
│   └── service/
├── support-pkg/
│   ├── ai/
│   ├── auth/
│   ├── cache/
│   ├── db/
│   ├── logger/
│   └── queue/
└── deployments/docker/
```

## Services

### Root entrypoint

The repository supports `go run .` through the root [`main.go`](/Users/dineshyadav/dev/p-pharma/backend/main.go).  
Select the service with `APP_SERVICE`.

Supported values:

- `api`
- `cron`
- `worker`

Example:

```bash
APP_SERVICE=api go run .
APP_SERVICE=cron go run .
APP_SERVICE=worker go run .
```

### Dedicated entrypoints

You can also run each service directly:

```bash
go run ./cmd/api
go run ./cmd/cron
go run ./cmd/worker
```

## Configuration

Configuration is loaded by [`internal/config/config.go`](/Users/dineshyadav/dev/p-pharma/backend/internal/config/config.go).

Precedence:

1. hardcoded defaults
2. `config.json` via `CONFIG_FILE`
3. environment variables

### `.env` example

Use [`.env.example`](/Users/dineshyadav/dev/p-pharma/backend/.env.example) as reference.

Key points:

- config keys use `UPPER_SNAKE_CASE`
- database config is passed as a `DB` JSON object
- MongoDB is expected to be external, not started by Docker Compose

Example:

```env
APP_ENV=development
APP_SERVICE=api
CONFIG_FILE=config.json
PORT=4545
JWT_SECRET=dev-secret
DB={"MONGO_URI":"mongodb://127.0.0.1:27017","MONGO_DB_NAME":"ppharma"}
INTERNAL_API_KEYS=int1:internal-secret:inventory.write|orders.item_status.write|products.write|queue.write
QUEUE_DIR=/tmp/ppharma-queue
QUEUE_TOPIC=inventory.sync
CRON_TICK_SECONDS=30
WORKER_POLL_SECONDS=2
WORKER_CONSUMER_ID=worker-1
```

### `config.json` example

Use [`config.example.json`](/Users/dineshyadav/dev/p-pharma/backend/config.example.json) as the template.

Example:

```json
{
  "APP_ENV": "development",
  "PORT": "4545",
  "JWT_SECRET": "dev-secret",
  "DB": {
    "MONGO_URI": "mongodb://127.0.0.1:27017",
    "MONGO_DB_NAME": "ppharma"
  },
  "INTERNAL_API_KEYS": [
    {
      "id": "int1",
      "key": "internal-secret",
      "scopes": [
        "inventory.write",
        "orders.item_status.write",
        "products.write",
        "queue.write"
      ]
    }
  ],
  "QUEUE_DIR": "/tmp/ppharma-queue",
  "QUEUE_TOPIC": "inventory.sync",
  "CRON_TICK_SECONDS": 30,
  "WORKER_POLL_SECONDS": 2,
  "WORKER_CONSUMER_ID": "worker-1"
}
```

## Local Development

### Run API

```bash
go run .
```

By default, root `main.go` runs the API if `APP_SERVICE` is not set.

### Run all services locally

Use the helper script:

```bash
make run-all-local
```

This runs:

- API
- cron
- worker

The helper script is at [`scripts/run-all-local.sh`](/Users/dineshyadav/dev/p-pharma/backend/scripts/run-all-local.sh).

### Useful make targets

Defined in [`Makefile`](/Users/dineshyadav/dev/p-pharma/backend/Makefile):

- `make run-api`
- `make run-cron`
- `make run-worker`
- `make run-all-local`
- `make up`
- `make down`
- `make logs`
- `make test`

## Docker

Docker files are split per service:

- [`Dockerfile.api`](/Users/dineshyadav/dev/p-pharma/backend/deployments/docker/Dockerfile.api)
- [`Dockerfile.cron`](/Users/dineshyadav/dev/p-pharma/backend/deployments/docker/Dockerfile.cron)
- [`Dockerfile.worker`](/Users/dineshyadav/dev/p-pharma/backend/deployments/docker/Dockerfile.worker)

Compose file:

- [`docker-compose.yml`](/Users/dineshyadav/dev/p-pharma/backend/deployments/docker/docker-compose.yml)

Start all services:

```bash
DB='{"MONGO_URI":"mongodb://10.0.0.5:27017","MONGO_DB_NAME":"ppharma"}' \
docker compose -f deployments/docker/docker-compose.yml up --build
```

Notes:

- Compose does not start MongoDB.
- You must provide a reachable external MongoDB URI through `DB`.
- Queue storage uses a shared volume mounted at `/shared-queue`.

## API Overview

### Health

- `GET /health/live`
- `GET /health/ready`

### Swagger

Available only when `APP_ENV != production`.

- `GET /swagger`
- `GET /swagger/openapi.yaml`

### Versioned routing

Versioned API routes live under [`internal/http/routes/v1`](/Users/dineshyadav/dev/p-pharma/backend/internal/http/routes/v1).

Current structure:

- auth: [`auth.go`](/Users/dineshyadav/dev/p-pharma/backend/internal/http/routes/v1/auth.go)
- customer: [`customer.go`](/Users/dineshyadav/dev/p-pharma/backend/internal/http/routes/v1/customer.go)
- admin: [`admin.go`](/Users/dineshyadav/dev/p-pharma/backend/internal/http/routes/v1/admin.go)
- internal admin: [`internal.go`](/Users/dineshyadav/dev/p-pharma/backend/internal/http/routes/v1/internal.go)

### Auth routes

Base path: `/api/v1/auth`

- `POST /login`
- `POST /refresh`
- `POST /logout`
- `GET /sessions`
- `DELETE /sessions/:sessionId`

### Customer routes

Base path: `/api/v1/customer`

- `GET /orders/:orderId`
- `GET /orders/:orderId/payments`
- `GET /products`
- `GET /products/:productId`
- `POST /consultations`
- `GET /consultations`
- `GET /consultations/:consultationId`

### Admin routes

Base path: `/api/v1/admin`

- `PATCH /orders/:orderId/items/:itemId/status`
- `POST /products`
- `PATCH /products/:productId`
- `PATCH /inventory/:productId/stock`
- `PATCH /payments/:paymentId/status`
- `PATCH /consultations/:consultationId/status`

### Admin internal routes

Base path: `/api/v1/admin/internal`

Authenticated with `X-API-Key` and scope checks.

- `POST /inventory/sync`
- `PATCH /orders/:orderId/items/:itemId/status`
- `POST /products/bulk-upsert`
- `POST /queue/publish`

## Authentication

### Customer/admin API

- JWT bearer token
- customer and admin routes are separated by route group and middleware

### Internal admin API

- `X-API-Key` header
- scope-based authorization

Default scopes used today:

- `inventory.write`
- `orders.item_status.write`
- `products.write`
- `queue.write`

## Queue, Cron, Worker

The queue is currently implemented with a simple file-backed adapter in [`support-pkg/queue/filequeue`](/Users/dineshyadav/dev/p-pharma/backend/support-pkg/queue/filequeue).

Current behavior:

- cron publishes messages to a queue topic on an interval
- worker polls and consumes those messages
- API internal route can publish a test message

This is a good development scaffold, but not yet a production-grade broker.

## Support Package

Concrete integrations live under [`support-pkg`](/Users/dineshyadav/dev/p-pharma/backend/support-pkg):

- auth adapters
- logger
- cache
- Mongo wrappers
- queue implementation
- future AI providers/factories

This keeps domain and usecase code loosely coupled to implementations.

## Testing

Run the test suite:

```bash
go test ./...
```

Current tests mainly cover:

- order status derivation and transitions
- middleware behavior
- app bootstrap routing assumptions

## Next Practical Steps

1. Replace in-memory repositories with Mongo-backed repositories in the API bootstrap path.
2. Replace stub handlers with real usecase logic for auth, products, payments, and consultations.
3. Add request validation, rate limiting, and structured error codes across all endpoints.
4. Introduce real queue/broker and real Mongo client lifecycle management for production.
