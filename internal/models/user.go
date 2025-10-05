package models

import "time"

type User struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Age          string    `db:"age" json:"age"`
	Email        string    `db:"email" json:"email"`
	Password     string    `db:"-" json:"password,omitempty"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
