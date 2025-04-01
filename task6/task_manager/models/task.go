package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task represents a task in the system.
type Task struct {
    ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Title       string             `json:"title" bson:"title"`
    Description string             `json:"description" bson:"description"`
    DueDate     string             `json:"due_date" bson:"due_date"`
    Status      string             `json:"status" bson:"status"` // e.g., "pending", "completed"
}
