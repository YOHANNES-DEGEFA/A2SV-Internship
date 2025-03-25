package main

import (
    "fmt"
    "library_management/controllers"
    "library_management/models"
    "library_management/services"
)

func main() {
    library := services.NewLibrary()
    controller := controllers.LibraryController{LibraryService: library}

    // Example usage
    book1 := models.Book{ID: 1, Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Status: "Available"}
    book2 := models.Book{ID: 2, Title: "Clean Code", Author: "Robert C. Martin", Status: "Available"}
    member1 := models.Member{ID: 1, Name: "John Doe"}

    controller.AddBook(book1)
    controller.AddBook(book2)
    controller.BorrowBook(1, 1)
    controller.ListAvailableBooks()
    controller.ListBorrowedBooks(1)
    controller.ReturnBook(1, 1)
    controller.ListAvailableBooks()
}
