package domain

import "time"

type User struct {
	UserID     int64     `json:"user_id"`
	TelegramID int64     `json:"telegram_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserAnswer struct {
	AnswerID   int64     `json:"answer_id"`
	TelegramID int64     `json:"telegram_id"`
	TaskID     int64     `json:"task_id"`
	UserText   string    `json:"user_text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
