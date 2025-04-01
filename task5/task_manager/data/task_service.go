package data

import (
	"errors"
	"sync"

	"task_manager/models"
)

// TaskService holds the in-memory data store and a mutex for safe concurrent access.
type TaskService struct {
	tasks      map[int]models.Task
	nextTaskID int
	mutex      sync.Mutex
}

// NewTaskService initializes a new TaskService.
func NewTaskService() *TaskService {
	return &TaskService{
		tasks:      make(map[int]models.Task),
		nextTaskID: 1,
	}
}

// GetAll returns all tasks.
func (s *TaskService) GetAll() []models.Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	tasks := []models.Task{}
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetByID returns a task by its ID.
func (s *TaskService) GetByID(id int) (models.Task, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

// Create adds a new task to the store.
func (s *TaskService) Create(task models.Task) models.Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task.ID = s.nextTaskID
	s.nextTaskID++
	s.tasks[task.ID] = task
	return task
}

// Update modifies an existing task.
func (s *TaskService) Update(id int, updatedTask models.Task) (models.Task, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	// Update fields.
	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	task.DueDate = updatedTask.DueDate
	task.Status = updatedTask.Status

	s.tasks[id] = task
	return task, nil
}

// Delete removes a task by its ID.
func (s *TaskService) Delete(id int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return errors.New("task not found")
	}
	delete(s.tasks, id)
	return nil
}
