# Noddit GCP Deployment Checklist

Use this checklist to ensure a smooth deployment to GCP.

## Pre-Deployment

- [ ] **GCP Account Setup**
  - [ ] GCP account created
  - [ ] Billing enabled (required for free tier)
  - [ ] Credit card on file
  - [ ] Project created

- [ ] **Local Environment**
  - [ ] `gcloud` CLI installed
  - [ ] `gcloud auth login` completed
  - [ ] Git repository pushed to GitHub

- [ ] **Secrets Prepared**
  - [ ] Database password generated (`openssl rand -base64 24`)
  - [ ] JWT secret generated (`openssl rand -base64 32`)
  - [ ] Clerk API keys obtained (if using Clerk)
  - [ ] Secrets documented securely (password manager)

---

## Initial Deployment

- [ ] **Step 1: Create GCP VM**
  - [ ] Run `bash scripts/gcp-setup.sh`
  - [ ] VM created successfully
  - [ ] Firewall rules configured
  - [ ] Docker installed on VM
  - [ ] Deploy user created
  - [ ] SSH private key saved for GitHub Actions

- [ ] **Step 2: Configure GitHub Secrets**
  - [ ] Navigate to GitHub repo settings → Secrets
  - [ ] Add `SSH_PRIVATE_KEY` (from setup script)
  - [ ] Add `VM_HOST` (VM external IP)
  - [ ] Add `VM_USER` (`deploy`)
  - [ ] Add `VM_PORT` (`22`)

- [ ] **Step 3: Configure VM Environment**
  - [ ] SSH into VM: `gcloud compute ssh noddit-vm --zone=us-central1-a`
  - [ ] Switch to deploy user: `sudo su - deploy`
  - [ ] Create `.env` file in `~/noddit/`
  - [ ] Add `DB_PASSWORD`
  - [ ] Add `JWT_SECRET`
  - [ ] Add `CLERK_SECRET_KEY`
  - [ ] Add `NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY`
  - [ ] Add `EXTERNAL_IP` (VM IP address)
  - [ ] File saved and permissions correct (`chmod 600 .env`)

- [ ] **Step 4: Initial Manual Deployment**
  - [ ] Copy production compose: `cp docker-compose.production.yml docker-compose.yml`
  - [ ] Build containers: `docker-compose build`
  - [ ] Start containers: `docker-compose up -d`
  - [ ] All containers running: `docker-compose ps`
  - [ ] No errors in logs: `docker-compose logs`

- [ ] **Step 5: Verify Deployment**
  - [ ] Backend accessible: `curl http://VM_IP:8080/api/public/subnoddits`
  - [ ] Frontend accessible in browser: `http://VM_IP:8081`
  - [ ] Can register new user
  - [ ] Can login
  - [ ] Database persists data

- [ ] **Step 6: Setup Auto-Start**
  - [ ] Run `bash ~/noddit/scripts/vm-systemd-setup.sh`
  - [ ] Service enabled: `sudo systemctl status noddit`
  - [ ] Test reboot: `sudo reboot` (wait 2 min, verify containers auto-start)

- [ ] **Step 7: Setup Backups**
  - [ ] Run `bash ~/noddit/scripts/vm-backup-setup.sh`
  - [ ] Test manual backup: `~/backup-db.sh`
  - [ ] Verify backup created: `ls ~/backups/`
  - [ ] Cron job configured: `crontab -l`

---

## Optional Enhancements

- [ ] **Nginx Reverse Proxy** (serve on port 80)
  - [ ] Run `bash ~/noddit/scripts/vm-nginx-setup.sh`
  - [ ] Update `.env` with new URLs
  - [ ] Rebuild frontend: `docker-compose up -d --build frontend`
  - [ ] Test: `http://VM_IP/` (port 80)

- [ ] **SSL Certificate** (requires domain)
  - [ ] Domain DNS pointing to VM IP
  - [ ] DNS propagated (check with `dig yourdomain.com`)
  - [ ] Install certbot: `sudo apt install certbot python3-certbot-nginx`
  - [ ] Obtain certificate: `sudo certbot --nginx -d yourdomain.com`
  - [ ] Auto-renewal enabled: `sudo certbot renew --dry-run`
  - [ ] Test: `https://yourdomain.com`

- [ ] **Security Hardening**
  - [ ] UFW firewall configured
  - [ ] SSH password authentication disabled
  - [ ] Unattended security updates enabled
  - [ ] Non-standard SSH port (optional)
  - [ ] Fail2ban installed (optional)

---

## Post-Deployment Testing

- [ ] **Functionality Tests**
  - [ ] Homepage loads
  - [ ] User registration works
  - [ ] User login works
  - [ ] Create subnoddit
  - [ ] Create post
  - [ ] Add comment
  - [ ] Upvote/downvote
  - [ ] User profile page
  - [ ] Search functionality

- [ ] **Performance Tests**
  - [ ] Page load times acceptable (<2 seconds)
  - [ ] API response times acceptable (<500ms)
  - [ ] Multiple concurrent users (test with friends)
  - [ ] Memory usage stable: `docker stats`

- [ ] **Reliability Tests**
  - [ ] Containers restart on failure
  - [ ] System reboot recovery (auto-start)
  - [ ] Database backup/restore
  - [ ] GitHub Actions deployment works

---

## Continuous Deployment

- [ ] **GitHub Actions Setup**
  - [ ] Workflow file at `.github/workflows/deploy.yml`
  - [ ] Secrets configured in GitHub
  - [ ] Test push triggers deployment
  - [ ] Health checks pass
  - [ ] Can view deployment logs in Actions tab

- [ ] **Deployment Process**
  - [ ] Make code changes locally
  - [ ] Commit changes
  - [ ] Push to `main` or `deployment` branch
  - [ ] GitHub Actions deploys automatically
  - [ ] Verify deployment succeeded
  - [ ] Test changes on live site

---

## Monitoring & Maintenance

- [ ] **Regular Checks** (weekly)
  - [ ] Check VM resource usage: `free -h && df -h`
  - [ ] Review container logs: `docker-compose logs --tail=100`
  - [ ] Verify backups are running: `ls -lh ~/backups/`
  - [ ] Check for security updates: `sudo apt list --upgradable`

- [ ] **Monthly Tasks**
  - [ ] Review GCP billing (should be $0)
  - [ ] Test backup restoration
  - [ ] Review and rotate logs
  - [ ] Update dependencies (if needed)

- [ ] **Monitoring Setup** (optional)
  - [ ] Setup monitoring alerts
  - [ ] Configure uptime monitoring (UptimeRobot, etc.)
  - [ ] Setup error tracking (Sentry, etc.)

---

## Troubleshooting Reference

### Common Issues

**Containers not starting:**
```bash
docker-compose ps
docker-compose logs
free -h  # Check memory
```

**Out of memory:**
```bash
docker stats --no-stream
docker system prune -a
```

**Database connection errors:**
```bash
docker-compose logs postgres
docker exec -it noddit-postgres psql -U postgres -d userdb
```

**GitHub Actions failing:**
- Verify SSH key in secrets
- Check VM is accessible
- Review Actions logs

**Can't access application:**
```bash
# Check firewall
gcloud compute firewall-rules list

# Check containers
docker-compose ps

# Check ports
sudo netstat -tlnp | grep -E '8080|8081'
```

---

## Rollback Plan

If deployment fails:

1. **Quick rollback on VM:**
   ```bash
   cd ~/noddit
   docker-compose down
   git checkout PREVIOUS_COMMIT_HASH
   docker-compose build
   docker-compose up -d
   ```

2. **Database restoration:**
   ```bash
   gunzip -c ~/backups/noddit_YYYYMMDD_HHMMSS.sql.gz | \
     docker exec -i noddit-postgres psql -U postgres -d userdb
   ```

3. **Revert GitHub commit:**
   ```bash
   git revert COMMIT_HASH
   git push origin deployment
   ```

---

## Success Criteria

Deployment is successful when:

- ✅ All containers running and healthy
- ✅ Frontend accessible in browser
- ✅ Backend API responding
- ✅ User can register and login
- ✅ Database persists data across restarts
- ✅ Auto-start works after VM reboot
- ✅ Backups running daily
- ✅ GitHub Actions deployments working
- ✅ Memory usage <900MB (leaves buffer)
- ✅ GCP billing shows $0.00

---

## Documentation

After deployment, ensure you have:

- [ ] VM IP address documented
- [ ] All passwords in password manager
- [ ] GitHub secrets documented
- [ ] Domain configuration (if applicable)
- [ ] Deployment runbook (this checklist)
- [ ] Backup restoration procedure tested

---

## Next Steps

After successful deployment:

1. **Share with users** - Send them the URL
2. **Monitor usage** - Check metrics and logs
3. **Gather feedback** - Identify issues early
4. **Plan improvements** - Based on user feedback
5. **Regular maintenance** - Follow monitoring schedule

---

## Emergency Contacts

**GCP Support:** [https://console.cloud.google.com/support](https://console.cloud.google.com/support)

**GitHub Support:** [https://support.github.com](https://support.github.com)

**Your documentation:** See DEPLOYMENT.md for detailed procedures

---

## Completion

- [ ] All items in this checklist completed
- [ ] Application tested and working
- [ ] Documentation updated
- [ ] Team notified of deployment
- [ ] Monitoring configured

**Deployed by:** ________________
**Date:** ________________
**VM IP:** ________________
**Domain:** ________________ (if applicable)

---

**🎉 Congratulations on your successful deployment!**
