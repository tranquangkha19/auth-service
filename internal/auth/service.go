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
