.PHONY: build run test clean tidy

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application
run:
	go run ./cmd/server

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download dependencies
tidy:
	go mod tidy

# Run database migrations
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/quizme_go_db?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/quizme_go_db?sslmode=disable" down

# Generate mocks (requires mockery)
mocks:
	mockery --all --output ./internal/mocks

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Development server with hot reload (requires air)
dev:
	air

# Docker commands
docker-build:
	docker build -t quizme-backend .

docker-run:
	docker run -p 8080:8080 quizme-backend
