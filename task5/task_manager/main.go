package main

import (
	"task_manager/router"
)

func main() {
	r := router.SetupRouter()
	// Run on default port 8080.
	r.Run(":8080")
}
