package usecases

import (
	"task_manager/Domain"
	"task_manager/Infrastructure" // Depends on Infrastructure interfaces
	"task_manager/Repositories"  // Depends on Repository interfaces
)

// --- User Usecase Interface ---
// Defines application-specific user operations.
type UserUsecase interface {
	Register(username, password string) (domain.User, error) // Returns domain user (no hash)
	Login(username, password string) (string, error)         // Returns JWT token string
	Promote(usernameToPromote string) error                  // Caller authorization is handled by delivery layer
}

// --- User Usecase Implementation ---
type userUsecase struct {
	userRepo      repositories.UserRepository    // Interface dependency
	passwordSvc   infrastructure.PasswordService // Interface dependency
	jwtSvc        infrastructure.JWTService      // Interface dependency
}

// NewUserUsecase creates a new UserUsecase instance.
func NewUserUsecase(
	userRepo repositories.UserRepository,
	passwordSvc infrastructure.PasswordService,
	jwtSvc infrastructure.JWTService,
) UserUsecase {
	return &userUsecase{
		userRepo:      userRepo,
		passwordSvc:   passwordSvc,
		jwtSvc:        jwtSvc,
	}
}

// Register handles the user registration business logic.
func (uc *userUsecase) Register(username, password string) (domain.User, error) {
	// 1. Basic Input Validation
	if username == "" || password == "" {
		return domain.User{}, domain.ErrBadRequest // Use domain error for bad input
	}

	// 2. Check if username already exists (primary check)
	_, _, err := uc.userRepo.FindByUsername(username)
	if err == nil { // User *was* found
		return domain.User{}, domain.ErrUsernameExists
	}
	if err != domain.ErrUserNotFound { // Unexpected repository error
		// Log internal error details (not shown)
		return domain.User{}, domain.ErrInternalServer
	}
	// If err is ErrUserNotFound, we can proceed.

	// 3. Hash the password
	hashedPassword, err := uc.passwordSvc.Hash(password)
	if err != nil {
		// Log internal error details (not shown)
		return domain.User{}, domain.ErrInternalServer
	}

	// 4. Determine the role (first user is admin)
	count, err := uc.userRepo.CountUsers()
	if err != nil {
		// Log internal error details (not shown)
		return domain.User{}, domain.ErrInternalServer
	}
	role := domain.RoleUser
	if count == 0 {
		role = domain.RoleAdmin
	}

	// 5. Prepare domain user object
	userToSave := domain.User{
		Username: username,
		Role:     role,
		// ID will be assigned by the repository
	}

	// 6. Save the user via the repository
	savedUser, err := uc.userRepo.Save(userToSave, hashedPassword)
	if err != nil {
		// Could be ErrUsernameExists (if race condition) or other repo error
		if err == domain.ErrUsernameExists {
			return domain.User{}, err // Return specific error
		}
		// Log internal error details (not shown)
		return domain.User{}, domain.ErrInternalServer
	}

	// 7. Return the created domain user (without the hash)
	return savedUser, nil
}

// Login handles the user login business logic.
func (uc *userUsecase) Login(username, password string) (string, error) {
	// 1. Find user by username
	user, hashedPassword, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		if err == domain.ErrUserNotFound {
			// Return generic invalid credentials error for security
			return "", domain.ErrInvalidCredentials
		}
		// Log internal error details (not shown)
		return "", domain.ErrInternalServer
	}

	// 2. Compare the provided password with the stored hash
	err = uc.passwordSvc.Compare(hashedPassword, password)
	if err != nil {
		// Password does not match (bcrypt.CompareHashAndPassword returns error on mismatch)
		return "", domain.ErrInvalidCredentials // Generic error
	}

	// 3. Credentials are valid, generate JWT token
	token, err := uc.jwtSvc.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		// Log internal error details (not shown)
		return "", domain.ErrInternalServer
	}

	// 4. Return the generated token string
	return token, nil
}

// Promote handles the logic for promoting a user to admin.
// Assumes the *caller* has already been authorized as admin by middleware.
func (uc *userUsecase) Promote(usernameToPromote string) error {
	// 1. Find the user to be promoted
	userToPromote, _, err := uc.userRepo.FindByUsername(usernameToPromote)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return err // Specific error is okay here
		}
		// Log internal error details (not shown)
		return domain.ErrInternalServer
	}

	// 2. Check if user is already admin (optional, maybe just idempotent)
	if userToPromote.Role == domain.RoleAdmin {
		return nil // Already admin, no error
	}

	// 3. Update the user's role in the repository
	err = uc.userRepo.UpdateRole(userToPromote.ID, domain.RoleAdmin)
	if err != nil {
		// Could be ErrUserNotFound if user deleted between find and update (race condition)
		if err == domain.ErrUserNotFound {
			return err
		}
		// Log internal error details (not shown)
		return domain.ErrInternalServer
	}

	return nil // Success
}