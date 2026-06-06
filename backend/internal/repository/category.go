package repository

import (
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/jmoiron/sqlx"
)

// CategoryRepository defines the persistence interface for category operations.
type CategoryRepository interface {
	FindAllByUser(userID int) ([]models.Category, error)
	Create(userID int, name string) (*models.Category, error)
	Delete(userID, categoryID int) (bool, error)
}

// pgCategoryRepository is a PostgreSQL-backed implementation of CategoryRepository.
type pgCategoryRepository struct {
	db *sqlx.DB
}

// NewCategoryRepository initializes and returns a new pgCategoryRepository.
func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &pgCategoryRepository{db: db}
}

// FindAllByUser retrieves all categories belonging to the given user.
func (r *pgCategoryRepository) FindAllByUser(userID int) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Select(&categories,
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
func (r *pgCategoryRepository) Create(userID int, name string) (*models.Category, error) {
	var category models.Category
	err := r.db.QueryRowx(
		"INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id, name",
		userID, name,
	).StructScan(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Delete removes a category by ID if it belongs to the given user.
func (r *pgCategoryRepository) Delete(userID, categoryID int) (bool, error) {
	res, err := r.db.Exec(
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
