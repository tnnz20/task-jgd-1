# Category & Product API - Go Clean Architecture

A RESTful CRUD API built with Go using Clean Architecture principles. This project demonstrates best practices for structuring Go applications with separation of concerns, dependency injection, and testability. Includes multi-entity support with categories and products.

## Tech Stack

- **Go 1.25+** - Programming language
- **net/http** - Standard library HTTP server (Go 1.22+ enhanced routing)
- **log/slog** - Structured logging (Go 1.21+)
- **testing** - Standard library testing

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        HTTP Request                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Delivery Layer                               │
│                  (HTTP Controllers)                              │
│         - Handles HTTP requests/responses                        │
│         - Input validation                                       │
│         - Route handling                                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Use Case Layer                               │
│                   (Business Logic)                               │
│         - Business rules                                         │
│         - Orchestrates data flow                                 │
│         - Domain validation                                      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Repository Layer                              │
│                    (Data Access)                                 │
│         - Data persistence                                       │
│         - In-memory storage                                      │
└─────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
task-1/
├── cmd/
│   └── http/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   ├── app.go                     # Bootstrap & dependency injection
│   │   └── logger.go                  # Structured logger configuration
│   ├── delivery/
│   │   └── http/
│   │       ├── category_controller.go # Category HTTP handlers
│   │       ├── product_controller.go  # Product HTTP handlers
│   │       ├── helper.go              # Shared HTTP utilities
│   │       └── route/
│   │           └── route.go           # Route definitions
│   ├── entity/
│   │   ├── category_entity.go         # Category domain entity
│   │   └── product_entity.go          # Product domain entity
│   ├── model/
│   │   ├── model.go                   # Generic response wrapper
│   │   ├── category_model.go          # Category Request/Response DTOs
│   │   ├── product_model.go           # Product Request/Response DTOs
│   │   └── converter/
│   │       ├── category_converter.go  # Category Entity ↔ Model converters
│   │       └── product_converter.go   # Product Entity ↔ Model converters
│   ├── repository/
│   │   ├── interface.go               # Repository interfaces
│   │   ├── memory/
│   │   │   ├── category.go            # In-memory category repository
│   │   │   └── product.go             # In-memory product repository
│   │   └── postgres/
│   │       ├── category.go            # PostgreSQL category repository
│   │       └── product.go             # PostgreSQL product repository
│   └── usecase/
│       ├── category_usecase.go        # Category business logic
│       ├── category_usecase_test.go   # Category usecase tests
│       ├── product_usecase.go         # Product business logic
│       └── product_usecase_test.go    # Product usecase tests
├── test/
│   └── category_test.go               # Integration tests
├── go.mod
├── Makefile
└── README.md
```

## Features

- ✅ Clean Architecture pattern
- ✅ RESTful API design
- ✅ Graceful shutdown
- ✅ Structured JSON logging (slog)
- ✅ In-memory data storage
- ✅ Thread-safe operations
- ✅ Comprehensive unit tests
- ✅ Integration tests
- ✅ Go 1.22+ enhanced routing

## API Endpoints

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Check server health status |

**Response:**
```json
{"status":"healthy"}
```

### Categories

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/categories` | Create a new category |
| GET | `/api/categories` | Get all categories |
| GET | `/api/categories/{id}` | Get category by ID |
| PUT | `/api/categories/{id}` | Update category by ID |
| DELETE | `/api/categories/{id}` | Delete category by ID |

### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/products` | Create a new product |
| GET | `/api/products` | Get all products |
| GET | `/api/products/{id}` | Get product by ID |
| PUT | `/api/products/{id}` | Update product by ID |
| DELETE | `/api/products/{id}` | Delete product by ID |

#### Create Category

**Request:**
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and gadgets"
  }'
```

**Response (201 Created):**
```json
{
  "data": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and gadgets",
    "created_at": 1737783600000,
    "updated_at": 1737783600000
  }
}
```

#### Get All Categories

**Request:**
```bash
curl http://localhost:8080/api/categories
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Electronics",
      "description": "Electronic devices and gadgets",
      "created_at": 1737783600000,
      "updated_at": 1737783600000
    }
  ]
}
```

#### Get Category by ID

**Request:**
```bash
curl http://localhost:8080/api/categories/1
```

**Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and gadgets",
    "created_at": 1737783600000,
    "updated_at": 1737783600000
  }
}
```

**Response (404 Not Found):**
```json
{
  "errors": "Category not found"
}
```

#### Update Category

**Request:**
```bash
curl -X PUT http://localhost:8080/api/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Electronics",
    "description": "Updated description"
  }'
```

**Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "name": "Updated Electronics",
    "description": "Updated description",
    "created_at": 1737783600000,
    "updated_at": 1737783660000
  }
}
```

#### Delete Category

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/categories/1
```

**Response (200 OK):**
```json
{
  "data": true
}
```

## Products API Examples

#### Create Product

**Request:**
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price": 999.99,
    "stock": 50,
    "category_id": 1
  }'
```

**Response (201 Created):**
```json
{
  "data": {
    "id": 1,
    "name": "Laptop",
    "price": 999.99,
    "stock": 50,
    "category_id": 1,
    "created_at": 1737783600000,
    "updated_at": 1737783600000
  }
}
```

#### Get All Products

**Request:**
```bash
curl http://localhost:8080/api/products
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Laptop",
      "price": 999.99,
      "stock": 50,
      "category_id": 1,
      "created_at": 1737783600000,
      "updated_at": 1737783600000
    }
  ]
}
```

#### Get Product by ID

**Request:**
```bash
curl http://localhost:8080/api/products/1
```

**Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "name": "Laptop",
    "price": 999.99,
    "stock": 50,
    "category_id": 1,
    "created_at": 1737783600000,
    "updated_at": 1737783600000
  }
}
```

**Response (404 Not Found):**
```json
{
  "errors": "Product not found"
}
```

#### Update Product

**Request:**
```bash
curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Laptop",
    "price": 1099.99,
    "stock": 45,
    "category_id": 1
  }'
```

**Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "name": "Updated Laptop",
    "price": 1099.99,
    "stock": 45,
    "category_id": 1,
    "created_at": 1737783600000,
    "updated_at": 1737783660000
  }
}
```

#### Delete Product

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/products/1
```

**Response (200 OK):**
```json
{
  "data": true
}
```

## Getting Started

### Prerequisites

- Go 1.22 or higher

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd task-1
```

2. Download dependencies:
```bash
go mod tidy
```

### Running the Application

```bash
# Run directly
go run ./cmd/http

# Or build and run
go build -o app ./cmd/http
./app
```

The server will start at `http://localhost:8080`

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `LOG_LEVEL` | Logging level (DEBUG, INFO, WARN, ERROR) | `INFO` |

**Example:**
```bash
PORT=3000 LOG_LEVEL=DEBUG go run ./cmd/http
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with coverage
go test ./... -cover

# Run with coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Error Responses

All error responses follow this format:

```json
{
  "errors": "Error message here"
}
```

| Status Code | Description |
|-------------|-------------|
| 400 | Bad Request - Invalid input |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error |

## Graceful Shutdown

The application supports graceful shutdown:

- Listens for `SIGINT` (Ctrl+C) and `SIGTERM` signals
- Waits up to 30 seconds for active connections to complete
- Logs shutdown progress

## License

This project is for educational purposes.
