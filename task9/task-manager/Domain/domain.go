package domain

import "errors"

// Constants for user roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User represents the core user entity.
// No json tags or implementation details like PasswordHash.
type User struct {
	ID       int
	Username string
	Role     string
}

// Task represents the core task entity.
// No json tags.
type Task struct {
	ID          int
	Title       string
	Description string
	DueDate     string // Consider using time.Time in a real application
	Status      string
}

// --- Domain specific errors ---
// These represent business rule violations or specific data states.
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrTaskNotFound       = errors.New("task not found")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials") // Keep generic for security
	ErrForbidden          = errors.New("forbidden access")
	ErrInternalServer     = errors.New("internal server error") // For unexpected issues
	ErrBadRequest         = errors.New("bad request")           // For invalid input structure/format
)