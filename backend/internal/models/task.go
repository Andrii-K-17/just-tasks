package models

import "time"

type Task struct {
	ID          int       `db:"id"           json:"id"`
	UserID      int       `db:"user_id"      json:"user_id"`
	TaskText    string    `db:"task_text"    json:"task_text"`
	Priority    int       `db:"priority"     json:"priority"`
	Deadline    *string   `db:"deadline"     json:"deadline"`
	IsCompleted bool      `db:"is_completed" json:"is_completed"`
	Position    int       `db:"position"     json:"position"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}
