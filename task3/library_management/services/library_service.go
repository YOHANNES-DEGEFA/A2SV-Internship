package services

import (
	"errors"
	"fmt"
	"library_management/models"
)

// LibraryManager defines methods for managing the library.
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

// Library implements LibraryManager.
type Library struct {
	Books   map[int]models.Book  // Keyed by book ID
	Members map[int]models.Member // Keyed by member ID
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
	book.Status = "Available"
	l.Books[book.ID] = book
}

// RemoveBook removes a book from the library by its ID.
func (l *Library) RemoveBook(bookID int) {
	delete(l.Books, bookID)
}

// BorrowBook allows a member to borrow a book if it is available.
func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}
	if book.Status != "Available" {
		return errors.New("book is already borrowed")
	}
	member, exists := l.Members[memberID]
	if !exists {
		return errors.New("member not found")
	}
	// Update book status and add to member's borrowed books.
	book.Status = "Borrowed"
	l.Books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member
	return nil
}

// ReturnBook allows a member to return a borrowed book.
func (l *Library) ReturnBook(bookID int, memberID int) error {
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
	member, exists := l.Members[memberID]
	if !exists {
		return []models.Book{}
	}
	return member.BorrowedBooks
}

// AddMember is a helper method to add a new member.
func (l *Library) AddMember(member models.Member) {
	l.Members[member.ID] = member
}

// ListAllBooks prints all books (for debugging or full listing).
func (l *Library) ListAllBooks() {
	fmt.Println("All Books in Library:")
	for _, book := range l.Books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Status: %s\n", book.ID, book.Title, book.Author, book.Status)
	}
}
