package concurrency

import (
	"fmt"
	"library_management/services"
)

// ReservationRequest encapsulates a reservation request.
type ReservationRequest struct {
	BookID   int
	MemberID int
	Response chan error
}

// StartReservationWorker starts a goroutine that processes reservation requests from the channel.
func StartReservationWorker(library *services.Library, requests chan ReservationRequest) {
	go func() {
		for req := range requests {
			err := library.ReserveBook(req.BookID, req.MemberID)
			req.Response <- err
		}
	}()
	fmt.Println("Reservation worker started...")
}
