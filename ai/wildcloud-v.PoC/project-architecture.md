# Wild Cloud Project Architecture

Wild Cloud consists of two main directory structures: the **Wild Cloud Repository** (source code and templates) and **User Cloud Directories** (individual deployments). Understanding this architecture is essential for working with Wild Cloud effectively.

## Architecture Overview

```
Wild Cloud Repository (/path/to/wild-cloud-repo)  ←  Source code, templates, scripts
                    ↓
User Cloud Directory (/path/to/my-cloud)          ←  Individual deployment instance
                    ↓
Kubernetes Cluster                                ←  Running infrastructure
```

## Wild Cloud Repository Structure

The Wild Cloud repository (`WC_ROOT`) contains the source code, templates, and tools:

### `/bin/` - Command Line Interface
**Purpose**: All Wild Cloud CLI commands and utilities
```
bin/
├── wild-*                    # All user-facing commands (34+ scripts)
├── wild-common.sh            # Common utilities and functions
├── README.md                 # CLI documentation
└── helm-chart-to-kustomize   # Utility for converting Helm charts
```

**Key Commands**:
- **Setup**: `wild-init`, `wild-setup`, `wild-setup-cluster`, `wild-setup-services`
- **Apps**: `wild-app-*`, `wild-apps-list`
- **Config**: `wild-config*`, `wild-secret*`
- **Operations**: `wild-backup`, `wild-health`, `wild-dashboard-token`

### `/apps/` - Application Templates
**Purpose**: Pre-built applications ready for deployment
```
apps/
├── README.md                 # Apps system documentation
├── ghost/                    # Blog publishing platform
│   ├── manifest.yaml         # App metadata and defaults
│   ├── kustomization.yaml    # Kustomize configuration
│   ├── deployment.yaml       # Kubernetes deployment
│   ├── service.yaml          # Service definition
│   ├── ingress.yaml          # HTTPS ingress
│   └── ...
├── immich/                   # Photo management
├── gitea/                    # Git hosting
├── postgresql/               # Database service
├── vllm/                     # AI/LLM inference
└── ...
```

**Application Categories**:
- **Content Management**: Ghost, Discourse
- **Media**: Immich
- **Development**: Gitea, Docker Registry
- **Databases**: PostgreSQL, MySQL, Redis
- **AI/ML**: vLLM
- **Infrastructure**: Memcached, NFS

### `/setup/` - Infrastructure Templates
**Purpose**: Cluster and service deployment templates
```
setup/
├── README.md
├── cluster-nodes/            # Talos node configuration
│   ├── init-cluster.sh       # Cluster initialization script
│   ├── patch.templates/      # Node-specific config templates
│   │   ├── controlplane.yaml # Control plane template
│   │   └── worker.yaml       # Worker node template
│   └── talos-schemas.yaml    # Version mappings
├── cluster-services/         # Core Kubernetes services
│   ├── README.md
│   ├── metallb/              # Load balancer
│   ├── traefik/              # Ingress controller
│   ├── cert-manager/         # Certificate management
│   ├── longhorn/             # Distributed storage
│   ├── coredns/              # DNS resolution
│   ├── externaldns/          # DNS record management
│   ├── kubernetes-dashboard/ # Web UI
│   └── ...
├── dnsmasq/                  # DNS and PXE boot server
├── home-scaffold/            # User directory templates
└── operator/                 # Additional operator tools
```

### `/experimental/` - Development Projects
**Purpose**: Experimental features and development tools
```
experimental/
├── daemon/                   # Go API daemon
│   ├── main.go               # API server
│   ├── Makefile              # Build automation
│   └── README.md
└── app/                      # React dashboard
    ├── src/                  # React source code
    ├── package.json          # Dependencies
    ├── pnpm-lock.yaml        # Lock file
    └── README.md
```

### `/scripts/` - Utility Scripts
**Purpose**: Installation and utility scripts
```
scripts/
├── setup-utils.sh            # Install dependencies
└── install-wild-cloud-dependencies.sh
```

### `/docs/` - Documentation
**Purpose**: User guides and documentation
```
docs/
├── guides/                   # Setup and usage guides
├── agent-context/           # Agent documentation
│   └── wildcloud/           # Context files for AI agents
└── *.md                     # Various documentation files
```

### `/test/` - Test Suite
**Purpose**: Automated testing with Bats
```
test/
├── bats/                     # Bats testing framework
├── fixtures/                 # Test data and configurations
├── run_bats_tests.sh         # Test runner
└── *.bats                    # Individual test files
```

### Root Files
```
/
├── README.md                 # Project overview
├── CLAUDE.md                 # AI assistant context
├── LICENSE                   # GNU AGPLv3
├── CONTRIBUTING.md           # Contribution guidelines
├── env.sh                    # Environment setup
├── .gitignore               # Git exclusions
└── .gitmodules              # Git submodules
```

## User Cloud Directory Structure

Each user deployment (`WC_HOME`) is an independent cloud instance:

### Directory Layout
```
my-cloud/                     # User's cloud directory
├── .wildcloud/               # Project marker and cache
│   ├── cache/                # Downloaded templates
│   │   ├── apps/             # Cached app templates
│   │   └── services/         # Cached service templates
│   └── logs/                 # Operation logs
├── config.yaml               # Main configuration
├── secrets.yaml              # Sensitive data (600 permissions)
├── env.sh                    # Environment setup (auto-generated)
├── apps/                     # Deployed application configs
│   ├── ghost/                # Compiled ghost configuration
│   ├── postgresql/           # Database configuration
│   └── ...
├── setup/                    # Infrastructure configurations
│   ├── cluster-nodes/        # Node-specific configurations
│   │   └── generated/        # Generated Talos configs
│   └── cluster-services/     # Compiled service configurations
├── docs/                     # Project-specific documentation
├── .kube/                    # Kubernetes configuration
│   └── config                # kubectl configuration
├── .talos/                   # Talos configuration
│   └── config                # talosctl configuration
└── backups/                  # Local backup staging
```

### Configuration Files

**`config.yaml`** - Main configuration (version controlled):
```yaml
cloud:
  domain: "example.com"
  email: "admin@example.com"
cluster:
  name: "my-cluster"
  nodeCount: 3
apps:
  ghost:
    domain: "blog.example.com"
```

**`secrets.yaml`** - Sensitive data (not version controlled):
```yaml
apps:
  ghost:
    dbPassword: "generated-password"
  postgresql:
    rootPassword: "generated-password"
cluster:
  talos:
    secrets: "base64-encoded-secrets"
```

**`.wildcloud/`** - Project metadata:
- Marks directory as Wild Cloud project
- Contains cached templates and temporary files
- Used for project detection by scripts

### Generated Directories

**`apps/`** - Compiled application configurations:
- Created by `wild-app-add` command
- Contains ready-to-deploy Kubernetes manifests
- Templates processed with user configuration
- Each app in separate subdirectory

**`setup/cluster-nodes/generated/`** - Talos configurations:
- Base cluster configuration (`controlplane.yaml`, `worker.yaml`)
- Node-specific patches and final configs
- Cluster secrets and certificates
- Generated by `wild-cluster-config-generate`

**`setup/cluster-services/`** - Kubernetes services:
- Compiled service configurations
- Generated by `wild-cluster-services-configure`
- Ready for deployment to cluster

## Template Processing Flow

### From Repository to Deployment

1. **Template Storage**: Templates stored in repository with placeholder variables
2. **Configuration Merge**: `wild-app-add` reads templates directly from repository and merges app defaults with user config
3. **Template Compilation**: gomplate processes templates with user data
4. **Manifest Generation**: Final Kubernetes manifests created in user directory
5. **Deployment**: `wild-app-deploy` applies manifests to cluster

### Template Variables

**Repository Templates** (before processing):
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ghost
  namespace: {{ .apps.ghost.namespace }}
spec:
  replicas: {{ .apps.ghost.replicas | default 1 }}
  template:
    spec:
      containers:
      - name: ghost
        image: "{{ .apps.ghost.image }}"
        env:
        - name: url
          value: "https://{{ .apps.ghost.domain }}"
```

**User Directory** (after processing):
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ghost
  namespace: ghost
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: ghost
        image: "ghost:5.0.0"
        env:
        - name: url
          value: "https://blog.example.com"
```

## File Permissions and Security

### Security Model

**Configuration Security**:
```bash
config.yaml         # 644 (readable by group)
secrets.yaml        # 600 (owner only)
.wildcloud/         # 755 (standard directory)
apps/               # 755 (standard directory)
```

**Git Integration**:
```gitignore
# Automatically excluded from version control
secrets.yaml                    # Never commit secrets
.wildcloud/cache/              # Temporary files
.wildcloud/logs/               # Operation logs
setup/cluster-nodes/generated/ # Generated configs
.kube/                         # Kubernetes configs
.talos/                        # Talos configs
backups/                       # Backup files
```

### Access Patterns

**Read Operations**:
- Scripts read config and secrets via `wild-config` and `wild-secret`
- Template processor accesses both files for compilation
- Kubernetes tools read generated manifests

**Write Operations**:
- Only Wild Cloud commands modify config and secrets
- Manual editing supported but not recommended
- Backup processes create read-only copies

## Development Workflow

### Repository Development

**Setup Development Environment**:
```bash
git clone https://github.com/username/wild-cloud
cd wild-cloud
source env.sh                   # Set up environment
scripts/setup-utils.sh          # Install dependencies
```

**Testing Changes**:
```bash
# Test specific functionality
test/run_bats_tests.sh

# Test with real cloud directory
cd /path/to/test-cloud
wild-app-add myapp              # Test app changes
wild-setup-cluster --dry-run    # Test cluster changes
```

### User Workflow

**Initial Setup**:
```bash
mkdir my-cloud && cd my-cloud
wild-init                       # Initialize project
wild-setup                      # Deploy infrastructure
```

**Daily Operations**:
```bash
wild-apps-list                  # Browse available apps
wild-app-add ghost              # Configure app
wild-app-deploy ghost           # Deploy to cluster
```

**Configuration Management**:
```bash
wild-config apps.ghost.domain   # Read configuration
wild-config-set apps.ghost.storage "20Gi"  # Update configuration
wild-app-deploy ghost           # Apply changes
```

## Integration Points

### External Systems

**DNS Providers**:
- Cloudflare API for DNS record management
- Route53 support for AWS domains
- Generic webhook support for other providers

**Certificate Authorities**:
- Let's Encrypt (primary)
- Custom CA support
- Manual certificate import

**Storage Backends**:
- Local storage via Longhorn
- NFS network storage
- Cloud storage integration (S3, etc.)

**Backup Systems**:
- Restic for deduplication and encryption
- S3-compatible storage backends
- Local and remote backup targets

### Kubernetes Integration

**Custom Resources**:
- Traefik IngressRoute and Middleware
- cert-manager Certificate and Issuer
- Longhorn Volume and Engine
- ExternalDNS DNSEndpoint

**Standard Resources**:
- Deployments, Services, ConfigMaps
- Ingress, PersistentVolumes, Secrets
- NetworkPolicies, ServiceAccounts
- Jobs, CronJobs, DaemonSets

## Extensibility Points

### Custom Applications

**Create New Apps**:
1. Create directory in `apps/`
2. Define `manifest.yaml` with metadata
3. Create Kubernetes resource templates
4. Test with `wild-app-add` and `wild-app-deploy`

**App Requirements**:
- Follow Wild Cloud labeling standards
- Use gomplate template syntax
- Include external-dns annotations
- Implement proper security contexts

### Custom Services

**Add Infrastructure Services**:
1. Create directory in `setup/cluster-services/`
2. Define installation and configuration scripts
3. Create Kubernetes manifests with templates
4. Integrate with service deployment pipeline

### Script Extensions

**Extend CLI**:
- Add scripts to `bin/` directory with `wild-` prefix
- Follow common script patterns (error handling, help text)
- Source `wild-common.sh` for utilities
- Use configuration system for customization

## Deployment Patterns

### Single-Node Development

**Configuration**:
```yaml
cluster:
  nodeCount: 1
  nodes:
    all-in-one:
      roles: ["controlplane", "worker"]
```

**Suitable For**:
- Development and testing
- Learning Kubernetes concepts
- Small personal deployments
- Resource-constrained environments

### Multi-Node Production

**Configuration**:
```yaml
cluster:
  nodeCount: 5
  nodes:
    control-1: { role: "controlplane" }
    control-2: { role: "controlplane" }
    control-3: { role: "controlplane" }
    worker-1: { role: "worker" }
    worker-2: { role: "worker" }
```

**Suitable For**:
- Production workloads
- High availability requirements
- Scalable application hosting
- Enterprise-grade deployments

### Hybrid Deployments

**Configuration**:
```yaml
cluster:
  nodes:
    control-1:
      role: "controlplane"
      taints: []                # Allow workloads on control plane
    worker-gpu:
      role: "worker"
      labels:
        nvidia.com/gpu: "true"  # GPU-enabled node
```

**Use Cases**:
- Mixed workload requirements
- Specialized hardware (GPU, storage)
- Cost optimization
- Gradual scaling

The Wild Cloud architecture provides a solid foundation for personal cloud infrastructure while maintaining flexibility for customization and extension.