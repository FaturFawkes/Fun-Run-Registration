# API Contracts: Public API

**Feature**: 001-event-registration-system  
**Base URL**: `/api/v1/public`  
**Authentication**: None (public endpoints)  
**Date**: 2025-12-31

---

## Endpoints

### 1. Register Participant

**Purpose**: Allow public users to register for the Tau-Tau Run Fun Run 5K event.

**Endpoint**: `POST /api/v1/public/register`

**Authentication**: None

**Request Headers**:
```
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "+6281234567890",
  "instagram_handle": "johndoe",
  "address": "Jl. Example Street No. 123, Jakarta, Indonesia 12345"
}
```

**Request Body Schema**:

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `name` | string | Yes | 2-255 characters | Full name of participant |
| `email` | string | Yes | Valid email format, unique | Email address |
| `phone` | string | Yes | 10-50 characters | Phone number (any format) |
| `instagram_handle` | string | No | Max 100 characters, alphanumeric + underscore | Instagram username (without @) |
| `address` | string | Yes | Min 10 characters | Full address |

**Success Response** (HTTP 201 Created):

```json
{
  "success": true,
  "message": "Registration successful! Your payment status is pending.",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john.doe@example.com",
    "registration_status": "PENDING",
    "payment_status": "UNPAID"
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
        "message": "Email must be a valid email address"
      },
      {
        "field": "name",
        "message": "Name is required and must be 2-255 characters"
      }
    ]
  }
}
```

**409 Conflict** - Duplicate email:
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

**500 Internal Server Error** - Server error:
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

### 2. Health Check (Optional)

**Purpose**: Check if the public API is running and healthy.

**Endpoint**: `GET /api/v1/public/health`

**Authentication**: None

**Request Headers**: None

**Success Response** (HTTP 200 OK):

```json
{
  "status": "healthy",
  "timestamp": "2025-12-31T17:46:37Z",
  "version": "1.0.0"
}
```

---

## Validation Rules

### Email Validation
- Must match regex: `^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$`
- Case-insensitive duplicate check
- Trimmed before validation

### Name Validation
- Minimum length: 2 characters
- Maximum length: 255 characters
- Trimmed before validation
- Allowed characters: letters, spaces, hyphens, apostrophes

### Phone Validation
- Minimum length: 10 characters
- Maximum length: 50 characters
- Allowed characters: digits, spaces, hyphens, parentheses, plus sign
- Example formats: `+6281234567890`, `0812-3456-7890`, `(021) 12345678`

### Instagram Handle Validation
- Optional field
- Maximum length: 100 characters
- Allowed characters: alphanumeric, underscores, periods
- Must not start with @ (stripped if provided)

### Address Validation
- Minimum length: 10 characters
- Maximum length: 1000 characters (TEXT field)
- Trimmed before validation

---

## Rate Limiting (Future Consideration)

**Recommendation for Production**:
- Limit: 5 registrations per IP per hour
- Prevents spam and abuse
- Not implemented in MVP (add in production)

---

## CORS Configuration

**Allowed Origins**: 
- Development: `http://localhost:3000`
- Production: `https://tautaurun.com` (configure via environment variable)

**Allowed Methods**: `GET, POST, OPTIONS`

**Allowed Headers**: `Content-Type, Accept`

**Max Age**: 3600 seconds (1 hour)

---

## Example Requests

### Valid Registration (cURL)

```bash
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane.smith@example.com",
    "phone": "+6281234567890",
    "instagram_handle": "janesmith",
    "address": "Jl. Merdeka No. 45, Bandung, Indonesia 40111"
  }'
```

### Invalid Registration - Missing Email (cURL)

```bash
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "phone": "+6281234567890",
    "address": "Jl. Merdeka No. 45, Bandung, Indonesia 40111"
  }'
```

Expected Response: `400 Bad Request` with validation error for missing email.

### Duplicate Registration (cURL)

```bash
# First registration
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "duplicate@example.com",
    "phone": "+6281234567890",
    "address": "Jl. Example Street No. 123, Jakarta, Indonesia 12345"
  }'

# Second registration with same email
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Different Name",
    "email": "duplicate@example.com",
    "phone": "+6289999999999",
    "address": "Different address"
  }'
```

Expected Second Response: `409 Conflict` with duplicate email error.

---

## Security Considerations

1. **Input Sanitization**: All inputs must be sanitized to prevent XSS and SQL injection
2. **SQL Injection Prevention**: Use parameterized queries only
3. **XSS Prevention**: Escape HTML characters in responses (if any user data is rendered)
4. **HTTPS Only**: All API calls must be over HTTPS in production
5. **No Sensitive Data**: Never expose database internals or server details in error messages

---

## Frontend Integration

**Example: Registration Form Submit (TypeScript/React)**

```typescript
interface RegistrationData {
  name: string;
  email: string;
  phone: string;
  instagram_handle?: string;
  address: string;
}

async function registerParticipant(data: RegistrationData) {
  try {
    const response = await fetch('/api/v1/public/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    const result = await response.json();

    if (!response.ok) {
      // Handle error
      if (result.error.code === 'DUPLICATE_EMAIL') {
        throw new Error('This email is already registered.');
      } else if (result.error.code === 'VALIDATION_ERROR') {
        throw new Error(result.error.details.map(d => d.message).join(', '));
      } else {
        throw new Error('Registration failed. Please try again.');
      }
    }

    // Success
    return result.data;
  } catch (error) {
    console.error('Registration error:', error);
    throw error;
  }
}
```

---

**Document Status**: Ready for implementation  
**Review Required**: Yes (validate response formats and error codes)  
**Next Step**: Implement backend handler in Golang
