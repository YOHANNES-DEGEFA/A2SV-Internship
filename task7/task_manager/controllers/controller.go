package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"task_manager/data"
	"task_manager/models"
	// "task_manager/middleware"
)

// jwtSecret should match the secret in the middleware.
var jwtSecret = []byte("your_secret_key_here")

// Controller holds references to both user and task services.
type Controller struct {
	UserService *data.UserService
	TaskService *data.TaskService
}

// NewController creates a new Controller instance.
func NewController(userService *data.UserService, taskService *data.TaskService) *Controller {
	return &Controller{
		UserService: userService,
		TaskService: taskService,
	}
}

// Register handles POST /register to create a new user.
func (ctr *Controller) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, err := ctr.UserService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles POST /login to authenticate a user and return a JWT token.
func (ctr *Controller) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, err := ctr.UserService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Promote handles POST /promote to promote a user to admin.
// Only an admin can perform this action.
func (ctr *Controller) Promote(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	err := ctr.UserService.Promote(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}

// GetTasks handles GET /tasks. Accessible by any authenticated user.
func (ctr *Controller) GetTasks(c *gin.Context) {
	tasks := ctr.TaskService.GetAll()
	c.JSON(http.StatusOK, tasks)
}

// GetTask handles GET /tasks/:id. Accessible by any authenticated user.
func (ctr *Controller) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	task, err := ctr.TaskService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask handles POST /tasks. Accessible only by admin.
func (ctr *Controller) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	createdTask := ctr.TaskService.Create(task)
	c.JSON(http.StatusCreated, createdTask)
}

// UpdateTask handles PUT /tasks/:id. Accessible only by admin.
func (ctr *Controller) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	task, err := ctr.TaskService.Update(id, updatedTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /tasks/:id. Accessible only by admin.
func (ctr *Controller) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	err = ctr.TaskService.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
