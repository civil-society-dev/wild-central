# Wild Cloud Central

A Go-based web service for managing wild-cloud infrastructure, providing DNS, DHCP, and PXE boot services for Talos Linux clusters.

## Features

- **Web Management Interface** - Browser-based configuration and monitoring
- **REST API** - JSON API for programmatic management
- **DNS/DHCP Services** - Integrated dnsmasq configuration management
- **PXE Boot Support** - Automatic Talos Linux asset downloading and serving
- **Containerized Testing** - Docker support for development and testing

## Quick Start

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

> ⚠️ **Note**: These Docker scripts test the installation process only. In production, use `sudo apt install wild-cloud-central` and manage via systemd.

Choose the testing approach that fits your needs:

### 1. Automated Verification - `./test-docker.sh`
- **When to use**: Verify the installation works correctly
- **What it does**: Builds .deb package, installs it, tests all endpoints automatically
- **Best for**: CI/CD, quick verification that everything works

### 2. Background Testing - `./start-background.sh` / `./stop-background.sh`  
- **When to use**: You want to test APIs while doing other work
- **What it does**: Starts services silently in background, gives you your terminal back
- **Example workflow**: Start services, test in another terminal, stop when done
```bash
./start-background.sh           # Services start, terminal returns immediately
curl http://localhost:9081/api/v1/health  # Test in same or different terminal
# Continue working while services run...
./stop-background.sh            # Clean shutdown when finished
```

### 3. Interactive Development - `./start-interactive.sh`
- **When to use**: You want to see what's happening as you test
- **What it does**: Starts services with live logs, takes over your terminal
- **Example workflow**: Start services, watch logs in real-time, Ctrl+C to stop
```bash
./start-interactive.sh          # Services start, shows live logs
# You see all HTTP requests, errors, debug info in real-time
# Press Ctrl+C when done - terminal is "busy" until then
```

### 4. Shell Access - `./debug-container.sh`
- **When to use**: Deep debugging, manual service control, file inspection
- **What it does**: Drops you into the container shell
- **Best for**: Investigating issues, manually starting/stopping services

### Test Access Points
All services bind to localhost (127.0.0.1) on non-standard ports, so they won't interfere with your local services:

- Management UI: http://localhost:9080
- API: http://localhost:9081
- DNS: localhost:9053 (UDP) - test with `dig @localhost -p 9053 wildcloud.local`
- DHCP: localhost:9067 (UDP)
- TFTP: localhost:9069 (UDP)
- Container logs: `docker logs wild-central-bg`

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