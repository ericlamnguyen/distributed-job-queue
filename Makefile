.PHONY: fmt tidy lint test build run-api run-worker check clean

BINARY_DIR := bin

fmt:
	go fmt ./...

tidy:
	go mod tidy

lint:
	go vet ./...

test:
	go test ./...

build:
	mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/api ./cmd/api
	go build -o $(BINARY_DIR)/worker ./cmd/worker

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

check: fmt lint

clean:
	rm -rf $(BINARY_DIR)