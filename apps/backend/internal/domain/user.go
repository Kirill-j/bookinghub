package domain

import "time"

type UserRole string

const (
	RoleUser    UserRole = "USER"
	RoleManager UserRole = "MANAGER"
	RoleAdmin   UserRole = "ADMIN"
)

type User struct {
	ID           uint64    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	Role         UserRole  `json:"role" db:"role"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
