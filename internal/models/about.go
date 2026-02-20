package models

import "time"

type About struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}
