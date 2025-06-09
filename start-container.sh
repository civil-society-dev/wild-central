#!/bin/bash
set -e

echo "Starting Wild Cloud Central container..."

# Start nginx
echo "Starting nginx..."
nginx

# Start dnsmasq (will fail initially, but that's expected until config is generated)
echo "Starting dnsmasq..."
dnsmasq --keep-in-foreground --log-facility=- &
DNSMASQ_PID=$!

# Function to handle shutdown
shutdown() {
    echo "Shutting down services..."
    kill $DNSMASQ_PID 2>/dev/null || true
    nginx -s quit 2>/dev/null || true
    exit 0
}

# Set up signal handlers
trap shutdown SIGTERM SIGINT

# Start wild-cloud-central service
echo "Starting wild-cloud-central service..."
exec /usr/bin/wild-cloud-central