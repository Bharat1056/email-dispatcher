.DEFAULT_GOAL := help

.PHONY: help run build lint test setup migrate-up migrate-down

help:
	@echo "Available targets:"
	@echo "  run"
	@echo "  build"
	@echo "  lint"
	@echo "  test"
	@echo "  setup"
	@echo "  migrate-up"
	@echo "  migrate-down"


run:
	go run cmd/server/main.go

build:
	@mkdir -p bin
	go build -o bin/server cmd/server/main.go

lint:
	golangci-lint run ./...

test:
	go test ./...

setup:
	go install github.com/evilmartians/lefthook@latest
	lefthook install
	go mod tidy

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down
