package auth

import (
	"fmt"
	"time"

	"github.com/tranquangkha19/auth-service/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Fullname    string     `json:"fullname"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	Email       string     `json:"email,omitempty"`
	Username    string     `json:"username,omitempty"`
	Password    string     `json:"password"`
	Birthday    *time.Time `json:"birthday,omitempty"`
}

type Service struct {
	db *database.Database
}

func NewService(db *database.Database) *Service {
	return &Service{db: db}
}

func (s *Service) Login(account, password string) (string, error) {
	user, err := s.db.GetUserByAccount(account)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Update latest login
	if err := s.db.UpdateLatestLogin(user.ID); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update latest login: %v\n", err)
	}

	// TODO: Generate JWT token here
	return "mock-jwt-token", nil
}

func (s *Service) Register(req RegisterRequest) error {
	// Validate required fields
	if req.Fullname == "" || req.Password == "" {
		return fmt.Errorf("fullname and password are required")
	}

	// Check if at least one contact method is provided
	if req.Email == "" && req.Username == "" && req.PhoneNumber == "" {
		return fmt.Errorf("at least one of email, username, or phone number is required")
	}

	// Check if user already exists
	if req.Email != "" && s.db.UserExists(req.Email) {
		return fmt.Errorf("email already exists")
	}
	if req.Username != "" && s.db.UserExists(req.Username) {
		return fmt.Errorf("username already exists")
	}
	if req.PhoneNumber != "" && s.db.UserExists(req.PhoneNumber) {
		return fmt.Errorf("phone number already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &database.User{
		Fullname:     req.Fullname,
		PasswordHash: string(hashedPassword),
		Birthday:     req.Birthday,
	}

	// Set optional fields only if provided
	if req.PhoneNumber != "" {
		user.PhoneNumber = &req.PhoneNumber
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Username != "" {
		user.Username = &req.Username
	}

	return s.db.CreateUser(user)
}

// ValidateToken validates a JWT token and returns user information
// For now, this is a simple mock implementation
func (s *Service) ValidateToken(token string) (uint, *database.User, error) {
	// TODO: Implement proper JWT token validation
	// For now, we'll use a simple mock token validation

	if token == "" {
		return 0, nil, fmt.Errorf("token is required")
	}

	// Mock token validation - in production, you'd decode and verify the JWT
	if token == "mock-jwt-token" {
		// Return a mock user for testing
		user, err := s.db.GetUserByID(1) // Assuming user ID 1 exists
		if err != nil {
			return 0, nil, fmt.Errorf("user not found")
		}
		return user.ID, user, nil
	}

	return 0, nil, fmt.Errorf("invalid token")
}
