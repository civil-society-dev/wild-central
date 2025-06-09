# Wild Cloud Central

A Go-based web service for managing wild-cloud infrastructure, providing DNS, DHCP, and PXE boot services for Talos Linux clusters.

## Features

- **Web Management Interface** - Browser-based configuration and monitoring
- **REST API** - JSON API for programmatic management
- **DNS/DHCP Services** - Integrated dnsmasq configuration management
- **PXE Boot Support** - Automatic Talos Linux asset downloading and serving
- **Containerized Testing** - Docker support for development and testing

## Quick Start

### Using Docker (Recommended for Testing)

1. **Build and run with Docker Compose:**
   ```bash
   docker-compose up --build
   ```

2. **Access the web interface:**
   - Management UI: http://localhost:8080
   - API directly: http://localhost:8081

### Building Locally

1. **Build the application:**
   ```bash
   make build
   ```

2. **Run locally:**
   ```bash
   make run
   ```

### Creating Debian Package

```bash
make deb
```

This creates `build/wild-cloud-central_0.1.0_amd64.deb` ready for installation.

## API Endpoints

- `GET /api/v1/health` - Service health check
- `GET /api/v1/config` - Get current configuration
- `PUT /api/v1/config` - Update configuration
- `GET /api/v1/dnsmasq/config` - Generate dnsmasq configuration
- `POST /api/v1/dnsmasq/restart` - Restart dnsmasq service
- `POST /api/v1/pxe/assets` - Download/update PXE boot assets

## Configuration

Edit `config.yaml` to customize your deployment:

```yaml
server:
  port: 8081
  host: "0.0.0.0"

cloud:
  domain: "wildcloud.local"
  dns:
    ip: "192.168.8.50"
  dhcpRange: "192.168.8.100,192.168.8.200"

cluster:
  endpointIp: "192.168.8.60"
  nodes:
    talos:
      version: "v1.8.0"
```

## Testing

The Docker container includes:
- dnsmasq for DNS/DHCP services
- nginx for serving PXE assets and web interface
- wild-cloud-central management service

### Container Testing Commands

```bash
# Build and start services
docker-compose up --build

# Check service health
curl http://localhost:8081/api/v1/health

# Download PXE assets
curl -X POST http://localhost:8081/api/v1/pxe/assets

# Generate dnsmasq config
curl http://localhost:8081/api/v1/dnsmasq/config
```

## Production Installation

```bash
# Install the .deb package
sudo dpkg -i build/wild-cloud-central_0.1.0_amd64.deb

# Configure
sudo cp config.yaml /etc/wild-cloud-central/config.yaml

# Enable and start service
sudo systemctl enable wild-cloud-central
sudo systemctl start wild-cloud-central
```

## Development

- **Language:** Go 1.21+
- **Dependencies:** gorilla/mux, gopkg.in/yaml.v3
- **Target:** Debian-based systems (Ubuntu/Debian)

## Architecture

This service replaces the original bash script implementation with:
- Unified configuration management
- Real-time dnsmasq configuration generation
- Integrated Talos factory asset downloading
- Web-based management interface
- Proper systemd service integration