package services

import (
	"errors"
	"fmt"
	"library_management/models"
	"sync"
	"time"
)

// LibraryManager defines methods for managing the library.
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
	ReserveBook(bookID int, memberID int) error
	AddMember(member models.Member)
}

// Library implements LibraryManager.
type Library struct {
	Books   map[int]models.Book   // Keyed by book ID
	Members map[int]models.Member // Keyed by member ID
	mu      sync.Mutex            // Protects access to Books and Members
}

// NewLibrary creates a new Library instance.
func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

// AddBook adds a new book to the library.
func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()
	book.Status = "Available"
	book.ReservedBy = 0
	l.Books[book.ID] = book
}

// RemoveBook removes a book from the library by its ID.
func (l *Library) RemoveBook(bookID int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.Books, bookID)
}

// BorrowBook allows a member to borrow a book if it is available or reserved for them.
func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	// If the book is reserved, only the member who reserved it can borrow it.
	if book.Status == "Reserved" {
		if book.ReservedBy != memberID {
			return errors.New("book is reserved by another member")
		}
		// Reservation will be cleared upon borrowing.
	}

	if book.Status != "Available" && book.Status != "Reserved" {
		return errors.New("book is already borrowed")
	}

	// Update book status to Borrowed.
	book.Status = "Borrowed"
	book.ReservedBy = 0
	l.Books[bookID] = book

	member, exists := l.Members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	return nil
}

// ReturnBook allows a member to return a borrowed book.
func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	member, exists := l.Members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	// Check if the member has borrowed this book.
	found := false
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			// Remove the book from the member's slice.
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("this book is not borrowed by the member")
	}

	// Update the book status.
	book.Status = "Available"
	l.Books[bookID] = book
	l.Members[memberID] = member
	return nil
}

// ListAvailableBooks lists all books that are currently available.
func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	available := []models.Book{}
	for _, book := range l.Books {
		if book.Status == "Available" {
			available = append(available, book)
		}
	}
	return available
}

// ListBorrowedBooks lists all books borrowed by a specific member.
func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	member, exists := l.Members[memberID]
	if !exists {
		return []models.Book{}
	}
	return member.BorrowedBooks
}

// ReserveBook reserves a book for a member if it is available.
// It also starts a timer that will auto-cancel the reservation after 5 seconds if not borrowed.
func (l *Library) ReserveBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}
	if book.Status != "Available" {
		return errors.New("book is not available for reservation")
	}

	// Reserve the book.
	book.Status = "Reserved"
	book.ReservedBy = memberID
	l.Books[bookID] = book

	// Start auto-cancellation in a separate goroutine.
	go l.autoCancelReservation(bookID, memberID)

	return nil
}

// autoCancelReservation will cancel a reservation if the book is not borrowed within 5 seconds.
func (l *Library) autoCancelReservation(bookID int, memberID int) {
	timer := time.NewTimer(5 * time.Second)
	<-timer.C

	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.Books[bookID]
	if !exists {
		return
	}
	// If still reserved by the same member, cancel the reservation.
	if book.Status == "Reserved" && book.ReservedBy == memberID {
		book.Status = "Available"
		book.ReservedBy = 0
		l.Books[bookID] = book
		fmt.Printf("Auto-cancellation: Reservation for book %d by member %d has timed out.\n", bookID, memberID)
	}
}

// AddMember adds a new member to the library.
func (l *Library) AddMember(member models.Member) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Members[member.ID] = member
}

// ListAllBooks prints all books (for debugging or full listing).
func (l *Library) ListAllBooks() {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Println("All Books in Library:")
	for _, book := range l.Books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Status: %s, ReservedBy: %d\n", book.ID, book.Title, book.Author, book.Status, book.ReservedBy)
	}
}
