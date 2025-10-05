package models

import "time"

type Card struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Position    float64   `db:"position" json:"position"`
	ListID      int       `db:"list_id" json:"list_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type List struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Position  float64   `db:"position" json:"position"`
	BoardID   int       `db:"board_id" json:"board_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Cards []Card `json:"cards,omitempty"`
}

type Board struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	OwnerID   int       `db:"owner_id" json:"owner_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Lists   []List `json:"lists,omitempty"`
	Members []User `json:"members,omitempty"`
}

type WebSocketMessage struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}
