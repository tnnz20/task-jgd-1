.PHONY: run build start clean test test-cover test-verbose docker-build docker-run docker-stop migrate migrate-create migrate-up migrate-down help

# Application name
APP_NAME=app
APP_PATH=./cmd/http/main.go
BUILD_DIR=bin
DOCKER_IMAGE=category-api
MIGRATE_PATH=./cmd/migrate

# Detect Windows only for the .exe extension
ifeq ($(OS),Windows_NT)
	EXE_EXT = .exe
else
	EXE_EXT = 
endif

# Run the application directly (development)
run:
	@go run $(APP_PATH)

# Build the application
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME)$(EXE_EXT) $(APP_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)$(EXE_EXT)"

# Build and run the application
start: build
	@echo "Starting application..."
	./$(BUILD_DIR)/$(APP_NAME)$(EXE_EXT)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
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

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-run:
	@echo "Running Docker container..."
	docker run -d -p 8080:8080 --name $(DOCKER_IMAGE) $(DOCKER_IMAGE)
	@echo "Container started at http://localhost:8080"

docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(DOCKER_IMAGE)
	docker rm $(DOCKER_IMAGE)
	@echo "Container stopped and removed"

# --- MIGRATION HELPERS ---

# 1. Helper: If user types just "make migrate", show options
migrate:
	@echo "⚠️  Command incomplete. Choose one of the following:"
	@echo "  make migrate-create name=my_table   - Create a new migration file"
	@echo "  make migrate-up                     - Run all pending migrations"
	@echo "  make migrate-down                   - Rollback the last migration"

# 2. Safety Check: Ensure 'name' is provided for create
migrate-create:
ifndef name
	@echo "❌ Error: Missing migration name."
	@echo "👉 Usage: make migrate-create name=create_users_table"
	@exit 1
else
	@echo "Creating migration: $(name)..."
	@go run $(MIGRATE_PATH) -create=$(name)
endif

migrate-up:
	@echo "Applying migrations..."
	@go run $(MIGRATE_PATH) -up

migrate-down:
	@echo "Rolling back last migration..."
	@go run $(MIGRATE_PATH) -down

# Help command
help:
	@echo "Available commands:"
	@echo "  make run         - Run the application (development)"
	@echo "  make build       - Build the application"
	@echo "  make start       - Build and run the application"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make migrate     - Show migration help menu"