package domain

import "time"

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      string    `json:"userID"`
	CreatedAt   time.Time `json:"createdAt"`
}
