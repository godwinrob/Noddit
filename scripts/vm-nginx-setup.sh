#!/bin/bash
# Optional Nginx Reverse Proxy Setup
# Run this ON THE VM if you want to serve frontend on port 80 with /api routing

set -e

echo "Installing and configuring Nginx..."

# Install Nginx
sudo apt update
sudo apt install nginx -y

# Create Nginx configuration
sudo tee /etc/nginx/sites-available/noddit << 'EOF'
server {
    listen 80;
    server_name _;

    # Frontend
    location / {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Backend API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

# Enable site
sudo ln -sf /etc/nginx/sites-available/noddit /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default

# Test configuration
sudo nginx -t

# Restart Nginx
sudo systemctl restart nginx
sudo systemctl enable nginx

echo "✅ Nginx setup complete!"
echo ""
echo "Your application is now available on port 80:"
echo "  Frontend: http://YOUR_VM_IP/"
echo "  Backend: http://YOUR_VM_IP/api"
echo ""
echo "⚠️  IMPORTANT: Update your .env file:"
echo "  NEXT_PUBLIC_API_URL=http://YOUR_VM_IP/api"
echo "  FRONTEND_URL=http://YOUR_VM_IP"
echo ""
echo "Then rebuild the frontend:"
echo "  cd ~/noddit"
echo "  docker-compose down"
echo "  docker-compose up -d --build frontend"
