package models

type Category struct {
	ID     int    `db:"id"      json:"id"`
	UserID int    `db:"user_id" json:"-"`
	Name   string `db:"name"    json:"name"`
}
