# Wild Cloud Apps System

The Wild Cloud apps system provides a streamlined way to deploy and manage applications on your Kubernetes cluster. It uses Kustomize for configuration management and follows a standardized structure for consistent deployment patterns.

## App Structure and Components

### Directory Structure
Each subdirectory represents a Wild Cloud app. Each app directory contains:

**Required Files:**
- `manifest.yaml` - App metadata and configuration
- `kustomization.yaml` - Kustomize configuration with Wild Cloud labels

**Standard Configuration Files (one or more YAML files containing Kubernetes resource definitions):**
```
apps/myapp/
├── manifest.yaml          # Required: App metadata and configuration
├── kustomization.yaml     # Required: Kustomize configuration with Wild Cloud labels
├── namespace.yaml         # Kubernetes namespace definition
├── deployment.yaml        # Application deployment
├── service.yaml           # Kubernetes service definition
├── ingress.yaml           # HTTPS ingress with external DNS
├── pvc.yaml               # Persistent volume claims (if needed)
├── db-init-job.yaml       # Database initialization (if needed)
└── configmap.yaml         # Configuration data (if needed)
```

### App Manifest (`manifest.yaml`)

The required `manifest.yaml` file contains metadata about the app. Here's an example `manifest.yaml` file:

```yaml
name: myapp
description: A brief description of the application and its purpose.
version: 1.0.0
icon: https://example.com/icon.png
requires:
  - name: postgres
defaultConfig:
  image: myapp/server:1.0.0
  domain: myapp.{{ .cloud.domain }}
  timezone: UTC
  storage: 10Gi
  dbHostname: postgres.postgres.svc.cluster.local
  dbUsername: myapp
requiredSecrets:
  - apps.myapp.dbPassword
  - apps.postgres.password
```

**Manifest Fields**:
- `name` - The name of the app, used for identification (must match directory name)
- `description` - A brief description of the app
- `version` - The version of the app (should generally follow the versioning scheme of the app itself)
- `icon` - A URL to an icon representing the app
- `requires` - A list of other apps that this app depends on (each entry should be the name of another app)
- `defaultConfig` - A set of default configuration values for the app (when an app is added using `wild-app-add`, these values will be added to the Wild Cloud `config.yaml` file)
- `requiredSecrets` - A list of secrets that must be set in the Wild Cloud `secrets.yaml` file for the app to function properly (these secrets are typically sensitive information like database passwords or API keys; keys with random values will be generated automatically when the app is added)

### Kustomization Configuration

Wild Cloud apps use standard Kustomize with required Wild Cloud labels:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: myapp
labels:
  - includeSelectors: true
    pairs:
      app: myapp
      managedBy: kustomize
      partOf: wild-cloud
resources:
  - namespace.yaml
  - deployment.yaml
  - service.yaml
  - ingress.yaml
  - pvc.yaml
  - db-init-job.yaml
```

**Kustomization Requirements**:
- Every Wild Cloud kustomization should include the Wild Cloud labels in its `kustomization.yaml` file (this allows Wild Cloud to identify and manage the app correctly)
- The `app` label and `namespace` should match the app's name/directory
- **includeSelectors: true** - Automatically applies labels to all resources AND their selectors

#### Standard Wild Cloud Labels

Wild Cloud uses a consistent labeling strategy across all apps:

```yaml
labels:
  - includeSelectors: true
    pairs:
      app: myapp              # The app name (matches directory)
      managedBy: kustomize    # Managed by Kustomize
      partOf: wild-cloud      # Part of Wild Cloud ecosystem
```

The `includeSelectors: true` setting automatically applies these labels to all resources AND their selectors, which means:

1. **Resource labels** - All resources get the standard Wild Cloud labels
2. **Selector labels** - All selectors automatically include these labels for robust selection

This allows individual resources to use simple, component-specific selectors:

```yaml
selector:
  matchLabels:
    component: web
```

Which Kustomize automatically expands to:

```yaml
selector:
  matchLabels:
    app: myapp
    component: web
    managedBy: kustomize
    partOf: wild-cloud
```

### Template System

Wild Cloud apps are actually **templates** that get compiled with your specific configuration when you run `wild-app-add`. This allows for:

- **Dynamic Configuration** - Reference user settings via `{{ .apps.appname.key }}`
- **Gomplate Processing** - Full template capabilities including conditionals and loops
- **Secret Integration** - Automatic secret generation and referencing
- **Domain Management** - Automatic subdomain assignment based on your domain

**Template Variable Examples**:
```yaml
# Configuration references
image: "{{ .apps.myapp.image }}"
domain: "{{ .apps.myapp.domain }}"
namespace: "{{ .apps.myapp.namespace }}"

# Cloud-wide settings
timezone: "{{ .cloud.timezone }}"
domain_suffix: "{{ .cloud.domain }}"

# Conditional logic
{{- if .apps.myapp.enableSSL }}
- name: ENABLE_SSL
  value: "true"
{{- end }}
```

## App Lifecycle Management

### 1. Discovery Phase
**Command**: `wild-apps-list`

Lists all available applications with metadata:
```bash
wild-apps-list --verbose    # Detailed view with descriptions
wild-apps-list --json       # JSON output for automation
```

Shows:
- App name and description
- Version and dependencies
- Installation status
- Required configuration

### 2. Configuration Phase
**Command**: `wild-app-add <app-name>`

Processes app templates and prepares for deployment:

**What it does**:
1. Reads app manifest directly from Wild Cloud repository
2. Merges default configuration with existing `config.yaml`
3. Generates required secrets automatically
4. Compiles templates with gomplate using your configuration
5. Creates ready-to-deploy Kustomize files in `apps/<app-name>/`

**Generated Files**:
- Compiled Kubernetes manifests (no more template variables)
- Standard Kustomize configuration
- App-specific configuration merged into your `config.yaml`
- Required secrets added to your `secrets.yaml`

### 3. Deployment Phase
**Command**: `wild-app-deploy <app-name>`

Deploys the app to your Kubernetes cluster:

**Deployment Process**:
1. Creates namespace if it doesn't exist
2. Handles app dependencies (deploys required apps first)
3. Creates secrets from your `secrets.yaml`
4. Applies Kustomize configuration to cluster
5. Copies TLS certificates to app namespace
6. Validates deployment success

**Options**:
- `--force` - Overwrite existing resources
- `--dry-run` - Preview changes without applying

### 4. Operations Phase

**Monitoring**: `wild-app-doctor <app-name>`
- Runs app-specific diagnostic tests
- Checks pod status, resource usage, connectivity
- Options: `--keep`, `--follow`, `--timeout`

**Updates**: Re-run `wild-app-add` then `wild-app-deploy`
- Use `--force` flag to overwrite existing configuration
- Updates configuration changes
- Handles image updates
- Preserves persistent data

**Removal**: `wild-app-delete <app-name>`
- Deletes namespace and all resources
- Removes local configuration files
- Options: `--force` for no confirmation

## Configuration System

### Configuration Storage

**Global Configuration** (`config.yaml`):
```yaml
cloud:
  domain: example.com
  timezone: America/New_York
apps:
  myapp:
    domain: app.example.com
    image: myapp:1.0.0
    storage: 20Gi
    timezone: UTC
```

**Secrets Management** (`secrets.yaml`):
```yaml
apps:
  myapp:
    dbPassword: "randomly-generated-password"
    adminPassword: "user-set-password"
  postgres:
    password: "randomly-generated-password"
```

### Secret Generation

When you run `wild-app-add`, required secrets are automatically generated:
- **Random Generation**: 32-character base64 strings for passwords/keys
- **User Prompts**: For secrets that need specific values
- **Preservation**: Existing secrets are never overwritten
- **Permissions**: `secrets.yaml` has 600 permissions (owner-only)

### Configuration Commands
```bash
# Read app configuration
wild-config apps.myapp.domain

# Set app configuration
wild-config-set apps.myapp.storage "50Gi"

# Read app secrets
wild-secret apps.myapp.dbPassword

# Set app secrets
wild-secret-set apps.myapp.adminPassword "my-secure-password"
```

## Networking and DNS

### External DNS Integration

Wild Cloud apps automatically manage DNS records through ingress annotations:

```yaml
metadata:
  annotations:
    external-dns.alpha.kubernetes.io/target: {{ .cloud.domain }}
    external-dns.alpha.kubernetes.io/cloudflare-proxied: "false"
```

**How it works**:
1. App ingress created with external-dns annotations
2. ExternalDNS controller detects new ingress
3. Creates CNAME record: `app.yourdomain.com` → `yourdomain.com`
4. DNS resolves to MetalLB load balancer IP
5. Traefik routes traffic to appropriate service

### HTTPS Certificate Management

Automatic TLS certificates via cert-manager:

```yaml
metadata:
  annotations:
    traefik.ingress.kubernetes.io/router.tls: "true"
    traefik.ingress.kubernetes.io/router.tls.certresolver: letsencrypt
spec:
  tls:
    - hosts:
        - {{ .apps.myapp.domain }}
      secretName: myapp-tls
```

**Certificate Lifecycle**:
1. Ingress created with TLS configuration
2. cert-manager detects certificate requirement
3. Let's Encrypt challenge initiated automatically
4. Certificate issued and stored in Kubernetes secret
5. Traefik uses certificate for TLS termination
6. Automatic renewal before expiration

## Database Integration

### Database Initialization Jobs

Apps that require databases use initialization jobs to set up the database before the main application starts:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: myapp-db-init
spec:
  template:
    spec:
      containers:
      - name: db-init
        image: postgres:15
        command:
          - /bin/bash
          - -c
          - |
            PGPASSWORD=$ROOT_PASSWORD psql -h $DB_HOST -U postgres -c "
              CREATE DATABASE IF NOT EXISTS $DB_NAME;
              CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
              GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
            "
        env:
        - name: DB_HOST
          value: {{ .apps.myapp.dbHostname }}
        - name: ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: myapp-secrets
              key: apps.postgres.password
      restartPolicy: OnFailure
```

**Database URL Secrets**: For apps requiring database URLs with embedded credentials, always use dedicated secrets:

```yaml
# In manifest.yaml
requiredSecrets:
  - apps.myapp.dbUrl

# Generated secret (by wild-app-add)
apps:
  myapp:
    dbUrl: "postgresql://myapp:password123@postgres.postgres.svc.cluster.local/myapp"
```

### Supported Databases

Wild Cloud apps commonly integrate with:
- **PostgreSQL** - Via `postgres` app dependency
- **MySQL** - Via `mysql` app dependency
- **Redis** - Via `redis` app dependency
- **SQLite** - For apps with embedded database needs

## Storage Management

### Persistent Volume Claims

Apps requiring persistent storage define PVCs:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: myapp-data
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: longhorn
  resources:
    requests:
      storage: {{ .apps.myapp.storage }}
```

**Storage Integration**:
- **Longhorn Storage Class** - Distributed, replicated storage
- **Dynamic Provisioning** - Automatic volume creation
- **Backup Support** - Via `wild-app-backup` command
- **Expansion** - Update storage size in configuration

### Backup and Restore

**Application Backup**: `wild-app-backup <app-name>`
- Discovers databases and PVCs automatically
- Creates restic snapshots with deduplication
- Supports PostgreSQL and MySQL database backups
- Streams PVC data for efficient storage

**Application Restore**: `wild-app-restore <app-name> <snapshot-id>`
- Restores from restic snapshots
- Options: `--db-only`, `--pvc-only`, `--skip-globals`
- Creates safety snapshots before destructive operations

## Security Considerations

### Pod Security Standards

All Wild Cloud apps comply with Pod Security Standards:

```yaml
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 999
        runAsGroup: 999
        seccompProfile:
          type: RuntimeDefault
      containers:
      - name: app
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          readOnlyRootFilesystem: false  # Set to true when possible
```

### Secret Management

- **Kubernetes Secrets** - All sensitive data stored as Kubernetes secrets
- **Secret References** - Apps reference secrets via `secretKeyRef`, never inline
- **Full Dotted Paths** - Always use complete secret paths (e.g., `apps.myapp.dbPassword`)
- **No Plaintext** - Secrets never stored in manifests or config files

### Network Policies

Apps can define network policies for traffic isolation:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: myapp-network-policy
spec:
  podSelector:
    matchLabels:
      app: myapp
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: traefik
```

## Available Applications

Wild Cloud includes apps for common self-hosted services:

### Content Management
- **Ghost** - Publishing platform for blogs and websites
- **Discourse** - Community discussion platform

### Development & Project Management Tools
- **Gitea** - Self-hosted Git service with web interface
- **OpenProject** - Open-source project management software
- **Docker Registry** - Private container image registry

### Media & File Management
- **Immich** - Self-hosted photo and video backup solution

### Communication
- **Keila** - Newsletter and email marketing platform
- **Listmonk** - Newsletter and mailing list manager

### Databases
- **PostgreSQL** - Relational database service
- **MySQL** - Relational database service
- **Redis** - In-memory data structure store
- **Memcached** - Distributed memory caching system

### AI/ML
- **vLLM** - Fast LLM inference server with OpenAI-compatible API

### Examples & Templates
- **example-admin** - Example admin interface application
- **example-app** - Template application for development reference

## Creating Custom Apps

### App Development Process

1. **Create Directory**: `apps/myapp/`
2. **Write Manifest**: Define metadata and configuration
3. **Create Resources**: Kubernetes manifests with templates
4. **Test Locally**: Use `wild-app-add` and `wild-app-deploy`
5. **Validate**: Ensure all resources deploy correctly

### Best Practices

**Manifest Design**:
- Include comprehensive `defaultConfig` for all configurable values
- List all `requiredSecrets` the app needs
- Specify dependencies in `requires` field
- Use semantic versioning

**Template Usage**:
- Reference configuration via `{{ .apps.myapp.key }}`
- Use conditionals for optional features
- Include proper gomplate syntax for lists and objects
- Test template compilation

**Resource Configuration**:
- Always include Wild Cloud standard labels
- Use appropriate security contexts
- Define resource requests and limits
- Include health checks and probes

**Storage and Networking**:
- Use Longhorn storage class for persistence
- Include external-dns annotations for automatic DNS
- Configure TLS certificates via cert-manager annotations
- Follow database initialization patterns for data apps

### Converting from Helm Charts

Wild Cloud provides tooling to convert Helm charts to Wild Cloud apps:

```bash
# Convert Helm chart to Kustomize base
helm fetch --untar --untardir charts stable/mysql
helm template --output-dir base --namespace mysql mysql charts/mysql
cd base/mysql
kustomize create --autodetect

# Then customize for Wild Cloud:
# 1. Add manifest.yaml
# 2. Replace hardcoded values with templates
# 3. Update labels to Wild Cloud standard
# 4. Configure secrets properly
```

## Troubleshooting Applications

### Common Issues

**App Won't Start**:
- Check pod logs: `kubectl logs -n <app-namespace> deployment/<app-name>`
- Verify secrets exist: `kubectl get secrets -n <app-namespace>`
- Check resource constraints: `kubectl describe pod -n <app-namespace>`

**Database Connection Issues**:
- Verify database is running: `kubectl get pods -n <db-namespace>`
- Check database initialization job: `kubectl logs job/<app>-db-init -n <app-namespace>`
- Validate database credentials in secrets

**DNS/Certificate Issues**:
- Check ingress status: `kubectl get ingress -n <app-namespace>`
- Verify certificate creation: `kubectl get certificates -n <app-namespace>`
- Check external-dns logs: `kubectl logs -n external-dns deployment/external-dns`

**Storage Issues**:
- Check PVC status: `kubectl get pvc -n <app-namespace>`
- Verify Longhorn cluster health: Access Longhorn UI
- Check storage class availability: `kubectl get storageclass`

### Diagnostic Tools

```bash
# App-specific diagnostics
wild-app-doctor <app-name>

# Resource inspection
kubectl get all -n <app-namespace>
kubectl describe deployment/<app-name> -n <app-namespace>

# Log analysis
kubectl logs -f deployment/<app-name> -n <app-namespace>
kubectl logs job/<app>-db-init -n <app-namespace>

# Configuration verification
wild-config apps.<app-name>
wild-secret apps.<app-name>
```

The Wild Cloud apps system provides a powerful, consistent way to deploy and manage self-hosted applications with enterprise-grade features like automatic HTTPS, DNS management, backup/restore, and integrated security.