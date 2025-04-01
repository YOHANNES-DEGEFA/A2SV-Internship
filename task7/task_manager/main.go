package main

import (
	"task_manager/router"
)

func main() {
	r := router.SetupRouter()
	// Run the server on port 8080.
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
