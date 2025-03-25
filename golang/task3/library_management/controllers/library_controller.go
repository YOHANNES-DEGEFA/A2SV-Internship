
package controllers

import (
    "fmt"
    "library_management/models"
    "library_management/services"
)

type LibraryController struct {
    LibraryService services.LibraryManager
}

func (lc *LibraryController) AddBook(book models.Book) {
    lc.LibraryService.AddBook(book)
    fmt.Println("Book added successfully!")
}

func (lc *LibraryController) RemoveBook(bookID int) {
    err := lc.LibraryService.RemoveBook(bookID)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Book removed successfully!")
    }
}

func (lc *LibraryController) BorrowBook(bookID int, memberID int) {
    err := lc.LibraryService.BorrowBook(bookID, memberID)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Book borrowed successfully!")
    }
}

func (lc *LibraryController) ReturnBook(bookID int, memberID int) {
    err := lc.LibraryService.ReturnBook(bookID, memberID)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Book returned successfully!")
    }
}

func (lc *LibraryController) ListAvailableBooks() {
    books := lc.LibraryService.ListAvailableBooks()
    fmt.Println("Available Books:")
    for _, book := range books {
        fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
    }
}

func (lc *LibraryController) ListBorrowedBooks(memberID int) {
    books := lc.LibraryService.ListBorrowedBooks(memberID)
    fmt.Printf("Books borrowed by Member ID %d:\n", memberID)
    for _, book := range books {
        fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
    }
}




// package main

// import (
// 	"fmt"
// 	"strings"
// )


// func promptUser() {

// 	 var command string
// 	 var userId int
// 	 fmt.Println("If you are a member Enter your ID if not Enter -1")
// 	 fmt.Scan(&userId)

// 	 if userId == -1 {
// 		// your id is ...  and add a new member
		
// 	 }

// 	 fmt.Print("Enter one of the following options: \n\t",
// 					"1. AddBook\n\t",
// 					"2. RemoveBook\n\t",
// 					"3. BorrowBook\n\t",
// 					"4. ReturnBook(bookID int, memberID int) error\n\t", 
// 					"6. ListBorrowedBooks(memberID int)\n\t")

// 	fmt.Scan(&command)

// 	if (len(command) != 1) || (!strings.Contains("123456", command)){
	
// 	} else if command == "1" {
		
// 	} else if command == "2" {

// 	} else if command == "3" {

// 	} else if command == "4" {
 
// 	} else if command == "5" {

// 	} else if command == "4" {
 
// 	} else if command == "5" {

// 	} else {

// }

// }