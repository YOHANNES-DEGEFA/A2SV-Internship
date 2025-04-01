package data

import (
	"errors"
	"sync"

	"task_manager/models"

	"golang.org/x/crypto/bcrypt"
)

// UserService provides user management functions.
type UserService struct {
	users      map[string]models.User // key is username
	nextUserID int
	mutex      sync.Mutex
}

// NewUserService creates a new UserService.
func NewUserService() *UserService {
	return &UserService{
		users:      make(map[string]models.User),
		nextUserID: 1,
	}
}

// Register registers a new user with a username and password.
// If no user exists, the first user becomes an admin.
func (us *UserService) Register(username, password string) (models.User, error) {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	if _, exists := us.users[username]; exists {
		return models.User{}, errors.New("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	role := "user"
	if len(us.users) == 0 {
		// First user becomes admin.
		role = "admin"
	}

	user := models.User{
		ID:           us.nextUserID,
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}
	us.users[username] = user
	us.nextUserID++
	return user, nil
}

// Login authenticates a user and returns the user if successful.
func (us *UserService) Login(username, password string) (models.User, error) {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	user, exists := us.users[username]
	if !exists {
		return models.User{}, errors.New("invalid credentials")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return models.User{}, errors.New("invalid credentials")
	}

	return user, nil
}

// Promote promotes a user to admin. Only admins can promote others.
func (us *UserService) Promote(username string) error {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	user, exists := us.users[username]
	if !exists {
		return errors.New("user not found")
	}
	user.Role = "admin"
	us.users[username] = user
	return nil
}
