package service

import (
	"context"
	"encoding/json"
	"fmt"
	"notes-project/internal/models"
	"notes-project/internal/repository"
)

type ListService interface {
	Create(ctx context.Context, list *models.List, boardID, userID int) error
}

type CacheInvalidator func(ctx context.Context, boardID int)

type listService struct {
	listRepo             repository.ListRepository
	boardRepo            repository.BoardRepository
	broadcaster          Broadcaster
	invalidateBoardCache CacheInvalidator
}

func NewListService(
	listRepo repository.ListRepository,
	boardRepo repository.BoardRepository,
	broadcaster Broadcaster,
	cacheInvalidator CacheInvalidator) ListService {
	return &listService{
		listRepo:             listRepo,
		boardRepo:            boardRepo,
		broadcaster:          broadcaster,
		invalidateBoardCache: cacheInvalidator}
}

func (s *listService) checkBoardPermissions(ctx context.Context, boardID, userID int) error {
	hasAccess, err := s.boardRepo.IsMemberOrOwner(ctx, boardID, userID)
	if err != nil {
		return fmt.Errorf("could not verify board permissions: %w", err)
	}
	if !hasAccess {
		return fmt.Errorf("access denied to this board")
	}
	return nil
}

func (s *listService) Create(ctx context.Context, list *models.List, boardID, userID int) error {
	if err := s.checkBoardPermissions(ctx, boardID, userID); err != nil {
		return err
	}
	maxPos, err := s.listRepo.GetMaxPositionForBoard(ctx, boardID)
	if err != nil {
		return fmt.Errorf("could not determine list position: %w", err)
	}
	list.Position = maxPos + 1.0
	list.BoardID = boardID
	if err := s.listRepo.Create(ctx, list); err != nil {
		return err
	}

	s.invalidateBoardCache(ctx, boardID)

	wsMessage := models.WebSocketMessage{Event: "LIST_CREATED", Payload: list}
	jsonMessage, _ := json.Marshal(wsMessage)
	s.broadcaster.BroadcastToBoard(boardID, jsonMessage)
	return nil
}
