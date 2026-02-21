.PHONY: run-api run-cron run-worker run-all-local up down logs test

run-api:
	go run ./cmd/api

run-cron:
	go run ./cmd/cron

run-worker:
	go run ./cmd/worker

run-all-local:
	./scripts/run-all-local.sh

up:
	docker compose -f deployments/docker/docker-compose.yml up --build

down:
	docker compose -f deployments/docker/docker-compose.yml down

logs:
	docker compose -f deployments/docker/docker-compose.yml logs -f api cron worker

test:
	go test ./...
