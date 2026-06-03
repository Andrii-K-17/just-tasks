package services

import (
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/jmoiron/sqlx"
)

// CategoryService handles category CRUD operations.
type CategoryService struct {
	db *sqlx.DB
}

// NewCategoryService initializes and returns a new CategoryService.
func NewCategoryService(db *sqlx.DB) *CategoryService {
	return &CategoryService{db: db}
}

// GetAll retrieves all categories belonging to the given user.
func (s *CategoryService) GetAll(userID int) ([]models.Category, error) {
	var categories []models.Category
	if err := s.db.Select(&categories,
		"SELECT id, name FROM categories WHERE user_id=$1 ORDER BY id ASC",
		userID,
	); err != nil {
		return nil, err
	}

	if categories == nil {
		categories = []models.Category{}
	}

	return categories, nil
}

// Create inserts a new category for the given user and returns it.
func (s *CategoryService) Create(userID int, name string) (*models.Category, error) {
	var category models.Category
	err := s.db.QueryRowx(
		"INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id, name",
		userID, name,
	).StructScan(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Delete removes a category by ID if it belongs to the given user.
func (s *CategoryService) Delete(userID, categoryID int) (bool, error) {
	res, err := s.db.Exec(
		"DELETE FROM categories WHERE id=$1 AND user_id=$2",
		categoryID, userID,
	)
	if err != nil {
		return false, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}
