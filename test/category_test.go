package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tnnz20/jgd-task-1/internal/config"
	"github.com/tnnz20/jgd-task-1/internal/model"
)

// setupTestServer creates a test server with all dependencies
func setupTestServer() *http.ServeMux {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	app := http.NewServeMux()

	config.Bootstrap(&config.BootstrapConfig{
		App:    app,
		Logger: logger,
	})

	return app
}

func TestHealthEndpoint(t *testing.T) {
	app := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestCreateCategory(t *testing.T) {
	app := setupTestServer()

	t.Run("success", func(t *testing.T) {
		body := `{"name":"Electronics","description":"Electronic devices"}`
		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
		}

		var response model.WebResponse[*model.CategoryResponse]
		err := json.NewDecoder(rec.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Data == nil {
			t.Fatal("Expected data to not be nil")
		}

		if response.Data.Name != "Electronics" {
			t.Errorf("Expected name 'Electronics', got '%s'", response.Data.Name)
		}
	})

	t.Run("bad request - empty name", func(t *testing.T) {
		body := `{"name":"","description":"Test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		body := `invalid json`
		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestListCategories(t *testing.T) {
	app := setupTestServer()

	// Create some categories first
	categories := []string{`{"name":"Category 1"}`, `{"name":"Category 2"}`}
	for _, body := range categories {
		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
	}

	// Test list
	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var response model.WebResponse[[]*model.CategoryResponse]
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(response.Data))
	}
}

func TestGetCategory(t *testing.T) {
	app := setupTestServer()

	// Create a category first
	body := `{"name":"Test Category","description":"Test Description"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	app.ServeHTTP(createRec, createReq)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		var response model.WebResponse[*model.CategoryResponse]
		err := json.NewDecoder(rec.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Data.Name != "Test Category" {
			t.Errorf("Expected name 'Test Category', got '%s'", response.Data.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/categories/999", nil)
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/categories/invalid", nil)
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestUpdateCategory(t *testing.T) {
	app := setupTestServer()

	// Create a category first
	createBody := `{"name":"Original Name"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	app.ServeHTTP(createRec, createReq)

	t.Run("success", func(t *testing.T) {
		body := `{"name":"Updated Name","description":"Updated Description"}`
		req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		var response model.WebResponse[*model.CategoryResponse]
		err := json.NewDecoder(rec.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Data.Name != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got '%s'", response.Data.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		body := `{"name":"Test"}`
		req := httptest.NewRequest(http.MethodPut, "/api/categories/999", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("bad request - empty name", func(t *testing.T) {
		body := `{"name":""}`
		req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestDeleteCategory(t *testing.T) {
	app := setupTestServer()

	// Create a category first
	createBody := `{"name":"To Delete"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	app.ServeHTTP(createRec, createReq)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		// Verify deletion
		getReq := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
		getRec := httptest.NewRecorder()
		app.ServeHTTP(getRec, getReq)

		if getRec.Code != http.StatusNotFound {
			t.Errorf("Expected category to be deleted, got status %d", getRec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/categories/999", nil)
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}
