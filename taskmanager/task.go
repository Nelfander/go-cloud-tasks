package taskmanager

import "time"

type Task struct {
	ID     int    `json:"id"`      // Tells Go to use lowercase "id" in JSON
	Title  string `json:"title"`   // Tells Go to use lowercase "title"
	IsDone bool   `json:"is_done"` // Tells Go to use lowercase "is_done"
	UserID int    `json:"user_id"` // This links the task to a user
}

type User struct { // User represents  account system
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // The "-" means: never send this back in JSON for safety!
	CreatedAt    time.Time `json:"created_at"`
}
