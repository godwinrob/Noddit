#!/bin/bash
set -e

# GCP e2-micro VM Setup Script for Noddit
# This script creates and configures a free-tier GCP VM for deployment

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Noddit GCP Deployment Setup${NC}"
echo -e "${GREEN}================================${NC}"
echo ""

# Configuration
read -p "Enter your GCP Project ID: " PROJECT_ID
read -p "Enter region (default: us-central1): " REGION
REGION=${REGION:-us-central1}
read -p "Enter zone (default: us-central1-a): " ZONE
ZONE=${ZONE:-us-central1-a}
read -p "Enter VM instance name (default: noddit-vm): " INSTANCE_NAME
INSTANCE_NAME=${INSTANCE_NAME:-noddit-vm}

echo ""
echo -e "${YELLOW}Configuration:${NC}"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"
echo "Zone: $ZONE"
echo "Instance Name: $INSTANCE_NAME"
echo ""
read -p "Continue with this configuration? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

# Set project
echo -e "${GREEN}Setting GCP project...${NC}"
gcloud config set project $PROJECT_ID

# Create VM instance
echo -e "${GREEN}Creating e2-micro VM instance...${NC}"
gcloud compute instances create $INSTANCE_NAME \
  --zone=$ZONE \
  --machine-type=e2-micro \
  --network-tier=STANDARD \
  --subnet=default \
  --maintenance-policy=MIGRATE \
  --tags=http-server,https-server \
  --image=ubuntu-2204-jammy-v20240319 \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=30GB \
  --boot-disk-type=pd-standard \
  --boot-disk-device-name=$INSTANCE_NAME \
  || echo -e "${YELLOW}VM might already exist, continuing...${NC}"

# Create firewall rules
echo -e "${GREEN}Creating firewall rules...${NC}"

gcloud compute firewall-rules create allow-http \
  --direction=INGRESS \
  --network=default \
  --action=ALLOW \
  --rules=tcp:80 \
  --source-ranges=0.0.0.0/0 \
  --target-tags=http-server \
  || echo -e "${YELLOW}HTTP rule might already exist${NC}"

gcloud compute firewall-rules create allow-https \
  --direction=INGRESS \
  --network=default \
  --action=ALLOW \
  --rules=tcp:443 \
  --source-ranges=0.0.0.0/0 \
  --target-tags=https-server \
  || echo -e "${YELLOW}HTTPS rule might already exist${NC}"

gcloud compute firewall-rules create allow-app-ports \
  --direction=INGRESS \
  --network=default \
  --action=ALLOW \
  --rules=tcp:8080,tcp:8081 \
  --source-ranges=0.0.0.0/0 \
  --target-tags=http-server \
  || echo -e "${YELLOW}App ports rule might already exist${NC}"

# Get VM IP
echo -e "${GREEN}Getting VM external IP...${NC}"
VM_IP=$(gcloud compute instances describe $INSTANCE_NAME \
  --zone=$ZONE \
  --format="get(networkInterfaces[0].accessConfigs[0].natIP)")

echo -e "${GREEN}VM External IP: $VM_IP${NC}"

# Wait for VM to be ready
echo -e "${YELLOW}Waiting for VM to be ready (30 seconds)...${NC}"
sleep 30

# Install Docker on VM
echo -e "${GREEN}Installing Docker on VM...${NC}"
gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --command="
  set -e
  echo 'Updating system...'
  sudo apt update && sudo apt upgrade -y

  echo 'Installing Docker...'
  curl -fsSL https://get.docker.com -o get-docker.sh
  sudo sh get-docker.sh
  sudo usermod -aG docker \$USER

  echo 'Installing Docker Compose plugin...'
  sudo apt install docker-compose-plugin -y

  echo 'Enabling Docker...'
  sudo systemctl enable docker
  sudo systemctl start docker

  echo 'Docker installed successfully!'
"

# Create deploy user
echo -e "${GREEN}Creating deploy user...${NC}"
gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --command="
  set -e

  # Create deploy user
  sudo useradd -m -s /bin/bash -G docker deploy || echo 'User already exists'

  # Create project directory
  sudo mkdir -p /home/deploy/noddit
  sudo chown -R deploy:deploy /home/deploy/noddit

  echo 'Deploy user created!'
"

# Generate SSH key for GitHub Actions
echo -e "${GREEN}Setting up SSH keys for GitHub Actions...${NC}"
echo ""
echo -e "${YELLOW}Creating SSH key pair for GitHub Actions...${NC}"

gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --command="
  sudo -u deploy bash << 'DEPLOY_SCRIPT'
    cd ~
    mkdir -p ~/.ssh
    chmod 700 ~/.ssh

    # Generate SSH key
    ssh-keygen -t ed25519 -C 'github-actions-deploy' -f ~/.ssh/github_actions -N ''

    # Add to authorized_keys
    cat ~/.ssh/github_actions.pub >> ~/.ssh/authorized_keys
    chmod 600 ~/.ssh/authorized_keys

    echo ''
    echo '==================== GITHUB ACTIONS PRIVATE KEY ===================='
    echo 'Copy this ENTIRE output (including BEGIN/END lines) to GitHub Secrets as SSH_PRIVATE_KEY:'
    echo ''
    cat ~/.ssh/github_actions
    echo ''
    echo '===================================================================='
DEPLOY_SCRIPT
"

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Setup Complete!${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo ""
echo "1. Add GitHub Secrets (go to: https://github.com/YOUR_USERNAME/noddit/settings/secrets/actions):"
echo "   - SSH_PRIVATE_KEY: (the private key shown above)"
echo "   - VM_HOST: $VM_IP"
echo "   - VM_USER: deploy"
echo "   - VM_PORT: 22"
echo ""
echo "2. Create .env.production file on VM:"
echo "   gcloud compute ssh $INSTANCE_NAME --zone=$ZONE"
echo "   sudo su - deploy"
echo "   cd ~/noddit"
echo "   nano .env"
echo ""
echo "3. Copy the template from .env.production.template and fill in:"
echo "   - DB_PASSWORD (generate strong password)"
echo "   - JWT_SECRET (run: openssl rand -base64 32)"
echo "   - CLERK_SECRET_KEY (from Clerk dashboard)"
echo "   - NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY (from Clerk dashboard)"
echo "   - EXTERNAL_IP=$VM_IP"
echo ""
echo "4. Initial deployment:"
echo "   - Push to main branch, or"
echo "   - Trigger workflow manually in GitHub Actions"
echo ""
echo "5. Access your app:"
echo "   - Frontend: http://$VM_IP:8081"
echo "   - Backend: http://$VM_IP:8080"
echo ""
echo -e "${GREEN}VM is ready for deployment!${NC}"
