# Library Management System - Concurrent Book Reservation

## Overview

This version of the Library Management System has been extended to support concurrent book reservations using Goroutines, Channels, and Mutexes. The key enhancements include:

- **Goroutines:** Allow multiple reservation requests to be processed concurrently.
- **Channels:** Queue incoming reservation requests, enabling asynchronous processing.
- **Mutexes:** Protect shared data (books and members) from race conditions during concurrent updates.
- **Auto-Cancellation:** A timer-based Goroutine auto-cancels reservations if the book is not borrowed within 5 seconds.

## Concurrency Details

### Reservation Process

1. **Reservation Request:**
   - A reservation request (containing the `bookID` and `memberID`) is sent to the reservation worker via a channel.
2. **Worker Processing:**
   - The reservation worker, running in its own Goroutine, reads from the channel and calls the `ReserveBook` method on the library.
3. **Mutex Protection:**
   - The `ReserveBook` method uses a Mutex (`sync.Mutex`) to lock the library data structures during updates, ensuring safe concurrent access.
4. **Auto-Cancellation:**
   - Once a book is reserved, a separate Goroutine starts a timer. If the book is not borrowed within 5 seconds, the reservation is automatically canceled.
5. **Error Handling:**
   - If the book is not available or already reserved by another member, an error is returned.

### Simulating Concurrent Requests

- The system is designed to safely handle multiple reservation requests simultaneously, preventing double reservations and ensuring data consistency.

## Folder Structure
