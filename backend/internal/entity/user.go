package entity

import "time"

type User struct {
	ID        int       `json:"id"` // later > string (uuid)
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Coins     int       `json:"coins"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
