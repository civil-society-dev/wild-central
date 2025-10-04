# Packaging Wild Central

## Desired Experience

This is the desired experience for installing Wild Cloud Central on a fresh Debian/Ubuntu system:

### APT Repository (Recommended)

```bash
# Download and install GPG key
curl -fsSL https://mywildcloud.org/apt/wild-cloud-central.gpg | sudo tee /usr/share/keyrings/wild-cloud-central-archive-keyring.gpg > /dev/null

# Add repository (modern .sources format)
sudo tee /etc/apt/sources.list.d/wild-cloud-central.sources << 'EOF'
Types: deb
URIs: https://mywildcloud.org/apt
Suites: stable
Components: main
Signed-By: /usr/share/keyrings/wild-cloud-central-archive-keyring.gpg
EOF

# Update and install
sudo apt update
sudo apt install wild-cloud-central
```

### Manual Installation

Download the latest `.deb` package from the [releases page](https://github.com/wildcloud/wild-central/releases) and install:

```bash
sudo dpkg -i wild-cloud-central_*.deb
sudo apt-get install -f  # Fix any dependency issues
```

## Quick Start

1. **Configure the service** (optional):

   ```bash
   sudo cp /etc/wild-cloud-central/config.yaml.example /etc/wild-cloud-central/config.yaml
   sudo nano /etc/wild-cloud-central/config.yaml
   ```

2. **Start the service**:

   ```bash
   sudo systemctl enable wild-cloud-central
   sudo systemctl start wild-cloud-central
   ```

3. **Access the web interface**:
   Open http://your-server-ip in your browser

## Developer tooling

Makefile commands for packaging:

Build targets (compile binaries):

make build           - Build for current architecture
make build-arm64     - Build arm64 binary
make build-amd64     - Build amd64 binary
make build-all       - Build all architectures

Package targets (create .deb packages):

make package         - Create .deb package for current arch
make package-arm64   - Create arm64 .deb package
make package-amd64   - Create amd64 .deb package
make package-all     - Create all .deb packages

Repository targets:

make repo            - Build APT repository from packages
make deploy-repo     - Deploy repository to server

Quality assurance:

make check           - Run all checks (fmt + vet + test)
make fmt             - Format Go code
make vet             - Run go vet
make test            - Run tests

Development:

make run             - Run application locally
make clean           - Remove all build artifacts
make deps-check      - Verify and tidy dependencies
make version         - Show build information
make install         - Install to system

Directory structure:

build/          - Intermediate build artifacts
dist/bin/       - Final binaries for distribution
dist/packages/  - OS packages (.deb files)
dist/repositories/ - APT repository for deployment

Example workflows:
make check && make build     - Safe development build
make clean && make repo      - Full release build
