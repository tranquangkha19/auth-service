package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/tranquangkha19/auth-service/internal/database"
)

type JWTService struct {
	secretKey       []byte
	expirationHours int
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTService() (*JWTService, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file not found, continue with system environment variables
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable is required")
	}

	// Get expiration time from environment
	expirationStr := os.Getenv("JWT_EXPIRATION_HOURS")
	expirationHours := 24 // Default 24 hours

	if expirationStr != "" {
		if hours, err := strconv.Atoi(expirationStr); err == nil {
			expirationHours = hours
		}
	}

	return &JWTService{
		secretKey:       []byte(secretKey),
		expirationHours: expirationHours,
	}, nil
}

// GenerateToken creates a new JWT token for a user
func (j *JWTService) GenerateToken(user *database.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: getUsername(user),
		Email:    getEmail(user),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expirationHours) * time.Hour)),
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Second)), //MOCK: 30 seconds for testing
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// RefreshToken creates a new token with extended expiration
func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Create new claims with extended expiration
	newClaims := &Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   fmt.Sprintf("%d", claims.UserID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString(j.secretKey)
}

// Helper functions to safely get user fields
func getUsername(user *database.User) string {
	if user.Username != nil {
		return *user.Username
	}
	return ""
}

func getEmail(user *database.User) string {
	if user.Email != nil {
		return *user.Email
	}
	return ""
}
