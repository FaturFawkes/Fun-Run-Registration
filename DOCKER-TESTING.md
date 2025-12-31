# Panduan Testing End-to-End dengan Docker Compose

## 🚀 Quick Start - Testing Lokal dengan Docker

### Persiapan Awal

```bash
cd /home/fatur/Documents/Projects/Tau-TauRun

# Stop services yang masih running
ps aux | grep "go run" | awk '{print $2}' | xargs -r kill 2>/dev/null

# Clean up containers lama
docker-compose down 2>/dev/null || true
```

---

## Cara 1: Menggunakan Docker Compose Development (Recommended untuk Testing)

### Step 1: Start Semua Services

```bash
cd /home/fatur/Documents/Projects/Tau-TauRun

# Start database, backend, frontend
docker-compose up -d

# Lihat logs untuk memastikan semua berjalan
docker-compose logs -f
```

Output yang diharapkan:
```
✅ Database connected successfully
Server listening on port 8080
Ready on http://localhost:3000
```

Tekan `Ctrl+C` untuk keluar dari logs.

### Step 2: Tunggu Services Siap (30-60 detik)

```bash
# Cek status containers
docker-compose ps

# Harus menunjukkan semua "Up" dan "healthy"
```

### Step 3: Verifikasi Health Check

```bash
# Test backend
curl http://localhost:8080/health

# Harus return:
# {"success":true,"data":{"status":"healthy","version":"1.0.0"}}
```

### Step 4: Test Registration (Public Endpoint)

```bash
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Docker User",
    "email": "docker.test@example.com",
    "phone": "081234567890",
    "address": "Jakarta, Indonesia - Docker Test Street 123",
    "instagram_handle": "@dockertest"
  }'
```

Response yang diharapkan:
```json
{
  "success": true,
  "message": "Registration successful! Your payment status is pending.",
  "data": {
    "id": "uuid-here",
    "email": "docker.test@example.com",
    "registration_status": "PENDING",
    "payment_status": "UNPAID"
  }
}
```

### Step 5: Test Admin Login

```bash
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@tautaurun.com",
    "password": "Admin123!"
  }'
```

**Simpan token** dari response untuk digunakan selanjutnya.

### Step 6: Get Participants List

```bash
# Ganti YOUR_TOKEN_HERE dengan token dari step 5
TOKEN="YOUR_TOKEN_HERE"

curl -X GET http://localhost:8080/api/v1/admin/participants \
  -H "Authorization: Bearer $TOKEN"
```

### Step 7: Update Payment Status

```bash
# Ganti PARTICIPANT_ID dengan ID dari step 4
PARTICIPANT_ID="uuid-from-step-4"

curl -X PATCH http://localhost:8080/api/v1/admin/participants/$PARTICIPANT_ID/payment \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"payment_status": "PAID"}'
```

Response yang diharapkan:
```json
{
  "success": true,
  "message": "Payment status updated successfully",
  "data": {
    "id": "uuid-here",
    "payment_status": "PAID",
    "email_sent": true,
    "updated_at": "2026-01-01T..."
  }
}
```

### Step 8: Verifikasi di Database

```bash
# Masuk ke container database
docker exec -it tau-tau-run-db psql -U postgres -d tau_tau_run

# Lihat data participants
SELECT * FROM participants;

# Lihat email logs
SELECT * FROM email_logs ORDER BY sent_at DESC LIMIT 5;

# Keluar
\q
```

### Step 9: Test via Browser

1. Buka browser: **http://localhost:3000**
2. Isi form registrasi
3. Buka: **http://localhost:3000/admin/login**
4. Login dengan:
   - Email: `admin@tautaurun.com`
   - Password: `Admin123!`
5. Lihat dashboard dan coba update payment status

### Step 10: Stop Services

```bash
# Stop semua services (data tetap ada)
docker-compose stop

# Atau stop dan hapus data
docker-compose down -v
```

---

## Cara 2: Automated Testing Script

Buat file `test-docker.sh`:

```bash
#!/bin/bash
set -e

echo "🚀 E2E Test - Tau-Tau Run"
echo "=========================="

# 1. Start services
echo "1️⃣ Starting services..."
docker-compose up -d
sleep 30

# 2. Check health
echo "2️⃣ Checking health..."
curl -s http://localhost:8080/health | grep "healthy" && echo "✅ Backend OK"

# 3. Register participant
echo "3️⃣ Registering participant..."
REGISTER=$(curl -s -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Auto Test",
    "email": "auto.test@example.com",
    "phone": "081234567890",
    "address": "Jakarta, Indonesia - Auto Test Street 123"
  }')

PARTICIPANT_ID=$(echo $REGISTER | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "✅ Participant ID: $PARTICIPANT_ID"

# 4. Admin login
echo "4️⃣ Admin login..."
LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@tautaurun.com",
    "password": "Admin123!"
  }')

TOKEN=$(echo $LOGIN | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "✅ Token received"

# 5. Get participants
echo "5️⃣ Getting participants..."
curl -s -X GET http://localhost:8080/api/v1/admin/participants \
  -H "Authorization: Bearer $TOKEN" | grep "auto.test@example.com" && echo "✅ Participant found"

# 6. Update payment
echo "6️⃣ Updating payment status..."
curl -s -X PATCH http://localhost:8080/api/v1/admin/participants/$PARTICIPANT_ID/payment \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"payment_status": "PAID"}' | grep '"payment_status":"PAID"' && echo "✅ Payment updated"

echo ""
echo "=========================="
echo "✅ ALL TESTS PASSED!"
echo "=========================="
echo ""
echo "Services running at:"
echo "  Frontend: http://localhost:3000"
echo "  Backend:  http://localhost:8080"
echo "  Admin:    http://localhost:3000/admin/login"
echo ""
echo "To stop: docker-compose down"
```

Jalankan:

```bash
chmod +x test-docker.sh
./test-docker.sh
```

---

## Troubleshooting

### 1. Port sudah digunakan

```bash
# Cek port 8080
sudo lsof -i :8080

# Kill process
kill <PID>

# Atau ubah port di docker-compose.yml
```

### 2. Database tidak connect

```bash
# Lihat logs database
docker-compose logs db

# Restart database
docker-compose restart db

# Wait dan cek lagi
sleep 10
curl http://localhost:8080/health
```

### 3. Admin login gagal

```bash
# Pastikan admin user sudah ada
docker exec -it tau-tau-run-db psql -U postgres -d tau_tau_run -c "SELECT * FROM admins;"

# Jika kosong, insert manual:
docker exec -it tau-tau-run-db psql -U postgres -d tau_tau_run << 'SQL'
INSERT INTO admins (email, password_hash)
VALUES ('admin@tautaurun.com', 
        '$2a$12$B25zHJMMgXbVqxzIUzy9Vu0CdR8FphthNgoLkVfGkUS/v07oXyIf6')
ON CONFLICT (email) DO NOTHING;
SQL
```

### 4. Reset semua data

```bash
# Stop dan hapus semua
docker-compose down -v

# Hapus images juga
docker-compose down -v --rmi all

# Start fresh
docker-compose up -d --build
```

### 5. Lihat logs error

```bash
# Backend logs
docker-compose logs backend | tail -50

# Frontend logs  
docker-compose logs frontend | tail -50

# Database logs
docker-compose logs db | tail -50

# Semua logs
docker-compose logs -f
```

---

## Monitoring

### Lihat status containers

```bash
docker-compose ps
```

### Lihat resource usage

```bash
docker stats
```

### Exec command di container

```bash
# Masuk ke backend container
docker exec -it tau-tau-run-backend sh

# Masuk ke database
docker exec -it tau-tau-run-db psql -U postgres -d tau_tau_run

# Lihat logs frontend
docker logs tau-tau-run-frontend
```

---

## Clean Up

```bash
# Stop services (data tetap ada)
docker-compose stop

# Stop dan hapus containers (data volume tetap)
docker-compose down

# Hapus semua termasuk volumes (DATA HILANG!)
docker-compose down -v

# Hapus semua termasuk images
docker-compose down -v --rmi all
```

---

## Tips Testing

1. **Gunakan Postman/Insomnia** untuk test API lebih mudah
2. **Install jq** untuk format JSON: `sudo apt install jq`
   ```bash
   curl ... | jq .
   ```
3. **Simpan environment variables:**
   ```bash
   export API_URL="http://localhost:8080/api/v1"
   export TOKEN="your-token-here"
   ```
4. **Gunakan watch** untuk monitor:
   ```bash
   watch -n 2 'docker-compose ps'
   ```

---

## Next Steps

Setelah testing lokal berhasil:

1. ✅ Commit changes ke git
2. ✅ Push ke repository
3. ✅ Deploy ke staging server
4. ✅ Run security audit
5. ✅ Deploy ke production

---

**Happy Testing! 🚀**
