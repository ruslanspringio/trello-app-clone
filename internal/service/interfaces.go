package service

type Broadcaster interface {
	BroadcastToBoard(boardID int, message []byte)
}
