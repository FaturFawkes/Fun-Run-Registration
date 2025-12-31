# Deployment Guide - Tau-Tau Run Event Registration System

**Version:** 1.0.0  
**Last Updated:** 2026-01-01

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Production Environment Setup](#production-environment-setup)
- [Database Setup](#database-setup)
- [Backend Deployment](#backend-deployment)
- [Frontend Deployment](#frontend-deployment)
- [SMTP Configuration](#smtp-configuration)
- [Security Checklist](#security-checklist)
- [Monitoring & Logging](#monitoring--logging)
- [Backup & Recovery](#backup--recovery)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Software
- **Go:** 1.21 or higher
- **Node.js:** 18 or higher
- **PostgreSQL:** 15 or higher
- **Git:** For deployment
- **Systemd or PM2:** For process management

### Recommended Infrastructure
- **Server:** 2 CPU cores, 4GB RAM minimum
- **Storage:** 20GB SSD
- **OS:** Ubuntu 22.04 LTS or similar
- **Reverse Proxy:** Nginx or Apache
- **SSL:** Let's Encrypt or similar

---

## Production Environment Setup

### 1. Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y postgresql postgresql-contrib nginx certbot python3-certbot-nginx

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Node.js (via nvm)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 18
nvm use 18
```

### 2. Create Application User

```bash
# Create dedicated user
sudo useradd -m -s /bin/bash tautaurun
sudo mkdir -p /opt/tautaurun
sudo chown tautaurun:tautaurun /opt/tautaurun
```

### 3. Clone Repository

```bash
sudo su - tautaurun
cd /opt/tautaurun
git clone https://github.com/your-org/tau-tau-run.git
cd tau-tau-run
```

---

## Database Setup

### 1. Create Production Database

```bash
# Switch to postgres user
sudo -u postgres psql

-- Create database and user
CREATE DATABASE tau_tau_run_prod;
CREATE USER tautaurun WITH ENCRYPTED PASSWORD 'STRONG_PASSWORD_HERE';
GRANT ALL PRIVILEGES ON DATABASE tau_tau_run_prod TO tautaurun;
\q
```

### 2. Run Migrations

```bash
cd /opt/tautaurun/tau-tau-run

# Run schema migration
PGPASSWORD='STRONG_PASSWORD_HERE' psql -h localhost -U tautaurun -d tau_tau_run_prod \
  -f database/migrations/001_init.sql

# Create admin user (change password in seed file first!)
PGPASSWORD='STRONG_PASSWORD_HERE' psql -h localhost -U tautaurun -d tau_tau_run_prod \
  -f database/seeds/001_admin_seed.sql
```

### 3. Configure PostgreSQL for Production

Edit `/etc/postgresql/15/main/postgresql.conf`:

```ini
# Connection settings
max_connections = 100
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 2621kB
min_wal_size = 1GB
max_wal_size = 4GB
```

Restart PostgreSQL:
```bash
sudo systemctl restart postgresql
```

---

## Backend Deployment

### 1. Configure Environment

```bash
cd /opt/tautaurun/tau-tau-run/backend

# Create production .env
cp .env.example .env
nano .env
```

**Production `.env` Configuration:**

```env
# SERVER
PORT=8080
ENV=production

# DATABASE
DB_HOST=localhost
DB_PORT=5432
DB_USER=tautaurun
DB_PASSWORD=STRONG_DATABASE_PASSWORD
DB_NAME=tau_tau_run_prod
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=20
DB_MAX_IDLE_CONNECTIONS=10

# JWT - IMPORTANT: Generate a strong random secret!
JWT_SECRET=REPLACE_WITH_RANDOM_64_CHARACTER_STRING
JWT_EXPIRATION_HOURS=24

# SMTP (Gmail example)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@tautaurun.com
SMTP_FROM_NAME=Tau-Tau Run Team

# EVENT
EVENT_NAME=Tau-Tau Run Fun Run 5K
EVENT_DATE=2026-02-15
EVENT_LOCATION=Gelora Bung Karno Stadium, Jakarta
EVENT_DESCRIPTION=Join us for an exciting 5K fun run event!

# CORS
CORS_ALLOWED_ORIGINS=https://tautaurun.com,https://www.tautaurun.com
```

**Generate Strong JWT Secret:**
```bash
openssl rand -base64 48
```

### 2. Build Backend

```bash
cd /opt/tautaurun/tau-tau-run/backend

# Install dependencies
go mod download

# Build binary
go build -o tau-tau-run-api cmd/server/main.go

# Make executable
chmod +x tau-tau-run-api
```

### 3. Create Systemd Service

Create `/etc/systemd/system/tautaurun-api.service`:

```ini
[Unit]
Description=Tau-Tau Run API Server
After=network.target postgresql.service

[Service]
Type=simple
User=tautaurun
Group=tautaurun
WorkingDirectory=/opt/tautaurun/tau-tau-run/backend
ExecStart=/opt/tautaurun/tau-tau-run/backend/tau-tau-run-api
Restart=always
RestartSec=5
StandardOutput=append:/var/log/tautaurun/api.log
StandardError=append:/var/log/tautaurun/api-error.log

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/tautaurun

# Environment
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
```

### 4. Start Backend Service

```bash
# Create log directory
sudo mkdir -p /var/log/tautaurun
sudo chown tautaurun:tautaurun /var/log/tautaurun

# Reload systemd
sudo systemctl daemon-reload

# Enable and start service
sudo systemctl enable tautaurun-api
sudo systemctl start tautaurun-api

# Check status
sudo systemctl status tautaurun-api

# View logs
sudo journalctl -u tautaurun-api -f
```

---

## Frontend Deployment

### 1. Configure Environment

```bash
cd /opt/tautaurun/tau-tau-run/frontend

# Create production .env
cp .env.local.example .env.local
nano .env.local
```

**Production `.env.local`:**

```env
NEXT_PUBLIC_API_URL=https://api.tautaurun.com/api/v1
NEXT_PUBLIC_EVENT_NAME=Tau-Tau Run Fun Run 5K
NEXT_PUBLIC_EVENT_DATE=February 15, 2026
NEXT_PUBLIC_EVENT_LOCATION=Gelora Bung Karno Stadium, Jakarta
```

### 2. Build Frontend

```bash
cd /opt/tautaurun/tau-tau-run/frontend

# Install dependencies
npm ci --production

# Build for production
npm run build

# Test production build
npm run start
```

### 3. Create Systemd Service

Create `/etc/systemd/system/tautaurun-frontend.service`:

```ini
[Unit]
Description=Tau-Tau Run Frontend
After=network.target

[Service]
Type=simple
User=tautaurun
Group=tautaurun
WorkingDirectory=/opt/tautaurun/tau-tau-run/frontend
ExecStart=/home/tautaurun/.nvm/versions/node/v18.19.0/bin/npm run start
Restart=always
RestartSec=5
StandardOutput=append:/var/log/tautaurun/frontend.log
StandardError=append:/var/log/tautaurun/frontend-error.log

# Environment
Environment="NODE_ENV=production"
Environment="PORT=3000"

[Install]
WantedBy=multi-user.target
```

### 4. Start Frontend Service

```bash
sudo systemctl daemon-reload
sudo systemctl enable tautaurun-frontend
sudo systemctl start tautaurun-frontend
sudo systemctl status tautaurun-frontend
```

---

## Nginx Configuration

### 1. Backend Reverse Proxy

Create `/etc/nginx/sites-available/tautaurun-api`:

```nginx
server {
    listen 80;
    server_name api.tautaurun.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.tautaurun.com;

    # SSL certificates (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/api.tautaurun.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.tautaurun.com/privkey.pem;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Logging
    access_log /var/log/nginx/api.tautaurun.com-access.log;
    error_log /var/log/nginx/api.tautaurun.com-error.log;

    # Proxy to backend
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

### 2. Frontend Reverse Proxy

Create `/etc/nginx/sites-available/tautaurun`:

```nginx
server {
    listen 80;
    server_name tautaurun.com www.tautaurun.com;
    return 301 https://tautaurun.com$request_uri;
}

server {
    listen 443 ssl http2;
    server_name www.tautaurun.com;
    return 301 https://tautaurun.com$request_uri;
}

server {
    listen 443 ssl http2;
    server_name tautaurun.com;

    # SSL certificates
    ssl_certificate /etc/letsencrypt/live/tautaurun.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/tautaurun.com/privkey.pem;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Logging
    access_log /var/log/nginx/tautaurun.com-access.log;
    error_log /var/log/nginx/tautaurun.com-error.log;

    # Proxy to Next.js
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 3. Enable Sites and Get SSL

```bash
# Enable sites
sudo ln -s /etc/nginx/sites-available/tautaurun-api /etc/nginx/sites-enabled/
sudo ln -s /etc/nginx/sites-available/tautaurun /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Get SSL certificates
sudo certbot --nginx -d tautaurun.com -d www.tautaurun.com
sudo certbot --nginx -d api.tautaurun.com

# Restart Nginx
sudo systemctl restart nginx
```

---

## SMTP Configuration

### Option 1: Gmail (Development/Small Scale)

1. Enable 2-Factor Authentication on your Google account
2. Generate an App Password: https://myaccount.google.com/apppasswords
3. Use in `.env`:
   ```env
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-16-char-app-password
   ```

### Option 2: SendGrid (Production)

1. Sign up at https://sendgrid.com
2. Create API key
3. Use in `.env`:
   ```env
   SMTP_HOST=smtp.sendgrid.net
   SMTP_PORT=587
   SMTP_USERNAME=apikey
   SMTP_PASSWORD=your-sendgrid-api-key
   ```

### Option 3: AWS SES (Production)

1. Verify your domain in AWS SES
2. Create SMTP credentials
3. Use in `.env`:
   ```env
   SMTP_HOST=email-smtp.us-east-1.amazonaws.com
   SMTP_PORT=587
   SMTP_USERNAME=your-ses-username
   SMTP_PASSWORD=your-ses-password
   ```

---

## Security Checklist

### ✅ Pre-Deployment

- [ ] Change default admin password in seed file
- [ ] Generate strong JWT secret (64+ characters)
- [ ] Use strong database password
- [ ] Update all passwords in `.env` files
- [ ] Enable SSL/TLS for database connections
- [ ] Configure firewall (UFW or iptables)
- [ ] Disable SSH password authentication
- [ ] Enable automatic security updates

### ✅ Application Security

- [ ] Set `GIN_MODE=release` in production
- [ ] Configure CORS with specific origins only
- [ ] Use HTTPS for all connections
- [ ] Implement rate limiting (future)
- [ ] Regular security audits
- [ ] Keep dependencies updated

### ✅ Database Security

- [ ] PostgreSQL accessible only from localhost
- [ ] Use SSL mode (`sslmode=require`)
- [ ] Regular backups configured
- [ ] Backup encryption enabled
- [ ] Strong password policy enforced

### ✅ Server Security

```bash
# Configure firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow http
sudo ufw allow https
sudo ufw enable

# Fail2ban for brute force protection
sudo apt install fail2ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

---

## Monitoring & Logging

### Application Logs

```bash
# View backend logs
sudo journalctl -u tautaurun-api -f

# View frontend logs
sudo journalctl -u tautaurun-frontend -f

# View Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Database Monitoring

```bash
# Check database connections
sudo -u postgres psql -c "SELECT * FROM pg_stat_activity;"

# Check database size
sudo -u postgres psql -c "SELECT pg_database_size('tau_tau_run_prod');"
```

### Email Monitoring

Query email logs:
```sql
SELECT 
  status,
  COUNT(*) as count,
  DATE(sent_at) as date
FROM email_logs
GROUP BY status, DATE(sent_at)
ORDER BY date DESC;
```

---

## Backup & Recovery

### Automated Database Backups

Create `/opt/tautaurun/scripts/backup-db.sh`:

```bash
#!/bin/bash

BACKUP_DIR="/opt/tautaurun/backups"
DATE=$(date +%Y%m%d_%H%M%S)
FILENAME="tau_tau_run_prod_${DATE}.sql.gz"

mkdir -p $BACKUP_DIR

PGPASSWORD='STRONG_PASSWORD' pg_dump \
  -h localhost \
  -U tautaurun \
  -d tau_tau_run_prod \
  | gzip > "$BACKUP_DIR/$FILENAME"

# Keep only last 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

echo "Backup completed: $FILENAME"
```

Add to crontab:
```bash
# Daily backup at 2 AM
0 2 * * * /opt/tautaurun/scripts/backup-db.sh >> /var/log/tautaurun/backup.log 2>&1
```

### Restore from Backup

```bash
gunzip -c backup.sql.gz | PGPASSWORD='STRONG_PASSWORD' psql \
  -h localhost \
  -U tautaurun \
  -d tau_tau_run_prod
```

---

## Troubleshooting

### Backend Won't Start

```bash
# Check logs
sudo journalctl -u tautaurun-api -n 50

# Common issues:
# - Database connection failed: Check DB credentials
# - Port already in use: Check if another process is using port 8080
# - Permission denied: Check file ownership
```

### Email Not Sending

```bash
# Check email logs in database
PGPASSWORD='password' psql -U tautaurun -d tau_tau_run_prod \
  -c "SELECT * FROM email_logs ORDER BY sent_at DESC LIMIT 10;"

# Common issues:
# - SMTP credentials invalid
# - Firewall blocking port 587
# - Gmail blocking less secure apps (use App Password)
```

### High CPU/Memory Usage

```bash
# Check process stats
top
htop

# Restart services if needed
sudo systemctl restart tautaurun-api
sudo systemctl restart tautaurun-frontend
```

---

## Production Checklist

Before going live:

- [ ] All tests passing
- [ ] Security checklist completed
- [ ] SSL certificates installed
- [ ] Backups configured and tested
- [ ] Monitoring in place
- [ ] Documentation updated
- [ ] Admin password changed
- [ ] SMTP tested with real emails
- [ ] Load testing completed
- [ ] Disaster recovery plan documented

---

**Support:** For deployment issues, check logs first, then consult the main README.md or contact the development team.
