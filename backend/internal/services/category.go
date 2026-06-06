package services

import (
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/Andrii-K-17/just-tasks/internal/repository"
)

// CategoryService handles category CRUD operations.
type CategoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService initializes and returns a new CategoryService.
func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAll retrieves all categories belonging to the given user.
func (s *CategoryService) GetAll(userID int) ([]models.Category, error) {
	return s.repo.FindAllByUser(userID)
}

// Create inserts a new category for the given user and returns it.
func (s *CategoryService) Create(userID int, name string) (*models.Category, error) {
	return s.repo.Create(userID, name)
}

// Delete removes a category by ID if it belongs to the given user.
func (s *CategoryService) Delete(userID, categoryID int) (bool, error) {
	return s.repo.Delete(userID, categoryID)
}
