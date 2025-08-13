package domain

import (
	"time"
)

type Poem struct {
	PoemID    int64     `json:"poem_id"`
	Name      string    `json:"name"`
	AuthorID  int64     `json:"author_id"`
	Text      string    `json:"text"`
	Meta      string    `json:"meta"`
	IsSent    bool      `json:"is_sent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    Author    `json:"author"`
}
