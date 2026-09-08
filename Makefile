.PHONY: fmt tidy lint test build run-api run-worker check clean

BINARY_DIR := bin

fmt:
	go fmt ./...

lint:
	go vet ./...

tidy:
	go mod tidy

check: fmt lint tidy

test:
	go test ./...

build:
	mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/api ./cmd/api
	go build -o $(BINARY_DIR)/worker ./cmd/worker

clean:
	rm -rf $(BINARY_DIR)

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

create-job:
	curl -X POST http://localhost:8080/jobs \
		-H "Content-Type: application/json" \
		-d '{"type":"email","payload":{"to":"user@example.com","subject":"Welcome"}}'
		