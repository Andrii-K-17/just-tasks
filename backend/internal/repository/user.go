package repository

import (
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/jmoiron/sqlx"
)

// UserRepository defines the persistence interface for user operations.
type UserRepository interface {
	ExistsByUsername(username string) (bool, error)
	Create(username, passwordHash string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByID(id int) (*models.User, error)
	Delete(id int) error
}

// pgUserRepository is a PostgreSQL-backed implementation of UserRepository.
type pgUserRepository struct {
	db *sqlx.DB
}

// NewUserRepository initializes and returns a new pgUserRepository.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &pgUserRepository{db: db}
}

// ExistsByUsername checks whether a user with the given username already exists.
func (r *pgUserRepository) ExistsByUsername(username string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists,
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username,
	)
	return exists, err
}

// Create inserts a new user record and returns the created user.
func (r *pgUserRepository) Create(username, passwordHash string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowx(
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)
		 RETURNING id, username, created_at`,
		username, passwordHash,
	).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername retrieves a user by their username.
func (r *pgUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user,
		"SELECT id, username, password_hash FROM users WHERE username=$1", username,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID retrieves a user by their primary key.
func (r *pgUserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user,
		"SELECT id, username FROM users WHERE id=$1", id,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete removes a user record by their primary key.
func (r *pgUserRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}
