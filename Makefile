.PHONY: help proto sqlc gen test build run-local migrate-up migrate-down seed reset-db docker-build docker-up docker-down clean

# Variables
BINARY_NAME=job-applicants-api
DOCKER_IMAGE=job-applicants-api
DATABASE_URL?=postgres://applicants_user:applicants_pass@localhost:5432/applicants?sslmode=disable

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  make proto           - Generate protobuf code with buf"
	@echo "  make sqlc            - Generate sqlc database code"
	@echo "  make gen             - Generate all code (proto + sqlc)"
	@echo "  make tidy            - Run go mod tidy"
	@echo "  make test            - Run unit tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make coverage-html   - Generate HTML coverage report"
	@echo "  make build           - Build server binary"
	@echo "  make run             - Run the server locally"
	@echo "  make migrate-up      - Run database migrations up"
	@echo "  make migrate-down    - Run database migrations down"
	@echo "  make seed            - Seed the database with sample data"
	@echo "  make reset-db        - Drop, create, migrate, and seed database"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-up       - Start services with docker compose"
	@echo "  make docker-down     - Stop services with docker compose"
	@echo "  make clean           - Clean build artifacts"

## proto: Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	buf generate

## sqlc: Generate sqlc database code
sqlc:
	@echo "Generating sqlc code..."
	sqlc generate

## gen: Generate all code
gen: proto sqlc
	@echo "All code generated successfully!"

## tidy: Run go mod tidy
tidy:
	@echo "Running go mod tidy..."
	go mod tidy

## test: Run unit tests
test:
	@echo "Running unit tests..."
	go test -v -race -short ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "\nCoverage summary:"
	go tool cover -func=coverage.out | grep total:
	@echo "\nTo view detailed HTML coverage report, run: make coverage-html"

## coverage-html: Generate and open HTML coverage report
coverage-html:
	@echo "Generating HTML coverage report..."
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

## build: Build server binary
build:
	@echo "Building server..."
	CGO_ENABLED=0 go build -o bin/server ./cmd/server
	CGO_ENABLED=0 go build -o bin/migrate ./cmd/migrate
	CGO_ENABLED=0 go build -o bin/seed ./cmd/seed
	@echo "Binaries built in bin/"

## run: Run the server locally
run: build
	@echo "Starting server..."
	./bin/server

## migrate-up: Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	DATABASE_URL=$(DATABASE_URL) go run ./cmd/migrate -direction up

## migrate-down: Run database migrations down
migrate-down:
	@echo "Running migrations down..."
	DATABASE_URL=$(DATABASE_URL) go run ./cmd/migrate -direction down

## seed: Seed the database
seed:
	@echo "Seeding database..."
	DATABASE_URL=$(DATABASE_URL) go run ./cmd/seed --clear

## reset-db: Reset database (down, up, seed)
reset-db: migrate-down migrate-up seed
	@echo "Database reset complete!"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

## docker-up: Start services with docker compose
docker-up:
	@echo "Starting services..."
	docker compose up -d postgres
	@echo "Waiting for database..."
	@sleep 3
	docker compose --profile setup run --rm migrate
	docker compose --profile seed run --rm seed
	docker compose up -d api
	@echo "Services started!"
	@echo "  - REST API: http://localhost:8080"
	@echo "  - gRPC API: localhost:9090"
	@echo "  - API Docs: http://localhost:8080/docs/"

## docker-down: Stop services with docker compose
docker-down:
	@echo "Stopping services..."
	docker compose down

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf api/proto/v1/*.go
	rm -rf api/proto/v1/*.json
	rm -rf internal/db/sqlc/*.go
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

## install-tools: Install required development tools
install-tools:
	@echo "Installing tools..."
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed!"
