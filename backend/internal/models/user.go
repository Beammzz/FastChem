package models

import "time"

// User represents a registered user.
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	TotalPoints  int       `json:"totalPoints"`
	CreatedAt    time.Time `json:"createdAt"`
}

// RegisterRequest is the payload for user registration.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// LoginRequest is the payload for user login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is returned after successful login/register.
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UserPublic is a public-facing user representation (no sensitive fields).
//
// IsAdmin is set only on /api/auth/me, where the caller is the account itself.
// It tells the client whether to offer the admin console; it grants nothing —
// every /api/admin route re-reads the flag from the database.
type UserPublic struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	TotalPoints int    `json:"totalPoints"`
	IsAdmin     bool   `json:"isAdmin,omitempty"`
}
