package domain

import "time"

type Task struct {
	TaskID    int64     `json:"task_id"`
	PoemID    int64     `json:"poem_id"`
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
