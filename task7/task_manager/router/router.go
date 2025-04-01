package router

import (
	// "log"

	"github.com/gin-gonic/gin"
	"task_manager/controllers"
	"task_manager/data"
	"task_manager/middleware"
)

// SetupRouter initializes the Gin router and configures routes.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Initialize services.
	userService := data.NewUserService()
	taskService := data.NewTaskService()
	controller := controllers.NewController(userService, taskService)

	// Public endpoints.
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	// Endpoint to promote a user (admin only).
	// This route is protected by the auth middleware and then admin-only check.
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/promote", middleware.AdminOnly(), controller.Promote)

	// Task endpoints.
	// All task endpoints require authentication.
	tasks := r.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware())
	{
		tasks.GET("", controller.GetTasks)
		tasks.GET("/:id", controller.GetTask)
		// Admin-only actions.
		tasks.POST("", middleware.AdminOnly(), controller.CreateTask)
		tasks.PUT("/:id", middleware.AdminOnly(), controller.UpdateTask)
		tasks.DELETE("/:id", middleware.AdminOnly(), controller.DeleteTask)
	}

	return r
}
