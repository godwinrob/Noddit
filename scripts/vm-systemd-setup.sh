#!/bin/bash
# Systemd Service Setup for Noddit
# Run this ON THE VM after initial deployment to enable auto-start

set -e

echo "Setting up systemd service for Noddit..."

# Create systemd service file
sudo tee /etc/systemd/system/noddit.service << EOF
[Unit]
Description=Noddit Docker Compose Application
Requires=docker.service
After=docker.service network-online.target
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/deploy/noddit
ExecStart=/usr/bin/docker compose up -d
ExecStop=/usr/bin/docker compose down
TimeoutStartSec=300
User=deploy
Group=deploy

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable noddit.service

echo "✅ Systemd service setup complete!"
echo ""
echo "Service management commands:"
echo "  sudo systemctl start noddit     # Start the application"
echo "  sudo systemctl stop noddit      # Stop the application"
echo "  sudo systemctl restart noddit   # Restart the application"
echo "  sudo systemctl status noddit    # Check service status"
echo ""
echo "The application will now start automatically on system boot."
