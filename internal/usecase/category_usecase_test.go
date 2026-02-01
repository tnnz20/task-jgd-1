package usecase

import (
	"io"
	"log/slog"
	"testing"

	"github.com/tnnz20/jgd-task-1/internal/model"
	"github.com/tnnz20/jgd-task-1/internal/repository/memory"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestNewCategoryUseCase(t *testing.T) {
	repo := memory.NewCategoryRepository()
	logger := newTestLogger()

	useCase := NewCategoryUseCase(repo, logger)

	if useCase == nil {
		t.Fatal("Expected useCase to not be nil")
	}

	if useCase.CategoryRepository == nil {
		t.Fatal("Expected CategoryRepository to be set")
	}

	if useCase.Log == nil {
		t.Fatal("Expected Log to be set")
	}
}

func TestCategoryUseCaseCreate(t *testing.T) {
	repo := memory.NewCategoryRepository()
	logger := newTestLogger()
	useCase := NewCategoryUseCase(repo, logger)

	t.Run("success", func(t *testing.T) {
		testCreateSuccess(t, useCase)
	})

	t.Run("empty name", func(t *testing.T) {
		testCreateEmptyName(t, useCase)
	})
}

func testCreateSuccess(t *testing.T, useCase *CategoryUseCase) {
	request := &model.CreateCategoryRequest{
		Name:        "Test Category",
		Description: "Test Description",
	}

	response, err := useCase.Create(request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response to not be nil")
	}

	if response.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", response.ID)
	}

	if response.Name != "Test Category" {
		t.Errorf("Expected Name to be 'Test Category', got '%s'", response.Name)
	}

	if response.Description != "Test Description" {
		t.Errorf("Expected Description to be 'Test Description', got '%s'", response.Description)
	}
}

func testCreateEmptyName(t *testing.T, useCase *CategoryUseCase) {
	request := &model.CreateCategoryRequest{
		Name:        "",
		Description: "Test Description",
	}

	response, err := useCase.Create(request)
	if err == nil {
		t.Fatal("Expected error for empty name")
	}

	if err != ErrBadRequest {
		t.Errorf("Expected ErrBadRequest, got %v", err)
	}

	if response != nil {
		t.Error("Expected response to be nil")
	}
}

func TestCategoryUseCaseGet(t *testing.T) {
	repo := memory.NewCategoryRepository()
	logger := newTestLogger()
	useCase := NewCategoryUseCase(repo, logger)

	// Create a category first
	createReq := &model.CreateCategoryRequest{
		Name:        "Test Category",
		Description: "Test Description",
	}
	_, _ = useCase.Create(createReq)

	t.Run("success", func(t *testing.T) {
		request := &model.GetCategoryRequest{ID: 1}

		response, err := useCase.Get(request)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if response == nil {
			t.Fatal("Expected response to not be nil")
		}

		if response.ID != 1 {
			t.Errorf("Expected ID to be 1, got %d", response.ID)
		}

		if response.Name != "Test Category" {
			t.Errorf("Expected Name to be 'Test Category', got '%s'", response.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		request := &model.GetCategoryRequest{ID: 999}

		response, err := useCase.Get(request)
		if err == nil {
			t.Fatal("Expected error for non-existing category")
		}

		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}

		if response != nil {
			t.Error("Expected response to be nil")
		}
	})
}

func TestCategoryUseCaseList(t *testing.T) {
	repo := memory.NewCategoryRepository()
	logger := newTestLogger()
	useCase := NewCategoryUseCase(repo, logger)

	t.Run("empty list", func(t *testing.T) {
		responses, err := useCase.List()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(responses) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(responses))
		}
	})

	t.Run("with categories", func(t *testing.T) {
		// Create categories
		_, _ = useCase.Create(&model.CreateCategoryRequest{Name: "Category 1"})
		_, _ = useCase.Create(&model.CreateCategoryRequest{Name: "Category 2"})
		_, _ = useCase.Create(&model.CreateCategoryRequest{Name: "Category 3"})

		responses, err := useCase.List()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(responses) != 3 {
			t.Errorf("Expected 3 categories, got %d", len(responses))
		}
	})
}

func TestCategoryUseCaseUpdate(t *testing.T) {
	repo := memory.NewCategoryRepository()
	logger := newTestLogger()
	useCase := NewCategoryUseCase(repo, logger)

	// Create a category first
	_, _ = useCase.Create(&model.CreateCategoryRequest{
		Name:        "Original Name",
		Description: "Original Description",
	})

	t.Run("success", func(t *testing.T) {
		testUpdateSuccess(t, useCase)
	})

	t.Run("empty name", func(t *testing.T) {
		testUpdateEmptyName(t, useCase)
	})

	t.Run("not found", func(t *testing.T) {
		testUpdateNotFound(t, useCase)
	})
}

func testUpdateSuccess(t *testing.T, useCase *CategoryUseCase) {
	request := &model.UpdateCategoryRequest{
		ID:          1,
		Name:        "Updated Name",
		Description: "Updated Description",
	}

	response, err := useCase.Update(request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response to not be nil")
	}

	if response.Name != "Updated Name" {
		t.Errorf("Expected Name to be 'Updated Name', got '%s'", response.Name)
	}

	if response.Description != "Updated Description" {
		t.Errorf("Expected Description to be 'Updated Description', got '%s'", response.Description)
	}
}

func testUpdateEmptyName(t *testing.T, useCase *CategoryUseCase) {
	request := &model.UpdateCategoryRequest{
		ID:          1,
		Name:        "",
		Description: "Test",
	}

	response, err := useCase.Update(request)
	if err == nil {
		t.Fatal("Expected error for empty name")
	}

	if err != ErrBadRequest {
		t.Errorf("Expected ErrBadRequest, got %v", err)
	}

	if response != nil {
		t.Error("Expected response to be nil")
	}
}

func testUpdateNotFound(t *testing.T, useCase *CategoryUseCase) {
	request := &model.UpdateCategoryRequest{
		ID:   999,
		Name: "Test",
	}

	response, err := useCase.Update(request)
	if err == nil {
		t.Fatal("Expected error for non-existing category")
	}

	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	if response != nil {
		t.Error("Expected response to be nil")
	}
}

func TestCategoryUseCaseDelete(t *testing.T) {
	repo := memory.NewCategoryRepository()
	logger := newTestLogger()
	useCase := NewCategoryUseCase(repo, logger)

	// Create a category first
	_, _ = useCase.Create(&model.CreateCategoryRequest{Name: "Test Category"})

	t.Run("success", func(t *testing.T) {
		request := &model.DeleteCategoryRequest{ID: 1}

		err := useCase.Delete(request)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify deletion
		_, err = useCase.Get(&model.GetCategoryRequest{ID: 1})
		if err != ErrNotFound {
			t.Error("Expected category to be deleted")
		}
	})

	t.Run("not found", func(t *testing.T) {
		request := &model.DeleteCategoryRequest{ID: 999}

		err := useCase.Delete(request)
		if err == nil {
			t.Fatal("Expected error for non-existing category")
		}

		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})
}
