package handlers

import (
	"net/http"
	"notes-project/internal/models"
	"notes-project/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type BoardHandler struct {
	service service.BoardService
}

type AddMemberInput struct {
	Email string `json:"email" binding:"required"`
}

func NewBoardHandler(s service.BoardService) *BoardHandler {
	return &BoardHandler{service: s}
}

func (h *BoardHandler) RegisterBoardRoutes(rg *gin.RouterGroup) {
	boards := rg.Group("boards")
	{
		boards.POST("/", h.CreateBoard)
		boards.GET("/", h.GetAllBoardsForUser)
		boards.GET("/:boardId", h.GetBoardByID)
		boards.PUT("/:boardId", h.UpdateBoard)
		boards.DELETE("/:boardId", h.DeleteBoard)

		boards.POST("/:boardId/members", h.AddMemberToBoard)
	}
}

// @Summary      Создать новую доску
// @Description  Создает новую доску для авторизованного пользователя.
// @Tags         Boards
// @Accept       json
// @Produce      json
// @Param        board  body      models.Board  true  "Данные для создания доски (нужно только поле 'name')"
// @Success      201    {object}  models.Board
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Security     ApiKeyAuth
// @Router       /boards [post]

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	var input models.Board
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// Теперь мы уверены, что userID не nil.
	if err := h.service.Create(c.Request.Context(), &input, userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create board"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// @Summary      Получить все доски пользователя
// @Description  Возвращает список всех досок, где пользователь является владельцем или участником.
// @Tags         Boards
// @Produce      json
// @Success      200 {array}   models.Board
// @Failure      401 {object}  ErrorResponse
// @Security     ApiKeyAuth
// @Router       /boards [get]
func (h *BoardHandler) GetAllBoardsForUser(c *gin.Context) {
	userID, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	boards, err := h.service.GetAllForUser(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch boards"})
		return
	}

	c.JSON(http.StatusOK, boards)
}

// @Summary      Получить доску по ID
// @Description  Возвращает полную информацию о доске, включая списки и карточки.
// @Tags         Boards
// @Produce      json
// @Param        boardId  path      int  true  "ID Доски"
// @Success      200      {object}  models.Board
// @Failure      401      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Security     ApiKeyAuth
// @Router       /boards/{boardId} [get]
func (h *BoardHandler) GetBoardByID(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	boardID, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	board, err := h.service.GetByID(c.Request.Context(), boardID, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, board)
}

func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	boardID, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.service.Update(c.Request.Context(), boardID, userID.(int), input.Name); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "board updated successfully"})
}

func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	boardID, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), boardID, userID.(int)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "board deleted successfully"})
}

func (h *BoardHandler) AddMemberToBoard(c *gin.Context) {
	inviterID, exists := c.Get("userId")
	if !exists {
	}

	boardID, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
	}

	var input AddMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	err = h.service.AddMember(c.Request.Context(), boardID, inviterID.(int), input.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member added successfully"})
}
