package repositories

import (
	"sync"
	"task_manager/Domain" // Depends only on Domain
)

// --- User Repository Interface ---
// Defines data access methods for users.
type UserRepository interface {
	Save(user domain.User, passwordHash string) (domain.User, error)
	FindByUsername(username string) (domain.User, string, error) // Returns User and PasswordHash
	FindByID(id int) (domain.User, error)                     // Needed for consistency checks sometimes
	UpdateRole(userID int, newRole string) error
	CountUsers() (int, error) // Needed for assigning initial admin role
}

// --- In-Memory Implementation ---

// internalUser is used only within this repository implementation
// to store the password hash alongside the domain user.
type internalUser struct {
	domain.User
	PasswordHash string
}

type inMemoryUserRepository struct {
	users      map[string]internalUser // Map username to internalUser for quick lookup
	usersByID  map[int]string        // Map ID to username for quick role updates by ID
	nextUserID int
	mutex      sync.RWMutex // Use RWMutex for better read performance
}

// NewInMemoryUserRepository creates a new in-memory user repository.
func NewInMemoryUserRepository() UserRepository {
	return &inMemoryUserRepository{
		users:      make(map[string]internalUser),
		usersByID:  make(map[int]string),
		nextUserID: 1,
	}
}

// Save stores a new user. Assumes username uniqueness check happened in usecase.
func (r *inMemoryUserRepository) Save(user domain.User, passwordHash string) (domain.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Double-check existence in case of race condition (though usecase should check first)
	if _, exists := r.users[user.Username]; exists {
		return domain.User{}, domain.ErrUsernameExists
	}

	// Assign ID
	user.ID = r.nextUserID
	r.nextUserID++

	internal := internalUser{
		User:         user,
		PasswordHash: passwordHash,
	}

	// Store in both maps
	r.users[user.Username] = internal
	r.usersByID[user.ID] = user.Username

	// Return the domain user (without hash)
	return user, nil
}

// FindByUsername retrieves a user and their password hash by username.
func (r *inMemoryUserRepository) FindByUsername(username string) (domain.User, string, error) {
	r.mutex.RLock() // Use Read Lock for finding
	defer r.mutex.RUnlock()

	internal, exists := r.users[username]
	if !exists {
		return domain.User{}, "", domain.ErrUserNotFound
	}
	// Return domain user and the hash
	return internal.User, internal.PasswordHash, nil
}

// FindByID retrieves a user by ID.
func (r *inMemoryUserRepository) FindByID(id int) (domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	username, idExists := r.usersByID[id]
	if !idExists {
		return domain.User{}, domain.ErrUserNotFound
	}

	// Should always exist if usersByID entry exists, but check for safety
	internal, userExists := r.users[username]
	if !userExists {
		// Data inconsistency - log this!
		return domain.User{}, domain.ErrInternalServer
	}

	return internal.User, nil
}


// UpdateRole changes the role of a user identified by userID.
func (r *inMemoryUserRepository) UpdateRole(userID int, newRole string) error {
	r.mutex.Lock() // Need Write Lock for updating
	defer r.mutex.Unlock()

	// Find the username associated with the ID
	username, idExists := r.usersByID[userID]
	if !idExists {
		return domain.ErrUserNotFound
	}

	// Get the user data
	internal, userExists := r.users[username]
	if !userExists {
		// Data inconsistency - log this!
		return domain.ErrInternalServer
	}

	// Update the role
	internal.User.Role = newRole
	r.users[username] = internal // Put the modified user back into the map

	return nil
}

// CountUsers returns the total number of registered users.
func (r *inMemoryUserRepository) CountUsers() (int, error) {
	r.mutex.RLock() // Use Read Lock
	defer r.mutex.RUnlock()
	count := len(r.users)
	return count, nil
}