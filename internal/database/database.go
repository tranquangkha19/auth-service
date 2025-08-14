package database

import (
	"fmt"
	"log"
	"time"

	"github.com/tranquangkha19/auth-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Fullname     string     `gorm:"not null" json:"fullname"`
	PhoneNumber  *string    `gorm:"uniqueIndex;size:20" json:"phone_number,omitempty"`
	Email        *string    `gorm:"uniqueIndex;size:255" json:"email,omitempty"`
	Username     *string    `gorm:"uniqueIndex;size:100" json:"username,omitempty"`
	PasswordHash string     `gorm:"not null;column:password_hash" json:"-"`
	Birthday     *time.Time `gorm:"type:date" json:"birthday,omitempty"`
	LatestLogin  *time.Time `gorm:"column:latest_login" json:"latest_login,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	// Use URL format for more reliable connection
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	log.Printf("Connecting to database with DSN: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Skip auto-migration since we're using manual migrations
	// The table structure is already created via SQL migration

	log.Println("Database connected and migrated successfully")
	return &Database{DB: db}, nil
}

func (d *Database) CreateUser(user *User) error {
	return d.DB.Create(user).Error
}

func (d *Database) GetUserByAccount(account string) (*User, error) {
	var user User
	err := d.DB.Where("email = ? OR username = ? OR phone_number = ?", account, account, account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Database) UserExists(account string) bool {
	var count int64
	d.DB.Model(&User{}).Where("email = ? OR username = ? OR phone_number = ?", account, account, account).Count(&count)
	return count > 0
}

func (d *Database) UpdateLatestLogin(userID uint) error {
	now := time.Now()
	return d.DB.Model(&User{}).Where("id = ?", userID).Update("latest_login", now).Error
}

func (d *Database) GetUserByID(userID uint) (*User, error) {
	var user User
	err := d.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
