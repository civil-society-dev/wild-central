#!/bin/bash

set -e

echo "🔑 Setting up GPG key for Wild Cloud Central APT repository..."

# Check if GPG key already exists
if gpg --list-secret-keys | grep -q "Wild Cloud Central"; then
    echo "✅ GPG key already exists"
    KEY_ID=$(gpg --list-secret-keys --with-colons | grep "Wild Cloud Central" -B1 | grep "^sec" | cut -d: -f5)
    echo "Key ID: $KEY_ID"
else
    echo "🔧 Creating new GPG key..."
    
    # Create GPG key configuration
    cat > gpg-key-config << EOF
%echo Generating GPG key for Wild Cloud Central
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: Wild Cloud Central
Name-Comment: APT Repository Signing Key
Name-Email: apt@mywildcloud.org
Expire-Date: 2y
%no-protection
%commit
%echo GPG key created
EOF

    # Generate the key
    gpg --batch --generate-key gpg-key-config
    rm gpg-key-config
    
    KEY_ID=$(gpg --list-secret-keys --with-colons | grep "Wild Cloud Central" -B1 | grep "^sec" | cut -d: -f5)
    echo "✅ New GPG key created with ID: $KEY_ID"
fi

# Export public key
echo "📤 Exporting public key..."
gpg --armor --export $KEY_ID > wild-cloud-central.gpg

echo ""
echo "✅ GPG setup complete!"
echo ""
echo "📋 Next steps:"
echo "   1. Upload wild-cloud-central.gpg to https://mywildcloud.org/apt/"
echo "   2. Users will add this key with:"
echo "      curl -fsSL https://mywildcloud.org/apt/wild-cloud-central.gpg | sudo apt-key add -"
echo ""
echo "🔐 Key ID: $KEY_ID"
echo "📄 Public key saved to: wild-cloud-central.gpg"