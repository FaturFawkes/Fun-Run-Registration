package services

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/tau-tau-run/backend/config"
	"github.com/tau-tau-run/backend/internal/database"
	"github.com/tau-tau-run/backend/internal/models"
	"github.com/tau-tau-run/backend/internal/utils"
)

// EmailService handles email operations
type EmailService struct {
	config *config.Config
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendConfirmationEmail sends a payment confirmation email to the participant
func (s *EmailService) SendConfirmationEmail(participant *models.Participant) error {
	// Build email content
	subject := fmt.Sprintf("Payment Confirmed - %s", s.config.Event.Name)
	htmlBody, err := s.buildConfirmationEmailHTML(participant)
	if err != nil {
		return fmt.Errorf("failed to build email template: %w", err)
	}
	
	plainBody := s.buildConfirmationEmailPlain(participant)

	// Send email
	return s.sendEmail(participant.Email, subject, htmlBody, plainBody)
}

// sendEmail sends an email via SMTP
func (s *EmailService) sendEmail(to, subject, htmlBody, plainBody string) error {
	// Check if SMTP is configured
	if s.config.SMTP.Host == "" {
		return fmt.Errorf("SMTP not configured")
	}

	// Build email message
	from := s.config.SMTP.FromEmail
	fromName := s.config.SMTP.FromName

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, from)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["Date"] = time.Now().Format(time.RFC1123Z)

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// SMTP authentication
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.config.SMTP.Host, s.config.SMTP.Port)
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// buildConfirmationEmailHTML creates HTML email template
func (s *EmailService) buildConfirmationEmailHTML(participant *models.Participant) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #FF6B35; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9f9f9; padding: 30px; border: 1px solid #ddd; border-radius: 0 0 5px 5px; }
        .info-box { background-color: white; padding: 15px; margin: 20px 0; border-left: 4px solid #FF6B35; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
        .highlight { color: #FF6B35; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Payment Confirmed!</h1>
        </div>
        <div class="content">
            <p>Dear <strong>{{.Name}}</strong>,</p>
            
            <p>Great news! We have successfully received your payment for <span class="highlight">{{.EventName}}</span>.</p>
            
            <div class="info-box">
                <h3>Event Details:</h3>
                <p><strong>Event:</strong> {{.EventName}}</p>
                <p><strong>Date:</strong> {{.EventDate}}</p>
                <p><strong>Location:</strong> {{.EventLocation}}</p>
            </div>
            
            <div class="info-box">
                <h3>Your Registration:</h3>
                <p><strong>Name:</strong> {{.Name}}</p>
                <p><strong>Email:</strong> {{.Email}}</p>
                <p><strong>Phone:</strong> {{.Phone}}</p>
                {{if .InstagramHandle}}
                <p><strong>Instagram:</strong> {{.InstagramHandle}}</p>
                {{end}}
                <p><strong>Registration Status:</strong> <span class="highlight">CONFIRMED</span></p>
                <p><strong>Payment Status:</strong> <span class="highlight">PAID</span></p>
            </div>
            
            <p><strong>What's Next?</strong></p>
            <ul>
                <li>We'll send you more details about the event as we get closer to the date</li>
                <li>Please arrive at least 30 minutes before the event starts</li>
                <li>Bring your ID for registration check-in</li>
                <li>Get ready to have fun! 🏃‍♂️</li>
            </ul>
            
            <p>If you have any questions, please don't hesitate to contact us.</p>
            
            <p>See you at the event!</p>
            
            <p><strong>{{.EventTeam}}</strong></p>
        </div>
        <div class="footer">
            <p>This is an automated confirmation email. Please do not reply to this message.</p>
            <p>&copy; {{.Year}} {{.EventName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("confirmation").Parse(tmpl)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"Name":             participant.Name,
		"Email":            participant.Email,
		"Phone":            participant.Phone,
		"InstagramHandle":  participant.InstagramHandle,
		"EventName":        s.config.Event.Name,
		"EventDate":        s.config.Event.Date,
		"EventLocation":    s.config.Event.Location,
		"EventTeam":        s.config.SMTP.FromName,
		"Year":             time.Now().Year(),
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// buildConfirmationEmailPlain creates plain text email template
func (s *EmailService) buildConfirmationEmailPlain(participant *models.Participant) string {
	instagram := ""
	if participant.InstagramHandle != nil {
		instagram = fmt.Sprintf("\nInstagram: %s", *participant.InstagramHandle)
	}

	return fmt.Sprintf(`
🎉 Payment Confirmed!

Dear %s,

Great news! We have successfully received your payment for %s.

EVENT DETAILS:
- Event: %s
- Date: %s
- Location: %s

YOUR REGISTRATION:
- Name: %s
- Email: %s
- Phone: %s%s
- Registration Status: CONFIRMED
- Payment Status: PAID

WHAT'S NEXT?
- We'll send you more details about the event as we get closer to the date
- Please arrive at least 30 minutes before the event starts
- Bring your ID for registration check-in
- Get ready to have fun! 🏃‍♂️

If you have any questions, please don't hesitate to contact us.

See you at the event!

%s

---
This is an automated confirmation email. Please do not reply to this message.
© %d %s. All rights reserved.
`,
		participant.Name,
		s.config.Event.Name,
		s.config.Event.Name,
		s.config.Event.Date,
		s.config.Event.Location,
		participant.Name,
		participant.Email,
		participant.Phone,
		instagram,
		s.config.SMTP.FromName,
		time.Now().Year(),
		s.config.Event.Name,
	)
}

// LogEmail logs email sending attempts to the database
func (s *EmailService) LogEmail(participantID, recipientEmail, emailType, status, errorMessage string) error {
	query := `
		INSERT INTO email_logs (participant_id, recipient_email, email_type, status, error_message, sent_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
	`

	var errMsg sql.NullString
	if errorMessage != "" {
		errMsg = sql.NullString{String: errorMessage, Valid: true}
	}

	_, err := database.DB.Exec(query, participantID, recipientEmail, emailType, status, errMsg)
	if err != nil {
		return fmt.Errorf("failed to log email: %w", err)
	}

	return nil
}

// SendConfirmationEmailAsync sends confirmation email asynchronously
func (s *EmailService) SendConfirmationEmailAsync(participant *models.Participant) {
	go func() {
		utils.EmailLogger.Info("Sending confirmation email to %s (ID: %s)", participant.Email, participant.ID)
		
		err := s.SendConfirmationEmail(participant)
		
		if err != nil {
			utils.EmailLogger.Error("Failed to send email to %s: %v", participant.Email, err)
			// Log failure to database
			logErr := s.LogEmail(participant.ID, participant.Email, "PAYMENT_CONFIRMATION", "FAILED", err.Error())
			if logErr != nil {
				utils.EmailLogger.Error("Failed to log email failure: %v", logErr)
			}
		} else {
			utils.EmailLogger.Info("Successfully sent confirmation email to %s", participant.Email)
			// Log success to database
			logErr := s.LogEmail(participant.ID, participant.Email, "PAYMENT_CONFIRMATION", "SUCCESS", "")
			if logErr != nil {
				utils.EmailLogger.Error("Failed to log email success: %v", logErr)
			}
		}
	}()
}

// ValidateSMTPConfig checks if SMTP configuration is valid
func ValidateSMTPConfig(cfg *config.Config) error {
	if cfg.SMTP.Host == "" {
		return fmt.Errorf("SMTP_HOST is not configured")
	}
	if cfg.SMTP.Port == "" {
		return fmt.Errorf("SMTP_PORT is not configured")
	}
	if cfg.SMTP.FromEmail == "" {
		return fmt.Errorf("SMTP_FROM_EMAIL is not configured")
	}
	
	utils.EmailLogger.Info("SMTP configuration validated: %s:%s", cfg.SMTP.Host, cfg.SMTP.Port)
	return nil
}
