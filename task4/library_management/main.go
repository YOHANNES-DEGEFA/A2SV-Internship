package main

import (
	"library_management/controllers"
	"library_management/models"
	"library_management/services"
	"library_management/concurrency"
)

func main() {
	// Initialize the library.
	library := services.NewLibrary()

	// Add an initial member for testing.
	library.AddMember(models.Member{
		ID:            1,
		Name:          "Alice",
		BorrowedBooks: []models.Book{},
	})

	// Add some sample books.
	library.AddBook(models.Book{
		ID:     101,
		Title:  "The Go Programming Language",
		Author: "Alan A. A. Donovan",
	})
	library.AddBook(models.Book{
		ID:     102,
		Title:  "Introducing Go",
		Author: "Caleb Doxsey",
	})

	// Create a channel for reservation requests.
	reservationChan := make(chan concurrency.ReservationRequest)

	// Start the reservation worker that listens on the channel.
	concurrency.StartReservationWorker(library, reservationChan)

	// Start the console-based library controller.
	controllers.LibraryController(library, reservationChan)
}
