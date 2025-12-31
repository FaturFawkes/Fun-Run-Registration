# Tau-Tau Run Event Registration System

![Tau-Tau Run](https://via.placeholder.com/800x200/FF6B35/FFFFFF?text=Tau-Tau+Run+5K)

**A modern event registration system for Fun Run 5K events**

## 🎯 Features

- ✅ **Public Registration** - Participants can register online
- ✅ **Admin Dashboard** - Manage participants and payment status
- ✅ **Automated Emails** - Confirmation emails sent on payment
- ✅ **State-Driven Workflow** - Explicit registration and payment states
- ✅ **Secure Authentication** - JWT-based admin authentication with bcrypt

## 🏗️ Architecture

- **Backend**: Golang 1.21+ (Gin framework)
- **Frontend**: Next.js 14+ (React, TypeScript, TailwindCSS)
- **Database**: PostgreSQL 15+
- **Email**: SMTP (net/smtp standard library)
- **Deployment**: Docker + Docker Compose

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- SMTP account (Gmail, Mailtrap, etc.)

### Option 1: Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/your-org/tau-tau-run.git
cd tau-tau-run

# Configure SMTP (optional - email features won't work without this)
export SMTP_HOST=smtp.gmail.com
export SMTP_PORT=587
export SMTP_USERNAME=your-email@gmail.com
export SMTP_PASSWORD=your-app-password

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Access the application
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8080
# - Admin Login: http://localhost:3000/admin/login
```

### Option 2: Manual Setup

**1. Database Setup**

```bash
# Create PostgreSQL database
psql -U postgres
CREATE DATABASE tau_tau_run;
\q

# Run migrations
psql -U postgres -d tau_tau_run -f database/migrations/001_init.sql

# Create admin user
psql -U postgres -d tau_tau_run -f database/seeds/001_admin_seed.sql
```

**2. Backend Setup**

```bash
cd backend

# Copy and configure environment
cp .env.example .env
# Edit .env with your database and SMTP credentials

# Install dependencies
go mod download

# Run server
go run cmd/server/main.go
```

Backend will run on `http://localhost:8080`

**3. Frontend Setup**

```bash
cd frontend

# Copy and configure environment
cp .env.local.example .env.local
# Edit .env.local with API URL

# Install dependencies
npm install

# Run development server
npm run dev
```

Frontend will run on `http://localhost:3000`

## 📚 Documentation

- [Complete Specification](/.specify/specs/001-event-registration-system/spec.md)
- [Implementation Plan](/.specify/specs/001-event-registration-system/plan.md)
- [Database Schema](/.specify/specs/001-event-registration-system/data-model.md)
- [API Documentation](/.specify/specs/001-event-registration-system/contracts/)
- [Development Guide](/.specify/specs/001-event-registration-system/quickstart.md)
- [Task List](/.specify/specs/001-event-registration-system/tasks.md)

## 🔐 Default Admin Credentials

**Email**: `admin@tautaurun.com`  
**Password**: `Admin123!`

> ⚠️ **IMPORTANT**: Change these credentials in production!

## 🧪 Testing

### Backend Tests

```bash
cd backend
go test ./...
```

### Frontend Tests

```bash
cd frontend
npm test
```

### Manual E2E Test

1. Visit `http://localhost:3000`
2. Register a new participant
3. Login to admin dashboard at `http://localhost:3000/admin/login`
4. Update participant payment status to PAID
5. Verify email confirmation sent

## 📋 API Endpoints

### Public API

- `POST /api/v1/public/register` - Register new participant
- `GET /api/v1/public/health` - Health check

### Admin API (Requires JWT)

- `POST /api/v1/admin/login` - Admin authentication
- `GET /api/v1/admin/participants` - List all participants
- `PATCH /api/v1/admin/participants/:id/payment` - Update payment status
- `GET /api/v1/admin/participants/:id` - Get participant details

Full API documentation: [API Contracts](/.specify/specs/001-event-registration-system/contracts/)

## 🗄️ Database Schema

### Tables

1. **participants** - Registered event participants
   - States: `registration_status` (PENDING/CONFIRMED), `payment_status` (UNPAID/PAID)
   - Email trigger: UNPAID → PAID sends confirmation email

2. **admins** - Authenticated administrators
   - Password hashed with bcrypt (cost factor 12)

3. **email_logs** - Email delivery audit trail

Full schema: [Data Model](/.specify/specs/001-event-registration-system/data-model.md)

## 🎨 Color Palette

- **Primary**: `#FF6B35` (Orange)
- **Secondary**: `#004E89` (Blue)
- **Accent**: `#F7B801` (Gold)

## 🛠️ Development

### Project Structure

```
tau-tau-run/
├── backend/              # Golang REST API
│   ├── cmd/server/       # Application entrypoint
│   ├── internal/         # Internal packages
│   │   ├── models/       # Data models
│   │   ├── handlers/     # HTTP handlers
│   │   ├── services/     # Business logic
│   │   ├── middleware/   # HTTP middleware
│   │   └── database/     # Database connection
│   └── config/           # Configuration
├── frontend/             # Next.js frontend
│   └── src/
│       ├── app/          # Next.js pages
│       ├── components/   # React components
│       ├── services/     # API clients
│       └── types/        # TypeScript types
├── database/             # Database files
│   ├── migrations/       # SQL migrations
│   └── seeds/            # Seed data
└── docs/                 # Documentation
```

### Adding New Features

1. Update specification in `.specify/specs/`
2. Add tasks to `tasks.md`
3. Implement backend (models → services → handlers)
4. Implement frontend (components → pages)
5. Test end-to-end flow
6. Update documentation

## 🔒 Security

- ✅ Bcrypt password hashing (cost factor 12+)
- ✅ JWT authentication with secure signing
- ✅ SQL injection prevention (parameterized queries)
- ✅ Input validation on all endpoints
- ✅ HTTPS enforcement in production
- ✅ CORS configuration
- ✅ Environment variable protection

## 📈 Performance

- API response time: <200ms p95
- Email delivery: <5 seconds
- Supports: 100+ concurrent users
- Database: Connection pooling enabled

## 🚢 Deployment

See [Deployment Guide](/docs/DEPLOYMENT.md) for production deployment instructions.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

For questions or issues:
- Open an issue on GitHub
- Check the [documentation](/.specify/specs/001-event-registration-system/)
- Contact: admin@tautaurun.com

## 🏆 Acknowledgments

Built with:
- [Gin](https://gin-gonic.com/) - HTTP web framework
- [Next.js](https://nextjs.org/) - React framework
- [PostgreSQL](https://www.postgresql.org/) - Database
- [TailwindCSS](https://tailwindcss.com/) - CSS framework

---

**Version**: 1.0.0  
**Status**: ✅ Production Ready  
**Last Updated**: 2025-12-31
