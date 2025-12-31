# API Contracts: Admin API

**Feature**: 001-event-registration-system  
**Base URL**: `/api/v1/admin`  
**Authentication**: JWT Bearer Token (required for all endpoints except login)  
**Date**: 2025-12-31

---

## Authentication

All admin endpoints (except `/login`) require a valid JWT token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

**Token Expiration**: 24 hours  
**Token Refresh**: Not implemented in MVP (user must re-login after expiration)

---

## Endpoints

### 1. Admin Login

**Purpose**: Authenticate admin user and receive JWT token.

**Endpoint**: `POST /api/v1/admin/login`

**Authentication**: None (public endpoint)

**Request Headers**:
```
Content-Type: application/json
```

**Request Body**:
```json
{
  "email": "admin@tautaurun.com",
  "password": "Admin123!"
}
```

**Request Body Schema**:

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `email` | string | Yes | Valid email format | Admin email address |
| `password` | string | Yes | Min 8 characters | Admin password (plaintext, will be hashed for comparison) |

**Success Response** (HTTP 200 OK):

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "email": "admin@tautaurun.com"
    },
    "expires_at": "2026-01-01T17:46:37Z"
  }
}
```

**Error Responses**:

**400 Bad Request** - Validation error:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Email is required"
      }
    ]
  }
}
```

**401 Unauthorized** - Invalid credentials:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid email or password"
  }
}
```

**500 Internal Server Error**:
```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An unexpected error occurred. Please try again later."
  }
}
```

---

### 2. List All Participants

**Purpose**: Retrieve list of all registered participants for admin dashboard.

**Endpoint**: `GET /api/v1/admin/participants`

**Authentication**: Required (JWT token)

**Request Headers**:
```
Authorization: Bearer <jwt_token>
```

**Query Parameters** (Optional - for future pagination):

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number (future) |
| `limit` | integer | 100 | Items per page (future) |
| `payment_status` | string | - | Filter by UNPAID or PAID (future) |

**Success Response** (HTTP 200 OK):

```json
{
  "success": true,
  "data": {
    "participants": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "John Doe",
        "email": "john.doe@example.com",
        "phone": "+6281234567890",
        "instagram_handle": "johndoe",
        "address": "Jl. Example Street No. 123, Jakarta, Indonesia 12345",
        "registration_status": "PENDING",
        "payment_status": "UNPAID",
        "created_at": "2025-12-31T10:30:00Z",
        "updated_at": "2025-12-31T10:30:00Z"
      },
      {
        "id": "550e8400-e29b-41d4-a716-446655440002",
        "name": "Jane Smith",
        "email": "jane.smith@example.com",
        "phone": "+6281234567891",
        "instagram_handle": "janesmith",
        "address": "Jl. Merdeka No. 45, Bandung, Indonesia 40111",
        "registration_status": "PENDING",
        "payment_status": "PAID",
        "created_at": "2025-12-31T11:00:00Z",
        "updated_at": "2025-12-31T12:00:00Z"
      }
    ],
    "total": 2,
    "page": 1,
    "limit": 100
  }
}
```

**Error Responses**:

**401 Unauthorized** - Missing or invalid token:
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required. Please provide a valid token."
  }
}
```

**403 Forbidden** - Expired token:
```json
{
  "success": false,
  "error": {
    "code": "TOKEN_EXPIRED",
    "message": "Your session has expired. Please log in again."
  }
}
```

---

### 3. Update Participant Payment Status

**Purpose**: Update a participant's payment status from UNPAID to PAID. Triggers automatic email confirmation.

**Endpoint**: `PATCH /api/v1/admin/participants/:id/payment`

**Authentication**: Required (JWT token)

**Request Headers**:
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**URL Parameters**:

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Participant ID |

**Request Body**:
```json
{
  "payment_status": "PAID"
}
```

**Request Body Schema**:

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `payment_status` | string | Yes | Must be "PAID" or "UNPAID" | New payment status |

**Success Response** (HTTP 200 OK):

```json
{
  "success": true,
  "message": "Payment status updated successfully. Confirmation email sent.",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "payment_status": "PAID",
    "updated_at": "2025-12-31T17:46:37Z",
    "email_sent": true
  }
}
```

**Success Response - Email Failed** (HTTP 200 OK):

```json
{
  "success": true,
  "message": "Payment status updated successfully. Email sending failed.",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "payment_status": "PAID",
    "updated_at": "2025-12-31T17:46:37Z",
    "email_sent": false,
    "email_error": "SMTP connection timeout"
  }
}
```

**Note**: Payment status is updated even if email fails (decoupled for resilience).

**Error Responses**:

**400 Bad Request** - Invalid payment status:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_STATUS",
    "message": "Payment status must be either PAID or UNPAID"
  }
}
```

**404 Not Found** - Participant not found:
```json
{
  "success": false,
  "error": {
    "code": "PARTICIPANT_NOT_FOUND",
    "message": "Participant with the specified ID does not exist",
    "details": {
      "id": "550e8400-e29b-41d4-a716-446655440000"
    }
  }
}
```

**401 Unauthorized** - Missing or invalid token:
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required. Please provide a valid token."
  }
}
```

---

### 4. Get Single Participant (Optional for MVP)

**Purpose**: Retrieve detailed information about a specific participant.

**Endpoint**: `GET /api/v1/admin/participants/:id`

**Authentication**: Required (JWT token)

**Request Headers**:
```
Authorization: Bearer <jwt_token>
```

**URL Parameters**:

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Participant ID |

**Success Response** (HTTP 200 OK):

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+6281234567890",
    "instagram_handle": "johndoe",
    "address": "Jl. Example Street No. 123, Jakarta, Indonesia 12345",
    "registration_status": "PENDING",
    "payment_status": "UNPAID",
    "created_at": "2025-12-31T10:30:00Z",
    "updated_at": "2025-12-31T10:30:00Z"
  }
}
```

**Error Responses**:

**404 Not Found**:
```json
{
  "success": false,
  "error": {
    "code": "PARTICIPANT_NOT_FOUND",
    "message": "Participant with the specified ID does not exist"
  }
}
```

---

## JWT Token Structure

**Algorithm**: HS256 (HMAC SHA-256)  
**Secret**: Stored in environment variable `JWT_SECRET`

**Token Payload**:
```json
{
  "admin_id": "550e8400-e29b-41d4-a716-446655440001",
  "email": "admin@tautaurun.com",
  "exp": 1735679197,
  "iat": 1735592797
}
```

**Claims**:
- `admin_id` (string): Admin's UUID
- `email` (string): Admin's email address
- `exp` (integer): Token expiration timestamp (Unix epoch)
- `iat` (integer): Token issued at timestamp (Unix epoch)

---

## Email Trigger Logic

When payment status is updated from **UNPAID → PAID**:

1. **Update Database**: Set `payment_status = 'PAID'` and `updated_at = NOW()`
2. **Trigger Email** (asynchronous, non-blocking):
   - Fetch participant details (name, email)
   - Compose confirmation email using template
   - Send via SMTP
   - Log result to `email_logs` table (SUCCESS or FAILED)
3. **Return Response**: API responds immediately after database update (doesn't wait for email)

**Idempotency**: If status is already PAID, updating to PAID again does nothing (no email sent).

---

## Error Handling

### Common HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET, PATCH operations |
| 201 | Created | Successful POST (not used in admin API) |
| 400 | Bad Request | Validation errors, invalid input |
| 401 | Unauthorized | Missing or invalid authentication token |
| 403 | Forbidden | Expired token, insufficient permissions |
| 404 | Not Found | Participant ID doesn't exist |
| 500 | Internal Server Error | Unexpected server errors |

---

## Example Requests

### Admin Login (cURL)

```bash
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@tautaurun.com",
    "password": "Admin123!"
  }'
```

### List Participants (cURL)

```bash
curl -X GET http://localhost:8080/api/v1/admin/participants \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Update Payment Status to PAID (cURL)

```bash
curl -X PATCH http://localhost:8080/api/v1/admin/participants/550e8400-e29b-41d4-a716-446655440000/payment \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "payment_status": "PAID"
  }'
```

---

## Frontend Integration

**Example: Admin Dashboard API Client (TypeScript)**

```typescript
interface Admin {
  id: string;
  email: string;
}

interface LoginResponse {
  token: string;
  admin: Admin;
  expires_at: string;
}

interface Participant {
  id: string;
  name: string;
  email: string;
  phone: string;
  instagram_handle: string | null;
  address: string;
  registration_status: 'PENDING' | 'CONFIRMED';
  payment_status: 'UNPAID' | 'PAID';
  created_at: string;
  updated_at: string;
}

class AdminAPI {
  private baseURL = '/api/v1/admin';
  private token: string | null = null;

  async login(email: string, password: string): Promise<LoginResponse> {
    const response = await fetch(`${this.baseURL}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error.message);
    }

    const result = await response.json();
    this.token = result.data.token;
    return result.data;
  }

  async getParticipants(): Promise<Participant[]> {
    const response = await fetch(`${this.baseURL}/participants`, {
      headers: {
        'Authorization': `Bearer ${this.token}`,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error.message);
    }

    const result = await response.json();
    return result.data.participants;
  }

  async updatePaymentStatus(
    participantId: string,
    status: 'PAID' | 'UNPAID'
  ): Promise<void> {
    const response = await fetch(
      `${this.baseURL}/participants/${participantId}/payment`,
      {
        method: 'PATCH',
        headers: {
          'Authorization': `Bearer ${this.token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ payment_status: status }),
      }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error.message);
    }
  }
}

export default new AdminAPI();
```

---

## Security Considerations

1. **Token Storage**: Frontend should store JWT in HTTP-only cookies (preferred) or localStorage (acceptable for MVP)
2. **Token Expiration**: Frontend should handle 403 errors by redirecting to login
3. **HTTPS Only**: All admin API calls must be over HTTPS in production
4. **Password Hashing**: Passwords hashed with bcrypt (cost factor 12) before comparison
5. **SQL Injection**: Use parameterized queries only (never string concatenation)
6. **Rate Limiting**: Consider adding rate limiting to login endpoint (5 attempts per 15 minutes)

---

**Document Status**: Ready for implementation  
**Review Required**: Yes (validate JWT structure and email trigger logic)  
**Next Step**: Implement backend handlers with JWT middleware
