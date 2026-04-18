package models

import "time"

type FavoriteBook struct {
	UserID    int       `json:"user_id"    db:"user_id"`
	BookID    int       `json:"book_id"    db:"book_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
