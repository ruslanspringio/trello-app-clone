package handlers

import (
	"net/http"
	"notes-project/internal/models"
	"notes-project/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	service service.CardService
}

type MoveCardInput struct {
	NewListID   int     `json:"new_list_id" binding:"required"`
	NewPosition float64 `json:"new_position" binding:"required"`
}

func NewCardHandler(s service.CardService) *CardHandler {
	return &CardHandler{service: s}
}

func (h *CardHandler) RegisterCardRoutes(rg *gin.RouterGroup) {
	cardsGroup := rg.Group("/cards")
	{
		cardsGroup.PUT("/:cardId/move", h.MoveCard)
	}

	listsGroup := rg.Group("/lists/:listId/cards")
	{
		listsGroup.POST("/", h.CreateCard)
	}
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	listID, err := strconv.Atoi(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid list id"})
		return
	}

	var input models.Card
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.service.Create(c.Request.Context(), &input, listID, userID.(int)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (h *CardHandler) MoveCard(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	cardID, err := strconv.Atoi(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	var input MoveCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	err = h.service.Move(c.Request.Context(), cardID, input.NewListID, input.NewPosition, userID.(int))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "card moved successfully"})
}
