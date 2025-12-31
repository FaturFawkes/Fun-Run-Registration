# API Documentation - Tau-Tau Run Event Registration System

**Version:** 1.0.0  
**Base URL:** `http://localhost:8081/api/v1`  
**Production URL:** `https://api.tautaurun.com/api/v1`

---

## Table of Contents

- [Authentication](#authentication)
- [Response Format](#response-format)
- [Public Endpoints](#public-endpoints)
- [Admin Endpoints](#admin-endpoints)
- [Error Codes](#error-codes)
- [Rate Limiting](#rate-limiting)

---

## Authentication

### JWT Token Authentication

Admin endpoints require JWT authentication via the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

**Token Expiration:** 24 hours (configurable via `JWT_EXPIRATION_HOURS`)

---

## Response Format

All API responses follow a consistent JSON structure:

### Success Response
```json
{
  "success": true,
  "message": "Optional success message",
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      // Optional additional error details
    }
  }
}
```

---

## Public Endpoints

### Health Check

Check API server health status.

**Endpoint:** `GET /public/health`  
**Authentication:** None  

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "version": "1.0.0"
  }
}
```

---

### Register Participant

Register a new participant for the event.

**Endpoint:** `POST /public/register`  
**Authentication:** None  
**Content-Type:** `application/json`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "081234567890",
  "instagram_handle": "@johndoe",
  "address": "Jl. Sudirman No. 123, Jakarta, Indonesia"
}
```

**Field Validations:**
- `name` (required): 2-100 characters
- `email` (required): Valid email format
- `phone` (required): 10-15 digits
- `instagram_handle` (optional): Max 50 characters
- `address` (required): Min 10 characters

**Success Response (201):**
```json
{
  "success": true,
  "message": "Registration successful! Your payment status is pending.",
  "data": {
    "id": "uuid-here",
    "email": "john.doe@example.com",
    "registration_status": "PENDING",
    "payment_status": "UNPAID"
  }
}
```

**Error Response (409 - Duplicate Email):**
```json
{
  "success": false,
  "error": {
    "code": "DUPLICATE_EMAIL",
    "message": "Email address is already registered",
    "details": {
      "email": "john.doe@example.com"
    }
  }
}
```

**Error Response (400 - Validation Error):**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "address",
        "message": "address must be at least 10 characters long"
      }
    ]
  }
}
```

---

## Admin Endpoints

### Admin Login

Authenticate an admin user and receive JWT token.

**Endpoint:** `POST /admin/login`  
**Authentication:** None  
**Content-Type:** `application/json`

**Request Body:**
```json
{
  "email": "admin@tautaurun.com",
  "password": "Admin123!"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin": {
      "id": "uuid-here",
      "email": "admin@tautaurun.com"
    },
    "expires_at": "2026-01-02T10:30:00Z"
  }
}
```

**Error Response (401 - Invalid Credentials):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid email or password"
  }
}
```

---

### Get All Participants

Retrieve list of all registered participants.

**Endpoint:** `GET /admin/participants`  
**Authentication:** Required (JWT)  
**Query Parameters:** None (pagination to be added in future)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "participants": [
      {
        "id": "uuid-here",
        "name": "John Doe",
        "email": "john.doe@example.com",
        "phone": "081234567890",
        "instagram_handle": "johndoe",
        "address": "Jl. Sudirman No. 123, Jakarta, Indonesia",
        "registration_status": "PENDING",
        "payment_status": "PAID",
        "created_at": "2026-01-01T10:00:00Z",
        "updated_at": "2026-01-01T11:30:00Z"
      }
    ],
    "total": 42,
    "page": 1,
    "limit": 42
  }
}
```

**Error Response (401 - Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required"
  }
}
```

---

### Update Payment Status

Update participant's payment status.

**Endpoint:** `PATCH /admin/participants/:id/payment`  
**Authentication:** Required (JWT)  
**Content-Type:** `application/json`

**URL Parameters:**
- `id` (required): Participant UUID

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "payment_status": "PAID"
}
```

**Valid Values:** `PAID`, `UNPAID`

**Success Response (200):**
```json
{
  "success": true,
  "message": "Payment status updated successfully",
  "data": {
    "id": "uuid-here",
    "payment_status": "PAID",
    "updated_at": "2026-01-01T12:00:00Z",
    "email_sent": true
  }
}
```

**Behavior:**
- When status changes from `UNPAID` → `PAID`: Triggers confirmation email
- When status is already `PAID` → `PAID`: No email sent (idempotency)
- Email sending is asynchronous and non-blocking
- Payment update succeeds even if email fails

**Error Response (404 - Not Found):**
```json
{
  "success": false,
  "error": {
    "code": "PARTICIPANT_NOT_FOUND",
    "message": "Participant with the specified ID does not exist",
    "details": {
      "id": "invalid-uuid"
    }
  }
}
```

**Error Response (400 - Invalid Status):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_STATUS",
    "message": "Payment status must be either PAID or UNPAID"
  }
}
```

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Request data failed validation |
| `INVALID_STATUS` | 400 | Invalid payment status value |
| `INVALID_CREDENTIALS` | 401 | Wrong email or password |
| `UNAUTHORIZED` | 401 | Missing or invalid JWT token |
| `PARTICIPANT_NOT_FOUND` | 404 | Participant ID doesn't exist |
| `DUPLICATE_EMAIL` | 409 | Email already registered |
| `INTERNAL_ERROR` | 500 | Server error (check logs) |

---

## Rate Limiting

**Current Status:** Not implemented (planned for production)

**Recommended Limits:**
- Public endpoints: 100 requests/minute per IP
- Admin endpoints: 300 requests/minute per token
- Registration endpoint: 10 requests/minute per IP

---

## Email Automation

When a participant's payment status is updated to `PAID`, the system automatically:

1. Updates the database record
2. Triggers an async email send
3. Sends HTML confirmation email to participant
4. Logs email attempt to `email_logs` table (SUCCESS/FAILED)
5. Returns response immediately (non-blocking)

**Email will NOT be sent if:**
- Payment status is already `PAID` (idempotency)
- Status changes from `PAID` to `UNPAID`

---

## Curl Examples

### Register a Participant
```bash
curl -X POST http://localhost:8081/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane@example.com",
    "phone": "081234567890",
    "address": "Jakarta, Indonesia - Test Street 123"
  }'
```

### Admin Login
```bash
curl -X POST http://localhost:8081/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@tautaurun.com",
    "password": "Admin123!"
  }'
```

### Get Participants (with token)
```bash
TOKEN="your-jwt-token-here"

curl -X GET http://localhost:8081/api/v1/admin/participants \
  -H "Authorization: Bearer $TOKEN"
```

### Update Payment Status
```bash
TOKEN="your-jwt-token-here"
PARTICIPANT_ID="participant-uuid-here"

curl -X PATCH http://localhost:8081/api/v1/admin/participants/$PARTICIPANT_ID/payment \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"payment_status": "PAID"}'
```

---

## Database Schema Reference

### Participants Table
- `id` (UUID, PK)
- `name` (VARCHAR)
- `email` (VARCHAR, UNIQUE)
- `phone` (VARCHAR)
- `instagram_handle` (VARCHAR, nullable)
- `address` (TEXT)
- `registration_status` (VARCHAR) - PENDING, CONFIRMED
- `payment_status` (VARCHAR) - UNPAID, PAID
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Email Logs Table
- `id` (SERIAL, PK)
- `participant_id` (UUID, FK)
- `recipient_email` (VARCHAR)
- `email_type` (VARCHAR) - PAYMENT_CONFIRMATION
- `status` (VARCHAR) - SUCCESS, FAILED
- `error_message` (TEXT, nullable)
- `sent_at` (TIMESTAMP)

---

## Change Log

### Version 1.0.0 (2026-01-01)
- Initial API release
- Public registration endpoint
- Admin authentication
- Participant management
- Email automation

---

**Support:** For issues or questions, check the main README.md or contact admin@tautaurun.com
