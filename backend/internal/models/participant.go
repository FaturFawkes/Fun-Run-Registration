package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tau-tau-run/backend/internal/database"
)

// Participant represents a registered participant
type Participant struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	InstagramHandle    *string   `json:"instagram_handle"`
	Address            string    `json:"address"`
	RegistrationStatus string    `json:"registration_status"`
	PaymentStatus      string    `json:"payment_status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// CreateParticipantRequest represents registration request data
type CreateParticipantRequest struct {
	Name            string  `json:"name" binding:"required"`
	Email           string  `json:"email" binding:"required,email"`
	Phone           string  `json:"phone" binding:"required"`
	InstagramHandle *string `json:"instagram_handle"`
	Address         string  `json:"address" binding:"required"`
}

// Create creates a new participant in the database
func (p *Participant) Create() error {
	query := `
		INSERT INTO participants (name, email, phone, instagram_handle, address, registration_status, payment_status)
		VALUES ($1, $2, $3, $4, $5, 'PENDING', 'UNPAID')
		RETURNING id, created_at, updated_at
	`

	err := database.DB.QueryRow(
		query,
		p.Name,
		p.Email,
		p.Phone,
		p.InstagramHandle,
		p.Address,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create participant: %w", err)
	}

	p.RegistrationStatus = "PENDING"
	p.PaymentStatus = "UNPAID"

	return nil
}

// FindByEmail finds a participant by email
func FindParticipantByEmail(email string) (*Participant, error) {
	query := `
		SELECT id, name, email, phone, instagram_handle, address, 
		       registration_status, payment_status, created_at, updated_at
		FROM participants
		WHERE email = $1
	`

	participant := &Participant{}
	err := database.DB.QueryRow(query, email).Scan(
		&participant.ID,
		&participant.Name,
		&participant.Email,
		&participant.Phone,
		&participant.InstagramHandle,
		&participant.Address,
		&participant.RegistrationStatus,
		&participant.PaymentStatus,
		&participant.CreatedAt,
		&participant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find participant: %w", err)
	}

	return participant, nil
}

// FindByID finds a participant by ID
func FindParticipantByID(id string) (*Participant, error) {
	query := `
		SELECT id, name, email, phone, instagram_handle, address,
		       registration_status, payment_status, created_at, updated_at
		FROM participants
		WHERE id = $1
	`

	participant := &Participant{}
	err := database.DB.QueryRow(query, id).Scan(
		&participant.ID,
		&participant.Name,
		&participant.Email,
		&participant.Phone,
		&participant.InstagramHandle,
		&participant.Address,
		&participant.RegistrationStatus,
		&participant.PaymentStatus,
		&participant.CreatedAt,
		&participant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find participant: %w", err)
	}

	return participant, nil
}

// GetAll retrieves all participants
func GetAllParticipants() ([]Participant, error) {
	query := `
		SELECT id, name, email, phone, instagram_handle, address,
		       registration_status, payment_status, created_at, updated_at
		FROM participants
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	defer rows.Close()

	var participants []Participant
	for rows.Next() {
		var p Participant
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Email,
			&p.Phone,
			&p.InstagramHandle,
			&p.Address,
			&p.RegistrationStatus,
			&p.PaymentStatus,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		participants = append(participants, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating participants: %w", err)
	}

	return participants, nil
}

// UpdatePaymentStatus updates the payment status of a participant
func (p *Participant) UpdatePaymentStatus(status string) error {
	query := `
		UPDATE participants
		SET payment_status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at
	`

	err := database.DB.QueryRow(query, status, p.ID).Scan(&p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	p.PaymentStatus = status
	return nil
}
