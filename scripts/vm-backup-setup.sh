#!/bin/bash
# Database Backup Setup Script for GCP VM
# Run this ON THE VM after initial deployment

set -e

echo "Setting up automated database backups..."

# Create backup directory
mkdir -p ~/backups

# Create backup script
cat > ~/backup-db.sh << 'EOF'
#!/bin/bash
# Automated database backup script

BACKUP_DIR=~/backups
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/noddit_$DATE.sql.gz"

echo "Starting database backup at $(date)"

# Create backup
docker exec noddit-postgres pg_dump -U postgres userdb | gzip > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "Backup successful: $BACKUP_FILE"

    # Delete backups older than 7 days
    find $BACKUP_DIR -name "noddit_*.sql.gz" -mtime +7 -delete
    echo "Old backups cleaned up"
else
    echo "Backup failed!"
    exit 1
fi

echo "Backup completed at $(date)"
EOF

chmod +x ~/backup-db.sh

# Setup daily cron job at 2 AM
(crontab -l 2>/dev/null || true; echo "0 2 * * * ~/backup-db.sh >> ~/backups/backup.log 2>&1") | crontab -

echo "✅ Backup setup complete!"
echo ""
echo "Backup schedule: Daily at 2:00 AM"
echo "Backup retention: 7 days"
echo "Backup location: ~/backups/"
echo ""
echo "To run a manual backup now:"
echo "  ~/backup-db.sh"
echo ""
echo "To view backup logs:"
echo "  tail -f ~/backups/backup.log"
echo ""
echo "To list backups:"
echo "  ls -lh ~/backups/"
echo ""
echo "To restore from backup:"
echo "  gunzip -c ~/backups/noddit_YYYYMMDD_HHMMSS.sql.gz | docker exec -i noddit-postgres psql -U postgres -d userdb"
