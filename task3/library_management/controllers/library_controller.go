package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

// LibraryController provides a console interface to interact with the library.
func LibraryController(library *services.Library) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n--- Library Management System ---")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books by Member")
		fmt.Println("7. Add Member")
		fmt.Println("8. List All Books")
		fmt.Println("9. Exit")
		fmt.Print("Enter your choice: ")

		input, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("Invalid input, please enter a number.")
			continue
		}

		switch choice {
		case 1:
			addBook(reader, library)
		case 2:
			removeBook(reader, library)
		case 3:
			borrowBook(reader, library)
		case 4:
			returnBook(reader, library)
		case 5:
			listAvailableBooks(library)
		case 6:
			listBorrowedBooks(reader, library)
		case 7:
			addMember(reader, library)
		case 8:
			library.ListAllBooks()
		case 9:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func addBook(reader *bufio.Reader, library *services.Library) {
	fmt.Print("Enter Book ID: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	fmt.Print("Enter Book Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Enter Book Author: ")
	author, _ := reader.ReadString('\n')
	author = strings.TrimSpace(author)

	book := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: "Available",
	}
	library.AddBook(book)
	fmt.Println("Book added successfully!")
}

func removeBook(reader *bufio.Reader, library *services.Library) {
	fmt.Print("Enter Book ID to remove: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	library.RemoveBook(id)
	fmt.Println("Book removed successfully!")
}

func borrowBook(reader *bufio.Reader, library *services.Library) {
	fmt.Print("Enter Book ID to borrow: ")
	bookIDStr, _ := reader.ReadString('\n')
	bookID, err := strconv.Atoi(strings.TrimSpace(bookIDStr))
	if err != nil {
		fmt.Println("Invalid Book ID")
		return
	}

	fmt.Print("Enter Member ID: ")
	memberIDStr, _ := reader.ReadString('\n')
	memberID, err := strconv.Atoi(strings.TrimSpace(memberIDStr))
	if err != nil {
		fmt.Println("Invalid Member ID")
		return
	}

	if err := library.BorrowBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book borrowed successfully!")
	}
}

func returnBook(reader *bufio.Reader, library *services.Library) {
	fmt.Print("Enter Book ID to return: ")
	bookIDStr, _ := reader.ReadString('\n')
	bookID, err := strconv.Atoi(strings.TrimSpace(bookIDStr))
	if err != nil {
		fmt.Println("Invalid Book ID")
		return
	}

	fmt.Print("Enter Member ID: ")
	memberIDStr, _ := reader.ReadString('\n')
	memberID, err := strconv.Atoi(strings.TrimSpace(memberIDStr))
	if err != nil {
		fmt.Println("Invalid Member ID")
		return
	}

	if err := library.ReturnBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book returned successfully!")
	}
}

func listAvailableBooks(library *services.Library) {
	books := library.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}
	fmt.Println("Available Books:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func listBorrowedBooks(reader *bufio.Reader, library *services.Library) {
	fmt.Print("Enter Member ID: ")
	memberIDStr, _ := reader.ReadString('\n')
	memberID, err := strconv.Atoi(strings.TrimSpace(memberIDStr))
	if err != nil {
		fmt.Println("Invalid Member ID")
		return
	}
	books := library.ListBorrowedBooks(memberID)
	if len(books) == 0 {
		fmt.Println("No borrowed books for this member.")
		return
	}
	fmt.Println("Borrowed Books:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func addMember(reader *bufio.Reader, library *services.Library) {
	fmt.Print("Enter Member ID: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("Invalid Member ID")
		return
	}

	fmt.Print("Enter Member Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	member := models.Member{
		ID:            id,
		Name:          name,
		BorrowedBooks: []models.Book{},
	}
	library.AddMember(member)
	fmt.Println("Member added successfully!")
}
