package usecases_test // Use _test package for black-box testing of usecases

import (
	"errors"
	"testing"
	"task_manager/domain"
	mockinfra "task_manager/Mocks/infrastructure" // Aliased import for mocks
	mockrepo "task_manager/Mocks/repositories"   // Aliased import for mocks
	"task_manager/usecases"                      // The package being tested

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserUsecase_Register(t *testing.T) {
	mockUserRepo := mockrepo.NewMockUserRepository(t)
	mockPasswordSvc := mockinfra.NewMockPasswordService(t)
	mockJWTSvc := mockinfra.NewMockJWTService(t) // Not used in Register, but needed for constructor

	uc := usecases.NewUserUsecase(mockUserRepo, mockPasswordSvc, mockJWTSvc)

	username := "testuser"
	password := "password123"
	hashedPassword := "hashed_password"

	t.Run("Success - First user becomes admin", func(t *testing.T) {
		// Arrange
		mockUserRepo.On("FindByUsername", username).Return(domain.User{}, "", domain.ErrUserNotFound).Once()
		mockPasswordSvc.On("Hash", password).Return(hashedPassword, nil).Once()
		mockUserRepo.On("CountUsers").Return(0, nil).Once() // First user
		
		// Define expected user passed to Save (ID will be assigned by repo mock)
		expectedUserArg := domain.User{Username: username, Role: domain.RoleAdmin}
		// Define user returned by Save mock (with ID)
		expectedReturnedUser := domain.User{ID: 1, Username: username, Role: domain.RoleAdmin}
		mockUserRepo.On("Save", expectedUserArg, hashedPassword).Return(expectedReturnedUser, nil).Once()

		// Act
		user, err := uc.Register(username, password)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedReturnedUser, user)
		// Mocks are verified automatically by cleanup function
	})

	t.Run("Success - Subsequent user becomes user", func(t *testing.T) {
		// Arrange
		mockUserRepo.On("FindByUsername", username).Return(domain.User{}, "", domain.ErrUserNotFound).Once()
		mockPasswordSvc.On("Hash", password).Return(hashedPassword, nil).Once()
		mockUserRepo.On("CountUsers").Return(1, nil).Once() // Not the first user

		expectedUserArg := domain.User{Username: username, Role: domain.RoleUser}
		expectedReturnedUser := domain.User{ID: 2, Username: username, Role: domain.RoleUser}
		mockUserRepo.On("Save", expectedUserArg, hashedPassword).Return(expectedReturnedUser, nil).Once()

		// Act
		user, err := uc.Register(username, password)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedReturnedUser, user)
	})


	t.Run("Failure - Username already exists", func(t *testing.T) {
		// Arrange
		existingUser := domain.User{ID: 1, Username: username, Role: domain.RoleUser}
		mockUserRepo.On("FindByUsername", username).Return(existingUser, "somehash", nil).Once() // User found

		// Act
		user, err := uc.Register(username, password)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrUsernameExists, err)
		require.Empty(t, user) // Ensure no user is returned on error
		// Ensure other mocks were NOT called
		mockPasswordSvc.AssertNotCalled(t, "Hash", mock.Anything)
		mockUserRepo.AssertNotCalled(t, "CountUsers")
		mockUserRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Failure - Empty username", func(t *testing.T) {
        // Arrange - No mock setup needed as validation happens first

        // Act
        user, err := uc.Register("", password)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrBadRequest, err)
        require.Empty(t, user)
		// Ensure mocks were NOT called
		mockUserRepo.AssertNotCalled(t, "FindByUsername", mock.Anything)
		mockPasswordSvc.AssertNotCalled(t, "Hash", mock.Anything)
		mockUserRepo.AssertNotCalled(t, "CountUsers")
		mockUserRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

	t.Run("Failure - Password Hash error", func(t *testing.T) {
		// Arrange
		hashError := errors.New("bcrypt failure")
		mockUserRepo.On("FindByUsername", username).Return(domain.User{}, "", domain.ErrUserNotFound).Once()
		mockPasswordSvc.On("Hash", password).Return("", hashError).Once() // Simulate hashing error

		// Act
		user, err := uc.Register(username, password)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrInternalServer, err) // Should map to internal error
		require.Empty(t, user)
		mockUserRepo.AssertNotCalled(t, "CountUsers")
		mockUserRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Failure - CountUsers error", func(t *testing.T) {
		// Arrange
		countError := errors.New("db connection lost")
		mockUserRepo.On("FindByUsername", username).Return(domain.User{}, "", domain.ErrUserNotFound).Once()
		mockPasswordSvc.On("Hash", password).Return(hashedPassword, nil).Once()
		mockUserRepo.On("CountUsers").Return(0, countError).Once() // Simulate count error

		// Act
		user, err := uc.Register(username, password)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrInternalServer, err)
		require.Empty(t, user)
		mockUserRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Failure - Save error", func(t *testing.T) {
		// Arrange
		saveError := errors.New("db write failed")
		mockUserRepo.On("FindByUsername", username).Return(domain.User{}, "", domain.ErrUserNotFound).Once()
		mockPasswordSvc.On("Hash", password).Return(hashedPassword, nil).Once()
		mockUserRepo.On("CountUsers").Return(0, nil).Once()

		expectedUserArg := domain.User{Username: username, Role: domain.RoleAdmin}
		mockUserRepo.On("Save", expectedUserArg, hashedPassword).Return(domain.User{}, saveError).Once() // Simulate save error

		// Act
		user, err := uc.Register(username, password)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrInternalServer, err)
		require.Empty(t, user)
	})
}


func TestUserUsecase_Login(t *testing.T) {
    mockUserRepo := mockrepo.NewMockUserRepository(t)
	mockPasswordSvc := mockinfra.NewMockPasswordService(t)
	mockJWTSvc := mockinfra.NewMockJWTService(t)

	uc := usecases.NewUserUsecase(mockUserRepo, mockPasswordSvc, mockJWTSvc)

    username := "testuser"
    password := "password123"
    hashedPassword := "hashed_password"
    userID := 1
    userRole := domain.RoleUser
    expectedToken := "valid.jwt.token"

    t.Run("Success", func(t *testing.T) {
        // Arrange
        foundUser := domain.User{ID: userID, Username: username, Role: userRole}
        mockUserRepo.On("FindByUsername", username).Return(foundUser, hashedPassword, nil).Once()
        mockPasswordSvc.On("Compare", hashedPassword, password).Return(nil).Once() // Compare success
        mockJWTSvc.On("GenerateToken", userID, username, userRole).Return(expectedToken, nil).Once()

        // Act
        token, err := uc.Login(username, password)

        // Assert
        require.NoError(t, err)
        require.Equal(t, expectedToken, token)
    })

    t.Run("Failure - User not found", func(t *testing.T) {
        // Arrange
        mockUserRepo.On("FindByUsername", username).Return(domain.User{}, "", domain.ErrUserNotFound).Once()

        // Act
        token, err := uc.Login(username, password)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrInvalidCredentials, err) // Maps to invalid credentials
        require.Empty(t, token)
        mockPasswordSvc.AssertNotCalled(t, "Compare", mock.Anything, mock.Anything)
        mockJWTSvc.AssertNotCalled(t, "GenerateToken", mock.Anything, mock.Anything, mock.Anything)
    })

     t.Run("Failure - Incorrect password", func(t *testing.T) {
        // Arrange
        foundUser := domain.User{ID: userID, Username: username, Role: userRole}
		compareError := errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password") // Actual bcrypt error
        mockUserRepo.On("FindByUsername", username).Return(foundUser, hashedPassword, nil).Once()
        mockPasswordSvc.On("Compare", hashedPassword, password).Return(compareError).Once() // Compare fails

        // Act
        token, err := uc.Login(username, password)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrInvalidCredentials, err) // Maps to invalid credentials
        require.Empty(t, token)
        mockJWTSvc.AssertNotCalled(t, "GenerateToken", mock.Anything, mock.Anything, mock.Anything)
    })

     t.Run("Failure - FindByUsername repository error", func(t *testing.T) {
        // Arrange
		repoError := errors.New("db connection error")
        mockUserRepo.On("FindByUsername", username).Return(domain.User{}, "", repoError).Once()

        // Act
        token, err := uc.Login(username, password)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrInternalServer, err) // Maps to internal server error
        require.Empty(t, token)
        mockPasswordSvc.AssertNotCalled(t, "Compare", mock.Anything, mock.Anything)
        mockJWTSvc.AssertNotCalled(t, "GenerateToken", mock.Anything, mock.Anything, mock.Anything)
    })

	t.Run("Failure - JWT Generation error", func(t *testing.T) {
        // Arrange
		jwtError := errors.New("failed to sign token")
        foundUser := domain.User{ID: userID, Username: username, Role: userRole}
        mockUserRepo.On("FindByUsername", username).Return(foundUser, hashedPassword, nil).Once()
        mockPasswordSvc.On("Compare", hashedPassword, password).Return(nil).Once()
        mockJWTSvc.On("GenerateToken", userID, username, userRole).Return("", jwtError).Once() // Simulate JWT error

        // Act
        token, err := uc.Login(username, password)

        // Assert
        require.Error(t, err)
        require.Equal(t, domain.ErrInternalServer, err) // Maps to internal server error
        require.Empty(t, token)
    })
}


func TestUserUsecase_Promote(t *testing.T) {
	mockUserRepo := mockrepo.NewMockUserRepository(t)
	// Promote doesn't directly use password or JWT service, but need mocks for constructor
	mockPasswordSvc := mockinfra.NewMockPasswordService(t)
	mockJWTSvc := mockinfra.NewMockJWTService(t)

	uc := usecases.NewUserUsecase(mockUserRepo, mockPasswordSvc, mockJWTSvc)

	usernameToPromote := "normaluser"
	userIDToPromote := 5

	t.Run("Success", func(t *testing.T) {
		// Arrange
		userToPromote := domain.User{ID: userIDToPromote, Username: usernameToPromote, Role: domain.RoleUser}
		mockUserRepo.On("FindByUsername", usernameToPromote).Return(userToPromote, "somehash", nil).Once()
		mockUserRepo.On("UpdateRole", userIDToPromote, domain.RoleAdmin).Return(nil).Once()

		// Act
		err := uc.Promote(usernameToPromote)

		// Assert
		require.NoError(t, err)
	})

	t.Run("Success - User already admin (idempotent)", func(t *testing.T) {
		// Arrange
		userAlreadyAdmin := domain.User{ID: userIDToPromote, Username: usernameToPromote, Role: domain.RoleAdmin}
		mockUserRepo.On("FindByUsername", usernameToPromote).Return(userAlreadyAdmin, "somehash", nil).Once()
		// UpdateRole should NOT be called if already admin
		// mockUserRepo.AssertNotCalled(t, "UpdateRole", mock.Anything, mock.Anything) // Assert this after Act

		// Act
		err := uc.Promote(usernameToPromote)

		// Assert
		require.NoError(t, err)
		mockUserRepo.AssertNotCalled(t, "UpdateRole", mock.AnythingOfType("int"), mock.AnythingOfType("string"))
	})

	t.Run("Failure - User to promote not found", func(t *testing.T) {
		// Arrange
		mockUserRepo.On("FindByUsername", usernameToPromote).Return(domain.User{}, "", domain.ErrUserNotFound).Once()

		// Act
		err := uc.Promote(usernameToPromote)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrUserNotFound, err)
		mockUserRepo.AssertNotCalled(t, "UpdateRole", mock.Anything, mock.Anything)
	})

	t.Run("Failure - FindByUsername repository error", func(t *testing.T) {
		// Arrange
		repoError := errors.New("db connection error")
		mockUserRepo.On("FindByUsername", usernameToPromote).Return(domain.User{}, "", repoError).Once()

		// Act
		err := uc.Promote(usernameToPromote)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrInternalServer, err)
		mockUserRepo.AssertNotCalled(t, "UpdateRole", mock.Anything, mock.Anything)
	})

	t.Run("Failure - UpdateRole repository error", func(t *testing.T) {
		// Arrange
		updateError := errors.New("db write error")
		userToPromote := domain.User{ID: userIDToPromote, Username: usernameToPromote, Role: domain.RoleUser}
		mockUserRepo.On("FindByUsername", usernameToPromote).Return(userToPromote, "somehash", nil).Once()
		mockUserRepo.On("UpdateRole", userIDToPromote, domain.RoleAdmin).Return(updateError).Once()

		// Act
		err := uc.Promote(usernameToPromote)

		// Assert
		require.Error(t, err)
		require.Equal(t, domain.ErrInternalServer, err)
	})
}