package repositories

import (
	"sync"
	"task_manager/Domain" // Depends only on Domain
)

// --- Task Repository Interface ---
// Defines data access methods for tasks.
type TaskRepository interface {
	Save(task domain.Task) (domain.Task, error)
	FindByID(id int) (domain.Task, error)
	FindAll() ([]domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	Delete(id int) error
}

// --- In-Memory Implementation ---

type inMemoryTaskRepository struct {
	tasks      map[int]domain.Task
	nextTaskID int
	mutex      sync.RWMutex // Use RWMutex
}

// NewInMemoryTaskRepository creates a new in-memory task repository.
func NewInMemoryTaskRepository() TaskRepository {
	return &inMemoryTaskRepository{
		tasks:      make(map[int]domain.Task),
		nextTaskID: 1,
	}
}

// Save stores a new task.
func (r *inMemoryTaskRepository) Save(task domain.Task) (domain.Task, error) {
	r.mutex.Lock() // Write lock needed
	defer r.mutex.Unlock()

	task.ID = r.nextTaskID
	r.nextTaskID++
	r.tasks[task.ID] = task
	return task, nil
}

// FindByID retrieves a task by its ID.
func (r *inMemoryTaskRepository) FindByID(id int) (domain.Task, error) {
	r.mutex.RLock() // Read lock sufficient
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return domain.Task{}, domain.ErrTaskNotFound
	}
	return task, nil
}

// FindAll retrieves all stored tasks.
func (r *inMemoryTaskRepository) FindAll() ([]domain.Task, error) {
	r.mutex.RLock() // Read lock sufficient
	defer r.mutex.RUnlock()

	// Allocate slice with capacity for efficiency
	allTasks := make([]domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks, nil
}

// Update modifies an existing task.
func (r *inMemoryTaskRepository) Update(updatedTask domain.Task) (domain.Task, error) {
	r.mutex.Lock() // Write lock needed
	defer r.mutex.Unlock()

	// Check if task exists before updating
	_, exists := r.tasks[updatedTask.ID]
	if !exists {
		return domain.Task{}, domain.ErrTaskNotFound
	}

	// Replace the entire task object in the map
	r.tasks[updatedTask.ID] = updatedTask
	return updatedTask, nil
}

// Delete removes a task by its ID.
func (r *inMemoryTaskRepository) Delete(id int) error {
	r.mutex.Lock() // Write lock needed
	defer r.mutex.Unlock()

	// Check if task exists before deleting
	if _, exists := r.tasks[id]; !exists {
		return domain.ErrTaskNotFound
	}
	delete(r.tasks, id)
	return nil
}