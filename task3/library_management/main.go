package main

import (
	"fmt"
	"library_management/controllers"
	"library_management/models"
	"library_management/services"
)

func main() {
	library := services.NewLibrary()
	

	
	print([]models.Book{}) // output: []    // just empty slice or array of books
	// print([]models.Book)  // this is data type and  so error will be seen
	// Optionally, add an initial member for testing.
	library.AddMember(models.Member{
		ID:            1,
		Name:          "Alice",
		BorrowedBooks: []models.Book{},  // empty array of Book if to say []models.Book would be error cause this is type not a value 

	})

	// Start the console-based interface.
	controllers.LibraryController(library)
}


func print(x any) {
	fmt.Println(x)
}