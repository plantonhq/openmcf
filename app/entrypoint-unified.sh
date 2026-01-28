#!/bin/bash

set -e

echo "=========================================="
echo "🚀 OpenMCF Unified Container"
echo "=========================================="
echo ""

# Create necessary directories if they don't exist
echo "📁 Setting up directories..."
mkdir -p /data/db
mkdir -p /var/log/mongodb
mkdir -p /var/log/supervisor
mkdir -p /home/appuser/.pulumi/state
mkdir -p /home/appuser/go/cache
mkdir -p /home/appuser/go/tmp
mkdir -p /home/appuser/.openmcf

# Set proper permissions
echo "🔒 Setting permissions..."
chown -R mongodb:mongodb /data/db
chown -R mongodb:mongodb /var/log/mongodb
chown -R appuser:root /home/appuser
chown -R appuser:root /app/backend
chown -R appuser:root /app/frontend

echo "✅ Setup complete"
echo ""

# Configure Pulumi backend (required for backend service)
echo "🔧 Configuring Pulumi backend..."
export PULUMI_HOME=/home/appuser/.pulumi
export PULUMI_CONFIG_PASSPHRASE=${PULUMI_CONFIG_PASSPHRASE:-openmcf-default-passphrase}

# Automatically choose backend based on environment variables
if [ -n "$PULUMI_ACCESS_TOKEN" ]; then
  # Use Pulumi Cloud if access token is provided
  echo "🌐 Detected PULUMI_ACCESS_TOKEN - using Pulumi Cloud backend"
  if [ -n "$PULUMI_BACKEND_URL" ]; then
    echo "   Backend URL: $PULUMI_BACKEND_URL"
  else
    echo "   Backend URL: https://api.pulumi.com (default)"
  fi
  su -s /bin/sh appuser -c "pulumi login --non-interactive"
else
  # Use local file-based backend by default
  echo "📁 Using local file-based backend (no PULUMI_ACCESS_TOKEN found)"
  echo "   State storage: /home/appuser/.pulumi/state"
  su -s /bin/sh appuser -c "pulumi login --local --non-interactive"
fi

echo "✅ Pulumi backend configured successfully"
echo ""
echo "Starting services:"
echo "  - MongoDB (port 27017)"
echo "  - Backend gRPC Server (port 50051)"
echo "  - Frontend Next.js (port 3000)"
echo ""

# Execute the command passed to the entrypoint
exec "$@"

