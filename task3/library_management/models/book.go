package models

// Book represents a library book.
type Book struct {
	ID     int
	Title  string
	Author string
	Status string // "Available" or "Borrowed"
}
