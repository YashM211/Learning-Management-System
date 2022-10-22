package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// User represents a User schema
type User struct {
	gorm.Model
	Username string    `json:"username" gorm:"unique"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	Courses  []Courses `json:"courses"`
}

// UserErrors
type UserErrors struct {
	Err      bool   `json:"error"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims represent the structure of the JWT token
type Claims struct {
	jwt.StandardClaims
	ID uint `gorm:"primaryKey"`
}

var VerifiedUser string

type Courses struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `json:"index"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Enrollmentkey string         `json:"enrollmentkey"`
	Teacher       string         `json:"teacher"`
	UserID        uint           `json:"user_id"`
}
