#!/bin/bash

set -e

echo "🧪 Testing wild-cloud-central Docker installation..."

# Build the Docker image
echo "🔨 Building Docker image..."
docker build -t wild-cloud-central-test .

# Run the container to test installation
echo "🚀 Running installation test..."
echo "Access points after container starts:"
echo "  - Management UI: http://localhost:9080"
echo "  - API directly: http://localhost:9081"
echo ""
docker run --rm -p 9081:8081 -p 9080:80 wild-cloud-central-test