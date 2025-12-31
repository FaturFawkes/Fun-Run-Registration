# Quick Start Guide - Tau-Tau Run

## 🚀 How to Run the Frontend

### Prerequisites
- Node.js 18+ installed
- Backend already running on http://localhost:8080

### Steps

1. **Navigate to frontend directory:**
   ```bash
   cd /home/fatur/Documents/Projects/Tau-TauRun/frontend
   ```

2. **Install dependencies (first time only):**
   ```bash
   npm install
   ```
   ⏱️ This takes 1-2 minutes

3. **Configure environment (first time only):**
   ```bash
   cp .env.local.example .env.local
   ```
   The defaults work fine - no need to edit!

4. **Run development server:**
   ```bash
   npm run dev
   ```

5. **Open in browser:**
   - Public page: http://localhost:3000
   - Admin login: http://localhost:3000/admin/login

---

## 🎯 Complete System Startup

### Terminal 1: Database (One-Time Setup)
```bash
# Create database
psql -U postgres -c "CREATE DATABASE tau_tau_run;"

# Run migrations
psql -U postgres -d tau_tau_run -f database/migrations/001_init.sql

# Add admin user
psql -U postgres -d tau_tau_run -f database/seeds/001_admin_seed.sql
```

### Terminal 2: Backend
```bash
cd /home/fatur/Documents/Projects/Tau-TauRun/backend

# First time: setup environment
cp .env.example .env
# Edit .env and set DB_PASSWORD

# Run server
go run cmd/server/main.go
```
✅ Backend runs on: **http://localhost:8080**

### Terminal 3: Frontend
```bash
cd /home/fatur/Documents/Projects/Tau-TauRun/frontend

# First time: install dependencies
npm install
cp .env.local.example .env.local

# Run server
npm run dev
```
✅ Frontend runs on: **http://localhost:3000**

---

## 🌐 Available Pages

| Page | URL | Description |
|------|-----|-------------|
| **Public Registration** | http://localhost:3000 | Event landing page with registration form |
| **Admin Login** | http://localhost:3000/admin/login | Admin authentication |
| **Admin Dashboard** | http://localhost:3000/admin/dashboard | Protected admin area |

---

## 🔑 Default Admin Credentials

- **Email:** admin@tautaurun.com
- **Password:** Admin123!

---

## 🛠️ Troubleshooting

### "Cannot connect to API"
✅ Make sure backend is running on http://localhost:8080

### "Database connection failed"
✅ Check PostgreSQL is running: `sudo systemctl status postgresql`
✅ Verify .env file has correct DB_PASSWORD

### "Port 3000 already in use"
✅ Kill existing process: `lsof -ti:3000 | xargs kill -9`
✅ Or use different port: `npm run dev -- -p 3001`

### "Module not found"
✅ Run: `npm install` again
✅ Delete node_modules and reinstall: `rm -rf node_modules && npm install`

---

## 📝 Quick Commands Reference

```bash
# Frontend only (assumes backend is running)
cd frontend
npm install          # First time only
npm run dev          # Start dev server
npm run build        # Production build
npm run start        # Run production build

# Backend only
cd backend
go mod download      # First time only
go run cmd/server/main.go   # Start server

# Database
psql -U postgres -d tau_tau_run   # Connect to database
```

---

## ✅ What's Working Now

- ✅ **Phase 1-2:** Foundation complete
- ✅ **Phase 3:** Public registration working
- ✅ **Phase 4:** Admin authentication working
- ✅ **Phase 5:** Payment management working
- ✅ **Phase 6:** Email automation working
- ✅ **Phase 7:** Integration tested

---

## 📧 SMTP Configuration (Email Features)

The system sends confirmation emails when payment status is updated to PAID. Configure SMTP to enable this feature:

### Option 1: Gmail (Development/Testing)

1. **Enable 2-Factor Authentication** on your Google account
2. **Generate App Password**: https://myaccount.google.com/apppasswords
3. **Update backend/.env**:
   ```env
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-16-char-app-password
   SMTP_FROM_EMAIL=noreply@tautaurun.com
   SMTP_FROM_NAME=Tau-Tau Run Team
   ```

### Option 2: Mailtrap (Development/Testing)

Perfect for testing emails without sending to real addresses.

1. **Sign up** at https://mailtrap.io (free tier available)
2. **Get SMTP credentials** from your inbox
3. **Update backend/.env**:
   ```env
   SMTP_HOST=smtp.mailtrap.io
   SMTP_PORT=587
   SMTP_USERNAME=your-mailtrap-username
   SMTP_PASSWORD=your-mailtrap-password
   SMTP_FROM_EMAIL=noreply@tautaurun.com
   SMTP_FROM_NAME=Tau-Tau Run Team
   ```

### Option 3: SendGrid (Production)

Recommended for production use.

1. **Sign up** at https://sendgrid.com
2. **Create API Key** in Settings → API Keys
3. **Update backend/.env**:
   ```env
   SMTP_HOST=smtp.sendgrid.net
   SMTP_PORT=587
   SMTP_USERNAME=apikey
   SMTP_PASSWORD=your-sendgrid-api-key
   SMTP_FROM_EMAIL=noreply@tautaurun.com
   SMTP_FROM_NAME=Tau-Tau Run Team
   ```

### Verify Email Configuration

```bash
# Restart backend after updating .env
cd backend
go run cmd/server/main.go

# You should see:
# [EMAIL] INFO: SMTP configuration validated: smtp.gmail.com:587
```

**Note:** Without SMTP configuration, the system will still work but emails will be logged as FAILED in the database.

---

## 🎯 Test the System

1. Open http://localhost:3000
2. Fill out registration form and submit
3. Go to http://localhost:3000/admin/login
4. Login with admin@tautaurun.com / Admin123!
5. View participant list in dashboard
6. Update payment status to PAID
7. Check email inbox for confirmation email

---

**Need help?** Check the main README.md, [API Documentation](docs/API.md), or [Deployment Guide](docs/DEPLOYMENT.md)
