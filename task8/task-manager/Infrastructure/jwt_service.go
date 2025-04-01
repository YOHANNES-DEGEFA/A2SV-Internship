package infrastructure

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	// No direct domain dependency needed here if claims are self-contained
	// "task_manager/domain"
)

// --- JWT Claims ---
// These are the standard JWT claims plus custom ones for our app.
type JWTClaims struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// --- JWT Service Interface ---
// Defines operations for JWT handling.
type JWTService interface {
	GenerateToken(userID int, username, role string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

// --- JWT Service Implementation ---
// Uses github.com/golang-jwt/jwt/v4.
type jwtGoJWTService struct {
	secretKey []byte
	issuer    string
	expiryDur time.Duration
}

// NewJwtGoJWTService creates a new jwtGoJWTService.
// Secret key, issuer, and expiry should ideally come from config.
func NewJwtGoJWTService(secretKey, issuer string, expiry time.Duration) JWTService {
	return &jwtGoJWTService{
		secretKey: []byte(secretKey),
		issuer:    issuer,
		expiryDur: expiry,
	}
}

// GenerateToken creates a new JWT token for a user.
func (s *jwtGoJWTService) GenerateToken(userID int, username, role string) (string, error) {
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiryDur)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", userID), // User ID as subject
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		// Log internal error (not shown) and return a generic error indication
		return "", fmt.Errorf("could not sign token: %w", err)
	}
	return tokenString, nil
}

// ValidateToken parses and validates a JWT token string.
func (s *jwtGoJWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		return s.secretKey, nil
	})

	// Handle parsing/validation errors
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("malformed token")
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			// Consider specific errors if frontend needs to differentiate
			return nil, errors.New("token is expired or not valid yet")
		} else if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
			return nil, errors.New("invalid token signature")
		}
		// Generic invalid token error for other validation issues
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check if token is valid and claims can be extracted
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Optionally, verify issuer again (though ParseWithClaims does this if set)
		if !claims.VerifyIssuer(s.issuer, true) {
			return nil, fmt.Errorf("invalid token issuer")
		}
		return claims, nil
	}

	// Fallback invalid token error
	return nil, errors.New("invalid token")
}