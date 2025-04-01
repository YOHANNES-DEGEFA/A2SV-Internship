package infrastructure

import "golang.org/x/crypto/bcrypt"

// PasswordService defines the interface for password operations.
// This allows swapping the hashing algorithm later if needed.
type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

// bcryptPasswordService implements PasswordService using bcrypt.
type bcryptPasswordService struct{}

// NewBcryptPasswordService creates a new bcryptPasswordService.
func NewBcryptPasswordService() PasswordService {
	return &bcryptPasswordService{}
}

// Hash generates a bcrypt hash for the password.
func (s *bcryptPasswordService) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err // Propagate error
	}
	return string(hash), nil
}

// Compare compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on success, error on failure (including mismatch).
func (s *bcryptPasswordService) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}