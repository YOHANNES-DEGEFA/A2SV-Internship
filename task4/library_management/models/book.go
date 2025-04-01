package models

// Book represents a library book.
type Book struct {
	ID         int
	Title      string
	Author     string
	Status     string // "Available", "Reserved", or "Borrowed"
	ReservedBy int    // ID of the member who reserved the book (0 if not reserved)
}
