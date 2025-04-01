package controllers

import (
	"net/http"
	"strconv"
	"task_manager/Domain"   // Need Domain errors
	"task_manager/Usecases" // Need usecase interfaces

	"github.com/gin-gonic/gin"
)

// --- DTOs (Data Transfer Objects) ---
// Define structs for API request binding and response formatting.
// This decouples the API structure from the internal Domain models.

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// UserResponse excludes sensitive info like password hash.
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type PromoteRequest struct {
	Username string `json:"username" binding:"required"`
}

// TaskResponse is the structure returned for task endpoints.
type TaskResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}

// CreateTaskRequest defines the expected body for creating a task.
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // Add validation tags if needed (e.g., format)
	Status      string `json:"status"`   // Add validation tags if needed (e.g., enum)
}

// UpdateTaskRequest defines the expected body for updating a task.
type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}

// --- Controller ---
// Handles HTTP requests, calls Usecases, and formats responses.
type Controller struct {
	UserUsecase Usecases.UserUsecase // Interface dependency
	TaskUsecase Usecases.TaskUsecase // Interface dependency
}

// NewController creates a new Controller instance with dependencies.
func NewController(userUC Usecases.UserUsecase, taskUC Usecases.TaskUsecase) *Controller {
	return &Controller{
		UserUsecase: userUC,
		TaskUsecase: taskUC,
	}
}

// --- Error Mapping Helper ---
// Maps Domain/usecase errors to appropriate HTTP status codes and responses.
func mapDomainErrorToHTTP(err error) (int, gin.H) {
	switch err {
	case Domain.ErrInvalidCredentials:
		// Specific handling for login failure
		return http.StatusUnauthorized, gin.H{"error": err.Error()}
	case Domain.ErrUsernameExists:
		// Specific handling for registration conflict
		return http.StatusConflict, gin.H{"error": err.Error()} // 409 Conflict is often suitable
	case Domain.ErrUserNotFound, Domain.ErrTaskNotFound:
		return http.StatusNotFound, gin.H{"error": err.Error()}
	case Domain.ErrForbidden:
		return http.StatusForbidden, gin.H{"error": err.Error()}
	case Domain.ErrBadRequest:
		return http.StatusBadRequest, gin.H{"error": err.Error()}
	case Domain.ErrInternalServer:
		// Log the internal error details here (not shown for brevity)
		return http.StatusInternalServerError, gin.H{"error": "An unexpected internal error occurred"}
	default:
		// Log unexpected errors
		// logger.Error("Unexpected error in controller:", err)
		return http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"}
	}
}

// --- User Handlers ---

// Register handles POST /register requests.
func (ctr *Controller) Register(c *gin.Context) {
	var req RegisterRequest
	// Bind request body to DTO, handles basic validation (e.g., required fields)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Call the usecase
	user, err := ctr.UserUsecase.Register(req.Username, req.Password)
	if err != nil {
		// Map the Domain error to HTTP response
		status, body := mapDomainErrorToHTTP(err)
		c.JSON(status, body)
		return
	}

	// Map Domain.User to UserResponse DTO before sending
	resp := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	c.JSON(http.StatusCreated, resp) // 201 Created
}

// Login handles POST /login requests.
func (ctr *Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Call the usecase
	token, err := ctr.UserUsecase.Login(req.Username, req.Password)
	if err != nil {
		status, body := mapDomainErrorToHTTP(err)
		c.JSON(status, body)
		return
	}

	// Return the token in the response DTO
	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// Promote handles POST /promote requests (requires admin auth via middleware).
func (ctr *Controller) Promote(c *gin.Context) {
	// Authorization check (admin role) is handled by middleware before this runs.
	var req PromoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Call the usecase
	err := ctr.UserUsecase.Promote(req.Username)
	if err != nil {
		status, body := mapDomainErrorToHTTP(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User '%s' promoted to admin", req.Username)})
}

// --- Task Handlers ---

// mapDomainTaskToResponse converts a Domain.Task to a TaskResponse DTO.
func mapDomainTaskToResponse(task Domain.Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
	}
}

// mapDomainTasksToResponse converts a slice of Domain.Task to []TaskResponse.
func mapDomainTasksToResponse(tasks []Domain.Task) []TaskResponse {
	responses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = mapDomainTaskToResponse(task)
	}
	return responses
}

// GetTasks handles GET /tasks (requires authentication via middleware).
func (ctr *Controller) GetTasks(c *gin.Context) {
	tasks, err := ctr.TaskUsecase.GetAllTasks()
	if err != nil {
		status, body := mapDomainErrorToHTTP(err)
		c.JSON(status, body)
		return
	}
	// Map Domain tasks to response DTOs
	c.JSON(http.StatusOK, mapDomainTasksToResponse(tasks))
}

// GetTask handles GET /tasks/:id (requires authentication via middleware).
func (ctr *Controller) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	task, err := ctr.TaskUsecase.GetTaskByID(id)
	if err != nil {
		status, body := mapDomainErrorToHTTP(err)
		c.JSON(status, body)
		return
	}
	// Map Domain task to response DTO
	c.JSON(http.StatusOK, mapDomainTaskToResponse(task))
}

// CreateTask handles POST /tasks (requires admin auth via middleware).
func (ctr *Controller) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	createdTask, err := ctr.TaskUsecase.CreateTask(req.Title, req.Description, req.DueDate, req.Status)
	if err != nil {
		status, body := mapDomainErrorToHTTP(err)
		c.JSON(status, body)
		return
	}
	// Map Domain task to response DTO
	c.JSON(http.StatusCreated, mapDomainTaskToResponse(createdTask))
}

// UpdateTask handles PUT /tasks/:id (requires admin auth via middleware).
func (ctr *Controller) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	updatedTask, err := ctr.TaskUsecase.UpdateTask(id, req.Title, req.Description, req.DueDate, req.Status)
	if err != nil {
		status, body := mapDomainErrorToHTTP(err)
		c.JSON(status, body)
		return
	}
	// Map Domain task to response DTO
	c.JSON(http.StatusOK, mapDomainTaskToResponse(updatedTask))
}

// DeleteTask handles DELETE /tasks/:id (requires admin auth via middleware).
func (ctr *Controller) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	err = ctr.TaskUsecase.DeleteTask(id)
	if err != nil {
		status, body := mapDomainErrorToHTTP(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}