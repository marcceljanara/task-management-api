package web

import "time"

type TaskCreateRequest struct {
	Title       string    `validate:"required" json:"title"`
	Description string    `validate:"max=255" json:"description"`
	Priority    string    `validate:"oneof=low medium high" json:"priority"`
	Status      string    `validate:"oneof=pending in_progress completed" json:"status"`
	DueDate     time.Time `validate:"required" json:"due_date"`
}
