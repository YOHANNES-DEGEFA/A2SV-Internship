
package models

// User represents a system user.
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`      // Do not expose the hash
	Role         string `json:"role"`   // "admin" or "user"
}
