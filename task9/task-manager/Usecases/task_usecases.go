package usecases

import (
	"task_manager/Domain"
	"task_manager/Repositories" // Depends on Repository interface
)

// --- Task Usecase Interface ---
// Defines application-specific task operations.
type TaskUsecase interface {
	CreateTask(title, description, dueDate, status string) (domain.Task, error)
	GetTaskByID(id int) (domain.Task, error)
	GetAllTasks() ([]domain.Task, error)
	UpdateTask(id int, title, description, dueDate, status string) (domain.Task, error)
	DeleteTask(id int) error
	// Note: Authorization (who can call these) is handled by the delivery layer.
}

// --- Task Usecase Implementation ---
type taskUsecase struct {
	taskRepo repositories.TaskRepository // Interface dependency
}

// NewTaskUsecase creates a new TaskUsecase instance.
func NewTaskUsecase(taskRepo repositories.TaskRepository) TaskUsecase {
	return &taskUsecase{taskRepo: taskRepo}
}

// CreateTask handles the business logic for creating a new task.
func (uc *taskUsecase) CreateTask(title, description, dueDate, status string) (domain.Task, error) {
	// 1. Basic Input Validation (can be more complex)
	if title == "" {
		return domain.Task{}, domain.ErrBadRequest // Use domain error
	}
	// Add more validation for dueDate format, status enum etc. if needed

	// 2. Prepare domain task object
	taskToSave := domain.Task{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
		// ID assigned by repository
	}

	// 3. Save via repository
	savedTask, err := uc.taskRepo.Save(taskToSave)
	if err != nil {
		// Log internal error (not shown)
		return domain.Task{}, domain.ErrInternalServer
	}

	return savedTask, nil
}

// GetTaskByID handles retrieving a single task.
func (uc *taskUsecase) GetTaskByID(id int) (domain.Task, error) {
	task, err := uc.taskRepo.FindByID(id)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			return domain.Task{}, err // Propagate not found error
		}
		// Log internal error (not shown)
		return domain.Task{}, domain.ErrInternalServer
	}
	return task, nil
}

// GetAllTasks handles retrieving all tasks.
func (uc *taskUsecase) GetAllTasks() ([]domain.Task, error) {
	tasks, err := uc.taskRepo.FindAll()
	if err != nil {
		// Log internal error (not shown)
		return nil, domain.ErrInternalServer
	}
	return tasks, nil
}

// UpdateTask handles the business logic for updating an existing task.
func (uc *taskUsecase) UpdateTask(id int, title, description, dueDate, status string) (domain.Task, error) {
	// 1. Basic Input Validation
	if title == "" {
		return domain.Task{}, domain.ErrBadRequest
	}
	// Add more validation as needed

    // 2. Check if task exists (optional, depends if repo.Update checks)
    //    It's often good practice for the use case to ensure the entity exists first.
    _, err := uc.taskRepo.FindByID(id)
    if err != nil {
        if err == domain.ErrTaskNotFound {
            return domain.Task{}, err // Propagate not found
        }
        return domain.Task{}, domain.ErrInternalServer // Other FindByID error
    }


	// 3. Prepare the updated domain task object
	taskToUpdate := domain.Task{
		ID:          id, // Crucial: ID must be set for the update target
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
	}

	// 4. Update via repository
	updatedTask, err := uc.taskRepo.Update(taskToUpdate)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			// Repository might return this if update implementation checks existence again
			return domain.Task{}, err
		}
		// Log internal error (not shown)
		return domain.Task{}, domain.ErrInternalServer
	}

	return updatedTask, nil
}

// DeleteTask handles the business logic for deleting a task.
func (uc *taskUsecase) DeleteTask(id int) error {
	err := uc.taskRepo.Delete(id)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			return err // Propagate not found error
		}
		// Log internal error (not shown)
		return domain.ErrInternalServer
	}
	return nil // Success
}