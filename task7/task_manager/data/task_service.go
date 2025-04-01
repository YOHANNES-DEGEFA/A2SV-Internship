package data

import (
	"errors"
	"sync"

	"task_manager/models"
)

// TaskService provides methods to manage tasks.
type TaskService struct {
	tasks      map[int]models.Task
	nextTaskID int
	mutex      sync.Mutex
}

// NewTaskService creates a new TaskService.
func NewTaskService() *TaskService {
	return &TaskService{
		tasks:      make(map[int]models.Task),
		nextTaskID: 1,
	}
}

// GetAll returns all tasks.
func (ts *TaskService) GetAll() []models.Task {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	tasks := []models.Task{}
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetByID returns a task by its ID.
func (ts *TaskService) GetByID(id int) (models.Task, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

// Create adds a new task.
func (ts *TaskService) Create(task models.Task) models.Task {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	task.ID = ts.nextTaskID
	ts.nextTaskID++
	ts.tasks[task.ID] = task
	return task
}

// Update modifies an existing task.
func (ts *TaskService) Update(id int, updatedTask models.Task) (models.Task, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}

	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	task.DueDate = updatedTask.DueDate
	task.Status = updatedTask.Status
	ts.tasks[id] = task
	return task, nil
}

// Delete removes a task.
func (ts *TaskService) Delete(id int) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if _, exists := ts.tasks[id]; !exists {
		return errors.New("task not found")
	}
	delete(ts.tasks, id)
	return nil
}
