APP_NAME := kopesa-loan-platform-api
IMAGE := 010309/kopesa-loan-platform:latest
CONTAINER := kopesa-loan-api
PORT := 8080
HOST_PORT := 3352
COMPOSE_DIR := /apps/docker-compose-script/kopesa-loan
SERVER := safer

.PHONY: test fmt build migrate docker-build docker-push compose-up compose-down compose-pull server-dir deploy-files deploy-up

fmt:
	gofmt -w cmd internal

test:
	GOCACHE=$(PWD)/.cache/go-build GOMODCACHE=$(PWD)/.cache/gomod go test ./...

build:
	mkdir -p .bin
	GOCACHE=$(PWD)/.cache/go-build GOMODCACHE=$(PWD)/.cache/gomod go build -o .bin/kopesa-api ./cmd/api
	GOCACHE=$(PWD)/.cache/go-build GOMODCACHE=$(PWD)/.cache/gomod go build -o .bin/kopesa-migrate ./cmd/migrate

migrate:
	GOCACHE=$(PWD)/.cache/go-build GOMODCACHE=$(PWD)/.cache/gomod go run ./cmd/migrate

docker-build:
	docker build -t $(IMAGE) .

docker-push:
	docker push $(IMAGE)

compose-up:
	docker compose up -d

compose-down:
	docker compose down

compose-pull:
	docker compose pull

server-dir:
	ssh $(SERVER) "mkdir -p $(COMPOSE_DIR)"

deploy-files:
	scp docker-compose.yml $(SERVER):$(COMPOSE_DIR)/docker-compose.yml

deploy-up: deploy-files
	ssh $(SERVER) "cd $(COMPOSE_DIR) && sudo docker compose pull && sudo docker compose up -d"
