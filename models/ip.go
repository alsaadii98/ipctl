package models

import "time"

type IPAddress struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
