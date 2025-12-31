# Quickstart Guide: Event Registration System

**Feature**: 001-event-registration-system  
**Date**: 2025-12-31  
**Status**: Development Setup Guide

---

## Prerequisites

Before you begin, ensure you have the following installed on your system:

### Required Software

| Software | Version | Installation |
|----------|---------|--------------|
| **Go** | 1.21+ | [Download](https://golang.org/dl/) |
| **Node.js** | 18+ | [Download](https://nodejs.org/) |
| **PostgreSQL** | 15+ | [Download](https://www.postgresql.org/download/) |
| **Git** | Latest | [Download](https://git-scm.com/) |
| **Docker** (Optional) | Latest | [Download](https://www.docker.com/) |

### Verify Installations

```bash
go version        # Should show go1.21 or higher
node --version    # Should show v18 or higher
npm --version     # Should show v9 or higher
psql --version    # Should show 15 or higher
git --version     # Any recent version
```

---

## Quick Setup (5 Minutes)

### Option 1: Manual Setup

#### Step 1: Clone Repository

```bash
git clone https://github.com/your-org/tau-tau-run.git
cd tau-tau-run
```

#### Step 2: Setup Database

```bash
# Create PostgreSQL database
psql -U postgres
CREATE DATABASE tau_tau_run;
\q

# Run migration
psql -U postgres -d tau_tau_run -f database/migrations/001_init.sql

# Create admin user (optional seed)
psql -U postgres -d tau_tau_run -f database/seeds/001_admin_seed.sql
```

#### Step 3: Setup Backend

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env with your settings
# nano .env or vim .env or code .env

# Install dependencies
go mod download

# Run backend server
go run cmd/server/main.go
```

Backend should now be running on `http://localhost:8080`

#### Step 4: Setup Frontend

Open a new terminal:

```bash
cd frontend

# Copy environment file
cp .env.local.example .env.local

# Edit .env.local with API URL
# nano .env.local or vim .env.local or code .env.local

# Install dependencies
npm install

# Run development server
npm run dev
```

Frontend should now be running on `http://localhost:3000`

---

### Option 2: Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/your-org/tau-tau-run.git
cd tau-tau-run

# Start all services (database, backend, frontend)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

Services will be available at:
- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`
- PostgreSQL: `localhost:5432`

---

## Environment Configuration

### Backend `.env` File

Create `backend/.env` with the following:

```bash
# Server Configuration
PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=tau_tau_run
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION_HOURS=24

# SMTP Configuration (for email sending)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@tautaurun.com
SMTP_FROM_NAME=Tau-Tau Run Team

# Event Configuration
EVENT_NAME=Tau-Tau Run Fun Run 5K
EVENT_DATE=2026-02-15
EVENT_LOCATION=Gelora Bung Karno Stadium, Jakarta
```

**Important Security Notes**:
- **Never commit `.env` file to Git** (already in `.gitignore`)
- Change `JWT_SECRET` to a random 32+ character string in production
- For Gmail SMTP, use an [App Password](https://support.google.com/accounts/answer/185833), not your regular password

---

### Frontend `.env.local` File

Create `frontend/.env.local` with the following:

```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# Event Configuration (for display)
NEXT_PUBLIC_EVENT_NAME=Tau-Tau Run Fun Run 5K
NEXT_PUBLIC_EVENT_DATE=February 15, 2026
NEXT_PUBLIC_EVENT_LOCATION=Gelora Bung Karno Stadium, Jakarta
```

**Note**: All variables starting with `NEXT_PUBLIC_` are exposed to the browser.

---

## Database Setup (Detailed)

### Create Database and User

```sql
-- Connect to PostgreSQL as superuser
psql -U postgres

-- Create database
CREATE DATABASE tau_tau_run;

-- Create dedicated user (optional but recommended)
CREATE USER tau_tau_admin WITH PASSWORD 'secure_password_here';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE tau_tau_run TO tau_tau_admin;

-- Connect to the database
\c tau_tau_run

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO tau_tau_admin;

-- Exit
\q
```

### Run Migrations

```bash
# Apply initial schema
psql -U tau_tau_admin -d tau_tau_run -f database/migrations/001_init.sql

# Verify tables were created
psql -U tau_tau_admin -d tau_tau_run -c "\dt"
```

Expected output:
```
           List of relations
 Schema |      Name      | Type  |     Owner
--------+----------------+-------+---------------
 public | admins         | table | tau_tau_admin
 public | email_logs     | table | tau_tau_admin
 public | participants   | table | tau_tau_admin
```

### Create Admin User

**Option 1: Using Seed Script**

```bash
psql -U tau_tau_admin -d tau_tau_run -f database/seeds/001_admin_seed.sql
```

**Option 2: Manual SQL**

```sql
-- Connect to database
psql -U tau_tau_admin -d tau_tau_run

-- Insert admin (password: Admin123!)
INSERT INTO admins (email, password_hash)
VALUES (
  'admin@tautaurun.com',
  '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYlK4Qr1WZK'
);

-- Verify
SELECT id, email, created_at FROM admins;

\q
```

**Option 3: Using Golang Script** (recommended for production)

```bash
cd backend
go run scripts/create_admin.go --email admin@tautaurun.com --password YourSecurePassword
```

---

## SMTP Configuration Guide

### Gmail Setup (Development/Testing)

1. **Enable 2-Factor Authentication** on your Google account
2. **Create App Password**:
   - Go to [Google App Passwords](https://myaccount.google.com/apppasswords)
   - Select "Mail" and "Other (Custom name)"
   - Enter "Tau-Tau Run" and click Generate
   - Copy the 16-character password

3. **Update `.env`**:
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
SMTP_FROM_EMAIL=noreply@tautaurun.com
```

### Other SMTP Providers

**Mailtrap (Testing)**:
```bash
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USERNAME=your-mailtrap-username
SMTP_PASSWORD=your-mailtrap-password
```

**SendGrid (Production - if allowed)**:
Note: Constitution prohibits external email APIs, but SendGrid SMTP relay is acceptable.

```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
```

---

## Testing Your Setup

### 1. Test Backend Health

```bash
curl http://localhost:8080/api/v1/public/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2025-12-31T17:46:37Z",
  "version": "1.0.0"
}
```

### 2. Test Database Connection

```bash
# Backend should log successful connection on startup
# Check backend terminal for:
# [GIN] Database connected successfully
```

### 3. Test Public Registration

```bash
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "phone": "+6281234567890",
    "instagram_handle": "testuser",
    "address": "Jl. Test No. 123, Jakarta, Indonesia"
  }'
```

Expected response (201 Created):
```json
{
  "success": true,
  "message": "Registration successful! Your payment status is pending.",
  "data": {
    "id": "...",
    "email": "test@example.com",
    "registration_status": "PENDING",
    "payment_status": "UNPAID"
  }
}
```

### 4. Test Admin Login

```bash
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@tautaurun.com",
    "password": "Admin123!"
  }'
```

Expected response (200 OK):
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin": {
      "id": "...",
      "email": "admin@tautaurun.com"
    },
    "expires_at": "..."
  }
}
```

### 5. Test Frontend

Open browser and navigate to:
- Public page: `http://localhost:3000`
- Admin login: `http://localhost:3000/admin/login`

---

## Development Workflow

### 1. Starting Development Session

```bash
# Terminal 1: Backend
cd backend
go run cmd/server/main.go

# Terminal 2: Frontend
cd frontend
npm run dev

# Terminal 3: Database (if not using Docker)
# PostgreSQL should be running as a service
```

### 2. Making Changes

**Backend Changes**:
- Edit files in `backend/internal/`
- Server auto-restarts with `air` (live reload) or manually restart
- Run tests: `cd backend && go test ./...`

**Frontend Changes**:
- Edit files in `frontend/src/`
- Next.js hot-reloads automatically
- No restart needed

**Database Changes**:
- Create new migration file: `database/migrations/002_description.sql`
- Apply manually: `psql -U tau_tau_admin -d tau_tau_run -f database/migrations/002_description.sql`

### 3. Running Tests

```bash
# Backend tests
cd backend
go test ./... -v

# Frontend tests (if implemented)
cd frontend
npm test
```

---

## Common Issues & Troubleshooting

### Issue: "Database connection refused"

**Solution**:
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Start PostgreSQL
sudo systemctl start postgresql

# Verify connection
psql -U postgres -c "SELECT version();"
```

### Issue: "JWT token invalid"

**Solution**:
- Ensure `JWT_SECRET` in backend `.env` is at least 32 characters
- Clear browser cookies/localStorage
- Re-login to get a fresh token

### Issue: "Email not sending"

**Solution**:
- Check SMTP credentials in `.env`
- Test SMTP connection: `telnet smtp.gmail.com 587`
- Check backend logs for SMTP errors
- Verify Gmail App Password is correct (not regular password)

### Issue: "Port 8080 already in use"

**Solution**:
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in backend/.env
PORT=8081
```

### Issue: "npm install fails"

**Solution**:
```bash
# Clear npm cache
npm cache clean --force

# Delete node_modules and package-lock.json
rm -rf node_modules package-lock.json

# Reinstall
npm install
```

---

## Next Steps

1. ✅ Setup complete - all services running
2. 📖 Read API documentation: `contracts/public-api.md` and `contracts/admin-api.md`
3. 🏗️ Start implementing features (see `tasks.md` when generated)
4. 🧪 Test end-to-end flow:
   - Register participant
   - Login as admin
   - Update payment to PAID
   - Verify email received
5. 🚀 Deploy to production (see `docs/DEPLOYMENT.md`)

---

## Helpful Commands Cheatsheet

```bash
# Backend
cd backend
go run cmd/server/main.go          # Run server
go test ./...                       # Run tests
go mod tidy                         # Clean dependencies
go build -o bin/server cmd/server/main.go  # Build binary

# Frontend
cd frontend
npm run dev                         # Development server
npm run build                       # Production build
npm run start                       # Production server
npm run lint                        # Lint code

# Database
psql -U tau_tau_admin -d tau_tau_run         # Connect to DB
psql -U tau_tau_admin -d tau_tau_run -f migration.sql  # Run migration
psql -U tau_tau_admin -d tau_tau_run -c "SELECT * FROM participants;"  # Query

# Docker
docker-compose up -d                # Start all services
docker-compose down                 # Stop all services
docker-compose logs -f backend      # View backend logs
docker-compose restart backend      # Restart backend only
docker-compose exec db psql -U postgres  # Access database container
```

---

**Document Status**: Ready for use  
**Last Updated**: 2025-12-31  
**For Questions**: See main README.md or contact project maintainer
