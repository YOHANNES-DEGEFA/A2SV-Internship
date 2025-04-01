package models

// Member represents a library member.
type Member struct {
	ID            int
	Name          string
	BorrowedBooks []Book // Slice to hold borrowed books
}
