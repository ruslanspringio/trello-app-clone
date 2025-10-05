package repository

import (
	"context"
	"fmt"
	"notes-project/internal/models"

	"github.com/jmoiron/sqlx"
)

type BoardRepository interface {
	Create(ctx context.Context, board *models.Board) error
	GetByID(ctx context.Context, boardID int) (*models.Board, error)
	GetAllForUser(ctx context.Context, userID int) ([]models.Board, error)
	Update(ctx context.Context, boardID, ownerID int, name string) error
	Delete(ctx context.Context, boardID, ownerID int) error

	AddMember(ctx context.Context, boardID, userID int) error
	RemoveMember(ctx context.Context, boardID, userID int) error
	IsMemberOrOwner(ctx context.Context, boardID, userID int) (bool, error)
}

type boardRepository struct {
	db *sqlx.DB
}

func NewBoardRepository(db *sqlx.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) Create(ctx context.Context, board *models.Board) error {
	query := `INSERT INTO boards (name, owner_id) VALUES ($1, $2) 
					RETURNING id, created_at, updated_at`
	row := r.db.QueryRowxContext(ctx, query, board.Name, board.OwnerID)
	if err := row.Scan(&board.ID, &board.CreatedAt, &board.UpdatedAt); err != nil {
		return fmt.Errorf("boardRepository.Create: %w", err)
	}
	return nil
}

func (r *boardRepository) GetByID(ctx context.Context, boardID int) (*models.Board, error) {
	var board models.Board
	query := `SELECT * FROM boards WHERE id=$1`
	if err := r.db.GetContext(ctx, &board, query, boardID); err != nil {
		return nil, fmt.Errorf("boardRepository.GetByID: %w", err)
	}
	return &board, nil
}

func (r *boardRepository) GetAllForUser(ctx context.Context, userID int) ([]models.Board, error) {
	var boards []models.Board

	query := `SELECT DISTINCT b.* FROM boards b
			  LEFT JOIN board_members bm ON b.id = bm.board_id
			  WHERE b.owner_id = $1 OR bm.user_id = $1
			  ORDER BY b.updated_at DESC`
	if err := r.db.SelectContext(ctx, &boards, query, userID); err != nil {
		return nil, fmt.Errorf("boardRepository.GetAllForUser: %w", err)
	}
	return boards, nil
}

func (r *boardRepository) Update(ctx context.Context, boardID, ownerID int, name string) error {
	query := `UPDATE boards SET name=$1, updated_at=NOW() WHERE id=$2 AND owner_id=$3`
	result, err := r.db.ExecContext(ctx, query, name, boardID, ownerID)
	if err != nil {
		return fmt.Errorf("boardRepository.Update: %w", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("boardRepository.Update: failed to get rows affected: %w", err)
	}
	if rowAffected == 0 {
		return fmt.Errorf("board not found or user is not the owner")
	}
	return nil
}

func (r *boardRepository) Delete(ctx context.Context, boardID, ownerID int) error {
	query := `DELETE FROM boards WHERE id=$1 AND owner_id=$2`
	result, err := r.db.ExecContext(ctx, query, boardID, ownerID)
	if err != nil {
		return fmt.Errorf("boardRepository.Delete: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("boardRepository.Delete: failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("board not found or user is not the owner")
	}
	return nil
}

func (r *boardRepository) AddMember(ctx context.Context, boardID, userID int) error {
	query := `INSERT INTO board_members (board_id, user_id) VALUES ($1, $2)
			  ON CONFLICT (user_id, board_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, boardID, userID)
	return err
}

func (r *boardRepository) RemoveMember(ctx context.Context, boardID, userID int) error {
	query := `DELETE FROM board_members WHERE board_id=$1 AND user_id=$2`
	_, err := r.db.ExecContext(ctx, query, boardID, userID)
	return err
}

func (r *boardRepository) IsMemberOrOwner(ctx context.Context, boardID, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
				SELECT 1 FROM boards WHERE id=$1 AND owner_id=$2
				UNION ALL
				SELECT 1 FROM board_members WHERE board_id=$1 AND user_id=$2
			  )`
	err := r.db.GetContext(ctx, &exists, query, boardID, userID)
	if err != nil {
		return false, fmt.Errorf("IsMemberOrOwner check failed: %w", err)
	}
	return exists, nil
}
