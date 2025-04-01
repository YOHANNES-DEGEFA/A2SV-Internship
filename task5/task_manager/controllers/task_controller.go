package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"task_manager/data"
	"task_manager/models"
)

// TaskController holds the reference to the TaskService.
type TaskController struct {
	service *data.TaskService
}

// NewTaskController creates a new TaskController.
func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{service: service}
}

// GetTasks handles GET /tasks to list all tasks.
func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks := tc.service.GetAll()
	c.JSON(http.StatusOK, tasks)
}

// GetTask handles GET /tasks/:id to retrieve a task by ID.
func (tc *TaskController) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	task, err := tc.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask handles POST /tasks to create a new task.
func (tc *TaskController) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	createdTask := tc.service.Create(task)
	c.JSON(http.StatusCreated, createdTask)
}

// UpdateTask handles PUT /tasks/:id to update an existing task.
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	task, err := tc.service.Update(id, updatedTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /tasks/:id to remove a task.
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	err = tc.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
