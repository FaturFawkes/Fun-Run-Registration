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
- ⏳ **Phase 5:** Payment management (next)

---

## 🎯 Test the System

1. Open http://localhost:3000
2. Fill out registration form and submit
3. Go to http://localhost:3000/admin/login
4. Login with admin@tautaurun.com / Admin123!
5. You'll see the dashboard

---

**Need help?** Check the main README.md or the specification files in `.specify/specs/001-event-registration-system/`
