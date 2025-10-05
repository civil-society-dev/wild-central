# Wild Cloud Configuration System

Wild Cloud uses a comprehensive configuration management system that handles both non-sensitive configuration data and sensitive secrets through separate files and commands. The system supports YAML path-based access, template processing, and environment-specific customization.

## Configuration Architecture

### Core Components

1. **`config.yaml`** - Main configuration file for non-sensitive settings
2. **`secrets.yaml`** - Encrypted/protected storage for sensitive data
3. **`.wildcloud/`** - Project marker and cache directory
4. **`env.sh`** - Environment setup and path configuration
5. **Template System** - gomplate-based dynamic configuration processing

### File Structure of a Wild Cloud Project

```
your-cloud-directory/
├── .wildcloud/              # Project marker and cache
│   ├── cache/              # Downloaded templates and temporary files
│   └── logs/               # Operation logs
├── config.yaml             # Main configuration (tracked in git)
├── secrets.yaml            # Sensitive data (NOT tracked in git, 600 perms)
├── env.sh                  # Environment setup (auto-generated)
├── apps/                   # Deployed application configurations
├── setup/                  # Infrastructure setup files
└── docs/                   # Project documentation
```

## Configuration File (`config.yaml`)

### Structure and Organization

The configuration file uses a hierarchical YAML structure for organizing settings:

```yaml
# Cloud-wide settings
cloud:
  domain: "example.com"
  email: "admin@example.com"
  timezone: "America/New_York"

# Cluster infrastructure settings
cluster:
  name: "wild-cluster"
  nodeCount: 3
  network:
    subnet: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dnsServer: "192.168.1.50"
    metallbPool: "192.168.1.80-89"
    controlPlaneVip: "192.168.1.90"
  nodes:
    control-1:
      ip: "192.168.1.91"
      mac: "00:11:22:33:44:55"
      interface: "eth0"
      disk: "/dev/sda"
    control-2:
      ip: "192.168.1.92"
      mac: "00:11:22:33:44:56"
      interface: "eth0"
      disk: "/dev/sda"

# Application-specific settings
apps:
  ghost:
    domain: "blog.example.com"
    image: "ghost:5.0.0"
    storage: "10Gi"
    timezone: "UTC"
    namespace: "ghost"
  immich:
    domain: "photos.example.com"
    serverImage: "ghcr.io/immich-app/immich-server:release"
    storage: "250Gi"
    namespace: "immich"

# Service configurations
services:
  traefik:
    replicas: 2
    dashboard: true
  longhorn:
    defaultReplicas: 3
    storageClass: "longhorn"
```

### Configuration Commands

**Reading Configuration Values**:
```bash
# Read simple values
wild-config cloud.domain                    # "example.com"
wild-config cluster.name                    # "wild-cluster"

# Read nested values
wild-config apps.ghost.domain              # "blog.example.com"
wild-config cluster.nodes.control-1.ip     # "192.168.1.91"

# Check if key exists
wild-config --check apps.newapp.domain     # Returns exit code 0/1
```

**Writing Configuration Values**:
```bash
# Set simple values
wild-config-set cloud.domain "newdomain.com"
wild-config-set cluster.nodeCount 5

# Set nested values
wild-config-set apps.ghost.storage "20Gi"
wild-config-set cluster.nodes.worker-1.ip "192.168.1.94"

# Set complex values (JSON format)
wild-config-set apps.ghost '{"domain":"blog.com","storage":"50Gi"}'
```

### Configuration Sections

#### Cloud Settings (`cloud.*`)
Global settings that affect the entire Wild Cloud deployment:

```yaml
cloud:
  domain: "example.com"           # Primary domain for services
  email: "admin@example.com"      # Contact email for certificates
  timezone: "America/New_York"    # Default timezone for services
  backupLocation: "s3://backup"   # Backup storage location
  monitoring: true                # Enable monitoring services
```

#### Cluster Settings (`cluster.*`)
Infrastructure and node configuration:

```yaml
cluster:
  name: "production-cluster"
  version: "v1.28.0"
  network:
    subnet: "10.0.0.0/16"         # Cluster network range
    serviceCIDR: "10.96.0.0/12"   # Service network range
    podCIDR: "10.244.0.0/16"      # Pod network range
  nodes:
    control-1:
      ip: "10.0.0.10"
      role: "controlplane"
      taints: []
    worker-1:
      ip: "10.0.0.20"
      role: "worker"
      labels:
        node-type: "compute"
```

#### Application Settings (`apps.*`)
Per-application configuration that overrides defaults from app manifests:

```yaml
apps:
  postgresql:
    storage: "100Gi"
    maxConnections: 200
    sharedBuffers: "256MB"
  redis:
    memory: "1Gi"
    persistence: true
  ghost:
    domain: "blog.example.com"
    theme: "casper"
    storage: "10Gi"
    replicas: 2
```

## Secrets Management (`secrets.yaml`)

### Security Model

The `secrets.yaml` file stores all sensitive data with the following security measures:

- **File Permissions**: Automatically set to 600 (owner read/write only)
- **Git Exclusion**: Included in `.gitignore` by default
- **Encryption Support**: Can be encrypted at rest using tools like `age` or `gpg`
- **Access Control**: Only Wild Cloud commands can read/write secrets

### Secret Structure

```yaml
# Generated cluster secrets
cluster:
  talos:
    secrets: "base64-encoded-cluster-secrets"
    adminKey: "talos-admin-private-key"
  kubernetes:
    adminToken: "k8s-admin-service-account-token"

# Application secrets
apps:
  postgresql:
    rootPassword: "randomly-generated-32-char-string"
    replicationPassword: "randomly-generated-32-char-string"
  ghost:
    dbPassword: "randomly-generated-password"
    adminPassword: "user-set-password"
    jwtSecret: "randomly-generated-jwt-secret"
  immich:
    dbPassword: "randomly-generated-password"
    dbUrl: "postgresql://immich:password@postgres:5432/immich"
    jwtSecret: "jwt-signing-key"

# External service credentials
external:
  cloudflare:
    apiToken: "cloudflare-dns-api-token"
  letsencrypt:
    email: "admin@example.com"
  backup:
    s3AccessKey: "backup-s3-access-key"
    s3SecretKey: "backup-s3-secret-key"
```

### Secret Commands

**Reading Secrets**:
```bash
# Read secret values
wild-secret apps.postgresql.rootPassword
wild-secret cluster.kubernetes.adminToken

# Check if secret exists
wild-secret --check apps.newapp.apiKey
```

**Writing Secrets**:
```bash
# Set specific secret value
wild-secret-set apps.ghost.adminPassword "my-secure-password"

# Generate random secret (if no value provided)
wild-secret-set apps.newapp.apiKey        # Generates 32-char base64 string

# Set complex secret (JSON format)
wild-secret-set apps.database '{"user":"admin","password":"secret"}'
```

### Automatic Secret Generation

When you run `wild-app-add`, Wild Cloud automatically generates required secrets:

1. **Reads App Manifest**: Identifies `requiredSecrets` list
2. **Checks Existing Secrets**: Never overwrites existing values
3. **Generates Missing Secrets**: Creates secure random values
4. **Updates secrets.yaml**: Adds new secrets with proper structure

**Example App Manifest**:
```yaml
name: ghost
requiredSecrets:
  - apps.ghost.dbPassword      # Auto-generated if missing
  - apps.ghost.jwtSecret       # Auto-generated if missing
  - apps.postgresql.password   # Auto-generated if missing (dependency)
```

**Resulting secrets.yaml**:
```yaml
apps:
  ghost:
    dbPassword: "aB3kL9mN2pQ7rS8tU1vW4xY5zA6bC0dE"
    jwtSecret: "jF2gH5iJ8kL1mN4oP7qR0sT3uV6wX9yZ"
  postgresql:
    password: "eE8fF1gG4hH7iI0jJ3kK6lL9mM2nN5oO"
```

## Template System

### gomplate Integration

Wild Cloud uses [gomplate](https://gomplate.ca/) for dynamic configuration processing, allowing templates to access both configuration and secrets:

```yaml
# Template example (before processing)
apiVersion: v1
kind: ConfigMap
metadata:
  name: ghost-config
  namespace: {{ .apps.ghost.namespace }}
data:
  url: "https://{{ .apps.ghost.domain }}"
  timezone: "{{ .apps.ghost.timezone | default .cloud.timezone }}"
  database_host: "{{ .apps.postgresql.hostname }}"
  # Conditionals
  {{- if .apps.ghost.enableSSL }}
  ssl_enabled: "true"
  {{- end }}
  # Loops
  allowed_domains: |
    {{- range .apps.ghost.allowedDomains }}
    - {{ . }}
    {{- end }}
```

### Template Processing Commands

**Process Single Template**:
```bash
# From stdin
cat template.yaml | wild-compile-template > output.yaml

# With custom context
echo "domain: {{ .cloud.domain }}" | wild-compile-template
```

**Process Template Directory**:
```bash
# Recursively process all templates
wild-compile-template-dir source-dir output-dir

# Clean destination first
wild-compile-template-dir --clean source-dir output-dir
```

### Template Context

Templates have access to the complete configuration and secrets context:

```go
// Available template variables
.cloud.*          // All cloud configuration
.cluster.*        // All cluster configuration
.apps.*           // All application configuration
.services.*       // All service configuration

// Special functions
.cloud.domain             // Primary domain
default "fallback"        // Default value if key missing
env "VAR_NAME"           // Environment variable
file "path/to/file"      // File contents
```

**Template Examples**:
```yaml
# Basic variable substitution
domain: {{ .apps.myapp.domain }}

# Default values
timezone: {{ .apps.myapp.timezone | default .cloud.timezone }}

# Conditionals
{{- if .apps.myapp.enableFeature }}
feature_enabled: true
{{- else }}
feature_enabled: false
{{- end }}

# Lists and iteration
allowed_hosts:
{{- range .apps.myapp.allowedHosts }}
  - {{ . }}
{{- end }}

# Complex expressions
replicas: {{ if eq .cluster.environment "production" }}3{{ else }}1{{ end }}
```

## Environment Setup

### Environment Detection

Wild Cloud automatically detects and configures the environment through several mechanisms:

**Project Detection**:
- Searches for `.wildcloud` directory in current or parent directories
- Sets `WC_HOME` to the directory containing `.wildcloud`
- Fails if no Wild Cloud project found

**Repository Detection**:
- Locates Wild Cloud repository (source code)
- Sets `WC_ROOT` to repository location
- Used for accessing app templates and setup scripts

### Environment Variables

**Key Environment Variables**:
```bash
WC_HOME="/path/to/your-cloud"           # Your cloud directory
WC_ROOT="/path/to/wild-cloud-repo"      # Wild Cloud repository
PATH="$WC_ROOT/bin:$PATH"               # Wild Cloud commands available
KUBECONFIG="$WC_HOME/.kube/config"      # Kubernetes configuration
TALOSCONFIG="$WC_HOME/.talos/config"    # Talos configuration
```

**Environment Setup Script** (`env.sh`):
```bash
#!/bin/bash
# Auto-generated environment setup

export WC_HOME="/home/user/my-cloud"
export WC_ROOT="/opt/wild-cloud"
export PATH="$WC_ROOT/bin:$PATH"
export KUBECONFIG="$WC_HOME/.kubeconfig"
export TALOSCONFIG="$WC_HOME/setup/cluster-nodes/generated/talosconfig"

# Source this file to set up Wild Cloud environment
# source env.sh
```

### Common Script Pattern

Most Wild Cloud scripts follow this initialization pattern:

```bash
#!/bin/bash
set -e
set -o pipefail

# Initialize Wild Cloud environment
if [ -z "${WC_ROOT}" ]; then
    print "WC_ROOT is not set."
    exit 1
else
    source "${WC_ROOT}/scripts/common.sh"
    init_wild_env
fi

# Script logic here...
```

## Configuration Validation

### Schema Validation

Wild Cloud validates configuration against expected schemas:

**Cluster Configuration Validation**:
- Node IP addresses are valid and unique
- Network ranges don't overlap
- Required fields are present
- Hardware specifications meet minimums

**Application Configuration Validation**:
- Domain names are valid DNS names
- Storage sizes use valid Kubernetes formats
- Image references are valid container images
- Dependencies are satisfied

### Validation Commands

```bash
# Validate current configuration
wild-config --validate

# Check specific configuration sections
wild-config --validate --section cluster
wild-config --validate --section apps.ghost

# Test template compilation
wild-compile-template --validate < template.yaml
```

## Configuration Best Practices

### Organization

**Hierarchical Structure**:
- Group related settings under common prefixes
- Use consistent naming conventions
- Keep application configs under `apps.*`
- Separate infrastructure from application settings

**Documentation**:
```yaml
# Document complex configurations
cluster:
  # Node configuration - update IPs after hardware changes
  nodes:
    control-1:
      ip: "192.168.1.91"    # Main control plane node
      interface: "eth0"      # Primary network interface
```

### Security

**Configuration Security**:
- Never store secrets in `config.yaml`
- Use `wild-secret-set` for all sensitive data
- Regularly rotate generated secrets
- Backup `secrets.yaml` securely

**Access Control**:
```bash
# Ensure proper permissions
chmod 600 secrets.yaml
chmod 644 config.yaml

# Restrict directory access
chmod 755 your-cloud-directory
chmod 700 .wildcloud/
```

### Version Control

**Git Integration**:
```gitignore
# .gitignore for Wild Cloud projects
secrets.yaml          # Never commit secrets
.wildcloud/cache/      # Temporary files
.wildcloud/logs/       # Operation logs
setup/cluster-nodes/generated/  # Generated cluster configs
.kube/                 # Kubernetes configs
.talos/                # Talos configs
```

**Configuration Changes**:
- Commit `config.yaml` changes with descriptive messages
- Tag major configuration changes
- Use branches for experimental configurations
- Document configuration changes in commit messages

### Backup and Recovery

**Configuration Backup**:
```bash
# Backup configuration and secrets
wild-backup --home-only

# Export configuration for disaster recovery
cp config.yaml config-backup-$(date +%Y%m%d).yaml
cp secrets.yaml secrets-backup-$(date +%Y%m%d).yaml.gpg  # Encrypt first
```

**Recovery Process**:
1. Restore `config.yaml` from backup
2. Decrypt and restore `secrets.yaml`
3. Re-run `wild-setup` if needed
4. Validate configuration with `wild-config --validate`

## Advanced Configuration

### Multi-Environment Setup

**Development Environment**:
```yaml
cloud:
  domain: "dev.example.com"
cluster:
  name: "dev-cluster"
  nodeCount: 1
apps:
  ghost:
    domain: "blog.dev.example.com"
    replicas: 1
```

**Production Environment**:
```yaml
cloud:
  domain: "example.com"
cluster:
  name: "prod-cluster"
  nodeCount: 5
apps:
  ghost:
    domain: "blog.example.com"
    replicas: 3
```

### Configuration Inheritance

**Base Configuration**:
```yaml
# config.base.yaml
cloud:
  timezone: "UTC"
  email: "admin@example.com"
apps:
  postgresql:
    storage: "10Gi"
```

**Environment-Specific Override**:
```yaml
# config.prod.yaml (merged with base)
apps:
  postgresql:
    storage: "100Gi"    # Override for production
    replicas: 3         # Additional production setting
```

### Dynamic Configuration

**Runtime Configuration Updates**:
```bash
# Update configuration without restart
wild-config-set apps.ghost.replicas 3
wild-app-deploy ghost  # Apply changes

# Rolling updates
wild-config-set apps.ghost.image "ghost:5.1.0"
wild-app-deploy ghost --rolling-update
```

The Wild Cloud configuration system provides a powerful, secure, and flexible foundation for managing complex infrastructure deployments while maintaining simplicity for common use cases.