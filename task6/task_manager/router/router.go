package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"task_manager/controllers"
	"task_manager/data"
)

func SetupRouter() *gin.Engine {
	// MongoDB connection configuration.
	mongoURI := "mongodb://localhost:27017"
	dbName := "task_manager_db"
	collectionName := "tasks"

	taskService, err := data.NewTaskService(mongoURI, dbName, collectionName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	taskController := controllers.NewTaskController(taskService)
	r := gin.Default()

	// Task endpoints.
	r.GET("/tasks", taskController.GetTasks)
	r.GET("/tasks/:id", taskController.GetTask)
	r.POST("/tasks", taskController.CreateTask)
	r.PUT("/tasks/:id", taskController.UpdateTask)
	r.DELETE("/tasks/:id", taskController.DeleteTask)

	return r
}
