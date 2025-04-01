package models

// Task represents a task in the system.
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // Format: YYYY-MM-DD
	Status      string `json:"status"`   // e.g., "pending", "completed"
}
