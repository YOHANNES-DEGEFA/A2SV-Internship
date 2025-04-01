package services

import (
    "errors"
    "library_management/models"
)

type LibraryManager interface {
    AddBook(book models.Book)
    RemoveBook(bookID int) error
    BorrowBook(bookID int, memberID int) error
    ReturnBook(bookID int, memberID int) error
    ListAvailableBooks() []models.Book
    ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
    Books   map[int]models.Book
    Members map[int]models.Member
}

func NewLibrary() *Library {
    return &Library{
        Books:   make(map[int]models.Book),
        Members: make(map[int]models.Member),
    }
}

func (l *Library) AddBook(book models.Book) {
    l.Books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) error {
    if _, exists := l.Books[bookID]; !exists {
        return errors.New("book not found")
    }
    delete(l.Books, bookID)
    return nil
}

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

    book.Status = "Borrowed"
    l.Books[bookID] = book
    member.BorrowedBooks = append(member.BorrowedBooks, book)
    l.Members[memberID] = member
    return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
    book, exists := l.Books[bookID]
    if !exists {
        return errors.New("book not found")
    }
    if book.Status != "Borrowed" {
        return errors.New("book is not borrowed")
    }

    member, exists := l.Members[memberID]
    if !exists {
        return errors.New("member not found")
    }

    book.Status = "Available"
    l.Books[bookID] = book

    // Remove the book from the member's borrowed books
    for i, b := range member.BorrowedBooks {
        if b.ID == bookID {
            member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
            break
        }
    }
    l.Members[memberID] = member
    return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
    var availableBooks []models.Book
    for _, book := range l.Books {
        if book.Status == "Available" {
            availableBooks = append(availableBooks, book)
        }
    }
    return availableBooks
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
    member, exists := l.Members[memberID]
    if !exists {
        return nil
    }
    return member.BorrowedBooks
}
