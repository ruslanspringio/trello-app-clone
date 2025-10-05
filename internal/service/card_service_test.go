// Файл: internal/service/card_service_test.go
package service

import (
	"context"
	"notes-project/internal/models"
	"notes-project/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBroadcaster - простой мок для Broadcaster, если он вам нужен в других тестах.
type MockBroadcaster struct {
	mock.Mock
}

func (m *MockBroadcaster) BroadcastToBoard(boardID int, message []byte) {
	m.Called(boardID, message)
}

func TestCardService_Move(t *testing.T) {
	// --- ARRANGE (Подготовка) ---
	mockCardRepo := new(repository.MockCardRepository)
	mockListRepo := new(repository.MockListRepository)
	mockBoardRepo := new(repository.MockBoardRepository)
	mockBroadcaster := new(MockBroadcaster)
	// Для cacheInvalidator достаточно простой функции-заглушки.
	mockCacheInvalidator := func(ctx context.Context, boardID int) {}

	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBroadcaster, mockCacheInvalidator)

	ctx := context.Background()
	testUserID := 1
	testCardID := 10
	testOldListID := 100
	testNewListID := 101
	testOldBoardID := 1000
	testNewBoardID := 1000

	// Файл: internal/service/card_service_test.go

	// --- Финальная, правильная версия "обучения" ---
	mockCardRepo.On("GetByID", mock.Anything, testCardID).Return(&models.Card{ID: testCardID, ListID: testOldListID}, nil).Once()

	// Явно указываем, что при вызове с testOldListID, нужно вернуть oldList
	mockListRepo.On("GetByID", mock.Anything, testOldListID).Return(&models.List{ID: testOldListID, BoardID: testOldBoardID}, nil).Once()

	// Явно указываем, что при вызове с testNewListID, нужно вернуть newList
	mockListRepo.On("GetByID", mock.Anything, testNewListID).Return(&models.List{ID: testNewListID, BoardID: testNewBoardID}, nil).Once()

	// Остальные ожидания
	mockBoardRepo.On("IsMemberOrOwner", mock.Anything, testOldBoardID, testUserID).Return(true, nil).Once()
	mockBoardRepo.On("IsMemberOrOwner", mock.Anything, testNewBoardID, testUserID).Return(true, nil).Once()
	mockCardRepo.On("Move", mock.Anything, testCardID, testNewListID, 1.0).Return(nil).Once()
	mockBroadcaster.On("BroadcastToBoard", mock.Anything, mock.Anything).Return().Once()

	// --- ACT (Действие) ---
	err := cardService.Move(ctx, testCardID, testNewListID, 1.0, testUserID)

	// --- ASSERT (Проверка) ---
	assert.NoError(t, err)

	// Проверяем, что ВСЕ наши ожидания были выполнены.
	mockCardRepo.AssertExpectations(t)
	mockListRepo.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
	mockBroadcaster.AssertExpectations(t)
}
