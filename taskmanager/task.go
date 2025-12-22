package taskmanager

type Task struct {
	ID     int    `json:"id"`      // Tells Go to use lowercase "id" in JSON
	Title  string `json:"title"`   // Tells Go to use lowercase "title"
	IsDone bool   `json:"is_done"` // Tells Go to use lowercase "is_done"
}
