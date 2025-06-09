# Maintainer Guide

This guide covers package creation, repository management, and deployment for Wild Cloud Central.

## Package Management

### Creating Debian Package

```bash
make deb
```

This creates `build/wild-cloud-central_0.1.0_amd64.deb` with:
- Binary installed to `/usr/bin/wild-cloud-central`
- Systemd service file
- Configuration template
- Web interface files
- Nginx configuration

### Package Structure

The .deb package includes:
- `/usr/bin/wild-cloud-central` - Main binary
- `/etc/systemd/system/wild-cloud-central.service` - Systemd service
- `/etc/wild-cloud-central/config.yaml.example` - Configuration template
- `/var/www/html/wild-central/` - Web interface files
- `/etc/nginx/sites-available/wild-central` - Nginx configuration

### Post-installation Setup

The package automatically:
- Creates `wildcloud` system user
- Creates required directories with proper permissions
- Configures nginx
- Enables systemd service
- Sets up file ownership

## APT Repository Management

### Building Repository

```bash
make repo
```

This creates a complete APT repository structure in `apt-repo/`:
- `pool/main/w/wild-cloud-central/` - Package files
- `dists/stable/main/binary-amd64/` - Metadata and package lists
- `dists/stable/Release.asc` - GPG signature

### GPG Key Management

#### First-time Setup
```bash
./scripts/setup-gpg.sh
```

This creates:
- 4096-bit RSA GPG key pair
- Public key exported as `wild-cloud-central.gpg`
- Key configured for 2-year expiration

#### Key Renewal
When the key expires, regenerate with:
```bash
gpg --delete-secret-keys "Wild Cloud Central"
gpg --delete-keys "Wild Cloud Central"
rm wild-cloud-central.gpg
./scripts/setup-gpg.sh
```

### Repository Deployment

1. **Configure server details** in `scripts/deploy-repo.sh`:
   ```bash
   SERVER="user@mywildcloud.org"
   REMOTE_PATH="/var/www/html/apt"
   ```

2. **Deploy repository**:
   ```bash
   make deploy-repo
   ```

This uploads:
- Complete repository structure to server
- GPG public key for user verification
- Proper file permissions and structure

### Server Requirements

The target server needs:
- Web server (nginx/apache) serving `/var/www/html/apt`
- HTTPS support for `https://mywildcloud.org/apt`
- SSH access for deployment

### Repository Structure

```
/var/www/html/apt/
├── dists/
│   └── stable/
│       ├── Release
│       ├── Release.asc
│       └── main/
│           └── binary-amd64/
│               ├── Packages
│               ├── Packages.gz
│               └── Release
├── pool/
│   └── main/
│       └── w/
│           └── wild-cloud-central/
│               └── wild-cloud-central_0.1.0_amd64.deb
└── wild-cloud-central.gpg
```

## Release Process

1. **Update version** in `Makefile`:
   ```makefile
   VERSION?=0.2.0
   ```

2. **Build and test**:
   ```bash
   make clean
   make deb
   ./test-docker.sh
   ```

3. **Create repository**:
   ```bash
   make repo
   ```

4. **Deploy**:
   ```bash
   make deploy-repo
   ```

5. **Verify deployment**:
   ```bash
   curl -I https://mywildcloud.org/apt/dists/stable/Release
   curl -I https://mywildcloud.org/apt/wild-cloud-central.gpg
   ```

## Troubleshooting

### GPG Issues
- **"no default secret key"**: Run `./scripts/setup-gpg.sh`
- **Key conflicts**: Delete existing keys before recreating
- **Permission errors**: Ensure `~/.gnupg` has correct permissions (700)

### Repository Issues
- **Package not found**: Verify `dpkg-scanpackages` output
- **Signature verification failed**: Regenerate GPG key and re-sign
- **404 errors**: Check web server configuration and file permissions

### Deployment Issues
- **SSH failures**: Verify server credentials in `deploy-repo.sh`
- **Permission denied**: Ensure target directory is writable
- **rsync errors**: Check network connectivity and paths

## Monitoring

### Service Health
```bash
curl https://mywildcloud.org/apt/dists/stable/Release
curl https://mywildcloud.org/apt/wild-cloud-central.gpg
```

### Package Statistics
Monitor download statistics through web server logs:
```bash
grep "wild-cloud-central.*\.deb" /var/log/nginx/access.log | wc -l
```

### Repository Integrity
Verify signatures regularly:
```bash
gpg --verify Release.asc Release
```