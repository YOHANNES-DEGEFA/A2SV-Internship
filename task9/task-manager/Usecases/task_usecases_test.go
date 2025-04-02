package usecases_test

import (
	"errors"
	"testing"
	"task_manager/domain"
	mockrepo "task_manager/Mocks/repositories" // Aliased import for mocks
	"task_manager/usecases"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTaskUsecase_CreateTask(t *testing.T) {
	mockTaskRepo := mockrepo.NewMockTaskRepository(t)
	uc := usecases.NewTaskUsecase(mockTaskRepo)

	title := "New Task"
	description := "Task Description"
	dueDate := "2024-12-31"
	status := "pending"

	t.Run("Success", func(t *testing.T) {
		// Arrange
		taskToSave := domain.Task{
			Title:       title,
			Description: description,
			DueDate:     dueDate,
			Status:      status,
			// ID is assigned by repo
		}
		savedTask := domain.Task{
			ID:          1, // ID assigned by repo mock
			Title:       title,
			Description: description,
			DueDate:     dueDate,
			Status:      status,
		}
		mockTaskRepo.On("Save", taskToSave).Return(savedTask, nil).Once()

		// Act
		result, err := uc.CreateTask(title, description, dueDate, status)

		// Assert
		require.NoError(t, err)
		require.Equal(t, savedTask, result)
	})

	t.Run("Failure - Empty title", func(t *testing.T) {
		// Arrange - No mock setup needed

		// Act
		result, err := uc.CreateTask("", description, dueDate, status)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrBadRequest, err)
		require.Empty(t, result)
		mockTaskRepo.AssertNotCalled(t, "Save", mock.Anything)
	})

	t.Run("Failure - Repository Save error", func(t *testing.T) {
		// Arrange
		saveError := errors.New("db write error")
		taskToSave := domain.Task{
			Title:       title,
			Description: description,
			DueDate:     dueDate,
			Status:      status,
		}
		mockTaskRepo.On("Save", taskToSave).Return(domain.Task{}, saveError).Once()

		// Act
		result, err := uc.CreateTask(title, description, dueDate, status)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrInternalServer, err)
		require.Empty(t, result)
	})
}

func TestTaskUsecase_GetTaskByID(t *testing.T) {
	mockTaskRepo := mockrepo.NewMockTaskRepository(t)
	uc := usecases.NewTaskUsecase(mockTaskRepo)
	taskID := 1

	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedTask := domain.Task{ID: taskID, Title: "Test Task"}
		mockTaskRepo.On("FindByID", taskID).Return(expectedTask, nil).Once()

		// Act
		result, err := uc.GetTaskByID(taskID)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedTask, result)
	})

	t.Run("Failure - Task not found", func(t *testing.T) {
		// Arrange
		mockTaskRepo.On("FindByID", taskID).Return(domain.Task{}, domain.ErrTaskNotFound).Once()

		// Act
		result, err := uc.GetTaskByID(taskID)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrTaskNotFound, err)
		require.Empty(t, result)
	})

	t.Run("Failure - Repository error", func(t *testing.T) {
		// Arrange
		repoError := errors.New("db connection error")
		mockTaskRepo.On("FindByID", taskID).Return(domain.Task{}, repoError).Once()

		// Act
		result, err := uc.GetTaskByID(taskID)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrInternalServer, err)
		require.Empty(t, result)
	})
}

func TestTaskUsecase_GetAllTasks(t *testing.T) {
	mockTaskRepo := mockrepo.NewMockTaskRepository(t)
	uc := usecases.NewTaskUsecase(mockTaskRepo)

	t.Run("Success - Empty list", func(t *testing.T) {
		// Arrange
		mockTaskRepo.On("FindAll").Return([]domain.Task{}, nil).Once()

		// Act
		result, err := uc.GetAllTasks()

		// Assert
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("Success - Non-empty list", func(t *testing.T) {
		// Arrange
		expectedTasks := []domain.Task{
			{ID: 1, Title: "Task 1"},
			{ID: 2, Title: "Task 2"},
		}
		mockTaskRepo.On("FindAll").Return(expectedTasks, nil).Once()

		// Act
		result, err := uc.GetAllTasks()

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedTasks, result)
	})

	t.Run("Failure - Repository error", func(t *testing.T) {
		// Arrange
		repoError := errors.New("db read error")
		mockTaskRepo.On("FindAll").Return(nil, repoError).Once()

		// Act
		result, err := uc.GetAllTasks()

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrInternalServer, err)
		require.Nil(t, result)
	})
}


func TestTaskUsecase_UpdateTask(t *testing.T) {
    mockTaskRepo := mockrepo.NewMockTaskRepository(t)
    uc := usecases.NewTaskUsecase(mockTaskRepo)

    taskID := 1
    newTitle := "Updated Task Title"
    newDesc := "Updated description"
    newDueDate := "2025-01-01"
    newStatus := "in-progress"

	existingTask := domain.Task{ID: taskID, Title: "Old Title"} // Task returned by initial FindByID
	taskToUpdateArg := domain.Task{ // Task passed to repo.Update
        ID:          taskID,
        Title:       newTitle,
        Description: newDesc,
        DueDate:     newDueDate,
        Status:      newStatus,
    }
	returnedUpdatedTask := taskToUpdateArg // Task returned by repo.Update


    t.Run("Success", func(t *testing.T) {
        // Arrange
		mockTaskRepo.On("FindByID", taskID).Return(existingTask, nil).Once() // Initial check passes
        mockTaskRepo.On("Update", taskToUpdateArg).Return(returnedUpdatedTask, nil).Once()

        // Act
        result, err := uc.UpdateTask(taskID, newTitle, newDesc, newDueDate, newStatus)

        // Assert
        require.NoError(t, err)
        require.Equal(t, returnedUpdatedTask, result)
    })

    t.Run("Failure - Empty title", func(t *testing.T) {
        // Arrange - no mocks needed

        // Act
        result, err := uc.UpdateTask(taskID, "", newDesc, newDueDate, newStatus)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrBadRequest, err)
        require.Empty(t, result)
		mockTaskRepo.AssertNotCalled(t, "FindByID", mock.Anything)
        mockTaskRepo.AssertNotCalled(t, "Update", mock.Anything)
    })

	t.Run("Failure - Task not found on initial check", func(t *testing.T) {
        // Arrange
		mockTaskRepo.On("FindByID", taskID).Return(domain.Task{}, domain.ErrTaskNotFound).Once() // Initial check fails

        // Act
        result, err := uc.UpdateTask(taskID, newTitle, newDesc, newDueDate, newStatus)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrTaskNotFound, err)
        require.Empty(t, result)
        mockTaskRepo.AssertNotCalled(t, "Update", mock.Anything)
    })

	 t.Run("Failure - Repository error on initial check", func(t *testing.T) {
        // Arrange
		repoError := errors.New("db connection error")
		mockTaskRepo.On("FindByID", taskID).Return(domain.Task{}, repoError).Once() // Initial check fails

        // Act
        result, err := uc.UpdateTask(taskID, newTitle, newDesc, newDueDate, newStatus)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrInternalServer, err)
        require.Empty(t, result)
        mockTaskRepo.AssertNotCalled(t, "Update", mock.Anything)
    })

    t.Run("Failure - Task not found on Update call", func(t *testing.T) {
        // Arrange
		mockTaskRepo.On("FindByID", taskID).Return(existingTask, nil).Once() // Initial check passes
        mockTaskRepo.On("Update", taskToUpdateArg).Return(domain.Task{}, domain.ErrTaskNotFound).Once() // Update itself fails with not found

        // Act
        result, err := uc.UpdateTask(taskID, newTitle, newDesc, newDueDate, newStatus)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrTaskNotFound, err)
        require.Empty(t, result)
    })

    t.Run("Failure - Repository error on Update call", func(t *testing.T) {
        // Arrange
		repoError := errors.New("db write error")
		mockTaskRepo.On("FindByID", taskID).Return(existingTask, nil).Once() // Initial check passes
        mockTaskRepo.On("Update", taskToUpdateArg).Return(domain.Task{}, repoError).Once() // Update fails with repo error

        // Act
        result, err := uc.UpdateTask(taskID, newTitle, newDesc, newDueDate, newStatus)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrInternalServer, err)
        require.Empty(t, result)
    })
}


func TestTaskUsecase_DeleteTask(t *testing.T) {
	mockTaskRepo := mockrepo.NewMockTaskRepository(t)
	uc := usecases.NewTaskUsecase(mockTaskRepo)
	taskID := 1

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockTaskRepo.On("Delete", taskID).Return(nil).Once()

		// Act
		err := uc.DeleteTask(taskID)

		// Assert
		require.NoError(t, err)
	})

	t.Run("Failure - Task not found", func(t *testing.T) {
		// Arrange
		mockTaskRepo.On("Delete", taskID).Return(domain.ErrTaskNotFound).Once()

		// Act
		err := uc.DeleteTask(taskID)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrTaskNotFound, err)
	})

	t.Run("Failure - Repository error", func(t *testing.T) {
		// Arrange
		repoError := errors.New("db connection error")
		mockTaskRepo.On("Delete", taskID).Return(repoError).Once()

		// Act
		err := uc.DeleteTask(taskID)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrInternalServer, err)
	})
}