#!/bin/bash

set -e

# Configuration - update these for your server
SERVER="user@mywildcloud.org"
REMOTE_PATH="/var/www/html/apt"
LOCAL_REPO="apt-repo"

echo "🚀 Deploying APT repository to mywildcloud.org..."

# Check if repository exists
if [ ! -d "$LOCAL_REPO" ]; then
    echo "❌ Repository not found. Run './scripts/build-repo.sh' first."
    exit 1
fi

# Deploy repository
echo "📤 Uploading repository files..."
rsync -av --progress "$LOCAL_REPO/" "$SERVER:$REMOTE_PATH/"

# Deploy GPG public key
if [ -f "wild-cloud-central.gpg" ]; then
    echo "🔑 Uploading GPG public key..."
    scp wild-cloud-central.gpg "$SERVER:$REMOTE_PATH/"
else
    echo "⚠️  GPG public key not found. Run './scripts/setup-gpg.sh' first."
fi

echo ""
echo "✅ Deployment complete!"
echo ""
echo "🌐 Repository URL: https://mywildcloud.org/apt"
echo "🔑 GPG Key URL: https://mywildcloud.org/apt/wild-cloud-central.gpg"
echo ""
echo "👥 Users can now install with:"
echo "   curl -fsSL https://mywildcloud.org/apt/wild-cloud-central.gpg | sudo apt-key add -"
echo "   echo 'deb https://mywildcloud.org/apt stable main' | sudo tee /etc/apt/sources.list.d/wild-cloud-central.list"
echo "   sudo apt update"
echo "   sudo apt install wild-cloud-central"