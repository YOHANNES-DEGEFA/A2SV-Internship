package infrastructure

import (
	"net/http"
	"strings"
	"task_manager/Domain" // Need domain roles and errors

	"github.com/gin-gonic/gin"
)

// Constants for keys used in Gin context
const (
	ContextKeyUserID     = "userID"
	ContextKeyUsername   = "username"
	ContextKeyUserRole   = "role"
	ContextKeyIsVerified = "isVerified" // Flag to ensure Authenticate ran
)

// AuthMiddleware encapsulates dependencies needed for auth middleware.
type AuthMiddleware struct {
	jwtService JWTService // Depends on the JWTService interface
}

// NewAuthMiddleware creates a new AuthMiddleware instance.
func NewAuthMiddleware(jwtService JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// Authenticate is Gin middleware to verify the JWT token and set user info in context.
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := m.jwtService.ValidateToken(tokenStr)
		if err != nil {
			// Use the specific error message from the JWT service
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Token is valid, set user information into the context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyUserRole, claims.Role)
		c.Set(ContextKeyIsVerified, true) // Mark as verified

		c.Next() // Proceed to the next handler
	}
}

// AuthorizeAdmin is Gin middleware to ensure the authenticated user has the 'admin' role.
// It MUST run *after* the Authenticate middleware.
func AuthorizeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if authentication middleware ran successfully
		if verified, exists := c.Get(ContextKeyIsVerified); !exists || !verified.(bool) {
			// This indicates a configuration error (middleware order)
			// Log this scenario
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication context missing"})
			c.Abort()
			return
		}

		roleValue, exists := c.Get(ContextKeyUserRole)
		if !exists {
			// Should not happen if Authenticate succeeded
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found in context"})
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			// Should not happen
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role type in context"})
			c.Abort()
			return
		}

		// Check if the role is admin
		if role != domain.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": domain.ErrForbidden.Error()})
			c.Abort() // Stop processing
			return
		}

		// User is admin, allow proceeding
		c.Next()
	}
}