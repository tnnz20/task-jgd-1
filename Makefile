.PHONY: run build start clean test test-cover test-verbose docker-build docker-run docker-stop

# Application name
APP_NAME=app
APP_PATH=./cmd/http/main.go
BUILD_DIR=bin
DOCKER_IMAGE=category-api
MIGRATE_PATH=./cmd/migrate

# Check OS for cleanup command
ifeq ($(OS),Windows_NT)
    CLEAN_CMD = if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
    MKDIR_CMD = if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
    EXE_EXT = .exe
else
    CLEAN_CMD = rm -rf $(BUILD_DIR)
    MKDIR_CMD = mkdir -p $(BUILD_DIR)
    EXE_EXT = 
endif

# Run the application directly (development)
run:
	@go run $(APP_PATH)

# Build the application
build:
	@echo "Building application..."
	@$(MKDIR_CMD)
	@go build -o $(BUILD_DIR)/$(APP_NAME)$(EXE_EXT) $(APP_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)$(EXE_EXT)"

# Build and run the application
start: build
	@echo "Starting application..."
	./$(BUILD_DIR)/$(APP_NAME)$(EXE_EXT)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@$(CLEAN_CMD)
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

migrate-create:
	@go run $(MIGRATE_PATH) -create $(name)
