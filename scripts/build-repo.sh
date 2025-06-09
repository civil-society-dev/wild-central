#!/bin/bash

set -e

# Configuration
REPO_DIR="apt-repo"
DIST="stable"
COMPONENT="main"
ARCH="amd64"
PACKAGE_NAME="wild-cloud-central"

echo "🏗️  Building APT repository..."

# Create directory structure
mkdir -p "$REPO_DIR/pool/main/w/$PACKAGE_NAME"
mkdir -p "$REPO_DIR/dists/$DIST/$COMPONENT/binary-$ARCH"

# Copy .deb file to pool
cp build/wild-cloud-central_*.deb "$REPO_DIR/pool/main/w/$PACKAGE_NAME/"

# Generate Packages file
cd "$REPO_DIR"
dpkg-scanpackages pool/ /dev/null > "dists/$DIST/$COMPONENT/binary-$ARCH/Packages"
gzip -k -f "dists/$DIST/$COMPONENT/binary-$ARCH/Packages"

# Generate Release files
cd "dists/$DIST"

# Component Release file
cat > "$COMPONENT/binary-$ARCH/Release" << EOF
Archive: $DIST
Component: $COMPONENT
Origin: Wild Cloud Central
Label: Wild Cloud Central
Architecture: $ARCH
EOF

# Main Release file
cat > Release << EOF
Archive: $DIST
Codename: $DIST
Date: $(date -Ru)
Origin: Wild Cloud Central
Label: Wild Cloud Central Repository
Suite: $DIST
Architectures: $ARCH
Components: $COMPONENT
Description: Wild Cloud Central APT Repository
MD5Sum:
$(find . -name "Packages*" -exec md5sum {} \; | sed 's|\./||')
SHA1:
$(find . -name "Packages*" -exec sha1sum {} \; | sed 's|\./||')
SHA256:
$(find . -name "Packages*" -exec sha256sum {} \; | sed 's|\./||')
EOF

# Sign the Release file (requires GPG key)
if command -v gpg &> /dev/null; then
    gpg --armor --detach-sign --sign Release
    echo "✅ Repository signed with GPG"
else
    echo "⚠️  GPG not found - repository will not be signed"
fi

cd ../../..

echo "✅ APT repository built in $REPO_DIR/"
echo ""
echo "📦 To deploy:"
echo "   rsync -av $REPO_DIR/ user@mywildcloud.org:/var/www/html/apt/"
echo ""
echo "👥 Users can install with:"
echo "   curl -fsSL https://mywildcloud.org/apt/wild-cloud-central.gpg | sudo apt-key add -"
echo "   echo 'deb https://mywildcloud.org/apt stable main' | sudo tee /etc/apt/sources.list.d/wild-cloud-central.list"
echo "   sudo apt update"
echo "   sudo apt install wild-cloud-central"