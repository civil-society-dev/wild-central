# Wild-cloud Central

## Current bash implementation

These scripts help a "wild-cloud operator" set up the first component of a "wild-cloud", which is an inexpensive linux machine that will serve as "wild-cloud Central". Currently, it is a dnsmasq server that provides dns and pxe services to the cluster that will be set up next. Additional functionality and services will be added in the future.

### `install-dnsmasq`

- First phase of an installation that is intended to be run on a sysop's computer.

- Reads in config from .wildcloud/config.yaml and ./wildcloud/secrets.yaml
- Generates a set of customized templates for the user that will be used in the second phase.
- Downloads and caches some pxe files that.

### `dnsmasq\setup.sh`

This file as well as the rest of the files in the `dnsmasq` directory are copied to the wild-central server and used to set it up.

When run:

- installs dnsmasq and nginx
- configures dnsmasq to serve ipxe images (Talos linux) on the wild-cloud network

### Future work

In creating this, we note the system is overly complex and mixes up responsibilities between components and makes setting up difficult.

We want to re-implement the current functionality as a web service with a JSON API that will be installed with `sudo apt install wild-cloud-central`.

#### Implementation Details

The new wild-cloud-central service will be:

- **Written in Go** - for better performance, easier deployment, and cross-compilation capabilities
- **Debian-based systems only** - targeting Ubuntu/Debian distributions with apt package management
- **Single binary deployment** - packaged as a .deb file for simple installation and management
- **Web service with JSON API** - replacing bash scripts with HTTP endpoints for configuration and management
- **Systemd integration** - proper service management with start/stop/restart capabilities
- **Configuration management** - centralized config handling replacing scattered YAML template processing

#### Service Architecture

The Go-based service will provide:

- REST API endpoints for cluster configuration and management
- Built-in web UI for operator interaction
- Integrated dnsmasq configuration management
- Automatic PXE asset downloading and caching
- Real-time status monitoring and logging
- Simplified installation process removing multi-phase setup complexity

This approach will eliminate the current split between local preparation and remote setup phases, providing a unified management interface for wild-cloud infrastructure.
