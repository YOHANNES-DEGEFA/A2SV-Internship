package routers

import (
	"time"

	"github.com/gin-gonic/gin"
	"task_manager/Delivery/controllers"
	"task_manager/Infrastructure" // Need middleware + service implementations
	"task_manager/Repositories"  // Need repository implementations
	"task_manager/Usecases"      // Need usecase implementations
)

// SetupRouter initializes dependencies (DI) and configures Gin routes.
func SetupRouter() *gin.Engine {
	// --- Configuration (Move to config files/env vars in a real app) ---
	jwtSecret := "your_very_secret_key_here_CHANGE_ME_IN_PRODUCTION"
	jwtIssuer := "task-manager-api"
	jwtExpiry := 24 * time.Hour // Token valid for 24 hours

	// --- Dependency Injection ---

	// 1. Repositories (using in-memory implementations)
	userRepo := Repositories.NewInMemoryUserRepository()
	taskRepo := Repositories.NewInMemoryTaskRepository()

	// 2. Infrastructure Services
	passwordSvc := Infrastructure.NewBcryptPasswordService()
	jwtSvc := Infrastructure.NewJwtGoJWTService(jwtSecret, jwtIssuer, jwtExpiry)

	// 3. Usecases (injecting repository and Infrastructure dependencies)
	userUsecase := Usecases.NewUserUsecase(userRepo, passwordSvc, jwtSvc)
	taskUsecase := Usecases.NewTaskUsecase(taskRepo)

	// 4. Controllers (injecting usecase dependencies)
	appController := controllers.NewController(userUsecase, taskUsecase)

	// 5. Middleware (injecting JWT service dependency)
	authMiddleware := Infrastructure.NewAuthMiddleware(jwtSvc)

	// --- Gin Router Setup ---
	// Consider gin.ReleaseMode for production: gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// --- Public Routes (No Authentication Required) ---
	apiV1 := r.Group("/api/v1") // Optional: versioning
	{
		apiV1.POST("/register", appController.Register)
		apiV1.POST("/login", appController.Login)
	}

	// --- Protected Routes (Authentication Required) ---
	protected := apiV1.Group("/")
	// Apply authentication middleware to all routes in this group
	protected.Use(authMiddleware.Authenticate())
	{
		// Promote User (Requires Admin Role)
		// Apply AdminOnly middleware *after* Authenticate
		protected.POST("/promote", Infrastructure.AuthorizeAdmin(), appController.Promote)

		// Task Routes
		tasks := protected.Group("/tasks")
		{
			// Accessible by any authenticated user (Admin or User)
			tasks.GET("", appController.GetTasks)
			tasks.GET("/:id", appController.GetTask)

			// Accessible only by Admin users
			// Apply AdminOnly middleware *after* Authenticate
			tasks.POST("", Infrastructure.AuthorizeAdmin(), appController.CreateTask)
			tasks.PUT("/:id", Infrastructure.AuthorizeAdmin(), appController.UpdateTask)
			tasks.DELETE("/:id", Infrastructure.AuthorizeAdmin(), appController.DeleteTask)
		}

		// Add other protected routes here if needed
	}

	// Add Swagger/OpenAPI documentation endpoint setup here if desired

	return r
}