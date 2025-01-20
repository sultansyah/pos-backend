package notification

import "time"

type Notification struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
