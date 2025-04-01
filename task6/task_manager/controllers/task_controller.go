package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"task_manager/data"
	"task_manager/models"
)

// TaskController handles task-related HTTP requests.
type TaskController struct {
	service *data.TaskService
}

// NewTaskController creates a new TaskController.
func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{service: service}
}

// GetTasks handles GET /tasks.
func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks, err := tc.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask handles GET /tasks/:id.
func (tc *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := tc.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask handles POST /tasks.
func (tc *TaskController) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	createdTask, err := tc.service.Create(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdTask)
}

// UpdateTask handles PUT /tasks/:id.
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	task, err := tc.service.Update(id, updatedTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /tasks/:id.
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := tc.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
