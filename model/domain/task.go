package domain

import "time"

type Task struct {
	Id          string
	UserId      string
	Title       string
	Description string
	Status      string
	Priority    string
	DueDate    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}