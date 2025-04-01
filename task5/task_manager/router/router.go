package router

import (
	"github.com/gin-gonic/gin"
	"task_manager/controllers"
	"task_manager/data"
)

// SetupRouter initializes the Gin router with routes.
func SetupRouter() *gin.Engine {
	// Initialize the in-memory task service.
	taskService := data.NewTaskService()
	// Create the task controller.
	taskController := controllers.NewTaskController(taskService)

	// Initialize Gin router.
	r := gin.Default()

	// Task endpoints.
	r.GET("/tasks", taskController.GetTasks)
	r.GET("/tasks/:id", taskController.GetTask)
	r.POST("/tasks", taskController.CreateTask)
	r.PUT("/tasks/:id", taskController.UpdateTask)
	r.DELETE("/tasks/:id", taskController.DeleteTask)

	return r
}
