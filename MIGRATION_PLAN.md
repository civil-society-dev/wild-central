# Wild Cloud v.PoC to Wild Central Migration Plan

## Executive Summary

This document outlines the comprehensive migration plan from the bash-based Wild Cloud v.PoC system to the new Go-based Wild Central daemon architecture. The migration follows a **vertical slice approach**, implementing complete end-to-end functionality in five phases, ensuring working software at each stage.

**Key Changes:**
- **34+ bash scripts** → **Single Go daemon (wildd) + Go CLI (wild-cli)**
- **Single WC_HOME** → **Multiple instance management** via `$WILD_CENTRAL_DATA`
- **Direct script execution** → **REST API** with CLI and web clients
- **File-based operations** → **API-driven** with file storage backend
- **Local context** → **Multi-cloud context** switching

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

## Success Metrics

### Phase 1 Success
- [ ] Can manage 5+ instances simultaneously
- [ ] Config operations complete in < 100ms
- [ ] Template processing works for all existing apps
- [ ] Context switching works seamlessly

### Phase 2 Success
- [ ] Node discovery finds nodes in < 30 seconds
- [ ] Hardware detection is accurate
- [ ] Node setup completes successfully
- [ ] Async operations track correctly

### Phase 3 Success
- [ ] Cluster bootstraps in < 10 minutes
- [ ] All health checks pass
- [ ] kubectl commands work correctly
- [ ] Supports 3-5 node clusters

### Phase 4 Success
- [ ] All core services install successfully
- [ ] Service dependencies resolve correctly
- [ ] Certificates are issued automatically
- [ ] Services are accessible via ingress

### Phase 5 Success
- [ ] Can deploy all existing apps from v.PoC
- [ ] App dependencies work correctly
- [ ] Apps are accessible via configured domains
- [ ] App status reflects actual state

### Overall Success
- [ ] Complete feature parity with v.PoC bash scripts
- [ ] wild-cli provides all functionality
- [ ] wild-app provides intuitive UI
- [ ] Can manage multiple clusters from single daemon
- [ ] Setup time comparable to v.PoC
- [ ] No regressions in functionality

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
