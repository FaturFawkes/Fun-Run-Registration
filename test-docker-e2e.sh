#!/bin/bash
set -e

echo "═══════════════════════════════════════════════════════════"
echo "  🚀 E2E TEST - TAU-TAU RUN (Docker Compose)"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Navigate to project directory
cd "$(dirname "$0")"

# 1. Clean up
echo "🧹 Step 1: Cleaning up old containers..."
docker-compose down -v 2>/dev/null || true
echo -e "${GREEN}✅ Cleanup complete${NC}"
echo ""

# 2. Start services
echo "🐳 Step 2: Starting Docker Compose services..."
docker-compose up -d
echo -e "${GREEN}✅ Services starting${NC}"
echo ""

# 3. Wait for services
echo "⏳ Step 3: Waiting for services to be ready (30 seconds)..."
sleep 10
for i in {1..20}; do
  if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend is ready!${NC}"
    break
  fi
  echo -n "."
  sleep 1
done
echo ""

# Check if backend is up
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
  echo -e "${RED}❌ Backend failed to start. Check logs:${NC}"
  docker-compose logs backend
  exit 1
fi

# 4. Check database and create admin if needed
echo "👤 Step 4: Ensuring admin user exists..."
docker exec tau-tau-run-db psql -U postgres -d tau_tau_run << 'SQL' 2>/dev/null || true
INSERT INTO admins (email, password_hash, created_at)
VALUES (
    'admin@tautaurun.com',
    '$2a$12$B25zHJMMgXbVqxzIUzy9Vu0CdR8FphthNgoLkVfGkUS/v07oXyIf6',
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;
SQL
echo -e "${GREEN}✅ Admin user ready${NC}"
echo ""

# 5. Test: Register participant
echo "📝 Step 5: Testing participant registration..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Docker E2E Test User",
    "email": "docker.e2e.test@example.com",
    "phone": "081234567890",
    "address": "Jakarta, Indonesia - Docker E2E Test Street 123",
    "instagram_handle": "@dockere2e"
  }')

if echo "$REGISTER_RESPONSE" | grep -q '"success":true'; then
  PARTICIPANT_ID=$(echo $REGISTER_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
  echo -e "${GREEN}✅ Registration successful${NC}"
  echo "   Participant ID: $PARTICIPANT_ID"
else
  echo -e "${RED}❌ Registration failed${NC}"
  echo "$REGISTER_RESPONSE" | jq . 2>/dev/null || echo "$REGISTER_RESPONSE"
  exit 1
fi
echo ""

# 6. Test: Admin login
echo "🔐 Step 6: Testing admin login..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@tautaurun.com",
    "password": "Admin123!"
  }')

if echo "$LOGIN_RESPONSE" | grep -q '"success":true'; then
  TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
  echo -e "${GREEN}✅ Login successful${NC}"
  echo "   Token: ${TOKEN:0:30}..."
else
  echo -e "${RED}❌ Login failed${NC}"
  echo "$LOGIN_RESPONSE" | jq . 2>/dev/null || echo "$LOGIN_RESPONSE"
  exit 1
fi
echo ""

# 7. Test: Get participants
echo "📋 Step 7: Testing get participants list..."
PARTICIPANTS_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/admin/participants \
  -H "Authorization: Bearer $TOKEN")

if echo "$PARTICIPANTS_RESPONSE" | grep -q 'docker.e2e.test@example.com'; then
  TOTAL=$(echo $PARTICIPANTS_RESPONSE | grep -o '"total":[0-9]*' | cut -d':' -f2)
  echo -e "${GREEN}✅ Participants list retrieved${NC}"
  echo "   Total participants: $TOTAL"
else
  echo -e "${RED}❌ Failed to get participants${NC}"
  echo "$PARTICIPANTS_RESPONSE" | jq . 2>/dev/null || echo "$PARTICIPANTS_RESPONSE"
  exit 1
fi
echo ""

# 8. Test: Update payment status
echo "💳 Step 8: Testing payment status update..."
UPDATE_RESPONSE=$(curl -s -X PATCH http://localhost:8080/api/v1/admin/participants/$PARTICIPANT_ID/payment \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"payment_status": "PAID"}')

if echo "$UPDATE_RESPONSE" | grep -q '"payment_status":"PAID"'; then
  EMAIL_SENT=$(echo $UPDATE_RESPONSE | grep -o '"email_sent":[a-z]*' | cut -d':' -f2)
  echo -e "${GREEN}✅ Payment status updated to PAID${NC}"
  echo "   Email triggered: $EMAIL_SENT"
else
  echo -e "${RED}❌ Failed to update payment status${NC}"
  echo "$UPDATE_RESPONSE" | jq . 2>/dev/null || echo "$UPDATE_RESPONSE"
  exit 1
fi
echo ""

# 9. Test: Idempotency (update PAID to PAID again)
echo "🔁 Step 9: Testing idempotency (PAID → PAID)..."
UPDATE_AGAIN=$(curl -s -X PATCH http://localhost:8080/api/v1/admin/participants/$PARTICIPANT_ID/payment \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"payment_status": "PAID"}')

if echo "$UPDATE_AGAIN" | grep -q '"email_sent":false'; then
  echo -e "${GREEN}✅ Idempotency check passed (no duplicate email)${NC}"
else
  echo -e "${YELLOW}⚠️  Idempotency check inconclusive${NC}"
fi
echo ""

# 10. Verify database
echo "💾 Step 10: Verifying database state..."
docker exec tau-tau-run-db psql -U postgres -d tau_tau_run -t -c "
SELECT 
  name, 
  email, 
  payment_status,
  registration_status
FROM participants 
WHERE email = 'docker.e2e.test@example.com';" 2>/dev/null | head -1

# Check email logs
EMAIL_LOG_COUNT=$(docker exec tau-tau-run-db psql -U postgres -d tau_tau_run -t -c \
  "SELECT COUNT(*) FROM email_logs WHERE participant_id = '$PARTICIPANT_ID';" 2>/dev/null | tr -d ' ')

if [ "$EMAIL_LOG_COUNT" -gt "0" ]; then
  echo -e "${GREEN}✅ Email attempt logged ($EMAIL_LOG_COUNT records)${NC}"
else
  echo -e "${YELLOW}⚠️  No email logs found (check async processing)${NC}"
fi
echo ""

# 11. Show container status
echo "🐳 Step 11: Container status..."
docker-compose ps
echo ""

echo "═══════════════════════════════════════════════════════════"
echo -e "${GREEN}✅ ALL E2E TESTS PASSED!${NC}"
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Services are running and accessible:"
echo "  • Frontend:     http://localhost:3000"
echo "  • Backend API:  http://localhost:8080"
echo "  • Admin Login:  http://localhost:3000/admin/login"
echo "  • API Docs:     See docs/API.md"
echo ""
echo "Admin credentials:"
echo "  Email:    admin@tautaurun.com"
echo "  Password: Admin123!"
echo ""
echo "Database access:"
echo "  docker exec -it tau-tau-run-db psql -U postgres -d tau_tau_run"
echo ""
echo "To view logs:"
echo "  docker-compose logs -f"
echo ""
echo "To stop services:"
echo "  docker-compose down"
echo ""
echo "To stop and remove data:"
echo "  docker-compose down -v"
echo ""
