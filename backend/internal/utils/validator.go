package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// Validator provides validation utilities
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	email = strings.TrimSpace(strings.ToLower(email))
	
	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidateName validates participant name
func (v *Validator) ValidateName(name string) error {
	name = strings.TrimSpace(name)
	
	if name == "" {
		return errors.New("name is required")
	}

	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	if len(name) > 255 {
		return errors.New("name must not exceed 255 characters")
	}

	return nil
}

// ValidatePhone validates phone number
func (v *Validator) ValidatePhone(phone string) error {
	phone = strings.TrimSpace(phone)
	
	if phone == "" {
		return errors.New("phone number is required")
	}

	// Remove common separators for length check
	cleanPhone := strings.Map(func(r rune) rune {
		if r == ' ' || r == '-' || r == '(' || r == ')' || r == '+' {
			return -1
		}
		return r
	}, phone)

	if len(cleanPhone) < 10 {
		return errors.New("phone number must be at least 10 digits")
	}

	if len(phone) > 50 {
		return errors.New("phone number must not exceed 50 characters")
	}

	return nil
}

// ValidateInstagramHandle validates Instagram handle
func (v *Validator) ValidateInstagramHandle(handle *string) error {
	if handle == nil || *handle == "" {
		return nil // Optional field
	}

	cleaned := strings.TrimSpace(*handle)
	
	// Remove @ if present
	cleaned = strings.TrimPrefix(cleaned, "@")
	
	if len(cleaned) > 100 {
		return errors.New("Instagram handle must not exceed 100 characters")
	}

	// Check for valid Instagram handle format (alphanumeric, dots, underscores)
	for _, char := range cleaned {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '.' && char != '_' {
			return errors.New("Instagram handle can only contain letters, numbers, dots, and underscores")
		}
	}

	*handle = cleaned
	return nil
}

// ValidateAddress validates address
func (v *Validator) ValidateAddress(address string) error {
	address = strings.TrimSpace(address)
	
	if address == "" {
		return errors.New("address is required")
	}

	if len(address) < 10 {
		return errors.New("address must be at least 10 characters long")
	}

	if len(address) > 1000 {
		return errors.New("address must not exceed 1000 characters")
	}

	return nil
}

// SanitizeString removes potentially harmful characters
func (v *Validator) SanitizeString(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	return input
}

// ValidateRegistrationData validates all participant registration data
func (v *Validator) ValidateRegistrationData(name, email, phone string, instagramHandle *string, address string) []ValidationError {
	var errors []ValidationError

	// Validate name
	if err := v.ValidateName(name); err != nil {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: err.Error(),
		})
	}

	// Validate email
	if err := v.ValidateEmail(email); err != nil {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: err.Error(),
		})
	}

	// Validate phone
	if err := v.ValidatePhone(phone); err != nil {
		errors = append(errors, ValidationError{
			Field:   "phone",
			Message: err.Error(),
		})
	}

	// Validate Instagram handle (optional)
	if err := v.ValidateInstagramHandle(instagramHandle); err != nil {
		errors = append(errors, ValidationError{
			Field:   "instagram_handle",
			Message: err.Error(),
		})
	}

	// Validate address
	if err := v.ValidateAddress(address); err != nil {
		errors = append(errors, ValidationError{
			Field:   "address",
			Message: err.Error(),
		})
	}

	return errors
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
