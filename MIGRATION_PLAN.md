# Wild Cloud v.PoC to Wild Central Migration Plan

## Executive Summary

This document outlines the comprehensive migration plan from the bash-based Wild Cloud v.PoC system to the new Go-based Wild Central daemon architecture. The migration follows a **vertical slice approach**, implementing complete end-to-end functionality in five phases, ensuring working software at each stage.

**Key Changes:**
- **34+ bash scripts** → **Single Go daemon (wildd) + Go CLI (wild-cli)**
- **Single WC_HOME** → **Multiple instance management** via `$WILD_CENTRAL_DATA`
- **Direct script execution** → **REST API** with CLI and web clients
- **File-based operations** → **API-driven** with file storage backend
- **Local context** → **Multi-cloud context** switching

This is a large migration. Progress on the migration will be tracked in @MIGRATION_PROGRESS.md. We will update this document as we find plan changes necessary as we go.

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Domains & Functionality](#core-domains--functionality)
3. [Daemon Architecture Design](#daemon-architecture-design)
4. [REST API Specification](#rest-api-specification)
5. [Implementation Phases](#implementation-phases)
6. [Development Steps](#development-steps)
7. [Testing Strategy](#testing-strategy)
8. [Migration Path for Users](#migration-path-for-users)

---

## Architecture Overview

### Wild Cloud v.PoC (Current)

```
┌─────────────────────────────────────────────┐
│  User Machine                               │
│  ┌────────────────────────────────────┐     │
│  │  WC_ROOT (repository)              │     │
│  │  - bin/ (34+ bash scripts)         │     │
│  │  - apps/ (application templates)   │     │
│  │  - setup/ (infrastructure)         │     │
│  └────────────────────────────────────┘     │
│                 ↓                           │
│  ┌────────────────────────────────────┐     │
│  │  WC_HOME (deployment directory)    │     │
│  │  - config.yaml                     │     │
│  │  - secrets.yaml                    │     │
│  │  - apps/ (deployed configs)        │     │
│  └────────────────────────────────────┘     │
└─────────────────────────────────────────────┘
```

**Characteristics:**
- Scripts run directly on user's machine
- One WC_HOME per Wild Cloud instance
- Direct kubectl/talosctl/gomplate execution
- Environment variables (WC_ROOT, WC_HOME) for context

### Wild Central (Target)

```
┌─────────────────────────────────────────────┐
│  Wild Central Device (Raspberry Pi)         │
│  ┌────────────────────────────────────┐     │
│  │  wildd (Go daemon)                 │     │
│  │  - HTTP API Server                 │     │
│  │  - Multi-instance management       │     │
│  │  - dnsmasq integration             │     │
│  └────────────────────────────────────┘     │
│                 ↓                           │
│  ┌────────────────────────────────────┐     │
│  │  $WILD_CENTRAL_DATA                │     │
│  │  /var/lib/wild-central/            │     │
│  │  ├── instances/                    │     │
│  │  │   ├── prod-cluster/             │     │
│  │  │   │   ├── config.yaml           │     │
│  │  │   │   ├── secrets.yaml          │     │
│  │  │   │   └── apps/                 │     │
│  │  │   └── dev-cluster/              │     │
│  │  ├── apps/ (templates from repo)   │     │
│  │  └── logs/                         │     │
│  └────────────────────────────────────┘     │
└─────────────────────────────────────────────┘
         ↑                ↑
         │                │
    wild-cli          wild-app
   (Go CLI)          (React UI)
```

**Characteristics:**
- Daemon runs on dedicated device
- Manages multiple Wild Cloud instances
- REST-ish API for all operations
- Context-based instance selection (like kubectl)
- Integrated dnsmasq for LAN DNS/DHCP

---

## Core Domains & Functionality

### 1. Instance Management
**Purpose:** Create and manage multiple Wild Cloud instances

**v.PoC Scripts:**
- `wild-init` - Initialize new cloud directory

**Key Operations:**
- Create new Wild Cloud instance
- List available instances
- Get instance details
- Delete instance
- Set current context

### 2. Configuration Management
**Purpose:** Manage configuration and secrets for instances

**v.PoC Scripts:**
- `wild-config` - Read configuration values (via yq)
- `wild-config-set` - Write configuration values
- `wild-secret` - Read secret values
- `wild-secret-set` - Write secret values
- `wild-compile-template` - Process gomplate templates
- `wild-compile-template-dir` - Process template directories

**Key Operations:**
- Get/set config values (YAML path-based)
- Get/set secret values (with auto-generation)
- Process templates with gomplate
- Validate configuration
- Export/import configurations

### 3. Node Operations
**Purpose:** Discover and provision Talos nodes

**v.PoC Scripts:**
- `wild-node-detect` - Hardware detection (network, disks)
- `wild-node-setup` - Complete node lifecycle (detect → configure → deploy)
- `wild-cluster-node-ip` - Get node IP addresses
- `wild-cluster-node-boot-assets-download` - Download Talos assets

**Key Operations:**
- Scan network for available nodes
- Detect node hardware (NICs, disks)
- Configure node settings (role, IP, disk)
- Generate Talos configuration
- Deploy configuration to node
- Monitor node status
- Download Talos OS images

### 4. Cluster Lifecycle
**Purpose:** Deploy and manage Kubernetes clusters

**v.PoC Scripts:**
- `wild-setup-cluster` - Complete cluster setup
- `wild-cluster-config-generate` - Generate core Talos cluster configs
- `wild-health` - Comprehensive health checks

**Key Operations:**
- Generate cluster configuration
- Bootstrap Kubernetes cluster
- Add/remove control plane nodes
- Add/remove worker nodes
- Check cluster health
- Upgrade cluster version
- Reset cluster

### 5. Service Management
**Purpose:** Install and manage core cluster services

**v.PoC Scripts:**
- `wild-setup-services` - Install all services in order
- `wild-service-setup` - Setup individual service
- `wild-dashboard-token` - Get Kubernetes dashboard token

**Key Operations:**
- Configure service
- Install services in dependency order
- Check service status
- Update service

### 6. Application Management
**Purpose:** Deploy and manage cloud applications

**v.PoC Scripts:**
- `wild-apps-list` - List available apps
- `wild-app-add` - Configure app from template
- `wild-app-deploy` - Deploy app to cluster
- `wild-app-delete` - Remove app
- `wild-app-doctor` - Run diagnostics
- `wild-app-backup` - Backup app data
- `wild-app-restore` - Restore app from backup

**Key Operations:**
- Browse app catalog
- Configure app (merge defaults, generate secrets)
- Deploy app (with dependencies)
- Update app configuration
- Remove app
- Backup/restore app data
- Check app status

### 7. DNS/DHCP Management
**Purpose:** Manage dnsmasq for LAN services

**v.PoC Scripts:**
- `wild-dnsmasq-install` - Install dnsmasq services

**Key Operations:**
- Generate dnsmasq configuration
- Apply dnsmasq configuration
- Restart dnsmasq service
- Check dnsmasq status

### 8. Backup & Operations
**Purpose:** System maintenance and operations

**v.PoC Scripts:**
- `wild-backup` - Comprehensive backup
- `wild-cluster-secret-copy` - Copy secrets between namespaces
- `wild-talos-schema` - Talos schema management

**Key Operations:**
- Backup instance configuration
- Backup cluster resources
- Backup application data
- Restore from backup
- View logs
- Export configurations
- Copy secrets between namespaces

---

## Daemon Architecture Design

### Directory Structure

```
daemon/
├── cmd/
│   └── wildd/
│       └── main.go              # Entry point
│
├── internal/
│   ├── server/
│   │   ├── server.go            # HTTP server setup
│   │   ├── routes.go            # Route definitions
│   │   └── middleware.go        # CORS, auth, logging
│   │
│   ├── api/                     # API handlers (organized by domain)
│   │   ├── instances.go         # Instance CRUD
│   │   ├── config.go            # Configuration management
│   │   ├── nodes.go             # Node operations
│   │   ├── cluster.go           # Cluster lifecycle
│   │   ├── services.go          # Service management
│   │   ├── apps.go              # Application management
│   │   ├── operations.go        # Async operation tracking
│   │   └── health.go            # Health checks
│   │
│   ├── domain/                  # Business logic (organized by domain)
│   │   ├── instance/
│   │   │   ├── instance.go      # Instance management
│   │   │   ├── config.go        # Config operations
│   │   │   └── secrets.go       # Secret operations
│   │   ├── node/
│   │   │   ├── discover.go      # Node discovery
│   │   │   ├── provision.go     # Node provisioning
│   │   │   └── talos.go         # Talos operations
│   │   ├── cluster/
│   │   │   ├── bootstrap.go     # Cluster bootstrap
│   │   │   ├── config.go        # Cluster configuration
│   │   │   └── health.go        # Health checks
│   │   ├── service/
│   │   │   ├── catalog.go       # Service catalog
│   │   │   ├── deploy.go        # Service deployment
│   │   │   └── template.go      # Template processing
│   │   └── app/
│   │       ├── catalog.go       # App catalog
│   │       ├── configure.go     # App configuration
│   │       └── deploy.go        # App deployment
│   │
│   ├── tools/                   # External tool integration
│   │   ├── kubectl.go           # kubectl wrapper
│   │   ├── talosctl.go          # talosctl wrapper
│   │   ├── gomplate.go          # gomplate wrapper
│   │   ├── yq.go                # yq wrapper
│   │   └── kustomize.go         # kustomize wrapper
│   │
│   ├── storage/                 # File-based persistence
│   │   ├── storage.go           # Storage interface
│   │   ├── instance.go          # Instance storage
│   │   ├── config.go            # Config storage
│   │   └── template.go          # Template management
│   │
│   ├── context/                 # Multi-instance context
│   │   ├── context.go           # Context management
│   │   └── manager.go           # Context switching
│   │
│   ├── async/                   # Async operation handling
│   │   ├── operation.go         # Operation tracking
│   │   ├── executor.go          # Background execution
│   │   └── status.go            # Status management
│   │
│   ├── dnsmasq/                 # DNS/DHCP management
│   │   ├── config.go            # Configuration (existing)
│   │   └── service.go           # Service management
│   │
│   └── models/                  # Data models
│       ├── instance.go          # Instance model
│       ├── node.go              # Node model
│       ├── cluster.go           # Cluster model
│       ├── service.go           # Service model
│       ├── app.go               # App model
│       └── operation.go         # Operation model
│
└── pkg/                         # Public packages (if needed)
```

### Key Design Principles

1. **Direct Tool Integration**
   - Wrap external tools (kubectl, talosctl, gomplate) in thin adapters
   - Pass commands directly to tools via `exec.Command`
   - Capture stdout/stderr for logging and error handling
   - Use tool-specific kubeconfig/talosconfig for context

2. **File-Based Storage**
   - Store all instance data at `$WILD_CENTRAL_DATA/instances/{name}/`
   - Maintain same structure as v.PoC (config.yaml, secrets.yaml, apps/)
   - Allow direct file editing for power users
   - Validate on load

3. **Template Processing**
   - Use gomplate for all template processing
   - Process templates on-demand (not cached)
   - Pass config.yaml and secrets.yaml as data sources
   - Support same template syntax as v.PoC

4. **Async Operations**
   - Long-running operations return 202 Accepted with operation_id
   - Store operation state in memory (or simple file)
   - Client polls GET /api/v1/operations/{id} for status
   - Support log streaming via SSE

5. **Error Handling**
   - Consistent error response format
   - Include external tool errors in details
   - Log all errors with context
   - Return appropriate HTTP status codes

---

## REST API Specification

### API Version

All endpoints are prefixed with `/api/v1`

### Common Response Structures

#### Success Response
```json
{
  "success": true,
  "data": { ... }
}
```

#### Error Response
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "Additional context"
    }
  }
}
```

#### Async Operation Response
```json
{
  "operation_id": "op_abc123",
  "status": "running",
  "message": "Operation started"
}
```

### API Endpoints by Phase

## Phase 1: Instance Management

### Instances

#### Create Instance
```
POST /api/v1/instances
Request:
{
  "name": "prod-cluster",
  "description": "Production cluster",
  "config": {
    "cloud": {
      "domain": "example.com",
      "email": "admin@example.com"
    }
  }
}
Response: 201 Created
{
  "success": true,
  "data": {
    "name": "prod-cluster",
    "created_at": "2025-10-05T10:30:00Z",
    "status": "initialized"
  }
}
```

#### List Instances
```
GET /api/v1/instances
Response: 200 OK
{
  "success": true,
  "data": {
    "instances": [
      {
        "name": "prod-cluster",
        "description": "Production cluster",
        "status": "ready",
        "node_count": 3,
        "created_at": "2025-10-05T10:30:00Z"
      }
    ]
  }
}
```

#### Get Instance
```
GET /api/v1/instances/{name}
Response: 200 OK
{
  "success": true,
  "data": {
    "name": "prod-cluster",
    "description": "Production cluster",
    "status": "ready",
    "node_count": 3,
    "service_count": 7,
    "app_count": 2,
    "created_at": "2025-10-05T10:30:00Z",
    "config_path": "/var/lib/wild-central/instances/prod-cluster/config.yaml"
  }
}
```

#### Delete Instance
```
DELETE /api/v1/instances/{name}
Response: 200 OK
{
  "success": true,
  "message": "Instance deleted"
}
```

### Configuration

#### Get Configuration
```
GET /api/v1/instances/{name}/config
Query params: ?path=cluster.name (optional, for specific value)
Response: 200 OK
{
  "success": true,
  "data": {
    "cloud": {
      "domain": "example.com",
      "email": "admin@example.com"
    },
    "cluster": {
      "name": "prod-cluster",
      "nodeCount": 3
    }
  }
}
```

#### Update Configuration
```
PUT /api/v1/instances/{name}/config
Request:
{
  "path": "cluster.nodeCount",
  "value": 5
}
Response: 200 OK
{
  "success": true,
  "message": "Configuration updated"
}
```

#### Get Configuration YAML
```
GET /api/v1/instances/{name}/config/yaml
Response: 200 OK (text/plain)
cloud:
  domain: example.com
  email: admin@example.com
cluster:
  name: prod-cluster
```

#### Update Configuration YAML
```
PUT /api/v1/instances/{name}/config/yaml
Request: (text/plain)
cloud:
  domain: example.com
  ...
Response: 200 OK
{
  "success": true,
  "message": "Configuration updated"
}
```

### Secrets

#### Get Secret
```
GET /api/v1/instances/{name}/secrets
Query params: ?path=apps.ghost.dbPassword (optional)
Response: 200 OK
{
  "success": true,
  "data": {
    "apps": {
      "ghost": {
        "dbPassword": "***"  // Masked by default
      }
    }
  }
}
```

#### Set Secret
```
PUT /api/v1/instances/{name}/secrets
Request:
{
  "path": "apps.ghost.adminPassword",
  "value": "secure-password",
  "generate": false  // If true, generate random value
}
Response: 200 OK
{
  "success": true,
  "message": "Secret updated"
}
```

#### Generate Secret
```
POST /api/v1/instances/{name}/secrets/generate
Request:
{
  "path": "apps.newapp.apiKey",
  "length": 32
}
Response: 201 Created
{
  "success": true,
  "data": {
    "path": "apps.newapp.apiKey",
    "value": "generated-value"
  }
}
```

### Context

#### Get Current Context
```
GET /api/v1/context
Response: 200 OK
{
  "success": true,
  "data": {
    "current": "prod-cluster"
  }
}
```

#### Set Current Context
```
POST /api/v1/context
Request:
{
  "instance": "dev-cluster"
}
Response: 200 OK
{
  "success": true,
  "message": "Context switched to dev-cluster"
}
```

---

## Phase 2: Node Operations

### Node Discovery

#### Discover Nodes
```
POST /api/v1/instances/{name}/nodes/discover
Request:
{
  "subnet": "192.168.1.0/24",  // Optional, auto-detect if not provided
  "timeout": 30
}
Response: 202 Accepted
{
  "operation_id": "op_discover_abc123"
}

Poll: GET /api/v1/operations/op_discover_abc123
Response: 200 OK
{
  "status": "completed",
  "result": {
    "nodes": [
      {
        "mac": "00:11:22:33:44:55",
        "ip": "192.168.1.100",
        "interface": "eth0",
        "maintenance_mode": true,
        "disks": [
          {
            "path": "/dev/sda",
            "size": 500000000000
          }
        ]
      }
    ]
  }
}
```

#### Get Node Hardware Info
```
GET /api/v1/instances/{name}/nodes/{mac}/hardware
Response: 200 OK
{
  "success": true,
  "data": {
    "mac": "00:11:22:33:44:55",
    "ip": "192.168.1.100",
    "interface": "eth0",
    "maintenance_mode": true,
    "disks": [
      {
        "path": "/dev/sda",
        "size": 500000000000,
        "model": "Samsung SSD"
      }
    ]
  }
}
```

### Node Configuration

#### Add Node
```
POST /api/v1/instances/{name}/nodes
Request:
{
  "name": "control-1",
  "role": "controlplane",
  "mac": "00:11:22:33:44:55",
  "target_ip": "192.168.1.91",
  "interface": "eth0",
  "disk": "/dev/sda",
  "version": "v1.11.0",
  "schematic_id": "56774e0894..."
}
Response: 201 Created
{
  "success": true,
  "data": {
    "name": "control-1",
    "status": "configured"
  }
}
```

#### List Nodes
```
GET /api/v1/instances/{name}/nodes
Response: 200 OK
{
  "success": true,
  "data": {
    "nodes": [
      {
        "name": "control-1",
        "role": "controlplane",
        "status": "ready",
        "ip": "192.168.1.91",
        "version": "v1.11.0"
      }
    ]
  }
}
```

#### Get Node
```
GET /api/v1/instances/{name}/nodes/{node}
Response: 200 OK
{
  "success": true,
  "data": {
    "name": "control-1",
    "role": "controlplane",
    "status": "ready",
    "mac": "00:11:22:33:44:55",
    "target_ip": "192.168.1.91",
    "interface": "eth0",
    "disk": "/dev/sda",
    "version": "v1.11.0",
    "schematic_id": "56774e0894..."
  }
}
```

#### Setup Node
```
POST /api/v1/instances/{name}/nodes/{node}/setup
Request:
{
  "force_reconfigure": false,
  "deploy": true
}
Response: 202 Accepted
{
  "operation_id": "op_setup_abc123"
}
```

#### Delete Node
```
DELETE /api/v1/instances/{name}/nodes/{node}
Response: 200 OK
{
  "success": true,
  "message": "Node deleted"
}
```

---

## Phase 3: Cluster Lifecycle

### Cluster Configuration

#### Generate Cluster Config
```
POST /api/v1/instances/{name}/cluster/config/generate
Request:
{
  "cluster_name": "prod-cluster",
  "vip": "192.168.1.90",
  "version": "v1.11.0"
}
Response: 200 OK
{
  "success": true,
  "message": "Cluster configuration generated"
}
```

### Cluster Operations

#### Bootstrap Cluster
```
POST /api/v1/instances/{name}/cluster/bootstrap
Request:
{
  "node": "control-1"
}
Response: 202 Accepted
{
  "operation_id": "op_bootstrap_abc123"
}
```

#### Get Cluster Status
```
GET /api/v1/instances/{name}/cluster/status
Response: 200 OK
{
  "success": true,
  "data": {
    "status": "ready",
    "nodes": 3,
    "control_plane_nodes": 3,
    "worker_nodes": 0,
    "kubernetes_version": "v1.28.0",
    "talos_version": "v1.11.0",
    "services": {
      "api_server": "healthy",
      "etcd": "healthy",
      "scheduler": "healthy",
      "controller_manager": "healthy"
    }
  }
}
```

#### Check Cluster Health
```
GET /api/v1/instances/{name}/cluster/health
Response: 200 OK
{
  "success": true,
  "data": {
    "status": "healthy",
    "checks": [
      {
        "name": "API Server",
        "status": "passing",
        "message": "API server responding"
      },
      {
        "name": "etcd",
        "status": "passing",
        "message": "All etcd members healthy"
      },
      {
        "name": "Nodes",
        "status": "passing",
        "message": "All nodes ready"
      }
    ]
  }
}
```

#### Get Kubeconfig
```
GET /api/v1/instances/{name}/cluster/kubeconfig
Response: 200 OK (text/plain)
apiVersion: v1
kind: Config
clusters:
  - cluster:
      server: https://192.168.1.90:6443
```

#### Reset Cluster
```
POST /api/v1/instances/{name}/cluster/reset
Request:
{
  "confirm": true
}
Response: 202 Accepted
{
  "operation_id": "op_reset_abc123",
  "message": "Cluster reset initiated"
}
```

---

## Phase 4: Service Management

### Service Catalog

#### List Services
```
GET /api/v1/instances/{name}/services
Response: 200 OK
{
  "success": true,
  "data": {
    "services": [
      {
        "name": "metallb",
        "description": "Load balancer for bare metal",
        "status": "running",
        "version": "0.13.0",
        "namespace": "metallb-system"
      },
      {
        "name": "traefik",
        "description": "Ingress controller",
        "status": "running",
        "version": "2.10.0",
        "namespace": "traefik"
      }
    ]
  }
}
```

#### Get Service
```
GET /api/v1/instances/{name}/services/{service}
Response: 200 OK
{
  "success": true,
  "data": {
    "name": "traefik",
    "description": "Ingress controller",
    "status": "running",
    "version": "2.10.0",
    "namespace": "traefik",
    "pods": 2,
    "healthy_pods": 2,
    "dependencies": ["metallb"]
  }
}
```

### Service Operations

#### Install Service
```
POST /api/v1/instances/{name}/services
Request:
{
  "name": "cert-manager",
  "fetch": false,  // Fetch fresh templates
  "deploy": true   // Actually deploy (or just configure)
}
Response: 202 Accepted
{
  "operation_id": "op_service_abc123"
}
```

#### Install All Services
```
POST /api/v1/instances/{name}/services/install-all
Request:
{
  "fetch": false,
  "deploy": true
}
Response: 202 Accepted
{
  "operation_id": "op_services_abc123"
}
```

#### Update Service
```
PUT /api/v1/instances/{name}/services/{service}
Request:
{
  "version": "2.11.0"
}
Response: 202 Accepted
{
  "operation_id": "op_update_abc123"
}
```

#### Delete Service
```
DELETE /api/v1/instances/{name}/services/{service}
Response: 202 Accepted
{
  "operation_id": "op_delete_abc123"
}
```

#### Get Service Status
```
GET /api/v1/instances/{name}/services/{service}/status
Response: 200 OK
{
  "success": true,
  "data": {
    "name": "traefik",
    "status": "running",
    "pods": [
      {
        "name": "traefik-abc123",
        "status": "Running",
        "ready": true
      }
    ]
  }
}
```

---

## Phase 5: Application Management

### Application Catalog

#### List Available Apps
```
GET /api/v1/apps
Response: 200 OK
{
  "success": true,
  "data": {
    "apps": [
      {
        "name": "ghost",
        "description": "Publishing platform",
        "version": "5.0.0",
        "category": "content",
        "dependencies": ["postgresql"]
      },
      {
        "name": "immich",
        "description": "Photo management",
        "version": "1.0.0",
        "category": "media",
        "dependencies": ["postgresql", "redis"]
      }
    ]
  }
}
```

#### Get App Details
```
GET /api/v1/apps/{app}
Response: 200 OK
{
  "success": true,
  "data": {
    "name": "ghost",
    "description": "Publishing platform for blogs",
    "version": "5.0.0",
    "category": "content",
    "dependencies": ["postgresql"],
    "default_config": {
      "domain": "ghost.{{ .cloud.domain }}",
      "storage": "10Gi",
      "timezone": "UTC"
    },
    "required_secrets": [
      "apps.ghost.dbPassword",
      "apps.postgresql.password"
    ]
  }
}
```

### Application Operations

#### List Deployed Apps
```
GET /api/v1/instances/{name}/apps
Response: 200 OK
{
  "success": true,
  "data": {
    "apps": [
      {
        "name": "ghost",
        "status": "running",
        "version": "5.0.0",
        "namespace": "ghost",
        "url": "https://ghost.example.com"
      }
    ]
  }
}
```

#### Get Deployed App
```
GET /api/v1/instances/{name}/apps/{app}
Response: 200 OK
{
  "success": true,
  "data": {
    "name": "ghost",
    "status": "running",
    "version": "5.0.0",
    "namespace": "ghost",
    "url": "https://ghost.example.com",
    "pods": 1,
    "healthy_pods": 1,
    "config": {
      "domain": "ghost.example.com",
      "storage": "10Gi"
    }
  }
}
```

#### Configure App
```
POST /api/v1/instances/{name}/apps
Request:
{
  "name": "ghost",
  "force": false,
  "config": {
    "domain": "blog.example.com",
    "storage": "20Gi"
  }
}
Response: 200 OK
{
  "success": true,
  "message": "App configured",
  "data": {
    "name": "ghost",
    "status": "configured"
  }
}
```

#### Deploy App
```
POST /api/v1/instances/{name}/apps/{app}/deploy
Request:
{
  "force": false,
  "dry_run": false
}
Response: 202 Accepted
{
  "operation_id": "op_deploy_abc123"
}
```

#### Update App Configuration
```
PUT /api/v1/instances/{name}/apps/{app}
Request:
{
  "config": {
    "storage": "50Gi"
  }
}
Response: 200 OK
{
  "success": true,
  "message": "App configuration updated"
}
```

#### Delete App
```
DELETE /api/v1/instances/{name}/apps/{app}
Query params: ?force=true
Response: 202 Accepted
{
  "operation_id": "op_delete_abc123"
}
```

#### Run App Diagnostics
```
POST /api/v1/instances/{name}/apps/{app}/doctor
Response: 202 Accepted
{
  "operation_id": "op_doctor_abc123"
}
```

#### Get App Status
```
GET /api/v1/instances/{name}/apps/{app}/status
Response: 200 OK
{
  "success": true,
  "data": {
    "name": "ghost",
    "status": "running",
    "pods": [
      {
        "name": "ghost-abc123",
        "status": "Running",
        "ready": true,
        "restarts": 0
      }
    ],
    "dependencies": [
      {
        "name": "postgresql",
        "status": "running"
      }
    ]
  }
}
```

---

## System-Wide Endpoints

### Health & Status

#### System Health
```
GET /api/v1/health
Response: 200 OK
{
  "success": true,
  "data": {
    "status": "healthy",
    "version": "1.0.0",
    "uptime": "24h30m15s"
  }
}
```

#### System Status
```
GET /api/v1/status
Response: 200 OK
{
  "success": true,
  "data": {
    "daemon": "running",
    "version": "1.0.0",
    "uptime": "24h30m15s",
    "instances": 2,
    "dnsmasq": "running"
  }
}
```

### Operations

#### Get Operation Status
```
GET /api/v1/operations/{id}
Response: 200 OK
{
  "success": true,
  "data": {
    "operation_id": "op_abc123",
    "type": "cluster.bootstrap",
    "status": "running",
    "progress": 45,
    "message": "Waiting for API server...",
    "started_at": "2025-10-05T10:30:00Z",
    "logs": [
      "Step 1/6: Bootstrapping etcd cluster",
      "Step 2/6: Waiting for VIP assignment"
    ]
  }
}
```

#### List Operations
```
GET /api/v1/operations
Query params: ?instance=prod-cluster&status=running
Response: 200 OK
{
  "success": true,
  "data": {
    "operations": [
      {
        "operation_id": "op_abc123",
        "type": "cluster.bootstrap",
        "instance": "prod-cluster",
        "status": "running",
        "started_at": "2025-10-05T10:30:00Z"
      }
    ]
  }
}
```

#### Stream Operation Logs
```
GET /api/v1/operations/{id}/logs
Response: 200 OK (text/event-stream)
event: log
data: {"timestamp": "...", "message": "Step 1/6: ..."}

event: log
data: {"timestamp": "...", "message": "Step 2/6: ..."}

event: complete
data: {"status": "completed"}
```

### DNS/DHCP (dnsmasq)

#### Get dnsmasq Config
```
GET /api/v1/dnsmasq/config
Response: 200 OK
{
  "success": true,
  "data": {
    "interface": "eth0",
    "dhcp_range": "192.168.1.100,192.168.1.200",
    "dns_server": "1.1.1.1"
  }
}
```

#### Update dnsmasq Config
```
PUT /api/v1/dnsmasq/config
Request:
{
  "interface": "eth0",
  "dhcp_range": "192.168.1.100,192.168.1.200"
}
Response: 200 OK
{
  "success": true,
  "message": "dnsmasq configuration updated"
}
```

#### Restart dnsmasq
```
POST /api/v1/dnsmasq/restart
Response: 200 OK
{
  "success": true,
  "message": "dnsmasq restarted"
}
```

#### Get dnsmasq Status
```
GET /api/v1/dnsmasq/status
Response: 200 OK
{
  "success": true,
  "data": {
    "status": "running",
    "pid": 1234,
    "leases": 5
  }
}
```

### PXE Boot

#### Download PXE Assets
```
POST /api/v1/pxe/assets/download
Request:
{
  "version": "v1.11.0",
  "schematic_id": "56774e0894..."
}
Response: 202 Accepted
{
  "operation_id": "op_pxe_abc123"
}
```

---

## Implementation Phases

### Phase 1: Core Instance Management (Week 1-2)

**Goal:** Multi-instance configuration management with REST API

**Components to Build:**
1. `internal/storage/` - File-based storage for instances
2. `internal/domain/instance/` - Instance business logic
3. `internal/domain/instance/config.go` - Config operations (yq integration)
4. `internal/domain/instance/secrets.go` - Secret operations
5. `internal/tools/yq.go` - yq wrapper for YAML manipulation
6. `internal/tools/gomplate.go` - gomplate wrapper for templates
7. `internal/api/instances.go` - Instance CRUD handlers
8. `internal/api/config.go` - Config/secret handlers
9. `internal/context/` - Context management

**API Endpoints:**
- POST /api/v1/instances
- GET /api/v1/instances
- GET /api/v1/instances/{name}
- DELETE /api/v1/instances/{name}
- GET /api/v1/instances/{name}/config
- PUT /api/v1/instances/{name}/config
- GET /api/v1/instances/{name}/config/yaml
- PUT /api/v1/instances/{name}/config/yaml
- GET /api/v1/instances/{name}/secrets
- PUT /api/v1/instances/{name}/secrets
- POST /api/v1/instances/{name}/secrets/generate
- GET /api/v1/context
- POST /api/v1/context

**Validation:**
- Can create multiple instances
- Config values can be read/written via YAML path
- Secrets are generated and stored securely
- Context switching works
- Templates process with gomplate

**Success Criteria:**
- [ ] Can create 3+ test instances
- [ ] Config operations work with nested YAML paths
- [ ] Secrets are stored with 600 permissions
- [ ] Context switches between instances
- [ ] Templates render correctly

---

### Phase 2: Node Discovery & Setup (Week 3-4)

**Goal:** Discover and provision Talos nodes

**Components to Build:**
1. `internal/domain/node/discover.go` - Network scanning for nodes
2. `internal/domain/node/provision.go` - Node provisioning logic
3. `internal/domain/node/talos.go` - Talos configuration generation
4. `internal/tools/talosctl.go` - talosctl wrapper
5. `internal/api/nodes.go` - Node API handlers
6. `internal/async/` - Async operation tracking

**API Endpoints:**
- POST /api/v1/instances/{name}/nodes/discover
- GET /api/v1/instances/{name}/nodes/{mac}/hardware
- POST /api/v1/instances/{name}/nodes
- GET /api/v1/instances/{name}/nodes
- GET /api/v1/instances/{name}/nodes/{node}
- POST /api/v1/instances/{name}/nodes/{node}/setup
- DELETE /api/v1/instances/{name}/nodes/{node}
- GET /api/v1/operations/{id}
- GET /api/v1/operations/{id}/logs

**Validation:**
- Can discover nodes on network
- Hardware detection works
- Node configuration is generated correctly
- Talos configuration applies to nodes
- Async operations track status

**Success Criteria:**
- [ ] Discover 3 nodes on test network
- [ ] Hardware detection returns NICs and disks
- [ ] Generate Talos configuration for control plane
- [ ] Apply configuration to node in maintenance mode
- [ ] Track operation status via polling

---

### Phase 3: Cluster Deployment (Week 5-6)

**Goal:** Bootstrap and manage Kubernetes clusters

**Components to Build:**
1. `internal/domain/cluster/bootstrap.go` - Cluster bootstrap logic
2. `internal/domain/cluster/config.go` - Cluster configuration
3. `internal/domain/cluster/health.go` - Health checks
4. `internal/api/cluster.go` - Cluster API handlers

**API Endpoints:**
- POST /api/v1/instances/{name}/cluster/config/generate
- POST /api/v1/instances/{name}/cluster/bootstrap
- GET /api/v1/instances/{name}/cluster/status
- GET /api/v1/instances/{name}/cluster/health
- GET /api/v1/instances/{name}/cluster/kubeconfig
- POST /api/v1/instances/{name}/cluster/reset

**Validation:**
- Cluster configuration generates correctly
- Bootstrap process completes successfully
- VIP assignment works
- API server responds
- Kubeconfig is valid
- Health checks pass

**Success Criteria:**
- [ ] Generate base Talos configs (controlplane.yaml, worker.yaml)
- [ ] Bootstrap cluster on first control plane node
- [ ] VIP is assigned and API server responds
- [ ] Additional nodes join cluster
- [ ] Health checks validate all components
- [ ] kubectl commands work with generated kubeconfig

---

### Phase 4: Service Management (Week 7-8)

**Goal:** Install and manage core cluster services

**Components to Build:**
1. `internal/domain/service/catalog.go` - Service catalog
2. `internal/domain/service/deploy.go` - Service deployment
3. `internal/domain/service/template.go` - Template processing for services
4. `internal/tools/kubectl.go` - kubectl wrapper
5. `internal/tools/kustomize.go` - kustomize wrapper
6. `internal/api/services.go` - Service API handlers

**API Endpoints:**
- GET /api/v1/instances/{name}/services
- GET /api/v1/instances/{name}/services/{service}
- POST /api/v1/instances/{name}/services
- POST /api/v1/instances/{name}/services/install-all
- PUT /api/v1/instances/{name}/services/{service}
- DELETE /api/v1/instances/{name}/services/{service}
- GET /api/v1/instances/{name}/services/{service}/status

**Validation:**
- Service templates process correctly
- Dependencies are resolved
- Services install in correct order
- kubectl operations work with instance kubeconfig
- Service status reflects actual state

**Success Criteria:**
- [ ] Install MetalLB successfully
- [ ] Install Traefik with MetalLB dependency
- [ ] Install cert-manager and issue certificates
- [ ] All core services running and healthy
- [ ] Service dependency resolution works

---

### Phase 5: Application Deployment (Week 9-10)

**Goal:** Deploy and manage user applications

**Components to Build:**
1. `internal/domain/app/catalog.go` - App catalog from repo
2. `internal/domain/app/configure.go` - App configuration (wild-app-add logic)
3. `internal/domain/app/deploy.go` - App deployment (wild-app-deploy logic)
4. `internal/api/apps.go` - App API handlers

**API Endpoints:**
- GET /api/v1/apps
- GET /api/v1/apps/{app}
- GET /api/v1/instances/{name}/apps
- GET /api/v1/instances/{name}/apps/{app}
- POST /api/v1/instances/{name}/apps
- POST /api/v1/instances/{name}/apps/{app}/deploy
- PUT /api/v1/instances/{name}/apps/{app}
- DELETE /api/v1/instances/{name}/apps/{app}
- POST /api/v1/instances/{name}/apps/{app}/doctor
- GET /api/v1/instances/{name}/apps/{app}/status

**Validation:**
- App catalog reads from repository
- App configuration merges defaults correctly
- Required secrets are generated
- Templates process with gomplate
- Apps deploy with dependencies
- Kustomize builds work

**Success Criteria:**
- [ ] List apps from catalog
- [ ] Configure Ghost app (merge defaults, generate secrets)
- [ ] Deploy PostgreSQL dependency
- [ ] Deploy Ghost app
- [ ] App is accessible via ingress
- [ ] App status reflects actual state

---

## Development Steps

### Prerequisites

1. **Set up development environment:**
   ```bash
   cd /home/payne/repos/wild-central/daemon
   go mod init github.com/wild-cloud/wild-central/daemon
   go get github.com/gorilla/mux
   ```

2. **Install required tools locally:**
   - kubectl
   - talosctl
   - gomplate
   - yq
   - kustomize

3. **Set up test environment:**
   - Test Wild Cloud v.PoC instance
   - Test Talos nodes (VMs or hardware)
   - Network for testing

### Phase 1 Implementation Steps

1. **Create storage layer:**
   ```go
   // internal/storage/storage.go
   type Manager struct {
       dataDir string
   }
   func NewManager() *Manager
   func (m *Manager) GetInstancePath(name string) string
   func (m *Manager) ListInstances() ([]string, error)
   func (m *Manager) InstanceExists(name string) bool
   func (m *Manager) CreateInstance(name string) error
   func (m *Manager) DeleteInstance(name string) error
   ```

2. **Create configuration management:**
   ```go
   // internal/domain/instance/config.go
   func GetConfig(instance, path string) (interface{}, error)
   func SetConfig(instance, path string, value interface{}) error
   func GetConfigYAML(instance string) (string, error)
   func SetConfigYAML(instance, yaml string) error
   ```

3. **Create tool wrappers:**
   ```go
   // internal/tools/yq.go
   func Get(file, path string) (string, error)
   func Set(file, path, value string) error

   // internal/tools/gomplate.go
   func Process(template string, data map[string]interface{}) (string, error)
   ```

4. **Create API handlers:**
   ```go
   // internal/api/instances.go
   func CreateInstance(w http.ResponseWriter, r *http.Request)
   func ListInstances(w http.ResponseWriter, r *http.Request)
   func GetInstance(w http.ResponseWriter, r *http.Request)
   func DeleteInstance(w http.ResponseWriter, r *http.Request)
   ```

5. **Set up routing:**
   ```go
   // internal/server/routes.go
   func SetupRoutes(router *mux.Router) {
       api := router.PathPrefix("/api/v1").Subrouter()

       // Instances
       api.HandleFunc("/instances", api.CreateInstance).Methods("POST")
       api.HandleFunc("/instances", api.ListInstances).Methods("GET")
       api.HandleFunc("/instances/{name}", api.GetInstance).Methods("GET")
       api.HandleFunc("/instances/{name}", api.DeleteInstance).Methods("DELETE")

       // Config
       api.HandleFunc("/instances/{name}/config", api.GetConfig).Methods("GET")
       api.HandleFunc("/instances/{name}/config", api.UpdateConfig).Methods("PUT")
       // ... etc
   }
   ```

6. **Test Phase 1:**
   ```bash
   # Start daemon
   cd daemon
   go run cmd/wildd/main.go

   # Test with curl
   curl -X POST http://localhost:5055/api/v1/instances \
     -H "Content-Type: application/json" \
     -d '{"name":"test-cluster"}'

   curl http://localhost:5055/api/v1/instances

   curl http://localhost:5055/api/v1/instances/test-cluster/config
   ```

### Repeat for Phases 2-5

Each phase follows similar pattern:
1. Implement domain logic
2. Create tool wrappers
3. Build API handlers
4. Add routes
5. Test with curl/CLI
6. Validate against success criteria

---

## Testing Strategy

### Unit Tests

Test each package independently:

```go
// internal/storage/storage_test.go
func TestCreateInstance(t *testing.T) {
    manager := NewManager("/tmp/test-wild-central")
    err := manager.CreateInstance("test")
    assert.NoError(t, err)
    assert.True(t, manager.InstanceExists("test"))
}

// internal/tools/yq_test.go
func TestYQGet(t *testing.T) {
    value, err := Get("test-config.yaml", "cluster.name")
    assert.NoError(t, err)
    assert.Equal(t, "test-cluster", value)
}
```

### Integration Tests

Test API endpoints:

```go
// internal/api/instances_test.go
func TestCreateInstanceAPI(t *testing.T) {
    // Setup test server
    router := mux.NewRouter()
    SetupRoutes(router)
    server := httptest.NewServer(router)
    defer server.Close()

    // Test create instance
    resp, err := http.Post(
        server.URL+"/api/v1/instances",
        "application/json",
        strings.NewReader(`{"name":"test"}`),
    )
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

### End-to-End Tests

Test complete workflows:

```bash
#!/bin/bash
# test/e2e/phase1.sh

# Test Phase 1: Instance Management
echo "Testing Phase 1: Instance Management"

# Create instance
curl -X POST http://localhost:5055/api/v1/instances \
  -H "Content-Type: application/json" \
  -d '{"name":"e2e-test"}'

# Get instance
curl http://localhost:5055/api/v1/instances/e2e-test

# Set config
curl -X PUT http://localhost:5055/api/v1/instances/e2e-test/config \
  -H "Content-Type: application/json" \
  -d '{"path":"cluster.name","value":"e2e-cluster"}'

# Get config
curl http://localhost:5055/api/v1/instances/e2e-test/config?path=cluster.name

# Cleanup
curl -X DELETE http://localhost:5055/api/v1/instances/e2e-test

echo "Phase 1 tests completed"
```

### Manual Testing

Test with actual Talos nodes:

1. **Phase 2 test:**
   - Start daemon
   - Create instance
   - Discover nodes on network
   - Configure and setup node
   - Verify Talos configuration applied

2. **Phase 3 test:**
   - Generate cluster config
   - Bootstrap cluster
   - Verify API server responds
   - Check etcd health
   - Add additional nodes

3. **Phase 4 test:**
   - Install MetalLB
   - Install Traefik
   - Verify load balancer IP assigned
   - Install cert-manager
   - Verify certificate issued

4. **Phase 5 test:**
   - Browse app catalog
   - Configure Ghost app
   - Deploy PostgreSQL
   - Deploy Ghost
   - Access Ghost via ingress
   - Verify app running

---

### For New Users

**Recommended Path:**
1. Install Wild Central image on Raspberry Pi
2. Access web UI at `http://wild-central.local:5055`
3. Follow setup wizard:
   - Create first instance
   - Configure network settings
   - Discover nodes
   - Bootstrap cluster
   - Install core services
4. Deploy applications via web UI

### CLI vs Web UI

**wild-cli (for advanced users):**
```bash
# Create instance
wild-cli instance create prod-cluster

# Configure
wild-cli config set cluster.name prod-cluster

# Discover nodes
wild-cli node discover

# Setup nodes
wild-cli node setup control-1 control-2 control-3

# Bootstrap cluster
wild-cli cluster bootstrap

# Install services
wild-cli service install-all

# Deploy app
wild-cli app add ghost
wild-cli app deploy ghost
```

**wild-app (for everyone):**
- Visual setup wizard
- Node discovery with UI
- Click-to-deploy services
- App catalog with descriptions
- Real-time status dashboards
- Log streaming
- Error notifications

---

## Key Design Decisions

### 1. Multi-Instance Architecture

**Decision:** Store each instance in separate directory at `$WILD_CENTRAL_DATA/instances/{name}/`

**Rationale:**
- Clean separation of instances
- Easy to backup/restore individual instances
- Maintains same structure as v.PoC (config.yaml, secrets.yaml, apps/)
- Allows direct file editing for power users

**Implementation:**
```
/var/lib/wild-central/
├── instances/
│   ├── prod-cluster/
│   │   ├── config.yaml
│   │   ├── secrets.yaml
│   │   ├── apps/
│   │   ├── setup/
│   │   ├── .kube/config
│   │   └── .talos/config
│   └── dev-cluster/
│       └── ...
├── apps/              # Shared app catalog from repo
├── setup/             # Shared infrastructure templates
└── logs/
```

### 2. Template Processing

**Decision:** Process templates on-demand with gomplate, no caching

**Rationale:**
- Maintains compatibility with v.PoC approach
- Users can edit templates and see changes immediately
- No stale cache issues
- Gomplate is fast enough for on-demand processing

**Implementation:**
- Keep apps/ and setup/ directories from repository
- Copy to $WILD_CENTRAL_DATA for local modifications
- Process with gomplate when deploying
- Pass config.yaml and secrets.yaml as data sources

### 3. External Tool Integration

**Decision:** Direct execution via `exec.Command`, thin wrappers

**Rationale:**
- Simplest approach (KISS principle)
- No complex abstraction layers
- Captures full tool output for logging
- Easy to debug (see exact commands executed)

**Implementation:**
```go
// internal/tools/kubectl.go
func Apply(instance, manifest string) (string, error) {
    kubeconfig := getKubeconfig(instance)
    cmd := exec.Command("kubectl", "apply",
        "--kubeconfig", kubeconfig,
        "-f", manifest)
    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

### 4. Async Operations

**Decision:** Simple in-memory operation tracking, polling-based status

**Rationale:**
- Start simple, can enhance later
- Polling is sufficient for initial version
- No additional dependencies (Redis, etc.)
- Operations are infrequent enough that polling is acceptable

**Implementation:**
```go
// internal/async/operation.go
type Operation struct {
    ID        string
    Type      string
    Instance  string
    Status    string  // running, completed, failed
    Progress  int     // 0-100
    Message   string
    Logs      []string
    StartedAt time.Time
    Error     string
}

var operations = sync.Map{}
```

### 5. Authentication & Security

**Decision:** Start with basic security

- No authentication (trust LAN)
- CORS enabled for wild-app
- Secrets file has 600 permissions
- No API keys or tokens

**Rationale:**
- Wild Central runs on private LAN
- Initial target is single-user/family scenarios
- Can add security incrementally

### 6. Error Handling

**Decision:** Consistent error responses with external tool errors in details

**Format:**
```json
{
  "error": {
    "code": "KUBECTL_APPLY_FAILED",
    "message": "Failed to apply manifest",
    "details": {
      "command": "kubectl apply -f manifest.yaml",
      "exit_code": 1,
      "output": "error: unable to recognize manifest.yaml..."
    }
  }
}
```

**Rationale:**
- Clear for users
- Includes context for debugging
- External tool errors are surfaced
- Consistent across all endpoints

---

## Idempotency & 1:1 Bash Script Correspondence

**CRITICAL:** The Go daemon implementation must maintain exact 1:1 correspondence with bash script behavior, including idempotency patterns and progress tracking. Every bash script operation must have an equivalent daemon method that produces identical results.

### Principles

1. **Exact Command Correspondence**: Each bash script command sequence must map to daemon methods that execute the same external commands in the same order
2. **Idempotency Preservation**: Operations must be safely repeatable, checking state before acting
3. **Progress Tracking**: Multi-step operations must track completion and support resume from failure
4. **State Consistency**: File system state, configuration keys, and cluster state must match bash script expectations

### Idempotency Patterns from v.PoC

#### Pattern 1: Configuration Key Existence Checks

**Bash Pattern:**
```bash
if wild-config --check "cluster.nodes.active.${NODE_NAME}.interface"; then
    print_success "Node $NODE_NAME already configured"
    # Skip or update
else
    # Detect and configure
fi
```

**Go Implementation:**
```go
type StateManager struct {
    configPath string
    mu sync.RWMutex
}

func (s *StateManager) KeyExists(path string) bool {
    // Check YAML key existence
}

func (s *StateManager) GetKey(path string) (string, error) {
    // Read YAML value at path
}
```

#### Pattern 2: File/Directory Existence Guards

**Bash Pattern:**
```bash
if [ -f "${NODE_SETUP_DIR}/generated/secrets.yaml" ]; then
    print_success "Cluster configuration already exists"
    return 0
fi
```

**Go Implementation:**
```go
type FileSystemChecker struct {
    wcHome string
}

func (f *FileSystemChecker) ClusterInitialized() bool {
    secretsPath := filepath.Join(f.wcHome,
        "setup/cluster-nodes/generated/secrets.yaml")
    return fileExists(secretsPath)
}
```

#### Pattern 3: Kubernetes Resource Existence

**Bash Pattern:**
```bash
HAS_CONTEXT=$(talosctl config contexts | grep -c "$CLUSTER_NAME" || true)
if [ "$HAS_CONTEXT" -eq 0 ]; then
    talosctl config merge ${WC_HOME}/setup/cluster-nodes/generated/talosconfig
fi
```

**Go Implementation:**
```go
func (t *TalosctlWrapper) HasContext(clusterName string) (bool, error) {
    output, err := t.exec("config", "contexts")
    return strings.Contains(output, clusterName), err
}

func (t *TalosctlWrapper) EnsureContext(clusterName, configPath string) error {
    if has, _ := t.HasContext(clusterName); has {
        return nil // Already exists
    }
    return t.exec("config", "merge", configPath)
}
```

#### Pattern 4: Delete-and-Recreate for Jobs

**Bash Pattern:**
```bash
# Ensure idempotent job execution
kubectl delete job immich-db-init --namespace="${APP_NAME}" --ignore-not-found=true
kubectl wait --for=delete job/immich-db-init --namespace="${APP_NAME}" --timeout=30s || true
# Then create job
kubectl apply -f db-init-job.yaml
```

**Go Implementation:**
```go
func (k *KubectlWrapper) RecreateJob(namespace, name, manifestPath string) error {
    // Delete existing job
    if err := k.Delete("job", name, namespace, true); err != nil {
        return err
    }

    // Wait for deletion
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    k.WaitForDelete(ctx, "job", name, namespace)

    // Apply new job
    return k.Apply(manifestPath)
}
```

### Progress Tracking System

#### Phase Tracking Model

```go
type Phase string

const (
    PhaseSetup         Phase = "setup"
    PhaseInfrastructure Phase = "infrastructure"
    PhaseCluster       Phase = "cluster"
    PhaseServices      Phase = "services"
    PhaseApps          Phase = "apps"
)

type ProgressTracker struct {
    CurrentPhase    Phase             `yaml:"currentPhase"`
    CompletedPhases []Phase           `yaml:"completedPhases"`
    StepProgress    map[string]int    `yaml:"stepProgress"`    // step -> completion %
    LastError       string            `yaml:"lastError,omitempty"`
    UpdatedAt       time.Time         `yaml:"updatedAt"`
}

func (p *ProgressTracker) IsPhaseComplete(phase Phase) bool {
    for _, completed := range p.CompletedPhases {
        if completed == phase {
            return true
        }
    }
    return false
}

func (p *ProgressTracker) MarkPhaseComplete(phase Phase) {
    if !p.IsPhaseComplete(phase) {
        p.CompletedPhases = append(p.CompletedPhases, phase)
    }
}
```

#### Operation State Tracking (for API/UI feedback only)

**Important:** Operation state is for UI progress display only, not for resumption logic. Idempotency handles resumption automatically.

```go
type OperationState struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`     // "cluster.bootstrap", "app.deploy", etc.
    Instance    string                 `json:"instance"`
    Status      string                 `json:"status"`   // "running", "completed", "failed"
    Progress    int                    `json:"progress"` // 0-100 (for UI display)
    Message     string                 `json:"message"`  // Current step description
    StartedAt   time.Time              `json:"startedAt"`
    UpdatedAt   time.Time              `json:"updatedAt"`
    CompletedAt *time.Time             `json:"completedAt,omitempty"`
    Error       string                 `json:"error,omitempty"`
    Logs        []string               `json:"logs"`     // Recent log lines for UI
}

// Save for daemon restart recovery (so UI can show operation status)
func (o *OperationState) Save(dataDir string) error {
    path := filepath.Join(dataDir, "operations", o.ID+".json")
    return writeJSON(path, o)
}

// Note: If daemon restarts during operation, the operation will fail.
// User simply re-runs the same command - idempotency handles the rest.
```

### Progressive Waiting with Restarts

**Bash Pattern:**
```bash
# Wait for API server with periodic kubelet restarts
print_info -n "Waiting for API server to respond on VIP."
max_attempts=60
for attempt in $(seq 1 $max_attempts); do
    if curl -k -s --max-time 5 "https://$vip:6443/healthz" >/dev/null 2>&1; then
        print_success "API server responding"
        break
    fi
    # Restart kubelet every 15 attempts
    if [ $((attempt % 15)) -eq 0 ]; then
        talosctl -n "$TARGET_IP" service kubelet restart > /dev/null 2>&1
        sleep 30
    fi
    printf "."
    sleep 10
done
```

**Go Implementation:**
```go
type ProgressiveWaiter struct {
    checkFunc     func() error
    maxAttempts   int
    checkInterval time.Duration
    restartAt     []int                    // Attempt numbers for restart
    restartFunc   func() error
    restartDelay  time.Duration
}

func (w *ProgressiveWaiter) Wait(ctx context.Context) error {
    for attempt := 1; attempt <= w.maxAttempts; attempt++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // Try check
        if err := w.checkFunc(); err == nil {
            return nil // Success
        }

        // Check if we should restart service
        for _, restartAttempt := range w.restartAt {
            if attempt == restartAttempt {
                log.Printf("Attempt %d: Restarting service...", attempt)
                if err := w.restartFunc(); err != nil {
                    log.Printf("Restart failed: %v", err)
                }
                time.Sleep(w.restartDelay)
                break
            }
        }

        time.Sleep(w.checkInterval)
    }

    return fmt.Errorf("check failed after %d attempts", w.maxAttempts)
}

// Example usage matching bash behavior:
waiter := &ProgressiveWaiter{
    checkFunc: func() error {
        return checkAPIServer(vip)
    },
    maxAttempts:   60,
    checkInterval: 10 * time.Second,
    restartAt:     []int{15, 30, 45}, // Restart at these attempts
    restartFunc: func() error {
        return talosctl.RestartService(targetIP, "kubelet")
    },
    restartDelay: 30 * time.Second,
}
```

### Idempotent Operation Interface

```go
type IdempotentOperation interface {
    // Check if prerequisites are met
    CheckPreconditions(ctx context.Context) error

    // Check if operation is already complete
    IsAlreadyComplete(ctx context.Context) (bool, error)

    // Execute the operation
    Execute(ctx context.Context) error

    // Validate the result
    Validate(ctx context.Context) error

    // Rollback if needed
    Rollback(ctx context.Context) error
}

// Example implementation
type ClusterBootstrapOperation struct {
    instance     string
    firstNode    string
    vip          string
    stateManager *StateManager
    talosctl     *TalosctlWrapper
    kubectl      *KubectlWrapper
}

func (op *ClusterBootstrapOperation) IsAlreadyComplete(ctx context.Context) (bool, error) {
    // Check if kubeconfig exists and API server responds
    kubeconfigPath := op.stateManager.GetKubeconfigPath(op.instance)
    if !fileExists(kubeconfigPath) {
        return false, nil
    }

    // Test API server
    if err := op.kubectl.ServerVersion(); err == nil {
        return true, nil // Already bootstrapped
    }

    return false, nil
}
```

### External Command Execution Tracking

```go
type CommandLogger struct {
    operationID string
    logFile     *os.File
}

func (c *CommandLogger) ExecWithLog(name string, args ...string) (string, error) {
    // Log exact command
    cmdStr := fmt.Sprintf("%s %s", name, strings.Join(args, " "))
    c.log("EXEC: %s", cmdStr)

    cmd := exec.Command(name, args...)

    // Capture output
    output, err := cmd.CombinedOutput()
    outputStr := string(output)

    if err != nil {
        c.log("ERROR: %v\nOutput: %s", err, outputStr)
        return outputStr, fmt.Errorf("command failed: %w\nOutput: %s", err, outputStr)
    }

    c.log("SUCCESS: %s", outputStr)
    return outputStr, nil
}

func (c *CommandLogger) log(format string, args ...interface{}) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    line := fmt.Sprintf("[%s] "+format+"\n", append([]interface{}{timestamp}, args...)...)
    c.logFile.WriteString(line)
}
```

### Configuration Prompting Pattern

**Bash Pattern:**
```bash
prompt_if_unset_config() {
    local key="$1"
    local prompt="$2"
    local default="${3:-}"

    if wild-config --check "$key"; then
        print_info "Using existing $key"
        return
    fi

    read -p "$prompt [${default}]: " value
    value="${value:-$default}"
    wild-config-set "$key" "$value"
}
```

**Go Implementation:**
```go
type ConfigPrompter struct {
    stateManager *StateManager
    interactive  bool
    defaults     map[string]string
}

func (p *ConfigPrompter) PromptIfUnset(key, prompt, defaultVal string) (string, error) {
    // Check if already set
    if exists := p.stateManager.KeyExists(key); exists {
        val, _ := p.stateManager.GetKey(key)
        log.Printf("Using existing %s = %s", key, val)
        return val, nil
    }

    // Use default in non-interactive mode
    if !p.interactive {
        if defaultVal != "" {
            p.stateManager.SetKey(key, defaultVal)
            return defaultVal, nil
        }
        return "", fmt.Errorf("key %s not set and no default provided", key)
    }

    // Interactive prompt
    fmt.Printf("%s [%s]: ", prompt, defaultVal)
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    value := strings.TrimSpace(input)

    if value == "" {
        value = defaultVal
    }

    p.stateManager.SetKey(key, value)
    return value, nil
}
```

### State Validation

```go
type StateValidator struct {
    wcHome       string
    stateManager *StateManager
}

func (v *StateValidator) ValidatePhaseComplete(phase Phase) error {
    switch phase {
    case PhaseSetup:
        return v.validateSetupPhase()
    case PhaseInfrastructure:
        return v.validateInfrastructurePhase()
    case PhaseCluster:
        return v.validateClusterPhase()
    case PhaseServices:
        return v.validateServicesPhase()
    case PhaseApps:
        return v.validateAppsPhase()
    }
    return fmt.Errorf("unknown phase: %s", phase)
}

func (v *StateValidator) validateClusterPhase() error {
    // Check secrets exist
    if !fileExists(filepath.Join(v.wcHome, "setup/cluster-nodes/generated/secrets.yaml")) {
        return fmt.Errorf("cluster secrets not generated")
    }

    // Check kubeconfig exists
    if !fileExists(filepath.Join(v.wcHome, ".kube/config")) {
        return fmt.Errorf("kubeconfig not found")
    }

    // Check required config keys
    requiredKeys := []string{
        "cluster.name",
        "cluster.nodes.control.vip",
    }

    for _, key := range requiredKeys {
        if !v.stateManager.KeyExists(key) {
            return fmt.Errorf("required config key missing: %s", key)
        }
    }

    return nil
}
```

### Bash-to-Go Command Mapping

| Bash Script | Go Package | Key Methods | Idempotency Check |
|------------|------------|-------------|-------------------|
| `wild-init` | `domain/instance` | `CreateInstance()` | Check `.wildcloud/` marker |
| `wild-config` | `domain/instance/config` | `GetConfig(path)` | Direct YAML read |
| `wild-config-set` | `domain/instance/config` | `SetConfig(path, value)` | Atomic write with backup |
| `wild-cluster-config-generate` | `domain/cluster` | `GenerateConfig()` | Check `generated/secrets.yaml` |
| `wild-node-setup` | `domain/node` | `SetupNode(name)` | Check node config keys + final/*.yaml |
| `wild-setup-cluster` | `domain/cluster` | `BootstrapCluster()` | Check kubeconfig + API response |
| `wild-setup-services` | `domain/service` | `InstallServices()` | Check pod status per service |
| `wild-app-add` | `domain/app` | `ConfigureApp(name)` | Check `apps/{name}/manifest.yaml` |
| `wild-app-deploy` | `domain/app` | `DeployApp(name)` | Check deployment status |

### Testing Requirements for 1:1 Correspondence

**Important:** The bash v.PoC scripts have no integration tests, so validation must be done through:
- Manual verification against real Talos/Kubernetes clusters
- State inspection of files and configuration
- Behavioral observation of actual operations
- Side-by-side manual testing when needed

#### 1. State-Based Validation Tests

Verify operations produce the expected file system state and configuration:

```go
func TestClusterBootstrapState(t *testing.T) {
    instance := "test-cluster"

    // Bootstrap cluster
    err := cluster.Bootstrap(ctx, instance)
    require.NoError(t, err)

    wcHome := getInstancePath(instance)

    // Verify expected files exist at correct paths
    assertFileExists(t, filepath.Join(wcHome, "setup/cluster-nodes/generated/secrets.yaml"))
    assertFileExists(t, filepath.Join(wcHome, "setup/cluster-nodes/generated/controlplane.yaml"))
    assertFileExists(t, filepath.Join(wcHome, "setup/cluster-nodes/generated/worker.yaml"))
    assertFileExists(t, filepath.Join(wcHome, "setup/cluster-nodes/generated/talosconfig"))
    assertFileExists(t, filepath.Join(wcHome, ".kube/config"))

    // Verify config keys were set (matches bash behavior)
    assertConfigKeyExists(t, instance, "cluster.name")
    assertConfigKeyExists(t, instance, "cluster.nodes.control.vip")

    // Verify secrets.yaml has correct permissions
    info, err := os.Stat(filepath.Join(wcHome, "secrets.yaml"))
    require.NoError(t, err)
    assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

    // Verify kubeconfig is valid and API server responds
    kubectl := NewKubectlWrapper(instance)
    _, err = kubectl.ServerVersion()
    assert.NoError(t, err, "API server should be accessible")
}
```

#### 2. Idempotency Tests

Verify operations can be run multiple times safely:

```go
func TestClusterBootstrapIdempotency(t *testing.T) {
    instance := "test-cluster"

    // Bootstrap once
    err := cluster.Bootstrap(ctx, instance)
    require.NoError(t, err)

    // Capture initial state
    state1 := captureInstanceState(instance)

    // Bootstrap again - should be no-op
    err = cluster.Bootstrap(ctx, instance)
    assert.NoError(t, err)

    // Capture state after second run
    state2 := captureInstanceState(instance)

    // States should be identical (no duplicate resources)
    assert.Equal(t, state1, state2)

    // Verify no duplicate Kubernetes resources
    nodes, _ := kubectl.GetNodes()
    assert.Equal(t, 3, len(nodes), "Should still have 3 nodes, not 6")
}

func captureInstanceState(instance string) map[string]interface{} {
    return map[string]interface{}{
        "config":     readYAML(getConfigPath(instance)),
        "secrets":    getSecretKeys(instance), // Don't expose values
        "files":      listFiles(getInstancePath(instance)),
        "k8s_nodes":  countKubernetesNodes(instance),
        "k8s_pods":   countKubernetesPods(instance),
    }
}
```

#### 3. Automatic Resume via Idempotency Tests

**Important:** True idempotency means no explicit "resume" logic is needed. If an operation fails, just re-run it - the idempotency checks will skip completed steps automatically.

Verify that re-running after failure continues correctly:

```go
func TestClusterBootstrapReRunAfterFailure(t *testing.T) {
    instance := "test-cluster"

    // Simulate partial completion (e.g., secrets generated but cluster not bootstrapped)
    setupPartialClusterState(instance)

    // First attempt - will complete remaining steps
    err := cluster.Bootstrap(ctx, instance)
    require.NoError(t, err)

    // Verify it detected existing state and skipped completed steps
    // (secrets generation would be skipped, bootstrap would execute)
    assertFileExists(t, wcHome+"/setup/cluster-nodes/generated/secrets.yaml")
    assertAPIServerResponds(t, instance)
}

func TestIdempotentStepSkipping(t *testing.T) {
    instance := "test-cluster"

    // Complete bootstrap
    err := cluster.Bootstrap(ctx, instance)
    require.NoError(t, err)

    // Track which external commands are executed on second run
    commandLog := captureCommandLog()

    // Run again - should skip all steps
    err = cluster.Bootstrap(ctx, instance)
    require.NoError(t, err)

    // Verify no kubectl/talosctl commands were executed
    // (all checks returned "already complete")
    assert.Empty(t, commandLog.Commands, "Should not execute any commands when already complete")
}
```

**Design Principle:** Instead of explicit "resume from step N" logic, each step checks:
1. Is my prerequisite complete? (idempotency check)
2. If yes, skip
3. If no, execute

Example:
```go
func (b *ClusterBootstrap) Execute(ctx context.Context) error {
    // Step 1: Generate secrets
    if !fileExists(b.secretsPath) {
        if err := b.generateSecrets(); err != nil {
            return err
        }
    }

    // Step 2: Bootstrap etcd
    if !b.isEtcdHealthy() {
        if err := b.bootstrapEtcd(); err != nil {
            return err
        }
    }

    // Step 3: Wait for API server
    if !b.isAPIServerResponding() {
        if err := b.waitForAPIServer(); err != nil {
            return err
        }
    }

    // Each step checks state first - no explicit "resume" needed
    return nil
}
```

#### 4. Manual Verification Checklist

Since bash scripts have no automated tests, each operation must be manually verified:

**For Cluster Bootstrap:**
- [ ] Run `wild-cli cluster bootstrap` on real hardware
- [ ] Verify control plane nodes boot and join
- [ ] Check VIP is assigned and API server responds
- [ ] Verify etcd cluster is healthy
- [ ] Confirm kubectl commands work
- [ ] Compare file structure with v.PoC WC_HOME
- [ ] Verify config.yaml has same keys as bash scripts would create
- [ ] Check logs match bash script verbosity and patterns

**For Node Setup:**
- [ ] Run `wild-cli node setup control-1` on maintenance mode node
- [ ] Verify hardware detection matches `wild-node-detect` output
- [ ] Check configuration files created in correct locations
- [ ] Verify Talos configuration applies successfully
- [ ] Confirm node reboots and joins cluster

**For Service Installation:**
- [ ] Run `wild-cli service install metallb`
- [ ] Verify service deploys in correct order
- [ ] Check dependencies are installed first
- [ ] Confirm pods reach Running state
- [ ] Verify service configuration matches v.PoC templates

**For App Deployment:**
- [ ] Run `wild-cli app add ghost`
- [ ] Verify manifest.yaml copied and processed
- [ ] Check defaultConfig merged into config.yaml
- [ ] Confirm secrets generated with correct values
- [ ] Verify templates compiled with gomplate
- [ ] Run `wild-cli app deploy ghost`
- [ ] Check dependencies (PostgreSQL) deploy first
- [ ] Confirm app accessible via ingress

#### 5. Behavioral Equivalence Tests

Test that Go implementation behaves identically to bash:

```go
func TestConfigSetBehavior(t *testing.T) {
    instance := "test-cluster"

    // Set nested key (bash: wild-config-set "cluster.nodes.active.control-1.ip" "192.168.1.91")
    err := config.SetKey(instance, "cluster.nodes.active.control-1.ip", "192.168.1.91")
    require.NoError(t, err)

    // Verify key is set
    value, err := config.GetKey(instance, "cluster.nodes.active.control-1.ip")
    assert.NoError(t, err)
    assert.Equal(t, "192.168.1.91", value)

    // Verify YAML structure matches bash output
    configYAML := readYAML(getConfigPath(instance))
    assert.Equal(t, "192.168.1.91",
        configYAML["cluster"].(map[string]interface{})
            ["nodes"].(map[string]interface{})
            ["active"].(map[string]interface{})
            ["control-1"].(map[string]interface{})
            ["ip"])
}

func TestSecretGenerationBehavior(t *testing.T) {
    instance := "test-cluster"

    // Generate secret (bash: wild-secret-set generates 32-char base64)
    err := secrets.GenerateSecret(instance, "apps.ghost.dbPassword")
    require.NoError(t, err)

    // Verify secret format matches bash
    value, err := secrets.GetSecret(instance, "apps.ghost.dbPassword")
    require.NoError(t, err)

    // Should be 32 characters, alphanumeric
    assert.Equal(t, 32, len(value))
    assert.Regexp(t, "^[a-zA-Z0-9]+$", value)
}
```

#### 6. Real Infrastructure Tests

**Critical:** Must be tested against actual Talos nodes and Kubernetes clusters:

```go
// +build integration

func TestRealClusterBootstrap(t *testing.T) {
    if os.Getenv("INTEGRATION_TEST") != "true" {
        t.Skip("Skipping integration test")
    }

    // This test requires:
    // - 3 physical machines or VMs with Talos
    // - Machines in maintenance mode
    // - Network configuration from environment

    instance := "integration-test-cluster"

    // Create instance
    err := instance.Create(instance)
    require.NoError(t, err)

    // Set required config
    config.SetKey(instance, "cluster.name", "integration-test")
    config.SetKey(instance, "cluster.nodes.control.vip", os.Getenv("TEST_VIP"))
    // ... set other required config

    // Bootstrap cluster
    err = cluster.Bootstrap(ctx, instance)
    require.NoError(t, err)

    // Verify real cluster is functional
    nodes, err := kubectl.GetNodes()
    require.NoError(t, err)
    assert.GreaterOrEqual(t, len(nodes), 1)

    // Cleanup
    t.Cleanup(func() {
        cluster.Destroy(ctx, instance)
    })
}
```

### Implementation Checklist

For each bash script being migrated:

- [ ] **Document exact command sequence** from bash script
- [ ] **Identify all state checks** (config keys, files, kubectl queries)
- [ ] **Map to Go methods** with same external commands
- [ ] **Implement idempotency checks** before operations
- [ ] **Add progress tracking** for multi-step operations
- [ ] **Create rollback logic** for failures
- [ ] **Write integration test** comparing bash vs Go state
- [ ] **Verify log output** matches bash script patterns
- [ ] **Test resume from failure** scenarios
- [ ] **Validate with real cluster** deployment

### Additional Idempotency Patterns

#### Pattern 5: Network Retry Logic

**Bash Pattern:**
```bash
# Retry curl with exponential backoff
for i in {1..5}; do
    if curl -s "https://api.example.com" >/dev/null 2>&1; then
        break
    fi
    sleep $((2 ** i))
done
```

**Go Implementation:**
```go
func RetryWithBackoff(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        if err := fn(); err == nil {
            return nil
        }
        if i < maxRetries-1 {
            backoff := time.Duration(1<<uint(i)) * time.Second
            time.Sleep(backoff)
        }
    }
    return fmt.Errorf("operation failed after %d retries", maxRetries)
}
```

#### Pattern 6: Lock File Management

**Bash Pattern:**
```bash
LOCK_FILE="/tmp/wild-cluster-${CLUSTER_NAME}.lock"
exec 200>"${LOCK_FILE}"
if ! flock -n 200; then
    echo "Another operation is in progress"
    exit 1
fi
# Operations here
flock -u 200
```

**Go Implementation:**
```go
type OperationLock struct {
    lockDir string
}

func (l *OperationLock) Acquire(instance, operation string) (*os.File, error) {
    lockPath := filepath.Join(l.lockDir, fmt.Sprintf("%s-%s.lock", instance, operation))

    // Create lock file
    f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
    if err != nil {
        return nil, err
    }

    // Try to acquire exclusive lock
    if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
        f.Close()
        return nil, fmt.Errorf("another operation is in progress: %w", err)
    }

    return f, nil
}

func (l *OperationLock) Release(f *os.File) error {
    defer f.Close()
    return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
```

#### Pattern 7: Environment Variable Handling

**Bash Pattern:**
```bash
export KUBECONFIG="${WC_HOME}/.kube/config"
export TALOSCONFIG="${WC_HOME}/.talos/config"
kubectl get nodes  # Uses KUBECONFIG
```

**Go Implementation:**
```go
type ToolExecutor struct {
    instance string
    wcHome   string
}

func (t *ToolExecutor) execWithEnv(name string, args ...string) (string, error) {
    cmd := exec.Command(name, args...)

    // Set environment variables like bash scripts
    cmd.Env = append(os.Environ(),
        fmt.Sprintf("KUBECONFIG=%s/.kube/config", t.wcHome),
        fmt.Sprintf("TALOSCONFIG=%s/.talos/config", t.wcHome),
    )

    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

### Bash Construct Mappings

#### Pipes and Process Substitution

**Bash:**
```bash
kubectl get pods -A | grep Running | wc -l
```

**Go:**
```go
func CountRunningPods() (int, error) {
    output, err := kubectl.Get("pods", "-A", "-o", "json")
    if err != nil {
        return 0, err
    }

    var podList struct {
        Items []struct {
            Status struct {
                Phase string `json:"phase"`
            } `json:"status"`
        } `json:"items"`
    }

    json.Unmarshal([]byte(output), &podList)

    count := 0
    for _, pod := range podList.Items {
        if pod.Status.Phase == "Running" {
            count++
        }
    }
    return count, nil
}
```

#### Background Processes

**Bash:**
```bash
long_running_command &
bg_pid=$!
# Do other work
wait $bg_pid
```

**Go:**
```go
func RunInBackground(fn func() error) <-chan error {
    errCh := make(chan error, 1)
    go func() {
        errCh <- fn()
    }()
    return errCh
}

// Usage
errCh := RunInBackground(longRunningOperation)
// Do other work
if err := <-errCh; err != nil {
    return err
}
```

### Signal Handling

**Bash Pattern:**
```bash
trap cleanup EXIT INT TERM

cleanup() {
    echo "Cleaning up..."
    rm -f "${LOCK_FILE}"
}
```

**Go Implementation:**
```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle signals
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        sig := <-sigCh
        log.Printf("Received signal %v, shutting down gracefully", sig)
        cancel()
    }()

    // Run operations with context
    if err := runOperations(ctx); err != nil {
        if err == context.Canceled {
            log.Println("Operation canceled")
        } else {
            log.Fatalf("Operation failed: %v", err)
        }
    }

    // Cleanup
    cleanup()
}
```

### Critical Implementation Notes

1. **Exact Command Replication**: Use same kubectl/talosctl arguments as bash scripts
2. **State File Locations**: Maintain exact same paths as v.PoC (`setup/`, `apps/`, etc.)
3. **Error Messages**: Match bash script error output for user familiarity
4. **Log Verbosity**: Match bash script `print_info`, `print_success` patterns
5. **User Prompts**: Maintain same interactive prompts as bash scripts
6. **File Permissions**: Match bash behavior (secrets.yaml = 600, etc.)
7. **Atomic Operations**: Use file locking for concurrent safety
8. **Backup Before Modify**: Keep backup of configs before updates
9. **Timeout Values**: Match exact timeout values from bash scripts
10. **Environment Variables**: Propagate same env vars as bash scripts
11. **Signal Handling**: Handle SIGINT/SIGTERM gracefully like bash trap
12. **Stdout/Stderr Separation**: Capture separately for accurate logging

### Enhanced Implementation Checklist

For each bash script being migrated:

- [ ] **Document exact command sequence** from bash script
- [ ] **Identify all state checks** (config keys, files, kubectl queries)
- [ ] **Map to Go methods** with same external commands
- [ ] **Implement idempotency checks** before each step
- [ ] **Add operation state tracking** for UI progress display
- [ ] **Create rollback logic** for failures (where applicable)
- [ ] **Write state-based validation tests** (verify files, config keys, permissions)
- [ ] **Write idempotency tests** (run twice, compare state; run after partial completion)
- [ ] **Verify log output** matches bash script patterns
- [ ] **Manual verification** with real cluster deployment
- [ ] **Verify timeout values** match bash scripts
- [ ] **Test concurrent operation safety** (lock files prevent conflicts)
- [ ] **Validate network partition handling** (operations fail gracefully, can re-run)
- [ ] **Check signal handling** (SIGTERM/SIGINT cleanup)
- [ ] **Test lock file behavior** under contention
- [ ] **Verify environment variable propagation** (KUBECONFIG, TALOSCONFIG)
- [ ] **Compare performance** with bash script timing
- [ ] **Test degraded network conditions** (retries work, idempotency handles partial completion)

---

## Appendix

### Bash Script to API Endpoint Mapping

| Bash Script | API Endpoint | Phase |
|------------|--------------|-------|
| wild-init | POST /api/v1/instances | 1 |
| wild-config | GET /api/v1/instances/{name}/config | 1 |
| wild-config-set | PUT /api/v1/instances/{name}/config | 1 |
| wild-secret | GET /api/v1/instances/{name}/secrets | 1 |
| wild-secret-set | PUT /api/v1/instances/{name}/secrets | 1 |
| wild-compile-template | (internal: tools/gomplate.go) | 1 |
| wild-node-detect | POST /api/v1/instances/{name}/nodes/discover | 2 |
| wild-node-setup | POST /api/v1/instances/{name}/nodes/{node}/setup | 2 |
| wild-cluster-config-generate | POST /api/v1/instances/{name}/cluster/config/generate | 3 |
| wild-setup-cluster | POST /api/v1/instances/{name}/cluster/bootstrap | 3 |
| wild-health | GET /api/v1/instances/{name}/cluster/health | 3 |
| wild-setup-services | POST /api/v1/instances/{name}/services/install-all | 4 |
| wild-service-setup | POST /api/v1/instances/{name}/services | 4 |
| wild-dashboard-token | (handled by Kubernetes secrets) | 4 |
| wild-apps-list | GET /api/v1/apps | 5 |
| wild-app-add | POST /api/v1/instances/{name}/apps | 5 |
| wild-app-deploy | POST /api/v1/instances/{name}/apps/{app}/deploy | 5 |
| wild-app-delete | DELETE /api/v1/instances/{name}/apps/{app} | 5 |
| wild-app-doctor | POST /api/v1/instances/{name}/apps/{app}/doctor | 5 |
| wild-backup | (future: backup endpoints) | 6 |
| wild-dnsmasq-install | PUT /api/v1/dnsmasq/config | System |

### Technology Stack

**Backend (wildd):**
- Go 1.21+
- gorilla/mux (HTTP routing)
- Standard library (exec, os, io)

**External Tools:**
- kubectl (Kubernetes operations)
- talosctl (Talos operations)
- gomplate (template processing)
- yq (YAML manipulation)
- kustomize (Kubernetes config management)

**Storage:**
- File-based (YAML files)
- No database required

**CLI (wild-cli):**
- Go 1.21+
- cobra (CLI framework)
- viper (configuration)

**Web UI (wild-app):**
- React 18+
- TypeScript
- TailwindCSS
- React Router

---

## Conclusion

This migration plan provides a clear path from the Wild Cloud v.PoC bash scripts to the new Wild Central daemon architecture. By following the vertical slice approach, we ensure working software at each phase while maintaining the core simplicity and directness that makes Wild Cloud effective.

The five-phase implementation strategy allows for:
- **Incremental validation** - Test each phase before moving forward
- **Parallel development** - CLI and Web UI can develop alongside daemon
- **Risk mitigation** - Identify issues early in simpler phases
- **User feedback** - Can release Phase 1-3 for early testing

The architecture preserves the best aspects of the v.PoC (template-driven, file-based, direct tool integration) while enabling the multi-instance, network-accessible management that Wild Central requires.
