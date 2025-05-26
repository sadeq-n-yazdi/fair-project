package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	// jwtSecret is the secret key used to sign JWT tokens
	// In a real application, this should be loaded from a secure configuration
	jwtSecret = []byte("your-secret-key")

	// tokenExpiration is the duration for which a token is valid
	tokenExpiration = 24 * time.Hour
)

// Claims represents the JWT claims
type Claims struct {
	Username string        `json:"username"`
	Roles    []models.Role `json:"roles"`
	jwt.RegisteredClaims
}

// SetJWTSecret sets the JWT secret key
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// SetTokenExpiration sets the token expiration duration
func SetTokenExpiration(duration time.Duration) {
	tokenExpiration = duration
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if a password matches a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a JWT token for a user
func GenerateToken(user models.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(tokenExpiration)

	claims := &Claims{
		Username: user.Username,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GetUserFromClaims gets a user from the database based on JWT claims
func GetUserFromClaims(claims *Claims) (*models.User, bool, error) {
	// Get the user from the database
	user, exists, err := storage.GetUser(claims.Username)
	if err != nil {
		return nil, false, err
	}

	if !exists {
		return nil, false, nil
	}

	// Convert to pointer for context
	userPtr := &user

	return userPtr, true, nil
}
