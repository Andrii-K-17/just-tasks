package models

import "time"

type TaskCollaborator struct {
	ID       int    `db:"id"       json:"id"`
	Username string `db:"username" json:"username"`
}

type Task struct {
	ID          int       `db:"id"           json:"id"`
	UserID      int       `db:"user_id"      json:"user_id"`
	TaskText    string    `db:"task_text"    json:"task_text"`
	Priority    int       `db:"priority"     json:"priority"`
	Deadline    *string   `db:"deadline"     json:"deadline"`
	IsCompleted bool      `db:"is_completed" json:"is_completed"`
	Position    int       `db:"position"     json:"position"`
	CategoryID  *int      `db:"category_id"  json:"category_id"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}

type TaskResponse struct {
	Task
	OwnerName     string             `json:"owner_name"`
	Collaborators []TaskCollaborator `json:"collaborators"`
}
