package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tau-tau-run/backend/internal/database"
)

// Admin represents an authenticated administrator
type Admin struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose in JSON
	CreatedAt    time.Time `json:"created_at"`
}

// LoginRequest represents admin login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents admin login response
type LoginResponse struct {
	Token     string    `json:"token"`
	Admin     AdminInfo `json:"admin"`
	ExpiresAt time.Time `json:"expires_at"`
}

// AdminInfo represents public admin information
type AdminInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// FindAdminByEmail finds an admin by email
func FindAdminByEmail(email string) (*Admin, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM admins
		WHERE email = $1
	`

	admin := &Admin{}

	// Use database.DB instead of sql.DB
	err := database.DB.QueryRow(query, email).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}

	return admin, nil
}

// FindAdminByID finds an admin by ID
func FindAdminByID(id string) (*Admin, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM admins
		WHERE id = $1
	`

	admin := &Admin{}
	err := database.DB.QueryRow(query, id).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}

	return admin, nil
}
