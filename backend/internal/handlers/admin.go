package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tau-tau-run/backend/internal/middleware"
	"github.com/tau-tau-run/backend/internal/models"
	"github.com/tau-tau-run/backend/internal/services"
	"github.com/tau-tau-run/backend/internal/utils"
)

// AdminHandler handles admin-related requests
type AdminHandler struct {
	authService  *services.AuthService
	emailService *services.EmailService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(authService *services.AuthService, emailService *services.EmailService) *AdminHandler {
	return &AdminHandler{
		authService:  authService,
		emailService: emailService,
	}
}

// Login handles admin login
func (h *AdminHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request data", nil)
		return
	}

	// Normalize email
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Find admin by email
	admin, err := models.FindAdminByEmail(req.Email)
	if err != nil {
		utils.AuthLogger.Error("Failed to find admin: %v", err)
		middleware.RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", nil)
		return
	}

	// Check if admin exists
	if admin == nil {
		utils.AuthLogger.Warning("Login attempt with non-existent email: %s", req.Email)
		middleware.RespondWithError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
		return
	}

	// Verify password
	if err := h.authService.ComparePassword(admin.PasswordHash, req.Password); err != nil {
		utils.AuthLogger.Warning("Failed login attempt for admin: %s", req.Email)
		middleware.RespondWithError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
		return
	}

	// Generate JWT token
	token, expiresAt, err := h.authService.GenerateToken(admin.ID, admin.Email)
	if err != nil {
		utils.AuthLogger.Error("Failed to generate token for admin %s: %v", admin.Email, err)
		middleware.RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate authentication token", nil)
		return
	}

	// Log successful login
	utils.AuthLogger.Info("Admin logged in: %s", admin.Email)

	// Return success response
	middleware.RespondWithSuccess(c, http.StatusOK, "Login successful", models.LoginResponse{
		Token: token,
		Admin: models.AdminInfo{
			ID:    admin.ID,
			Email: admin.Email,
		},
		ExpiresAt: expiresAt,
	})
}

// GetParticipants returns all participants (protected route)
func (h *AdminHandler) GetParticipants(c *gin.Context) {
	// Get admin info from context (set by auth middleware)
	adminEmail := middleware.GetAdminEmail(c)
	utils.AuthLogger.Info("Admin %s requested participant list", adminEmail)

	// Get all participants
	participants, err := models.GetAllParticipants()
	if err != nil {
		utils.DBLogger.Error("Failed to get participants: %v", err)
		middleware.RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve participants", nil)
		return
	}

	// Return success response
	middleware.RespondWithSuccess(c, http.StatusOK, "", gin.H{
		"participants": participants,
		"total":        len(participants),
		"page":         1,
		"limit":        len(participants),
	})
}

// UpdatePaymentStatus updates participant payment status (protected route)
func (h *AdminHandler) UpdatePaymentStatus(c *gin.Context) {
	participantID := c.Param("id")
	adminEmail := middleware.GetAdminEmail(c)

	var req struct {
		PaymentStatus string `json:"payment_status" binding:"required"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request data", nil)
		return
	}

	// Validate payment status
	req.PaymentStatus = strings.ToUpper(req.PaymentStatus)
	if req.PaymentStatus != "PAID" && req.PaymentStatus != "UNPAID" {
		middleware.RespondWithError(c, http.StatusBadRequest, "INVALID_STATUS", "Payment status must be either PAID or UNPAID", nil)
		return
	}

	// Find participant
	participant, err := models.FindParticipantByID(participantID)
	if err != nil {
		utils.DBLogger.Error("Failed to find participant: %v", err)
		middleware.RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", nil)
		return
	}

	if participant == nil {
		middleware.RespondWithError(c, http.StatusNotFound, "PARTICIPANT_NOT_FOUND", "Participant with the specified ID does not exist", gin.H{
			"id": participantID,
		})
		return
	}

	// Store old status for email trigger logic
	oldStatus := participant.PaymentStatus

	// Update payment status
	if err := participant.UpdatePaymentStatus(req.PaymentStatus); err != nil {
		utils.DBLogger.Error("Failed to update payment status: %v", err)
		middleware.RespondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update payment status", nil)
		return
	}

	// Log the update
	utils.AuthLogger.Info("Admin %s updated participant %s payment status: %s → %s",
		adminEmail, participant.Email, oldStatus, req.PaymentStatus)

	// Trigger email if UNPAID → PAID (idempotency check)
	emailSent := false
	if oldStatus == "UNPAID" && req.PaymentStatus == "PAID" {
		utils.EmailLogger.Info("Payment status changed to PAID for %s - triggering confirmation email", participant.Email)
		
		// Send email asynchronously (non-blocking)
		h.emailService.SendConfirmationEmailAsync(participant)
		emailSent = true
	} else if oldStatus == "PAID" && req.PaymentStatus == "PAID" {
		utils.EmailLogger.Info("Payment status already PAID for %s - skipping duplicate email", participant.Email)
	}

	// Return success response
	middleware.RespondWithSuccess(c, http.StatusOK, "Payment status updated successfully", gin.H{
		"id":             participant.ID,
		"payment_status": participant.PaymentStatus,
		"updated_at":     participant.UpdatedAt,
		"email_sent":     emailSent,
	})
}
