package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notes-project/internal/models"
	"notes-project/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type BoardService interface {
	Create(ctx context.Context, board *models.Board, ownerID int) error
	GetByID(ctx context.Context, boardID, userID int) (*models.Board, error)
	GetAllForUser(ctx context.Context, userID int) ([]models.Board, error)
	Update(ctx context.Context, boardID, userID int, name string) error
	Delete(ctx context.Context, boardID, userID int) error
	AddMember(ctx context.Context, boardID, inviterID int, inviteeEmail string) error
	InvalidateBoardCache(ctx context.Context, boardID int)
}

type boardService struct {
	repo        repository.BoardRepository
	listRepo    repository.ListRepository
	cardRepo    repository.CardRepository
	userRepo    repository.UserRepository
	broadcaster Broadcaster
	rdb         *redis.Client
}

func NewBoardService(
	repo repository.BoardRepository,
	listRepo repository.ListRepository,
	cardRepo repository.CardRepository,
	userRepo repository.UserRepository,
	broadcaster Broadcaster,
	rdb *redis.Client) BoardService {
	return &boardService{
		repo:        repo,
		listRepo:    listRepo,
		cardRepo:    cardRepo,
		userRepo:    userRepo,
		broadcaster: broadcaster,
		rdb:         rdb,
	}
}

func (s *boardService) InvalidateBoardCache(ctx context.Context, boardID int) {
	cacheKey := fmt.Sprintf("board:%d", boardID)
	if err := s.rdb.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to invalidate cache for board %d: %v", boardID, err)
	} else {
		log.Println("Cache invalidated for board:", boardID)
	}
}

func (s *boardService) Create(ctx context.Context, board *models.Board, ownerID int) error {
	board.OwnerID = ownerID
	if err := s.repo.Create(ctx, board); err != nil {
		return err
	}
	if err := s.repo.AddMember(ctx, board.ID, ownerID); err != nil {
		log.Printf("CRITICAL: could not add owner as member to board %d: %v", board.ID, err)
	}
	return nil
}

func (s *boardService) GetByID(ctx context.Context, boardID, userID int) (*models.Board, error) {
	cacheKey := fmt.Sprintf("board:%d", boardID)
	val, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		log.Println("Cache HIT for board:", boardID)
		var board models.Board
		if json.Unmarshal([]byte(val), &board) == nil {
			hasAccess, _ := s.repo.IsMemberOrOwner(ctx, board.ID, userID)
			if hasAccess {
				return &board, nil
			}
			return nil, fmt.Errorf("access denied")
		}
	}
	log.Println("Cache MISS for board:", boardID)
	hasAccess, err := s.repo.IsMemberOrOwner(ctx, boardID, userID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, fmt.Errorf("access denied")
	}
	board, err := s.repo.GetByID(ctx, boardID)
	if err != nil {
		return nil, err
	}
	lists, err := s.listRepo.GetAllByBoardID(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch lists: %w", err)
	}
	listIDs := make([]int, len(lists))
	for i, list := range lists {
		listIDs[i] = list.ID
	}
	cardsByListID, err := s.cardRepo.GetAllByListIDs(ctx, listIDs)
	if err != nil {
		return nil, fmt.Errorf("could not fetch cards: %w", err)
	}
	for i := range lists {
		if cards, ok := cardsByListID[lists[i].ID]; ok {
			lists[i].Cards = cards
		} else {
			lists[i].Cards = []models.Card{}
		}
	}
	board.Lists = lists
	jsonData, err := json.Marshal(board)
	if err == nil {
		s.rdb.Set(ctx, cacheKey, jsonData, 10*time.Minute)
	}
	return board, nil
}

func (s *boardService) GetAllForUser(ctx context.Context, userID int) ([]models.Board, error) {
	return s.repo.GetAllForUser(ctx, userID)
}

func (s *boardService) Update(ctx context.Context, boardID, userID int, name string) error {
	if err := s.repo.Update(ctx, boardID, userID, name); err != nil {
		return err
	}
	s.InvalidateBoardCache(ctx, boardID)
	return nil
}

func (s *boardService) Delete(ctx context.Context, boardID, userID int) error {
	s.InvalidateBoardCache(ctx, boardID)
	return s.repo.Delete(ctx, boardID, userID)
}

func (s *boardService) AddMember(ctx context.Context, boardID, inviterID int, inviteeEmail string) error {
	board, err := s.repo.GetByID(ctx, boardID)
	if err != nil {
		return fmt.Errorf("board not found")
	}
	if board.OwnerID != inviterID {
		return fmt.Errorf("only the board owner can invite members")
	}

	invitee, err := s.userRepo.GetByEmail(ctx, inviteeEmail)
	if err != nil {
		return fmt.Errorf("user with email %s not found", inviteeEmail)
	}

	err = s.repo.AddMember(ctx, boardID, invitee.ID)
	if err != nil {
		return err
	}

	s.InvalidateBoardCache(ctx, boardID)

	return nil
}
