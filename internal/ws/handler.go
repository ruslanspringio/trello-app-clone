package ws

import (
	"log"
	"net/http"
	"strconv"

	"notes-project/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Conn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WsHandler struct {
	hub          *Hub
	boardService service.BoardService
}

func NewWsHandler(h *Hub, bs service.BoardService) *WsHandler {
	return &WsHandler{hub: h, boardService: bs}
}

func (h *WsHandler) RegisterWsRoutes(rg *gin.RouterGroup) {
	rg.GET("/boards/:boardId/ws", h.ServeWs)
}

func (h *WsHandler) ServeWs(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
	}

	boardID, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
	}

	_, err = h.boardService.GetByID(c.Request.Context(), boardID, userID.(int))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		hub:  h.hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.hub.register <- &subscription{client: client, boardID: boardID}

	go client.writePump()
	go client.readPump(boardID)

	log.Printf("Client connected to board %d", boardID)
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (c *Client) readPump(boardID int) {
	defer func() {
		c.hub.unregister <- &subscription{client: c, boardID: boardID}
		c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}
