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
	db         *database.Database
	jwtService *JWTService
}

func NewService(db *database.Database) (*Service, error) {
	jwtService, err := NewJWTService()
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT service: %w", err)
	}

	return &Service{
		db:         db,
		jwtService: jwtService,
	}, nil
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

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
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
func (s *Service) ValidateToken(token string) (uint, *database.User, error) {
	if token == "" {
		return 0, nil, fmt.Errorf("token is required")
	}

	// Validate JWT token
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user from database
	user, err := s.db.GetUserByID(claims.UserID)
	if err != nil {
		return 0, nil, fmt.Errorf("user not found: %w", err)
	}

	return user.ID, user, nil
}
