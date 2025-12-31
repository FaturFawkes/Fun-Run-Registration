# Data Model: Event Registration System

**Feature**: 001-event-registration-system  
**Date**: 2025-12-31  
**Status**: Draft

## Database Schema

### Database: `tau_tau_run`

**Character Set**: UTF8  
**Collation**: utf8_general_ci (or utf8mb4 for full Unicode support)  
**Engine**: PostgreSQL 15+

---

## Tables

### 1. `participants`

Stores all registered participants for the Tau-Tau Run Fun Run 5K event.

```sql
CREATE TABLE participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50) NOT NULL,
    instagram_handle VARCHAR(100),
    address TEXT NOT NULL,
    registration_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    payment_status VARCHAR(20) NOT NULL DEFAULT 'UNPAID',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_registration_status CHECK (registration_status IN ('PENDING', 'CONFIRMED')),
    CONSTRAINT check_payment_status CHECK (payment_status IN ('UNPAID', 'PAID'))
);

CREATE INDEX idx_participants_email ON participants(email);
CREATE INDEX idx_participants_payment_status ON participants(payment_status);
CREATE INDEX idx_participants_created_at ON participants(created_at DESC);
```

**Fields**:

- `id` (UUID): Unique identifier, auto-generated
- `name` (VARCHAR 255): Full name of participant (required)
- `email` (VARCHAR 255): Email address, must be unique (required)
- `phone` (VARCHAR 50): Phone number (required)
- `instagram_handle` (VARCHAR 100): Instagram username (optional)
- `address` (TEXT): Full address (required)
- `registration_status` (VARCHAR 20): Current registration state
  - Values: `PENDING`, `CONFIRMED`
  - Default: `PENDING`
  - Used for future workflow (confirmation emails, manual review)
- `payment_status` (VARCHAR 20): Payment state
  - Values: `UNPAID`, `PAID`
  - Default: `UNPAID`
  - **Email trigger**: Transition from `UNPAID` → `PAID` sends confirmation email
- `created_at` (TIMESTAMP): Registration timestamp
- `updated_at` (TIMESTAMP): Last modification timestamp

**Constraints**:
- `email` must be unique (prevents duplicate registrations)
- `registration_status` must be one of: PENDING, CONFIRMED
- `payment_status` must be one of: UNPAID, PAID

**Indexes**:
- `email`: Fast lookup for duplicate detection
- `payment_status`: Fast filtering for unpaid participants
- `created_at`: Chronological sorting in admin dashboard

---

### 2. `admins`

Stores authenticated administrators who can access the admin dashboard.

```sql
CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$')
);

CREATE INDEX idx_admins_email ON admins(email);
```

**Fields**:

- `id` (UUID): Unique identifier
- `email` (VARCHAR 255): Admin email address (required, unique)
- `password_hash` (VARCHAR 255): Bcrypt hashed password (required)
  - **Security**: Never store plaintext passwords
  - **Hashing**: Use bcrypt with cost factor 12+
- `created_at` (TIMESTAMP): Admin account creation timestamp

**Constraints**:
- `email` must be unique
- `email` must match email format (basic regex validation)

**Security Notes**:
- Password must be hashed with bcrypt before storage
- Cost factor: minimum 12 rounds (recommended: 12-14)
- Never log or expose password_hash in API responses

---

### 3. `email_logs` (Optional for MVP, Recommended)

Logs all email sending attempts for debugging and audit purposes.

```sql
CREATE TABLE email_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    email_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_participant FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,
    CONSTRAINT check_status CHECK (status IN ('SUCCESS', 'FAILED'))
);

CREATE INDEX idx_email_logs_participant_id ON email_logs(participant_id);
CREATE INDEX idx_email_logs_status ON email_logs(status);
CREATE INDEX idx_email_logs_sent_at ON email_logs(sent_at DESC);
```

**Fields**:

- `id` (UUID): Unique identifier
- `participant_id` (UUID): Reference to participant who received email
- `recipient_email` (VARCHAR 255): Email address where email was sent
- `email_type` (VARCHAR 50): Type of email sent (e.g., "PAYMENT_CONFIRMATION")
- `status` (VARCHAR 20): Email delivery status
  - Values: `SUCCESS`, `FAILED`
- `error_message` (TEXT): SMTP error message if status is FAILED
- `sent_at` (TIMESTAMP): When email was sent/attempted

**Purpose**:
- Debug email delivery issues
- Audit trail for compliance
- Identify participants who didn't receive confirmation
- Support manual resend in future iterations

---

## Entity Relationships

```
admins (1) ----< (no direct relation) >---- participants (many)
   |                                              |
   |                                              |
   +--- Authenticate to update payment           |
                                                  |
participants (1) ----< (many) email_logs
   |
   +--- Each payment status change may generate email log entry
```

**Relationships**:
- Admins have no direct foreign key to participants (many-to-many implicit through actions)
- Participants have one-to-many relationship with email_logs
- Email logs are cascade deleted if participant is deleted (data cleanup)

---

## State Machine: Payment Status

```
┌─────────┐
│ UNPAID  │ (Initial state)
└────┬────┘
     │
     │ Admin updates payment_status to PAID
     │
     ▼
┌─────────┐
│  PAID   │ (Final state)
└─────────┘
     │
     └──> Trigger: Send confirmation email via SMTP
```

**State Transitions**:

1. **UNPAID → PAID** (only valid transition in MVP)
   - Triggered by: Admin updating payment status in dashboard
   - Side effect: Automatic email sent to participant
   - Idempotent: Updating PAID→PAID does nothing (no duplicate emails)

2. **PAID → UNPAID** (not implemented in MVP)
   - Future consideration: Refund scenario
   - Would require additional business logic

**Email Trigger Logic**:

```go
// Pseudocode
func UpdatePaymentStatus(participantID, newStatus) {
    oldStatus := GetCurrentStatus(participantID)
    
    // Update database first
    UpdateDatabase(participantID, newStatus)
    
    // Trigger email only on UNPAID → PAID transition
    if oldStatus == "UNPAID" && newStatus == "PAID" {
        SendConfirmationEmail(participantID) // Async, non-blocking
    }
}
```

---

## Data Validation Rules

### Participants Table

| Field | Validation | Error Message |
|-------|------------|---------------|
| `name` | Required, min 2 chars, max 255 chars | "Name is required and must be 2-255 characters" |
| `email` | Required, valid email format, unique | "Valid email is required" / "Email already registered" |
| `phone` | Required, min 10 chars, max 50 chars | "Phone number is required (10-50 characters)" |
| `instagram_handle` | Optional, max 100 chars, alphanumeric + underscore | "Instagram handle must be alphanumeric" |
| `address` | Required, min 10 chars | "Address is required (minimum 10 characters)" |
| `registration_status` | Must be PENDING or CONFIRMED | "Invalid registration status" |
| `payment_status` | Must be UNPAID or PAID | "Invalid payment status" |

### Admins Table

| Field | Validation | Error Message |
|-------|------------|---------------|
| `email` | Required, valid email format, unique | "Valid email is required" / "Admin email already exists" |
| `password` | Min 8 chars, must contain letter + number (pre-hash) | "Password must be 8+ characters with letters and numbers" |

---

## Migration Script

**File**: `database/migrations/001_init.sql`

```sql
-- Migration: 001_init
-- Description: Initial schema for Tau-Tau Run registration system
-- Date: 2025-12-31

BEGIN;

-- Create participants table
CREATE TABLE participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50) NOT NULL,
    instagram_handle VARCHAR(100),
    address TEXT NOT NULL,
    registration_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    payment_status VARCHAR(20) NOT NULL DEFAULT 'UNPAID',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_registration_status CHECK (registration_status IN ('PENDING', 'CONFIRMED')),
    CONSTRAINT check_payment_status CHECK (payment_status IN ('UNPAID', 'PAID'))
);

CREATE INDEX idx_participants_email ON participants(email);
CREATE INDEX idx_participants_payment_status ON participants(payment_status);
CREATE INDEX idx_participants_created_at ON participants(created_at DESC);

-- Create admins table
CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$')
);

CREATE INDEX idx_admins_email ON admins(email);

-- Create email_logs table (optional for MVP)
CREATE TABLE email_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    email_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_participant FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,
    CONSTRAINT check_status CHECK (status IN ('SUCCESS', 'FAILED'))
);

CREATE INDEX idx_email_logs_participant_id ON email_logs(participant_id);
CREATE INDEX idx_email_logs_status ON email_logs(status);
CREATE INDEX idx_email_logs_sent_at ON email_logs(sent_at DESC);

-- Create trigger to auto-update updated_at on participants
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_participants_updated_at
    BEFORE UPDATE ON participants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMIT;
```

---

## Seed Data (Development Only)

**File**: `database/seeds/001_admin_seed.sql`

```sql
-- Seed admin user for development
-- Email: admin@tautaurun.com
-- Password: Admin123! (bcrypt hash below)

INSERT INTO admins (email, password_hash, created_at)
VALUES (
    'admin@tautaurun.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYlK4Qr1WZK', -- bcrypt hash of "Admin123!"
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;
```

**Note**: Replace with real password hash generated via bcrypt in production. Never commit production passwords.

---

## Performance Considerations

1. **Indexes**: Added on frequently queried columns (email, payment_status, created_at)
2. **UUID vs. Serial**: UUID chosen for security (non-sequential IDs) and potential future distributed systems
3. **Connection Pooling**: Backend must implement connection pooling (recommended: max 10-20 connections for MVP)
4. **Cascading Deletes**: Email logs cascade delete with participants (cleanup)

---

## Security Considerations

1. **Password Storage**: Bcrypt with cost factor 12+ (enforced in backend, not database)
2. **Email Uniqueness**: Prevents duplicate registrations, enforced by unique constraint
3. **SQL Injection Prevention**: Backend uses parameterized queries (never string concatenation)
4. **Data Encryption**: HTTPS enforced for all API communication (not database-level encryption for MVP)

---

## Future Enhancements (Post-MVP)

1. Add `registration_status` workflow (PENDING → CONFIRMED flow)
2. Add `email_logs` for better debugging and resend capabilities
3. Add `events` table for multi-event support
4. Add `payment_method` and `payment_reference` fields to participants
5. Add soft delete (`deleted_at`) instead of hard delete
6. Add audit log table for all admin actions

---

**Document Status**: Ready for review  
**Review Required**: Yes (validate field types, constraints, indexes)  
**Next Step**: Review and approve, then create migration file in repository
