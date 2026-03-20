# Noddit GCP Deployment Quick Start

## TL;DR - Deploy in 10 Minutes

This is the fastest path to get Noddit running on GCP for **$0/month**.

### Prerequisites
- GCP account with billing enabled
- `gcloud` CLI installed
- GitHub repository

---

## Step 1: Create VM (2 minutes)

```bash
cd C:\Users\rob\Documents\workspace\noddit
bash scripts/gcp-setup.sh
```

When prompted:
- Project ID: `your-gcp-project-id`
- Region: `us-central1` (press Enter for default)
- Zone: `us-central1-a` (press Enter for default)
- Instance: `noddit-vm` (press Enter for default)

**Save the SSH private key output!** You'll need it in Step 2.

---

## Step 2: Configure GitHub (1 minute)

Go to: `https://github.com/YOUR_USERNAME/noddit/settings/secrets/actions`

Click "New repository secret" and add:

| Name | Value |
|------|-------|
| `SSH_PRIVATE_KEY` | The private key from Step 1 output |
| `VM_HOST` | Your VM IP from Step 1 output |
| `VM_USER` | `deploy` |
| `VM_PORT` | `22` |

---

## Step 3: Create Environment File (2 minutes)

```bash
# SSH into VM
gcloud compute ssh noddit-vm --zone=us-central1-a
sudo su - deploy
cd ~/noddit
nano .env
```

Paste and fill in:

```bash
DB_PASSWORD=$(openssl rand -base64 24)
JWT_SECRET=$(openssl rand -base64 32)
CLERK_SECRET_KEY=sk_live_YOUR_KEY
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_live_YOUR_KEY
EXTERNAL_IP=YOUR_VM_IP_HERE
```

**Generate secrets:**
```bash
# Generate strong password
openssl rand -base64 24

# Generate JWT secret
openssl rand -base64 32
```

Save: `Ctrl+X`, `Y`, `Enter`

---

## Step 4: Deploy (3 minutes)

```bash
# Still on VM
cp docker-compose.production.yml docker-compose.yml
docker-compose build
docker-compose up -d

# Wait ~2 minutes for build
docker-compose ps
```

You should see 3 containers running (postgres, backend, frontend).

---

## Step 5: Setup Auto-Start (1 minute)

```bash
bash ~/noddit/scripts/vm-systemd-setup.sh
bash ~/noddit/scripts/vm-backup-setup.sh
exit  # Exit deploy user
exit  # Exit VM
```

---

## Step 6: Test (1 minute)

```bash
# From local machine
export VM_IP=YOUR_VM_IP

# Test backend
curl http://$VM_IP:8080/api/public/subnoddits

# Test frontend (open in browser)
echo "http://$VM_IP:8081"
```

---

## ✅ Done!

Your app is now live at:
- **Frontend:** `http://YOUR_VM_IP:8081`
- **Backend:** `http://YOUR_VM_IP:8080`

**Future deployments:** Just push to GitHub!
```bash
git push origin deployment
```

---

## Optional Enhancements

### Serve on Port 80 (instead of 8081)

```bash
gcloud compute ssh noddit-vm --zone=us-central1-a
bash ~/noddit/scripts/vm-nginx-setup.sh

# Update .env
nano ~/noddit/.env
# Change NEXT_PUBLIC_API_URL to: http://YOUR_IP/api
# Change FRONTEND_URL to: http://YOUR_IP

# Rebuild frontend
cd ~/noddit
docker-compose down
docker-compose up -d --build frontend
```

Now accessible at: `http://YOUR_VM_IP/`

### Add SSL Certificate (requires domain)

```bash
# Point your domain to VM IP first
gcloud compute ssh noddit-vm --zone=us-central1-a
sudo apt install certbot python3-certbot-nginx -y
sudo certbot --nginx -d yourdomain.com
```

Now accessible at: `https://yourdomain.com`

---

## Troubleshooting

### Containers not starting
```bash
docker-compose logs
free -h  # Check memory
```

### GitHub Actions failing
- Verify secrets are set correctly
- Check VM is accessible: `ssh deploy@VM_IP`

### Out of memory
```bash
docker stats --no-stream  # Check usage
docker system prune -a     # Clean up
```

### Can't connect to app
```bash
# Check firewall
gcloud compute firewall-rules list

# Check containers
docker-compose ps
docker-compose logs
```

---

## Full Documentation

For detailed information, see:
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Complete deployment guide
- **[scripts/README.md](scripts/README.md)** - Script documentation

---

## Quick Commands Reference

```bash
# SSH into VM
gcloud compute ssh noddit-vm --zone=us-central1-a

# View logs
docker-compose logs -f

# Restart app
sudo systemctl restart noddit

# Check status
docker-compose ps

# Check memory
free -h
docker stats --no-stream

# Manual backup
~/backup-db.sh

# Deploy changes
git push origin deployment
```

---

## Cost

**$0/month** - Uses GCP Always Free Tier

Free tier includes:
- 1x e2-micro VM (0.25 vCPU, 1GB RAM)
- 30GB storage
- 1GB network egress/month

**Monitor usage:** [GCP Console → Billing](https://console.cloud.google.com/billing)
