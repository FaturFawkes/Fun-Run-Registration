# Rollback & Disaster Recovery Plan

**System:** Tau-Tau Run Event Registration System  
**Version:** 1.0.0  
**Last Updated:** 2026-01-01

---

## Table of Contents

- [Quick Rollback Procedures](#quick-rollback-procedures)
- [Disaster Scenarios](#disaster-scenarios)
- [Recovery Procedures](#recovery-procedures)
- [Contact Information](#contact-information)
- [Post-Incident Review](#post-incident-review)

---

## Quick Rollback Procedures

### Application Rollback (Docker)

**Time Required:** 2-5 minutes

```bash
# 1. Stop current services
docker-compose -f docker-compose.prod.yml down

# 2. Pull previous image tag
docker pull your-registry/tautaurun-backend:previous-version
docker pull your-registry/tautaurun-frontend:previous-version

# 3. Update docker-compose to use previous version
# Edit docker-compose.prod.yml and change image tags

# 4. Start services
docker-compose -f docker-compose.prod.yml up -d

# 5. Verify health
curl http://localhost:8080/health
curl http://localhost:3000
```

### Database Rollback

**Time Required:** 5-15 minutes (depending on database size)

```bash
# 1. Stop application
sudo systemctl stop tautaurun-api
sudo systemctl stop tautaurun-frontend

# 2. Backup current state (just in case)
PGPASSWORD=prod_password pg_dump \
  -h localhost -U tautaurun -d tau_tau_run_prod \
  | gzip > emergency_backup_$(date +%Y%m%d_%H%M%S).sql.gz

# 3. Restore from known good backup
gunzip -c /opt/tautaurun/backups/tau_tau_run_prod_YYYYMMDD.sql.gz | \
PGPASSWORD=prod_password psql \
  -h localhost -U tautaurun -d tau_tau_run_prod

# 4. Verify restoration
PGPASSWORD=prod_password psql -h localhost -U tautaurun \
  -d tau_tau_run_prod \
  -c "SELECT COUNT(*) FROM participants;"

# 5. Restart application
sudo systemctl start tautaurun-api
sudo systemctl start tautaurun-frontend
```

---

## Disaster Scenarios

### Scenario 1: Application Crash

**Symptoms:**
- API returns 502/503 errors
- Health check fails
- No response from backend

**Recovery Steps:**

1. **Check application status**
   ```bash
   sudo systemctl status tautaurun-api
   sudo journalctl -u tautaurun-api -n 100
   ```

2. **Restart application**
   ```bash
   sudo systemctl restart tautaurun-api
   ```

3. **If restart fails, check logs**
   ```bash
   tail -f /var/log/tautaurun/api-error.log
   ```

4. **Common fixes:**
   - Database connection failed → Check PostgreSQL
   - Port already in use → Kill conflicting process
   - Permission denied → Check file ownership

**Time to Recovery:** 2-5 minutes

---

### Scenario 2: Database Corruption

**Symptoms:**
- Database queries failing
- Data inconsistencies
- Application errors with database messages

**Recovery Steps:**

1. **Stop application immediately**
   ```bash
   sudo systemctl stop tautaurun-api
   sudo systemctl stop tautaurun-frontend
   ```

2. **Check database integrity**
   ```bash
   PGPASSWORD=password psql -h localhost -U tautaurun \
     -d tau_tau_run_prod \
     -c "SELECT pg_database_size('tau_tau_run_prod');"
   ```

3. **Restore from latest backup**
   - Follow database rollback procedure above

4. **Verify data integrity**
   ```bash
   PGPASSWORD=password psql -h localhost -U tautaurun \
     -d tau_tau_run_prod << EOF
   SELECT COUNT(*) FROM participants;
   SELECT COUNT(*) FROM admins;
   SELECT COUNT(*) FROM email_logs;
   EOF
   ```

5. **Restart application**

**Time to Recovery:** 10-30 minutes

---

### Scenario 3: Complete Server Failure

**Symptoms:**
- Server unreachable
- No SSH access
- All services down

**Recovery Steps:**

1. **Provision new server**
   - Use same OS and configuration
   - Follow deployment guide

2. **Restore from backups**
   - Database backup (most critical)
   - Configuration files
   - SSL certificates

3. **Update DNS** (if IP changed)
   ```bash
   # Point domain to new IP
   # Wait for propagation (up to 24 hours)
   ```

4. **Verify all services**
   - Database connectivity
   - API health check
   - Frontend accessibility
   - Email sending

**Time to Recovery:** 2-4 hours

---

### Scenario 4: Data Breach

**Symptoms:**
- Suspicious database activity
- Unauthorized admin access
- Data exfiltration detected

**Immediate Actions:**

1. **Isolate the system**
   ```bash
   # Block all incoming traffic
   sudo ufw default deny incoming
   
   # Stop application
   sudo systemctl stop tautaurun-api
   sudo systemctl stop tautaurun-frontend
   ```

2. **Preserve evidence**
   ```bash
   # Copy logs
   cp -r /var/log/tautaurun /incident-$(date +%Y%m%d)
   
   # Database snapshot
   PGPASSWORD=password pg_dump -h localhost -U tautaurun \
     -d tau_tau_run_prod > incident_db_$(date +%Y%m%d).sql
   ```

3. **Contact team**
   - Notify security team
   - Notify legal (if required)
   - Notify users (if data compromised)

4. **Investigation**
   - Review access logs
   - Check for unauthorized changes
   - Identify entry point

5. **Remediation**
   - Change all passwords
   - Rotate JWT secrets
   - Revoke all active sessions
   - Patch vulnerabilities

**Time to Recovery:** 1-7 days (depending on severity)

---

### Scenario 5: SMTP Service Down

**Symptoms:**
- Emails not being delivered
- SMTP connection errors in logs
- Email logs show all FAILED

**Recovery Steps:**

1. **Check SMTP service**
   ```bash
   # Test SMTP connectivity
   telnet smtp.gmail.com 587
   ```

2. **Verify credentials**
   - Check SMTP_USERNAME and SMTP_PASSWORD
   - Ensure account not locked

3. **Switch SMTP provider** (if needed)
   - Update .env with backup SMTP
   - Restart backend

4. **Resend failed emails**
   ```sql
   -- Get participants who didn't receive email
   SELECT p.* FROM participants p
   LEFT JOIN email_logs e ON p.id = e.participant_id
   WHERE p.payment_status = 'PAID' 
     AND (e.status = 'FAILED' OR e.status IS NULL);
   ```

**Time to Recovery:** 15-30 minutes

---

## Recovery Time Objectives (RTO)

| Scenario | RTO | RPO |
|----------|-----|-----|
| Application Crash | 5 minutes | 0 |
| Database Corruption | 30 minutes | 24 hours |
| Server Failure | 4 hours | 24 hours |
| Data Breach | Variable | Variable |
| SMTP Down | 30 minutes | N/A |

**RTO:** Recovery Time Objective (how quickly to restore)  
**RPO:** Recovery Point Objective (acceptable data loss)

---

## Backup Strategy

### Automated Daily Backups

```bash
# Location
/opt/tautaurun/backups/

# Retention
- Daily: 7 days
- Weekly: 4 weeks
- Monthly: 12 months

# Verification
# Test restore monthly
```

### Offsite Backups

```bash
# S3 backup (if configured)
aws s3 ls s3://tautaurun-backups/

# Verify latest backup
aws s3 ls s3://tautaurun-backups/ --recursive | tail -1
```

---

## Pre-Deployment Checklist

Before ANY deployment:

- [ ] Backup current database
- [ ] Test rollback procedure in staging
- [ ] Verify health checks working
- [ ] Document changes
- [ ] Notify team of deployment window
- [ ] Have rollback plan ready

---

## Post-Deployment Verification

After deployment:

- [ ] Check health endpoint: `curl https://api.tautaurun.com/health`
- [ ] Test registration flow
- [ ] Test admin login
- [ ] Test payment update
- [ ] Verify email sending
- [ ] Check logs for errors
- [ ] Monitor for 30 minutes

---

## Contact Information

### Emergency Contacts

```
Primary Admin: ____________________
Phone: ____________________
Email: ____________________

Backup Admin: ____________________
Phone: ____________________
Email: ____________________

Database Admin: ____________________
Phone: ____________________
Email: ____________________
```

### Service Providers

```
Hosting: ____________________
Support: ____________________

SMTP Provider: ____________________
Support: ____________________

Domain Registrar: ____________________
Support: ____________________
```

---

## Post-Incident Review

After any incident:

1. **Document what happened**
   - Timeline of events
   - Actions taken
   - Time to resolution

2. **Identify root cause**
   - What failed?
   - Why did it fail?
   - How was it detected?

3. **Implement improvements**
   - Prevent recurrence
   - Improve detection
   - Update procedures

4. **Update this document**
   - Add new scenarios
   - Refine procedures
   - Update contact info

---

## Testing Schedule

- [ ] Monthly: Test database restore
- [ ] Quarterly: Full disaster recovery drill
- [ ] Annually: Review and update plan

**Last Tested:** _________________  
**Next Test:** _________________

---

## Document History

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2026-01-01 | 1.0.0 | Initial version | System |

---

**Remember:** The best disaster recovery plan is the one you've tested!
