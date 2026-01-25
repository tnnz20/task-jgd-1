.PHONY: run build start clean test test-cover test-verbose

# Application name
APP_NAME=app
BUILD_DIR=bin

# Run the application directly (development)
run:
	go run ./cmd/http

# Build the application
build:
	@echo "Building application..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME).exe ./cmd/http
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME).exe"

# Build and run the application
start: build
	@echo "Starting application..."
	./$(BUILD_DIR)/$(APP_NAME).exe

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
	@echo "Clean complete"

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test ./... -v

# Run tests with coverage
test-cover:
	go test ./... -cover

# Run tests with coverage report
test-cover-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Help
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application (development)"
	@echo "  make build        - Build the application"
	@echo "  make start        - Build and run the application"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make test-verbose - Run tests with verbose output"
	@echo "  make test-cover   - Run tests with coverage"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Lint code"
	@echo "  make tidy         - Tidy dependencies"