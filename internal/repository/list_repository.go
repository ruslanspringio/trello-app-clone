package repository

import (
	"context"
	"fmt"
	"notes-project/internal/models"

	"github.com/jmoiron/sqlx"
)

type ListRepository interface {
	Create(ctx context.Context, list *models.List) error
	GetByID(ctx context.Context, listID int) (*models.List, error)
	GetAllByBoardID(ctx context.Context, boardID int) ([]models.List, error)
	Update(ctx context.Context, list *models.List) error
	Delete(ctx context.Context, listID int) error
	GetMaxPositionForBoard(ctx context.Context, boardID int) (float64, error)
}

type listRepository struct {
	db *sqlx.DB
}

func NewListRepository(db *sqlx.DB) ListRepository {
	return &listRepository{db: db}
}

func (r *listRepository) GetMaxPositionForBoard(ctx context.Context, boardID int) (float64, error) {
	var maxPos float64
	query := `SELECT COALESCE(MAX("position"), 0) FROM lists WHERE board_id=$1`
	err := r.db.GetContext(ctx, &maxPos, query, boardID)
	return maxPos, err
}

func (r *listRepository) Create(ctx context.Context, list *models.List) error {
	query := `INSERT INTO lists (title, "position", board_id) VALUES ($1, $2, $3) 
							RETURNING id, created_at, updated_at`
	row := r.db.QueryRowxContext(ctx, query, list.Title, list.Position, list.BoardID)
	return row.Scan(&list.ID, &list.CreatedAt, &list.UpdatedAt)
}

func (r *listRepository) GetByID(ctx context.Context, listID int) (*models.List, error) {
	var list models.List
	query := `SELECT * FROM lists WHERE id = $1`
	if err := r.db.GetContext(ctx, &list, query, listID); err != nil {
		return nil, fmt.Errorf("listRepository.GetByID: %w", err)
	}
	return &list, nil
}

func (r *listRepository) GetAllByBoardID(ctx context.Context, boardID int) ([]models.List, error) {
	var lists []models.List
	query := `SELECT * FROM lists WHERE board_id=$1 ORDER BY "position" ASC`
	if err := r.db.SelectContext(ctx, &lists, query, boardID); err != nil {
		return nil, fmt.Errorf("listRepository.GetAllByBoardID: %w", err)
	}
	return lists, nil
}

func (r *listRepository) Update(ctx context.Context, list *models.List) error {
	query := `UPDATE lists SET title=$1, "position"=$2, updated_at=NOW() WHERE id=$3`
	_, err := r.db.ExecContext(ctx, query, list.Title, list.Position, list.ID)
	return err
}

func (r *listRepository) Delete(ctx context.Context, listID int) error {
	query := `DELETE FROM lists WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, listID)
	return err
}
