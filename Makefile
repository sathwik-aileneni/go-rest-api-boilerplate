.PHONY: help build run test clean docker-up docker-down migrate

help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-up    - Start services with docker-compose"
	@echo "  docker-down  - Stop services with docker-compose"
	@echo "  migrate      - Run database migrations"

build:
	@echo "Building application..."
	go build -o bin/main ./cmd/api

run:
	@echo "Running application..."
	go run ./cmd/api/main.go

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

docker-up:
	@echo "Starting services..."
	docker-compose up -d

docker-down:
	@echo "Stopping services..."
	docker-compose down

migrate:
	@echo "Running migrations..."
	@echo "Please run migrations manually using your preferred migration tool"
	@echo "Migration files are in ./migrations directory"

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Linting code..."
	golangci-lint run
