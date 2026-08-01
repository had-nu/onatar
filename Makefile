GOBIN := $(shell go env GOPATH)/bin
MIGRATE ?= $(GOBIN)/migrate

# Load local .env (ignored in CI; vars there are already in the environment).
-include .env
export
DB_DSN ?= mysql://$${DB_USER}:$${DB_PASS}@tcp($${DB_HOST}:$${DB_PORT})/$${DB_NAME}?multiStatements=true

.PHONY: help dev migrate rollback seed test test-integration e2e lint build release clean

help:
	@grep -E '^[a-zA-Z_-]+:' Makefile | sed 's/:.*//' | sort

# ---- database ----

migrate: $(MIGRATE)
	$(MIGRATE) -path migrations -database "$(DB_DSN)" up

rollback: $(MIGRATE)
	$(MIGRATE) -path migrations -database "$(DB_DSN)" down 1

$(MIGRATE):
	@command -v migrate >/dev/null 2>&1 || go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
seed:
	go run ./cmd/seed

# ---- dev ----

dev: 
	@echo "backend on :8090, frontend on :5173"; \
	(go run ./cmd/server &) ; \
	(cd frontend && npm run dev)

# ---- quality ----

test:
	go test ./cmd/... ./internal/...
	cd frontend && npm run test

# Requer TEST_DB_DSN (DSN MySQL); precisa de backend de dev ou de uma MariaDB.
test-integration:
	TEST_DB_DSN="$${TEST_DB_DSN}" go test -count=1 ./internal/integration/

e2e:
	cd frontend && npm run test:e2e

lint:
	golangci-lint run
	cd frontend && npm run lint

# ---- build ----

build:
	go build -o bin/onatar-server ./cmd/server
	go build -o bin/onatar-seed ./cmd/seed

release:
	GOOS=linux GOARCH=amd64 go build -o bin/onatar-server-linux-amd64 ./cmd/server
	GOOS=linux GOARCH=arm64 go build -o bin/onatar-server-linux-arm64 ./cmd/server

clean:
	rm -rf bin
