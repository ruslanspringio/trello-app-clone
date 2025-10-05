package repository

import (
	"context"
	"database/sql"
	"fmt"
	"notes-project/internal/models"

	"github.com/jmoiron/sqlx"
)

type CardRepository interface {
	Create(ctx context.Context, card *models.Card) error
	GetMaxPositionForList(ctx context.Context, listID int) (float64, error)
	GetAllByListIDs(ctx context.Context, listIDs []int) (map[int][]models.Card, error)
	GetByID(ctx context.Context, cardID int) (*models.Card, error)
	Move(ctx context.Context, cardID, newListID int, newPosition float64) error
}

type cardRepository struct {
	db *sqlx.DB
}

func NewCardRepository(db *sqlx.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) GetByID(ctx context.Context, cardID int) (*models.Card, error) {
	var card models.Card
	query := `SELECT * FROM cards WHERE id=$1`
	if err := r.db.GetContext(ctx, &card, query, cardID); err != nil {
		return nil, fmt.Errorf("cardRepository.GetByID: %w", err)
	}
	return &card, nil
}

func (r cardRepository) Move(ctx context.Context, cardID, newListID int, newPosition float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer tx.Rollback()

	var oldListID int
	var oldPosition float64
	queryGet := `SELECT list_id, "position" FROM cards WHERE id=$1 FOR UPDATE`

	if err := tx.QueryRowContext(ctx, queryGet, cardID).Scan(&oldListID, &oldPosition); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("card with id %d not found", cardID)
		}
		return fmt.Errorf("could not get current card state: %w", err)
	}

	queryShiftOld := `UPDATE cards SET "position" = "position" - 1 WHERE list_id = $1 AND "position" > $2`
	if _, err := tx.ExecContext(ctx, queryShiftOld, oldListID, oldPosition); err != nil {
		return fmt.Errorf("coluld net shift card in old list: %w", err)
	}

	queryShiftNew := `UPDATE cards SET "position" = "position" + 1 WHERE list_id = $1 AND "position" >= $2`
	if _, err := tx.ExecContext(ctx, queryShiftNew, newListID, newPosition); err != nil {
		return fmt.Errorf("could not shift cards in new list: %w", err)
	}

	queryMove := `UPDATE cards SET list_id = $1, "position" = $2, updated_at = NOW() WHERE id = $3`
	if _, err := tx.ExecContext(ctx, queryMove, newListID, newPosition, cardID); err != nil {
		return fmt.Errorf("could not move card: %w", err)
	}

	return tx.Commit()
}

func (r *cardRepository) GetMaxPositionForList(ctx context.Context, listID int) (float64, error) {
	var maxPos float64
	query := `SELECT COALESCE(MAX("position"), 0) FROM cards WHERE list_id=$1`
	err := r.db.GetContext(ctx, &maxPos, query, listID)
	return maxPos, err
}

func (r *cardRepository) Create(ctx context.Context, card *models.Card) error {
	query := `INSERT INTO cards (title, description, "position", list_id) VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, updated_at`
	row := r.db.QueryRowxContext(ctx, query, card.Title, card.Description, card.Position, card.ListID)
	return row.Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)
}

func (r *cardRepository) GetAllByListIDs(ctx context.Context, listIDs []int) (map[int][]models.Card, error) {
	if len(listIDs) == 0 {
		return make(map[int][]models.Card), nil
	}

	query, args, err := sqlx.In(`SELECT * FROM cards WHERE list_id IN (?) ORDER BY "position" ASC`, listIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = r.db.Rebind(query)
	var cards []models.Card
	if err := r.db.SelectContext(ctx, &cards, query, args...); err != nil {
		return nil, fmt.Errorf("failed to select cards: %w", err)
	}

	cardsByListID := make(map[int][]models.Card)
	for _, card := range cards {
		cardsByListID[card.ListID] = append(cardsByListID[card.ListID], card)
	}

	return cardsByListID, nil
}
