package models

import "time"

type Certification struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Issuer    string    `json:"issuer"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
}
