package model

import (
	"time"
)

type Todo struct {
	ID        uint64    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Title     string    `json:"title" db:"title"`
	Memo      string    `json:"memo" db:"memo"`
	IsDone    bool      `json:"is_done" db:"is_done"`
	DueDate   time.Time `json:"due_date" db:"due_date"`
}
