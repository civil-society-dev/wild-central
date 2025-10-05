# Wild Cloud Setup Process & Infrastructure

Wild Cloud provides a complete, production-ready Kubernetes infrastructure designed for personal use. It combines enterprise-grade technologies to create a self-hosted cloud platform with automated deployment, HTTPS certificates, and web management interfaces.

## Setup Phases Overview

The Wild Cloud setup follows a sequential, dependency-aware process:

1. **Environment Setup** - Install required tools and dependencies
2. **DNS/Network Foundation** - Set up dnsmasq for DNS and PXE booting
3. **Cluster Infrastructure** - Deploy Talos Linux nodes and Kubernetes cluster
4. **Cluster Services** - Install core services (ingress, storage, certificates, etc.)

## Phase 1: Environment Setup

### Dependencies Installation
**Script**: `scripts/setup-utils.sh`

**Required Tools**:
- `kubectl` - Kubernetes CLI
- `gomplate` - Template processor for configuration
- `kustomize` - Kubernetes configuration management
- `yq` - YAML processor
- `restic` - Backup tool
- `talosctl` - Talos Linux cluster management

### Project Initialization
**Command**: `wild-init`

Creates the basic Wild Cloud directory structure:
- `.wildcloud/` - Project marker and cache
- `config.yaml` - Main configuration file
- `secrets.yaml` - Sensitive data storage
- Basic project scaffolding

## Phase 2: DNS/Network Foundation

### dnsmasq Infrastructure
**Location**: `setup/dnsmasq/`
**Requirements**: Dedicated Linux machine with static IP

**Services Provided**:
1. **LAN DNS Server**
   - Forwards internal domains (`*.internal.domain.com`) to cluster
   - Forwards external domains (`*.domain.com`) to cluster
   - Provides DNS resolution for entire network

2. **PXE Boot Server**
   - Enables network booting for cluster node installation
   - DHCP/TFTP services for Talos Linux deployment
   - Automated node provisioning

**Network Configuration Example**:
```yaml
network:
  subnet: 192.168.1.0/24
  gateway: 192.168.1.1
  dnsmasq_ip: 192.168.1.50
  dhcp_range: 192.168.1.100-200
  metallb_pool: 192.168.1.80-89
  control_plane_vip: 192.168.1.90
  node_ips: 192.168.1.91-93
```

## Phase 3: Cluster Infrastructure Setup

### Talos Linux Foundation
**Command**: `wild-setup-cluster`

**Talos Configuration**:
- **Version**: v1.11.0 (configurable)
- **Immutable OS**: Designed specifically for Kubernetes
- **System Extensions**:
  - Intel microcode updates
  - iSCSI tools for storage
  - gVisor container runtime
  - NVIDIA GPU support (optional)
  - Additional system utilities

### Cluster Setup Process

#### 1. Configuration Generation
**Script**: `wild-cluster-config-generate`

- Generates base Talos configurations (`controlplane.yaml`, `worker.yaml`)
- Creates cluster secrets using `talosctl gen config`
- Establishes foundation for all node configurations

#### 2. Node Setup (Atomic Operations)
**Script**: `wild-node-setup <node-name> [options]`

**Complete Node Lifecycle Management**:
- **Hardware Detection**: Discovers network interfaces and storage devices
- **Configuration Generation**: Creates node-specific patches and final configs
- **Deployment**: Applies Talos configuration to the node

**Options**:
- `--detect`: Force hardware re-detection
- `--no-deploy`: Generate configuration only, skip deployment

**Integration with Cluster Setup**:
- `wild-setup-cluster` automatically calls `wild-node-setup` for each node
- Individual node failures don't break cluster setup
- Clear retry instructions for failed nodes

### Cluster Architecture

**Control Plane**:
- 3 nodes for high availability
- Virtual IP (VIP) for load balancing
- etcd distributed across all control plane nodes

**Worker Nodes**:
- Variable count (configured during setup)
- Dedicated workload execution
- Storage participation via Longhorn

**Networking**:
- All nodes on same LAN segment
- Sequential IP assignment
- MetalLB integration for load balancing

## Phase 4: Cluster Services Installation

### Services Deployment Process
**Command**: `wild-setup-services [options]`
- **`--fetch`**: Fetch fresh templates before setup
- **`--no-deploy`**: Configure only, skip deployment

**New Architecture**: Per-service atomic operations
- Uses `wild-service-setup <service>` for each service in dependency order
- Each service handles complete lifecycle: fetch → configure → deploy
- Dependency validation before each service deployment
- Fail-fast with clear recovery instructions

**Individual Service Management**: `wild-service-setup <service> [options]`
- **Default**: Configure and deploy using existing templates
- **`--fetch`**: Fetch fresh templates before setup
- **`--no-deploy`**: Configure only, skip deployment

### Core Services (Installed in Order)

#### 1. MetalLB Load Balancer
**Location**: `setup/cluster-services/metallb/`

- **Purpose**: Provides load balancing for bare metal clusters
- **Functionality**: Assigns external IPs to Kubernetes services
- **Configuration**: IP address pool from local network range
- **Integration**: Foundation for ingress traffic routing

#### 2. Longhorn Distributed Storage
**Location**: `setup/cluster-services/longhorn/`

- **Purpose**: Distributed block storage for persistent volumes
- **Features**:
  - Cross-node data replication
  - Snapshot and backup capabilities
  - Volume expansion and management
  - Web-based management interface
- **Storage**: Uses local disks from all cluster nodes

#### 3. Traefik Ingress Controller
**Location**: `setup/cluster-services/traefik/`

- **Purpose**: HTTP/HTTPS reverse proxy and ingress controller
- **Features**:
  - Automatic service discovery
  - TLS termination
  - Load balancing and routing
  - Gateway API support
- **Integration**: Works with MetalLB for external traffic

#### 4. CoreDNS
**Location**: `setup/cluster-services/coredns/`

- **Purpose**: DNS resolution for cluster services
- **Integration**: Connects with external DNS providers
- **Functionality**: Service discovery and DNS forwarding

#### 5. cert-manager
**Location**: `setup/cluster-services/cert-manager/`

- **Purpose**: Automatic TLS certificate management
- **Features**:
  - Let's Encrypt integration
  - Automatic certificate issuance and renewal
  - Multiple certificate authorities support
  - Certificate lifecycle management

#### 6. ExternalDNS
**Location**: `setup/cluster-services/externaldns/`

- **Purpose**: Automatic DNS record management
- **Functionality**:
  - Syncs Kubernetes services with DNS providers
  - Automatic A/CNAME record creation
  - Supports multiple DNS providers (Cloudflare, Route53, etc.)

#### 7. Kubernetes Dashboard
**Location**: `setup/cluster-services/kubernetes-dashboard/`

- **Purpose**: Web UI for cluster management
- **Access**: `https://dashboard.internal.domain.com`
- **Authentication**: Token-based access via `wild-dashboard-token`
- **Features**: Resource management, monitoring, troubleshooting

#### 8. NFS Storage (Optional)
**Location**: `setup/cluster-services/nfs/`

- **Purpose**: Network file system for shared storage
- **Use Cases**: Media storage, backups, shared data
- **Integration**: Mounted as persistent volumes in applications

#### 9. Docker Registry
**Location**: `setup/cluster-services/docker-registry/`

- **Purpose**: Private container registry
- **Features**: Store custom images locally
- **Integration**: Used by applications and CI/CD pipelines

## Infrastructure Components Deep Dive

### DNS and Domain Architecture

```
Internet → External DNS → MetalLB LoadBalancer → Traefik → Kubernetes Services
                               ↑
                           Internal DNS (dnsmasq)
                               ↑
                         Internal Network
```

**Domain Types**:
- **External**: `app.domain.com` - Public-facing services
- **Internal**: `app.internal.domain.com` - Admin interfaces only
- **Resolution**: dnsmasq forwards all domain traffic to cluster

### Certificate and TLS Management

**Automatic Certificate Flow**:
1. Service deployed with ingress annotation
2. cert-manager detects certificate requirement
3. Let's Encrypt challenge initiated
4. Certificate issued and stored in Kubernetes secret
5. Traefik uses certificate for TLS termination
6. Automatic renewal before expiration

### Storage Architecture

**Longhorn Distributed Storage**:
- Block-level replication across nodes
- Default 3-replica policy for data durability
- Automatic failover and recovery
- Snapshot and backup capabilities
- Web UI for management and monitoring

**Storage Classes**:
- `longhorn` - Default replicated storage
- `longhorn-single` - Single replica for non-critical data
- `nfs` - Shared network storage (if configured)

### Network Traffic Flow

**External Request Flow**:
1. DNS resolution via dnsmasq → cluster IP
2. Traffic hits MetalLB load balancer
3. MetalLB forwards to Traefik ingress
4. Traefik terminates TLS and routes to service
5. Service forwards to appropriate pod
6. Response follows reverse path

### High Availability Features

**Control Plane HA**:
- 3 control plane nodes with leader election
- Virtual IP for API server access
- etcd cluster with automatic failover
- Distributed workload scheduling

**Storage HA**:
- Longhorn 3-way replication
- Automatic replica placement across nodes
- Node failure recovery
- Data integrity verification

**Networking HA**:
- MetalLB speaker pods on all nodes
- Automatic load balancer failover
- Multiple ingress controller replicas

## Hardware Requirements

### Minimum Specifications
- **Nodes**: 3 control plane + optional workers
- **RAM**: 8GB minimum per node (16GB+ recommended)
- **CPU**: 4 cores minimum per node
- **Storage**: 100GB+ local storage per node
- **Network**: Gigabit ethernet connectivity

### Network Requirements
- All nodes on same LAN segment
- Static IP assignments or DHCP reservations
- dnsmasq server accessible by all nodes
- Internet connectivity for image pulls and Let's Encrypt

### Recommended Hardware
- **Control Plane**: 16GB RAM, 8 cores, 200GB NVMe SSD
- **Workers**: 32GB RAM, 16 cores, 500GB NVMe SSD
- **Network**: Dedicated VLAN or network segment
- **Redundancy**: UPS protection, dual network interfaces

## Configuration Management

### Configuration Files
- `config.yaml` - Main configuration (domains, network, apps)
- `secrets.yaml` - Sensitive data (passwords, API keys, certificates)
- `.wildcloud/` - Cache and temporary files

### Template System
**gomplate Integration**:
- All configurations processed as templates
- Access to config and secrets via template variables
- Dynamic configuration generation
- Environment-specific customization

### Configuration Commands
```bash
# Read configuration values
wild-config cluster.name
wild-config apps.ghost.domain

# Set configuration values
wild-config-set cloud.domain "example.com"
wild-config-set cluster.nodeCount 5

# Secret management
wild-secret apps.database.password
wild-secret-set apps.api.key "secret-value"
```

## Setup Commands Reference

### Complete Setup
```bash
wild-init              # Initialize project
wild-setup             # Complete automated setup
```

### Phase-by-Phase Setup
```bash
wild-setup-cluster     # Cluster infrastructure only
wild-setup-services    # Cluster services only
```

### Individual Operations
```bash
wild-cluster-config-generate     # Generate base configs
wild-node-setup <node-name>      # Complete node setup (detect → configure → deploy)
wild-node-setup <node-name> --detect    # Force hardware re-detection
wild-node-setup <node-name> --no-deploy # Configuration only
wild-dashboard-token            # Get dashboard access
wild-health                     # System health check
```

## Troubleshooting and Validation

### Health Checks
```bash
wild-health                     # Overall system status
kubectl get nodes              # Node status
kubectl get pods -A            # All pod status
talosctl health                # Talos cluster health
```

### Service Validation
```bash
kubectl get svc -n metallb-system     # MetalLB status
kubectl get pods -n longhorn-system   # Storage status
kubectl get pods -n traefik          # Ingress status
kubectl get certificates -A           # Certificate status
```

### Log Analysis
```bash
talosctl logs -f machined             # Talos system logs
kubectl logs -n traefik deployment/traefik  # Ingress logs
kubectl logs -n cert-manager deployment/cert-manager  # Certificate logs
```

This comprehensive setup process creates a production-ready personal cloud infrastructure with enterprise-grade reliability, security, and management capabilities.