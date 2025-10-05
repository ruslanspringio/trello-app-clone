package service

import (
	"context"
	"encoding/json"
	"fmt"
	"notes-project/internal/models"
	"notes-project/internal/repository"
)

type CardService interface {
	Create(ctx context.Context, card *models.Card, listID, userID int) error
	Move(ctx context.Context, cardID, newListID int, newPosition float64, userID int) error
}

type cardService struct {
	cardRepo             repository.CardRepository
	listRepo             repository.ListRepository
	boardRepo            repository.BoardRepository
	broadcaster          Broadcaster
	invalidateBoardCache CacheInvalidator
}

func NewCardService(
	cardRepo repository.CardRepository,
	listRepo repository.ListRepository,
	boardRepo repository.BoardRepository,
	broadcaster Broadcaster,
	cacheInvalidator CacheInvalidator) CardService {
	return &cardService{
		cardRepo:             cardRepo,
		listRepo:             listRepo,
		boardRepo:            boardRepo,
		broadcaster:          broadcaster,
		invalidateBoardCache: cacheInvalidator}
}

func (s *cardService) Create(ctx context.Context, card *models.Card, listID, userID int) error {
	list, err := s.listRepo.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("list with id %d not found", listID)
	}

	if err := s.cardRepo.Create(ctx, card); err != nil {
		return err
	}

	s.invalidateBoardCache(ctx, list.BoardID)

	wsMessage := models.WebSocketMessage{Event: "CARD_CREATED", Payload: card}
	jsonMessage, _ := json.Marshal(wsMessage)
	s.broadcaster.BroadcastToBoard(list.BoardID, jsonMessage)
	return nil
}

func (s *cardService) Move(ctx context.Context, cardID, newListID int, newPosition float64, userID int) error {
	card, _ := s.cardRepo.GetByID(ctx, cardID)
	oldList, _ := s.listRepo.GetByID(ctx, card.ListID)
	newList, _ := s.listRepo.GetByID(ctx, newListID)

	if err := s.cardRepo.Move(ctx, cardID, newListID, newPosition); err != nil {
		return err
	}

	s.invalidateBoardCache(ctx, oldList.BoardID)
	if oldList.BoardID != newList.BoardID {
		s.invalidateBoardCache(ctx, newList.BoardID)
	}

	return nil
}
