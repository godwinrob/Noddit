# Noddit GCP Deployment Guide

This guide walks you through deploying Noddit to Google Cloud Platform using the **free tier e2-micro VM** for $0/month.

## Prerequisites

- Google Cloud Platform account with billing enabled (free tier requires card on file)
- `gcloud` CLI installed ([install guide](https://cloud.google.com/sdk/docs/install))
- GitHub repository for the project
- Domain name (optional, for SSL)

## Architecture

```
Internet → e2-micro VM (us-central1)
           └─ Docker Compose
              ├─ PostgreSQL (128MB limit)
              ├─ Go Backend (256MB limit)
              └─ Next.js Frontend (512MB limit)
```

**Total Cost:** $0/month (GCP Always Free Tier)

## Deployment Steps

### 1. Initial VM Setup

Run the automated setup script from your local machine:

```bash
cd C:\Users\rob\Documents\workspace\noddit
bash scripts/gcp-setup.sh
```

This script will:
- Create e2-micro VM in GCP
- Configure firewall rules
- Install Docker and Docker Compose
- Create deploy user
- Generate SSH keys for GitHub Actions
- Display setup instructions

**Save the SSH private key output** - you'll need it for GitHub Actions.

### 2. Configure GitHub Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions → New repository secret

Add these secrets:
- `SSH_PRIVATE_KEY`: The private key from setup script output
- `VM_HOST`: Your VM's external IP address
- `VM_USER`: `deploy`
- `VM_PORT`: `22`

### 3. Create Production Environment File

SSH into your VM:

```bash
gcloud compute ssh noddit-vm --zone=us-central1-a
sudo su - deploy
cd ~/noddit
```

Create `.env` file with production values:

```bash
nano .env
```

Paste this template and fill in the values:

```bash
# Database
DB_PASSWORD=YOUR_STRONG_PASSWORD_HERE

# JWT Authentication
# Generate with: openssl rand -base64 32
JWT_SECRET=YOUR_BASE64_SECRET_HERE

# Clerk (if using Clerk authentication)
CLERK_SECRET_KEY=sk_live_YOUR_KEY
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_live_YOUR_KEY

# VM External IP
EXTERNAL_IP=YOUR_VM_IP_HERE
```

**Generate secure secrets:**
```bash
# Strong password
openssl rand -base64 24

# JWT secret
openssl rand -base64 32
```

Save and exit (`Ctrl+X`, then `Y`, then `Enter`).

### 4. Initial Manual Deployment

While still on the VM:

```bash
# Copy production docker-compose
cp docker-compose.production.yml docker-compose.yml

# Build and start containers
docker-compose build
docker-compose up -d

# Check status
docker-compose ps
docker-compose logs -f
```

Press `Ctrl+C` to stop following logs. Exit VM:

```bash
exit  # Exit deploy user
exit  # Exit VM
```

### 5. Verify Deployment

From your local machine:

```bash
# Get VM IP
export VM_IP=$(gcloud compute instances describe noddit-vm \
  --zone=us-central1-a \
  --format="get(networkInterfaces[0].accessConfigs[0].natIP)")

# Test backend
curl http://$VM_IP:8080/api/public/subnoddits

# Test frontend (in browser)
echo "Frontend: http://$VM_IP:8081"
```

### 6. Setup Auto-Start on Boot

SSH back into VM and run:

```bash
gcloud compute ssh noddit-vm --zone=us-central1-a
sudo su - deploy
cd ~/noddit
bash ~/noddit/scripts/vm-systemd-setup.sh
```

This ensures Docker containers start automatically if the VM reboots.

### 7. Setup Automated Backups

While on the VM:

```bash
bash ~/noddit/scripts/vm-backup-setup.sh
```

This creates daily database backups at 2 AM, keeping the last 7 days.

### 8. (Optional) Setup Nginx Reverse Proxy

To serve the app on port 80 instead of 8081:

```bash
bash ~/noddit/scripts/vm-nginx-setup.sh
```

Then update `.env` and rebuild:

```bash
nano .env
# Change:
# NEXT_PUBLIC_API_URL=http://YOUR_IP/api
# FRONTEND_URL=http://YOUR_IP

docker-compose down
docker-compose up -d --build frontend
```

Now accessible at `http://YOUR_IP/` (port 80).

## Continuous Deployment with GitHub Actions

Once GitHub secrets are configured, deployments are automatic:

1. Make changes to your code
2. Commit and push to `main` or `deployment` branch:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   git push origin deployment
   ```
3. GitHub Actions automatically:
   - Syncs files to VM
   - Rebuilds containers
   - Restarts services
   - Runs health checks

**Monitor deployment:** Go to GitHub repository → Actions tab

**Manual deployment:** GitHub Actions tab → Deploy to GCP e2-micro → Run workflow

## Post-Deployment Tasks

### Security Hardening

```bash
# SSH into VM
gcloud compute ssh noddit-vm --zone=us-central1-a

# Setup UFW firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp
sudo ufw allow 8081/tcp
sudo ufw enable

# Auto security updates
sudo apt install unattended-upgrades -y
sudo dpkg-reconfigure -plow unattended-upgrades

# Disable password SSH (key-only)
sudo sed -i 's/PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo systemctl restart sshd
```

### SSL Certificate (if you have a domain)

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain certificate (replace with your domain)
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Auto-renewal is configured automatically
```

## Monitoring & Maintenance

### Check Application Status

```bash
# From local machine
gcloud compute ssh noddit-vm --zone=us-central1-a --command="docker-compose -f /home/deploy/noddit/docker-compose.yml ps"

# Check logs
gcloud compute ssh noddit-vm --zone=us-central1-a --command="docker-compose -f /home/deploy/noddit/docker-compose.yml logs --tail=100"
```

### Monitor Resources

```bash
# SSH into VM
gcloud compute ssh noddit-vm --zone=us-central1-a

# Memory usage
free -h

# Disk usage
df -h

# Container stats
docker stats --no-stream

# System load
htop  # or: top
```

### View Logs

```bash
# All containers
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres

# Last 100 lines
docker-compose logs --tail=100
```

### Restart Services

```bash
# Restart all
sudo systemctl restart noddit

# Or manually
cd /home/deploy/noddit
docker-compose restart

# Restart specific service
docker-compose restart backend
```

### Update Application

```bash
# Option 1: Push to GitHub (automated)
git push origin deployment

# Option 2: Manual update
gcloud compute ssh noddit-vm --zone=us-central1-a
sudo su - deploy
cd ~/noddit
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Backup & Restore

### Manual Backup

```bash
# SSH into VM
gcloud compute ssh noddit-vm --zone=us-central1-a
sudo su - deploy

# Run backup
~/backup-db.sh

# List backups
ls -lh ~/backups/
```

### Restore from Backup

```bash
# SSH into VM
gcloud compute ssh noddit-vm --zone=us-central1-a
sudo su - deploy
cd ~/backups

# Restore (replace with actual backup filename)
gunzip -c noddit_20260319_020000.sql.gz | docker exec -i noddit-postgres psql -U postgres -d userdb
```

### Download Backup to Local Machine

```bash
# From local machine
gcloud compute scp \
  --zone=us-central1-a \
  deploy@noddit-vm:~/backups/noddit_20260319_020000.sql.gz \
  ./local-backup.sql.gz
```

## Troubleshooting

### Containers Not Starting

```bash
# Check container status
docker-compose ps

# View logs
docker-compose logs

# Check available memory
free -h

# Restart
docker-compose down
docker-compose up -d
```

### Out of Memory

```bash
# Check memory usage
docker stats --no-stream

# Reduce memory limits in docker-compose.yml
# Restart containers
docker-compose down
docker-compose up -d
```

### Database Connection Issues

```bash
# Check postgres is running
docker-compose ps postgres

# Test connection
docker exec -it noddit-postgres psql -U postgres -d userdb

# View postgres logs
docker-compose logs postgres
```

### Disk Space Issues

```bash
# Check disk usage
df -h

# Clean up Docker
docker system prune -a --volumes

# Remove old backups
find ~/backups -name "noddit_*.sql.gz" -mtime +7 -delete
```

### GitHub Actions Deployment Failing

```bash
# Check secrets are set correctly
# Verify VM is accessible:
ssh -i ~/.ssh/deploy_key deploy@VM_IP

# Check deploy user permissions
ls -la /home/deploy/noddit

# View GitHub Actions logs in browser
```

## Cost Management

### Free Tier Limits

- **VM:** 1x e2-micro in us-central1/west1/east1 (always free)
- **Storage:** 30GB standard persistent disk (always free)
- **Network:** 1GB egress to Americas/month (always free)

### Monitor Usage

Go to: [GCP Console → Billing → Reports](https://console.cloud.google.com/billing)

Filter by:
- Service: Compute Engine
- Time range: Current month

**Expected cost:** $0.00 (if staying within free tier)

### Avoid Charges

- Don't upgrade VM size
- Don't add additional storage
- Keep network egress under 1GB/month
- Don't deploy in regions outside us-central1/west1/east1

## Scaling Beyond Free Tier

If your app grows beyond the e2-micro capacity:

### Option 1: Upgrade VM ($13-30/month)
```bash
gcloud compute instances stop noddit-vm --zone=us-central1-a
gcloud compute instances set-machine-type noddit-vm \
  --machine-type=e2-small \
  --zone=us-central1-a
gcloud compute instances start noddit-vm --zone=us-central1-a
```

### Option 2: Cloud Run + Neon ($0-5/month)
- Migrate backend to Cloud Run
- Migrate frontend to Cloud Run
- Use Neon.tech for free PostgreSQL
- See `DEPLOYMENT-CLOUDRUN.md` for guide

## Support

For issues:
1. Check logs: `docker-compose logs`
2. Check GitHub Actions logs
3. Review this guide
4. Check memory/disk: `free -h && df -h`
5. Open issue on GitHub repository

## Quick Reference

```bash
# SSH into VM
gcloud compute ssh noddit-vm --zone=us-central1-a

# View logs
docker-compose logs -f

# Restart app
sudo systemctl restart noddit

# Check status
docker-compose ps

# Manual deployment
git push origin deployment

# Run backup
~/backup-db.sh

# Check memory
free -h
docker stats --no-stream
```

## Next Steps

- [ ] Configure custom domain
- [ ] Setup SSL certificate
- [ ] Configure monitoring/alerts
- [ ] Setup log aggregation
- [ ] Configure CDN (if needed)
- [ ] Setup staging environment
