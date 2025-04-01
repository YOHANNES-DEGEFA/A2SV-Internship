package models

// Task represents a task in the system.
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // You can change to time.Time if needed.
	Status      string `json:"status"`   // e.g., "pending", "completed"
}
