# Deployment Scripts

This directory contains scripts for deploying Noddit to Google Cloud Platform.

## Scripts Overview

### 🚀 `gcp-setup.sh`
**Run from:** Local machine
**Purpose:** Initial GCP VM setup and configuration

Creates and configures a free-tier e2-micro VM with Docker, firewall rules, and SSH keys for GitHub Actions.

**Usage:**
```bash
bash scripts/gcp-setup.sh
```

**What it does:**
- Creates e2-micro VM in GCP
- Configures firewall rules (ports 80, 443, 8080, 8081)
- Installs Docker and Docker Compose
- Creates deploy user
- Generates SSH keys for GitHub Actions
- Outputs configuration instructions

---

### 💾 `vm-backup-setup.sh`
**Run from:** GCP VM
**Purpose:** Setup automated daily database backups

**Usage:**
```bash
# SSH into VM first
gcloud compute ssh noddit-vm --zone=us-central1-a
sudo su - deploy
bash ~/noddit/scripts/vm-backup-setup.sh
```

**What it does:**
- Creates backup directory
- Generates backup script
- Configures daily cron job (2 AM)
- Retention: 7 days

**Manual backup:**
```bash
~/backup-db.sh
```

---

### 🔄 `vm-systemd-setup.sh`
**Run from:** GCP VM
**Purpose:** Setup systemd service for auto-start on boot

**Usage:**
```bash
# SSH into VM first
gcloud compute ssh noddit-vm --zone=us-central1-a
sudo su - deploy
bash ~/noddit/scripts/vm-systemd-setup.sh
```

**What it does:**
- Creates systemd service file
- Enables auto-start on boot
- Allows service management via systemctl

**Service commands:**
```bash
sudo systemctl start noddit      # Start
sudo systemctl stop noddit       # Stop
sudo systemctl restart noddit    # Restart
sudo systemctl status noddit     # Status
```

---

### 🌐 `vm-nginx-setup.sh`
**Run from:** GCP VM
**Purpose:** Optional Nginx reverse proxy for port 80 access

**Usage:**
```bash
# SSH into VM first
gcloud compute ssh noddit-vm --zone=us-central1-a
bash ~/noddit/scripts/vm-nginx-setup.sh
```

**What it does:**
- Installs Nginx
- Configures reverse proxy
- Frontend on `/` → `localhost:8081`
- Backend on `/api` → `localhost:8080`
- Enables service

**After setup:**
1. Update `.env` file:
   ```bash
   NEXT_PUBLIC_API_URL=http://YOUR_IP/api
   FRONTEND_URL=http://YOUR_IP
   ```
2. Rebuild frontend:
   ```bash
   cd ~/noddit
   docker-compose down
   docker-compose up -d --build frontend
   ```

---

## Deployment Workflow

### Initial Setup (One-time)

1. **Create GCP VM** (from local machine):
   ```bash
   bash scripts/gcp-setup.sh
   ```

2. **Configure GitHub Secrets**:
   - Go to GitHub repo → Settings → Secrets
   - Add: `SSH_PRIVATE_KEY`, `VM_HOST`, `VM_USER`, `VM_PORT`

3. **Create `.env` on VM**:
   ```bash
   gcloud compute ssh noddit-vm --zone=us-central1-a
   sudo su - deploy
   cd ~/noddit
   nano .env
   # Fill in DB_PASSWORD, JWT_SECRET, CLERK keys, EXTERNAL_IP
   ```

4. **Initial deployment**:
   ```bash
   # Still on VM
   cp docker-compose.production.yml docker-compose.yml
   docker-compose build
   docker-compose up -d
   ```

5. **Setup auto-start**:
   ```bash
   bash ~/noddit/scripts/vm-systemd-setup.sh
   ```

6. **Setup backups**:
   ```bash
   bash ~/noddit/scripts/vm-backup-setup.sh
   ```

7. **(Optional) Setup Nginx**:
   ```bash
   bash ~/noddit/scripts/vm-nginx-setup.sh
   ```

### Continuous Deployment

After initial setup, deployments are automatic via GitHub Actions:

```bash
# Make changes
git add .
git commit -m "feat: new feature"
git push origin deployment  # or main
```

GitHub Actions will:
- Sync files to VM
- Rebuild containers
- Run health checks

---

## Troubleshooting

### Script permissions error
```bash
chmod +x scripts/*.sh
```

### gcloud command not found
Install gcloud CLI: https://cloud.google.com/sdk/docs/install

### SSH connection refused
Wait 30 seconds after VM creation, then try again.

### Docker permission denied
```bash
# On VM
sudo usermod -aG docker $USER
# Logout and login again
```

### Backup not running
```bash
# Check cron
crontab -l

# View logs
tail -f ~/backups/backup.log

# Test manually
~/backup-db.sh
```

---

## Additional Resources

- Full deployment guide: [`DEPLOYMENT.md`](../docs/DEPLOYMENT.md)
- Quick start guide: [`DEPLOYMENT-QUICKSTART.md`](../docs/DEPLOYMENT-QUICKSTART.md)
- Deployment checklist: [`DEPLOYMENT-CHECKLIST.md`](../docs/DEPLOYMENT-CHECKLIST.md)
- GCP Always Free Tier: https://cloud.google.com/free
- Docker Compose docs: https://docs.docker.com/compose/
