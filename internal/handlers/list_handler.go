package handlers

import (
	"net/http"
	"notes-project/internal/models"
	"notes-project/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	service service.ListService
}

func NewListHandler(s service.ListService) *ListHandler {
	return &ListHandler{service: s}
}

func (h *ListHandler) RegisterListRoutes(rg *gin.RouterGroup) {
	lists := rg.Group("/boards/:boardId/lists")
	{
		lists.POST("/", h.CreateList)
	}
}

func (h *ListHandler) CreateList(c *gin.Context) {
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

	var input models.List
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.service.Create(c.Request.Context(), &input, boardID, userID.(int)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}
