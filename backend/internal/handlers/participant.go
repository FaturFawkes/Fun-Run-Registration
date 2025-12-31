package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/tau-tau-run/backend/internal/middleware"
	"github.com/tau-tau-run/backend/internal/models"
	"github.com/tau-tau-run/backend/internal/utils"
)

// ParticipantHandler handles participant-related requests
type ParticipantHandler struct {
	validator *utils.Validator
}

// NewParticipantHandler creates a new participant handler
func NewParticipantHandler() *ParticipantHandler {
	return &ParticipantHandler{
		validator: utils.NewValidator(),
	}
}

// Register handles participant registration
func (h *ParticipantHandler) Register(c *gin.Context) {
	var req models.CreateParticipantRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request data", nil)
		return
	}

	// Sanitize inputs
	req.Name = h.validator.SanitizeString(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Phone = h.validator.SanitizeString(req.Phone)
	req.Address = h.validator.SanitizeString(req.Address)
	if req.InstagramHandle != nil {
		sanitized := h.validator.SanitizeString(*req.InstagramHandle)
		req.InstagramHandle = &sanitized
	}

	// Validate all fields
	validationErrors := h.validator.ValidateRegistrationData(
		req.Name,
		req.Email,
		req.Phone,
		req.InstagramHandle,
		req.Address,
	)

	if len(validationErrors) > 0 {
		middleware.RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input data", validationErrors)
		return
	}

	// Check for duplicate email
	existing, err := models.FindParticipantByEmail(req.Email)
	if err != nil {
		utils.DBLogger.Error("Failed to check duplicate email: %v", err)
		middleware.RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", nil)
		return
	}

	if existing != nil {
		middleware.RespondWithError(c, http.StatusConflict, "DUPLICATE_EMAIL", "Email address is already registered", gin.H{
			"email": req.Email,
		})
		return
	}

	// Create participant
	participant := &models.Participant{
		Name:            req.Name,
		Email:           req.Email,
		Phone:           req.Phone,
		InstagramHandle: req.InstagramHandle,
		Address:         req.Address,
	}

	if err := participant.Create(); err != nil {
		// Check for unique constraint violation (just in case of race condition)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			middleware.RespondWithError(c, http.StatusConflict, "DUPLICATE_EMAIL", "Email address is already registered", nil)
			return
		}

		utils.DBLogger.Error("Failed to create participant: %v", err)
		middleware.RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register participant", nil)
		return
	}

	// Log successful registration
	utils.ServerLogger.Info("New participant registered: %s (%s)", participant.Name, participant.Email)

	// Return success response
	middleware.RespondWithSuccess(c, http.StatusCreated, "Registration successful! Your payment status is pending.", gin.H{
		"id":                  participant.ID,
		"email":               participant.Email,
		"registration_status": participant.RegistrationStatus,
		"payment_status":      participant.PaymentStatus,
	})
}
