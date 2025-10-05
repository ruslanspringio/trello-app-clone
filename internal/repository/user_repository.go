package repository

import (
	"context"
	"fmt"
	"notes-project/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (name, age, email, password_hash)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, updated_at`

	row := r.db.QueryRowxContext(ctx, query, user.Name, user.Age, user.Email, user.PasswordHash)
	if err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return fmt.Errorf("userRepository.Create: %w", err)
	}
	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE email=$1"
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, fmt.Errorf("userRepository.GetByEmail: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := "SELECT * FROM users ORDER BY id DESC"
	err := r.db.SelectContext(ctx, &users, query)
	return users, err
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id=$1"
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users
			  SET name=$1, age=$2, updated_at=NOW()
			  WHERE id=$3
			  RETURNING updated_at`

	row := r.db.QueryRowxContext(ctx, query, user.Name, user.Age, user.ID)
	return row.Scan(&user.UpdatedAt)
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id=$1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
