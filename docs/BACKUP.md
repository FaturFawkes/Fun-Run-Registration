# Database Backup and Restore Procedures

## Table of Contents
- [Automated Backups](#automated-backups)
- [Manual Backup](#manual-backup)
- [Restore from Backup](#restore-from-backup)
- [Backup Verification](#backup-verification)
- [Best Practices](#best-practices)

---

## Automated Backups

### Setup Daily Automated Backups

1. **Create backup script**

```bash
#!/bin/bash
# File: /opt/tautaurun/scripts/backup-db.sh

BACKUP_DIR="/opt/tautaurun/backups"
DB_NAME="tau_tau_run"
DB_USER="postgres"
DB_PASSWORD="your_password_here"
DATE=$(date +%Y%m%d_%H%M%S)
FILENAME="${DB_NAME}_${DATE}.sql.gz"
RETENTION_DAYS=30

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Create backup
PGPASSWORD=$DB_PASSWORD pg_dump \
  -h localhost \
  -U $DB_USER \
  -d $DB_NAME \
  --no-owner \
  --no-acl \
  | gzip > "$BACKUP_DIR/$FILENAME"

# Verify backup was created
if [ -f "$BACKUP_DIR/$FILENAME" ]; then
    SIZE=$(du -h "$BACKUP_DIR/$FILENAME" | cut -f1)
    echo "[$(date)] Backup successful: $FILENAME ($SIZE)"
    
    # Delete backups older than retention period
    find $BACKUP_DIR -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete
    echo "[$(date)] Old backups cleaned up (retention: $RETENTION_DAYS days)"
else
    echo "[$(date)] ERROR: Backup failed!" >&2
    exit 1
fi
```

2. **Make script executable**

```bash
chmod +x /opt/tautaurun/scripts/backup-db.sh
```

3. **Add to crontab**

```bash
# Edit crontab
crontab -e

# Add this line (daily backup at 2 AM)
0 2 * * * /opt/tautaurun/scripts/backup-db.sh >> /var/log/tautaurun/backup.log 2>&1
```

4. **Test the backup script**

```bash
/opt/tautaurun/scripts/backup-db.sh
```

---

## Manual Backup

### Full Database Backup

```bash
# Create backup with timestamp
PGPASSWORD=your_password pg_dump \
  -h localhost \
  -U postgres \
  -d tau_tau_run \
  --no-owner \
  --no-acl \
  | gzip > tau_tau_run_$(date +%Y%m%d_%H%M%S).sql.gz
```

### Backup Specific Tables

```bash
# Backup only participants table
PGPASSWORD=your_password pg_dump \
  -h localhost \
  -U postgres \
  -d tau_tau_run \
  -t participants \
  | gzip > participants_$(date +%Y%m%d).sql.gz

# Backup multiple tables
PGPASSWORD=your_password pg_dump \
  -h localhost \
  -U postgres \
  -d tau_tau_run \
  -t participants \
  -t email_logs \
  | gzip > core_tables_$(date +%Y%m%d).sql.gz
```

### Schema-Only Backup

```bash
# Backup database structure without data
PGPASSWORD=your_password pg_dump \
  -h localhost \
  -U postgres \
  -d tau_tau_run \
  --schema-only \
  > tau_tau_run_schema.sql
```

---

## Restore from Backup

### Restore Full Database

```bash
# Decompress and restore
gunzip -c tau_tau_run_backup.sql.gz | \
PGPASSWORD=your_password psql \
  -h localhost \
  -U postgres \
  -d tau_tau_run

# Or in one step
PGPASSWORD=your_password pg_restore \
  -h localhost \
  -U postgres \
  -d tau_tau_run \
  -v \
  tau_tau_run_backup.dump
```

### Restore to New Database

```bash
# Create new database
PGPASSWORD=your_password psql -h localhost -U postgres -c "CREATE DATABASE tau_tau_run_restore;"

# Restore backup
gunzip -c tau_tau_run_backup.sql.gz | \
PGPASSWORD=your_password psql \
  -h localhost \
  -U postgres \
  -d tau_tau_run_restore

# Verify data
PGPASSWORD=your_password psql -h localhost -U postgres -d tau_tau_run_restore -c "SELECT COUNT(*) FROM participants;"
```

### Restore Specific Tables

```bash
# Restore only participants table
gunzip -c participants_backup.sql.gz | \
PGPASSWORD=your_password psql \
  -h localhost \
  -U postgres \
  -d tau_tau_run
```

---

## Backup Verification

### Verify Backup File

```bash
# Check if file is a valid gzip
gunzip -t tau_tau_run_backup.sql.gz

# View first few lines
gunzip -c tau_tau_run_backup.sql.gz | head -20

# Check backup size
du -h tau_tau_run_backup.sql.gz
```

### Verify Database State

```bash
# Count records
PGPASSWORD=your_password psql -h localhost -U postgres -d tau_tau_run << EOF
SELECT 
  'participants' as table_name, 
  COUNT(*) as count 
FROM participants
UNION ALL
SELECT 'admins', COUNT(*) FROM admins
UNION ALL
SELECT 'email_logs', COUNT(*) FROM email_logs;
EOF
```

### Test Restore (Dry Run)

```bash
# Create temporary database
PGPASSWORD=your_password psql -h localhost -U postgres -c "CREATE DATABASE tau_tau_run_test;"

# Restore to test database
gunzip -c tau_tau_run_backup.sql.gz | \
PGPASSWORD=your_password psql \
  -h localhost \
  -U postgres \
  -d tau_tau_run_test

# Verify
PGPASSWORD=your_password psql -h localhost -U postgres -d tau_tau_run_test -c "SELECT COUNT(*) FROM participants;"

# Clean up
PGPASSWORD=your_password psql -h localhost -U postgres -c "DROP DATABASE tau_tau_run_test;"
```

---

## Best Practices

### Backup Strategy

1. **3-2-1 Rule**
   - 3 copies of data
   - 2 different storage types
   - 1 offsite backup

2. **Retention Policy**
   - Daily backups: Keep 7 days
   - Weekly backups: Keep 4 weeks
   - Monthly backups: Keep 12 months

3. **Backup Schedule**
   - Daily: 2 AM (low traffic)
   - Weekly: Sunday 1 AM
   - Monthly: 1st of month

### Offsite Backup

```bash
#!/bin/bash
# Upload backup to S3 (requires AWS CLI)

BACKUP_FILE="tau_tau_run_$(date +%Y%m%d).sql.gz"
S3_BUCKET="s3://your-backup-bucket/tautaurun/"

# Create backup
PGPASSWORD=your_password pg_dump \
  -h localhost \
  -U postgres \
  -d tau_tau_run \
  | gzip > /tmp/$BACKUP_FILE

# Upload to S3
aws s3 cp /tmp/$BACKUP_FILE $S3_BUCKET

# Verify upload
aws s3 ls $S3_BUCKET$BACKUP_FILE

# Clean up local file
rm /tmp/$BACKUP_FILE
```

### Encryption

```bash
# Backup with encryption
PGPASSWORD=your_password pg_dump \
  -h localhost \
  -U postgres \
  -d tau_tau_run \
  | gzip \
  | openssl enc -aes-256-cbc -salt -out tau_tau_run_encrypted.sql.gz.enc

# Decrypt and restore
openssl enc -aes-256-cbc -d -in tau_tau_run_encrypted.sql.gz.enc \
  | gunzip \
  | PGPASSWORD=your_password psql -h localhost -U postgres -d tau_tau_run
```

### Monitoring

```bash
# Check last backup
ls -lht /opt/tautaurun/backups/ | head -5

# Verify backup log
tail -f /var/log/tautaurun/backup.log

# Alert if backup fails
#!/bin/bash
if ! /opt/tautaurun/scripts/backup-db.sh; then
    echo "Database backup failed!" | mail -s "URGENT: Backup Failure" admin@tautaurun.com
fi
```

---

## Emergency Recovery

### Quick Recovery Steps

1. **Stop application**
   ```bash
   sudo systemctl stop tautaurun-api
   sudo systemctl stop tautaurun-frontend
   ```

2. **Backup current state** (just in case)
   ```bash
   PGPASSWORD=your_password pg_dump -h localhost -U postgres -d tau_tau_run \
     | gzip > emergency_backup_$(date +%Y%m%d_%H%M%S).sql.gz
   ```

3. **Restore from backup**
   ```bash
   gunzip -c tau_tau_run_good_backup.sql.gz | \
   PGPASSWORD=your_password psql -h localhost -U postgres -d tau_tau_run
   ```

4. **Verify data**
   ```bash
   PGPASSWORD=your_password psql -h localhost -U postgres -d tau_tau_run \
     -c "SELECT COUNT(*) FROM participants;"
   ```

5. **Restart application**
   ```bash
   sudo systemctl start tautaurun-api
   sudo systemctl start tautaurun-frontend
   ```

---

## Backup Checklist

- [ ] Automated daily backups configured
- [ ] Backup retention policy set
- [ ] Backups stored in multiple locations
- [ ] Offsite backup configured
- [ ] Restore procedure tested monthly
- [ ] Backup monitoring in place
- [ ] Emergency recovery plan documented
- [ ] Team trained on restore procedure

---

**Remember:** A backup is only good if you can restore from it. Test your backups regularly!
