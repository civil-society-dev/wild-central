# Wild Cloud CLI Scripts Reference

Wild Cloud provides 34+ command-line tools (all prefixed with `wild-`) for managing your personal Kubernetes cloud infrastructure. These scripts handle everything from initial setup to day-to-day operations.

## Script Categories

### 🚀 Initial Setup & Scaffolding

**`wild-init`** - Initialize new Wild Cloud instance
- Creates `.wildcloud` directory structure
- Copies template files from repository
- Sets up basic configuration (email, domains, cluster name)
- **Usage**: `wild-init`
- **When to use**: First command to run in a new directory

**`wild-setup`** - Master setup orchestrator
- Runs complete Wild Cloud deployment sequence
- Options: `--skip-cluster`, `--skip-services`
- Executes: cluster setup → services setup
- **Usage**: `wild-setup [options]`
- **When to use**: After `wild-init` for complete automated setup

**`wild-update-docs`** - Copy documentation to cloud directory
- Options: `--force` to overwrite existing docs
- Copies `/docs` from repository to your cloud home
- **Usage**: `wild-update-docs [--force]`

### ⚙️ Configuration Management

**`wild-config`** - Read configuration values
- Access YAML paths from `config.yaml` (e.g., `cluster.name`, `cloud.domain`)
- Option: `--check` to test key existence
- **Usage**: `wild-config <key>` or `wild-config --check <key>`

**`wild-config-set`** - Write configuration values
- Sets values using YAML paths, creates config file if needed
- **Usage**: `wild-config-set <key> <value>`

**`wild-secret`** - Read secret values
- Similar to `wild-config` but for sensitive data in `secrets.yaml`
- File has restrictive permissions (600)
- **Usage**: `wild-secret <key>` or `wild-secret --check <key>`

**`wild-secret-set`** - Write secret values
- Generates random values if none provided (32-char base64)
- **Usage**: `wild-secret-set <key> [value]`

**`wild-compile-template`** - Process gomplate templates
- Uses `config.yaml` and `secrets.yaml` as template context
- **Usage**: `wild-compile-template < template.yaml`

**`wild-compile-template-dir`** - Process template directories
- Options: `--clean` to remove destination first
- **Usage**: `wild-compile-template-dir <source> <destination>`

### 🏗️ Cluster Infrastructure Management

**`wild-setup-cluster`** - Complete cluster setup (Phases 1-3)
- Automated control plane node setup and bootstrapping
- Configures Talos control plane nodes using wild-node-setup
- Options: `--skip-hardware`
- **Usage**: `wild-setup-cluster [options]`
- **Requires**: `wild-init` completed first

**`wild-cluster-config-generate`** - Generate Talos cluster config
- Creates base `controlplane.yaml` and `worker.yaml`
- Generates cluster secrets using `talosctl gen config`
- **Usage**: `wild-cluster-config-generate`

**`wild-node-setup`** - Complete node lifecycle management
- Handles detect → configure → patch → deploy for individual nodes
- Automatically detects maintenance mode and handles IP transitions
- Options: `--reconfigure`, `--no-deploy`
- **Usage**: `wild-node-setup <node-name> [options]`
- **Examples**:
  - `wild-node-setup control-1` (complete setup)
  - `wild-node-setup worker-1 --reconfigure` (force node reconfiguration)
  - `wild-node-setup control-2 --no-deploy` (configuration only)

**`wild-node-detect`** - Hardware detection utility
- Discovers network interfaces and disks from maintenance mode
- Returns JSON with hardware specifications and maintenance mode status
- **Usage**: `wild-node-detect <node-ip>`
- **Note**: Primarily used internally by `wild-node-setup`

**`wild-cluster-node-ip`** - Get node IP addresses
- Sources: config.yaml, kubectl, or talosctl
- Options: `--from-config`, `--from-talosctl`
- **Usage**: `wild-cluster-node-ip <node-name> [options]`

### 🔧 Cluster Services Management

**`wild-setup-services`** - Set up all cluster services (Phase 4)
- Manages MetalLB, Traefik, cert-manager, etc. in dependency order
- Options: `--fetch` for fresh templates, `--no-deploy` for config-only
- **Usage**: `wild-setup-services [options]`
- **Requires**: Working Kubernetes cluster

**`wild-service-setup`** - Complete service lifecycle management
- Handles fetch → configure → deploy for individual services
- Options: `--fetch` for fresh templates, `--no-deploy` for config-only
- **Usage**: `wild-service-setup <service> [--fetch] [--no-deploy]`
- **Examples**:
  - `wild-service-setup cert-manager` (configure + deploy)
  - `wild-service-setup cert-manager --fetch` (fetch + configure + deploy)
  - `wild-service-setup cert-manager --no-deploy` (configure only)

**`wild-dashboard-token`** - Get Kubernetes dashboard token
- Extracts token for dashboard authentication
- Copies to clipboard if available
- **Usage**: `wild-dashboard-token`

**`wild-cluster-secret-copy`** - Copy secrets between namespaces
- **Usage**: `wild-cluster-secret-copy <source-ns:secret> <target-ns1> [target-ns2]`

### 📱 Application Management

**`wild-apps-list`** - List available applications
- Shows metadata, installation status, dependencies
- Options: `--verbose`, `--json`, `--yaml`
- **Usage**: `wild-apps-list [options]`

**`wild-app-add`** - Configure app from repository
- Processes manifest.yaml with configuration
- Generates required secrets automatically
- Options: `--force` to overwrite existing app files
- **Usage**: `wild-app-add <app-name> [--force]`

**`wild-app-deploy`** - Deploy application to cluster
- Creates namespaces, handles dependencies
- Options: `--force`, `--dry-run`
- **Usage**: `wild-app-deploy <app-name> [options]`

**`wild-app-delete`** - Remove application
- Deletes namespace and all resources
- Options: `--force`, `--dry-run`
- **Usage**: `wild-app-delete <app-name> [options]`

**`wild-app-doctor`** - Run application diagnostics
- Executes app-specific diagnostic tests
- Options: `--keep`, `--follow`, `--timeout`
- **Usage**: `wild-app-doctor <app-name> [options]`

### 💾 Backup & Restore

**`wild-backup`** - Comprehensive backup system
- Backs up home directory, apps, and cluster resources
- Options: `--home-only`, `--apps-only`, `--cluster-only`
- Uses restic for deduplication
- **Usage**: `wild-backup [options]`

**`wild-app-backup`** - Application-specific backups
- Discovers databases and PVCs automatically
- Supports PostgreSQL and MySQL
- Options: `--all` for all applications
- **Usage**: `wild-app-backup <app-name> [--all]`

**`wild-app-restore`** - Application restore
- Restores databases and/or PVC data
- Options: `--db-only`, `--pvc-only`, `--skip-globals`
- **Usage**: `wild-app-restore <app-name> <snapshot-id> [options]`

### 🔍 Utilities & Helpers

**`wild-health`** - Comprehensive infrastructure validation
- Validates core components (MetalLB, Traefik, CoreDNS)
- Checks installed services (cert-manager, ExternalDNS, Kubernetes Dashboard)
- Tests DNS resolution, routing, certificates, and storage systems
- **Usage**: `wild-health`

**`wild-talos-schema`** - Talos schema management
- Handles configuration schema operations
- **Usage**: `wild-talos-schema [options]`

**`wild-cluster-node-boot-assets-download`** - Download Talos assets
- Downloads installation images for nodes
- **Usage**: `wild-cluster-node-boot-assets-download`

**`wild-dnsmasq-install`** - Install dnsmasq services
- Sets up DNS and DHCP for cluster networking
- **Usage**: `wild-dnsmasq-install`

## Common Usage Patterns

### Complete Setup from Scratch
```bash
wild-init                    # Initialize cloud directory
wild-setup                   # Complete automated setup
# or step by step:
wild-setup-cluster           # Just cluster infrastructure
wild-setup-services          # Just cluster services
```

### Individual Service Management
```bash
# Most common - reconfigure and deploy service
wild-service-setup cert-manager

# Get fresh templates and deploy (for updates)
wild-service-setup cert-manager --fetch

# Configure only, don't deploy (for planning)
wild-service-setup cert-manager --no-deploy

# Fix failed service and resume setup
wild-service-setup cert-manager --fetch
wild-setup-services  # Resume full setup if needed
```

### Application Management
```bash
wild-apps-list              # See available apps
wild-app-add ghost          # Configure app
wild-app-deploy ghost       # Deploy to cluster
wild-app-doctor ghost       # Troubleshoot issues
```

### Configuration Management
```bash
wild-config cluster.name                    # Read values
wild-config-set apps.ghost.domain "blog.example.com"  # Write values
wild-secret apps.ghost.adminPassword       # Read secrets
wild-secret-set apps.ghost.apiKey          # Generate random secret
```

### Cluster Operations
```bash
wild-cluster-node-ip control-1             # Get node IP
wild-dashboard-token                        # Get dashboard access
wild-health                                # Check system health
```

## Script Design Principles

1. **Consistent Interface**: All scripts use `--help` and follow common argument patterns
2. **Error Handling**: All scripts use `set -e` and `set -o pipefail` for robust error handling
3. **Idempotent**: Scripts check existing state before making changes
4. **Template-Driven**: Extensive use of gomplate for configuration flexibility
5. **Environment-Aware**: Scripts source `wild-common.sh` and initialize Wild Cloud environment
6. **Progressive Disclosure**: Complex operations broken into phases with individual controls

## Dependencies Between Scripts

### Setup Phase Dependencies
1. `wild-init` → creates basic structure
2. `wild-setup-cluster` → provisions infrastructure
3. `wild-setup-services` → installs cluster services
4. `wild-setup` → orchestrates all phases

### App Deployment Pipeline
1. `wild-apps-list` → discover applications
2. `wild-app-add` → configure and prepare application
3. `wild-app-deploy` → deploy to cluster

### Node Management Flow
1. `wild-cluster-config-generate` → base configurations
2. `wild-node-setup <node-name>` → atomic node operations (detect → patch → deploy)
   - Internally uses `wild-node-detect` for hardware discovery
   - Generates node-specific patches and final configurations
   - Deploys configuration to target node

All scripts are designed to work together as a cohesive Infrastructure as Code system for personal Kubernetes deployments.