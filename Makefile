.PHONY: run build lint test setup

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
