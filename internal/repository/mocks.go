// Файл: internal/repository/mocks.go
package repository

import (
	"context"
	"notes-project/internal/models"

	"github.com/stretchr/testify/mock"
)

// --- MockCardRepository ---
type MockCardRepository struct {
	mock.Mock
}

func (m *MockCardRepository) Create(ctx context.Context, card *models.Card) error {
	args := m.Called(ctx, card)
	return args.Error(0)
}

func (m *MockCardRepository) GetMaxPositionForList(ctx context.Context, listID int) (float64, error) {
	args := m.Called(ctx, listID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCardRepository) GetAllByListIDs(ctx context.Context, listIDs []int) (map[int][]models.Card, error) {
	args := m.Called(ctx, listIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int][]models.Card), args.Error(1)
}

func (m *MockCardRepository) GetByID(ctx context.Context, cardID int) (*models.Card, error) {
	args := m.Called(ctx, cardID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Card), args.Error(1)
}

func (m *MockCardRepository) Move(ctx context.Context, cardID, newListID int, newPosition float64) error {
	args := m.Called(ctx, cardID, newListID, newPosition)
	return args.Error(0)
}

// --- MockListRepository ---
type MockListRepository struct {
	mock.Mock
}

func (m *MockListRepository) Create(ctx context.Context, list *models.List) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockListRepository) GetByID(ctx context.Context, listID int) (*models.List, error) {
	args := m.Called(ctx, listID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.List), args.Error(1)
}

func (m *MockListRepository) GetAllByBoardID(ctx context.Context, boardID int) ([]models.List, error) {
	args := m.Called(ctx, boardID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.List), args.Error(1)
}

func (m *MockListRepository) Update(ctx context.Context, list *models.List) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockListRepository) Delete(ctx context.Context, listID int) error {
	args := m.Called(ctx, listID)
	return args.Error(0)
}

func (m *MockListRepository) GetMaxPositionForBoard(ctx context.Context, boardID int) (float64, error) {
	args := m.Called(ctx, boardID)
	return args.Get(0).(float64), args.Error(1)
}

// --- MockBoardRepository ---
type MockBoardRepository struct {
	mock.Mock
}

func (m *MockBoardRepository) Create(ctx context.Context, board *models.Board) error {
	args := m.Called(ctx, board)
	return args.Error(0)
}

func (m *MockBoardRepository) GetByID(ctx context.Context, boardID int) (*models.Board, error) {
	args := m.Called(ctx, boardID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Board), args.Error(1)
}

func (m *MockBoardRepository) GetAllForUser(ctx context.Context, userID int) ([]models.Board, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Board), args.Error(1)
}

func (m *MockBoardRepository) Update(ctx context.Context, boardID, ownerID int, name string) error {
	args := m.Called(ctx, boardID, ownerID, name)
	return args.Error(0)
}

func (m *MockBoardRepository) Delete(ctx context.Context, boardID, ownerID int) error {
	args := m.Called(ctx, boardID, ownerID)
	return args.Error(0)
}

func (m *MockBoardRepository) AddMember(ctx context.Context, boardID, userID int) error {
	args := m.Called(ctx, boardID, userID)
	return args.Error(0)
}

func (m *MockBoardRepository) RemoveMember(ctx context.Context, boardID, userID int) error {
	args := m.Called(ctx, boardID, userID)
	return args.Error(0)
}

func (m *MockBoardRepository) IsMemberOrOwner(ctx context.Context, boardID, userID int) (bool, error) {
	args := m.Called(ctx, boardID, userID)
	return args.Bool(0), args.Error(1)
}
