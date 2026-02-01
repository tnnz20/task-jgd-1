package memory

import (
	"testing"

	"github.com/tnnz20/jgd-task-1/internal/entity"
)

func TestNewCategoryRepository(t *testing.T) {
	repo := NewCategoryRepository()

	if repo == nil {
		t.Fatal("Expected repository to not be nil")
	}

	if repo.categories == nil {
		t.Fatal("Expected categories slice to be initialized")
	}

	if len(repo.categories) != 0 {
		t.Errorf("Expected empty categories slice, got %d", len(repo.categories))
	}
}

func TestCategoryRepositoryCreate(t *testing.T) {
	repo := NewCategoryRepository()

	category := &entity.Category{
		Name:        "Test Category",
		Description: "Test Description",
	}

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if category.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", category.ID)
	}

	if category.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if category.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	// Create another category
	category2 := &entity.Category{
		Name:        "Test Category 2",
		Description: "Test Description 2",
	}

	err = repo.Create(category2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if category2.ID != 2 {
		t.Errorf("Expected ID to be 2, got %d", category2.ID)
	}
}

func TestCategoryRepositoryFindById(t *testing.T) {
	repo := NewCategoryRepository()

	// Create a category first
	original := &entity.Category{
		Name:        "Test Category",
		Description: "Test Description",
	}
	_ = repo.Create(original)

	// Test finding existing category
	found := new(entity.Category)
	err := repo.FindById(found, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if found.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", found.ID)
	}

	if found.Name != "Test Category" {
		t.Errorf("Expected Name to be 'Test Category', got '%s'", found.Name)
	}

	// Test finding non-existing category
	notFound := new(entity.Category)
	err = repo.FindById(notFound, 999)
	if err == nil {
		t.Error("Expected error for non-existing category")
	}

	if err != ErrCategoryNotFound {
		t.Errorf("Expected ErrCategoryNotFound, got %v", err)
	}
}

func TestCategoryRepositoryFindAll(t *testing.T) {
	repo := NewCategoryRepository()

	// Test empty repository
	categories, err := repo.FindAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(categories) != 0 {
		t.Errorf("Expected 0 categories, got %d", len(categories))
	}

	// Add some categories
	_ = repo.Create(&entity.Category{Name: "Category 1"})
	_ = repo.Create(&entity.Category{Name: "Category 2"})
	_ = repo.Create(&entity.Category{Name: "Category 3"})

	categories, err = repo.FindAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(categories) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(categories))
	}
}

func TestCategoryRepositoryUpdate(t *testing.T) {
	repo := NewCategoryRepository()

	// Create a category first
	original := &entity.Category{
		Name:        "Original Name",
		Description: "Original Description",
	}
	_ = repo.Create(original)

	// Update the category
	updated := &entity.Category{
		ID:          1,
		Name:        "Updated Name",
		Description: "Updated Description",
	}

	err := repo.Update(updated)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the update
	found := new(entity.Category)
	_ = repo.FindById(found, 1)

	if found.Name != "Updated Name" {
		t.Errorf("Expected Name to be 'Updated Name', got '%s'", found.Name)
	}

	if found.Description != "Updated Description" {
		t.Errorf("Expected Description to be 'Updated Description', got '%s'", found.Description)
	}

	// Test updating non-existing category
	nonExisting := &entity.Category{
		ID:   999,
		Name: "Non-existing",
	}

	err = repo.Update(nonExisting)
	if err == nil {
		t.Error("Expected error for non-existing category")
	}

	if err != ErrCategoryNotFound {
		t.Errorf("Expected ErrCategoryNotFound, got %v", err)
	}
}

func TestCategoryRepositoryDelete(t *testing.T) {
	repo := NewCategoryRepository()

	// Create categories
	cat1 := &entity.Category{Name: "Category 1"}
	cat2 := &entity.Category{Name: "Category 2"}
	cat3 := &entity.Category{Name: "Category 3"}

	_ = repo.Create(cat1)
	_ = repo.Create(cat2)
	_ = repo.Create(cat3)

	// Delete the second category
	err := repo.Delete(cat2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify deletion
	categories, _ := repo.FindAll()
	if len(categories) != 2 {
		t.Errorf("Expected 2 categories after deletion, got %d", len(categories))
	}

	// Verify cat2 is not found
	notFound := new(entity.Category)
	err = repo.FindById(notFound, 2)
	if err == nil {
		t.Error("Expected error when finding deleted category")
	}

	// Test deleting non-existing category
	nonExisting := &entity.Category{ID: 999}
	err = repo.Delete(nonExisting)
	if err == nil {
		t.Error("Expected error for non-existing category")
	}

	if err != ErrCategoryNotFound {
		t.Errorf("Expected ErrCategoryNotFound, got %v", err)
	}
}

func TestCategoryRepositoryCountById(t *testing.T) {
	repo := NewCategoryRepository()

	// Test count for non-existing category
	count, err := repo.CountById(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count to be 0, got %d", count)
	}

	// Create a category
	_ = repo.Create(&entity.Category{Name: "Test"})

	// Test count for existing category
	count, err = repo.CountById(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count to be 1, got %d", count)
	}
}
