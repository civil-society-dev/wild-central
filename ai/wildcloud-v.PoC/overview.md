# Wild Cloud Overview

Wild Cloud is a complete, production-ready Kubernetes infrastructure designed for personal use. It combines enterprise-grade technologies to create a self-hosted cloud platform with automated deployment, HTTPS certificates, and web management interfaces.

## What is Wild Cloud?

### Vision
In a world where digital lives are increasingly controlled by large corporations, Wild Cloud puts you back in control by providing:

- **Privacy**: Your data stays on your hardware, under your control
- **Ownership**: No subscription fees or sudden price increases
- **Freedom**: Run the apps you want, the way you want them
- **Learning**: Gain valuable skills in modern cloud technologies
- **Resilience**: Reduce reliance on third-party services that can disappear

### Core Capabilities

**Complete Infrastructure Stack**:
- Kubernetes cluster with Talos Linux
- Automatic HTTPS certificates via Let's Encrypt
- Load balancing with MetalLB
- Ingress routing with Traefik
- Distributed storage with Longhorn
- DNS management with CoreDNS and ExternalDNS

**Application Platform**:
- One-command application deployment
- Pre-built apps for common self-hosted services
- Automatic database setup and configuration
- Integrated backup and restore system
- Web-based management interfaces

**Enterprise Features**:
- High availability and fault tolerance
- Automated certificate management
- Network policies and security contexts
- Monitoring and observability
- Infrastructure as code principles

## Technology Stack

### Core Infrastructure
- **Talos Linux** - Immutable OS designed for Kubernetes
- **Kubernetes** - Container orchestration platform
- **MetalLB** - Load balancer for bare metal deployments
- **Traefik** - HTTP reverse proxy and ingress controller
- **Longhorn** - Distributed block storage system
- **cert-manager** - Automatic TLS certificate management

### Supporting Services
- **CoreDNS** - DNS server for service discovery
- **ExternalDNS** - Automatic DNS record management
- **Kubernetes Dashboard** - Web UI for cluster management
- **restic** - Backup solution with deduplication
- **gomplate** - Template processor for configurations

### Development Tools
- **Kustomize** - Kubernetes configuration management
- **kubectl** - Kubernetes command line interface
- **talosctl** - Talos Linux management tool
- **Bats** - Testing framework for bash scripts

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                          Internet                               │
└─────────────────┬───────────────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────────────┐
│                      DNS Provider                               │
│              (Cloudflare, Route53, etc.)                        │
└─────────────────┬───────────────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────────────┐
│                    Your Network                                 │
│  ┌─────────────┐  ┌─────────────────────────────────────────┐   │
│  │   dnsmasq   │  │          Kubernetes Cluster             │   │
│  │   Server    │  │  ┌─────────────┐ ┌─────────────────┐    │   │
│  │             │  │  │  MetalLB    │ │    Traefik      │    │   │
│  │ DNS + DHCP  │  │  │ LoadBalancer│ │   Ingress       │    │   │
│  └─────────────┘  │  └─────────────┘ └─────────────────┘    │   │
│                   │  ┌───────────────────────────────────┐  │   │
│                   │  │           Applications            │  │   │
│                   │  │   Ghost, Immich, Gitea, vLLM...   │  │   │
│                   │  └───────────────────────────────────┘  │   │
│                   └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### Traffic Flow
1. **External Request** → DNS resolution via provider
2. **DNS Response** → Points to your cluster's external IP
3. **Network Request** → Hits MetalLB load balancer
4. **Load Balancer** → Routes to Traefik ingress controller
5. **Ingress Controller** → Terminates TLS and routes to application
6. **Application** → Serves content from Kubernetes pod

## Getting Started

### Prerequisites

**Hardware Requirements**:
- Minimum 3 nodes for high availability
- 8GB RAM per node (16GB+ recommended)
- 100GB+ storage per node
- Gigabit network connectivity
- x86_64 architecture

**Network Requirements**:
- All nodes on same network segment
- One dedicated machine for dnsmasq (can be lightweight)
- Static IP assignments or DHCP reservations
- Internet connectivity for downloads and certificates

### Quick Start Guide

#### 1. Install Dependencies
```bash
# Clone Wild Cloud repository
git clone https://github.com/your-org/wild-cloud
cd wild-cloud

# Install required tools
scripts/setup-utils.sh
```

#### 2. Initialize Your Cloud
```bash
# Create and initialize new cloud directory
mkdir my-cloud && cd my-cloud
wild-init

# Follow interactive setup prompts for:
# - Domain name configuration
# - Email for certificates
# - Network settings
```

#### 3. Deploy Infrastructure
```bash
# Complete automated setup
wild-setup

# Or step-by-step:
wild-setup-cluster     # Deploy Kubernetes cluster
wild-setup-services    # Install core services
```

#### 4. Deploy Your First App
```bash
# List available applications
wild-apps-list

# Deploy a blog
wild-app-add ghost
wild-app-deploy ghost

# Access at https://ghost.yourdomain.com
```

#### 5. Verify Deployment
```bash
# Check system health
wild-health

# Access Kubernetes dashboard
wild-dashboard-token
# Visit https://dashboard.internal.yourdomain.com
```

## Key Concepts

### Configuration Management

Wild Cloud uses a dual-file configuration system:

**`config.yaml`** - Non-sensitive settings:
```yaml
cloud:
  domain: "example.com"
  email: "admin@example.com"
apps:
  ghost:
    domain: "blog.example.com"
    storage: "10Gi"
```

**`secrets.yaml`** - Sensitive data (auto-generated):
```yaml
apps:
  ghost:
    dbPassword: "secure-random-password"
  postgresql:
    rootPassword: "another-secure-password"
```

### Template System

All configurations are templates processed with gomplate:

**Before Processing** (in repository):
```yaml
domain: {{ .apps.ghost.domain }}
storage: {{ .apps.ghost.storage | default "5Gi" }}
```

**After Processing** (in your cloud):
```yaml
domain: blog.example.com
storage: 10Gi
```

### Application Lifecycle

1. **Discovery**: `wild-apps-list` - Browse available apps
2. **Configuration**: `wild-app-add app-name` - Configure and prepare application
3. **Deployment**: `wild-app-deploy app-name` - Deploy to cluster
4. **Operations**: `wild-app-doctor app-name` - Monitor and troubleshoot

## Available Applications

### Content Management & Publishing
- **Ghost** - Modern publishing platform for blogs and websites
- **Discourse** - Community discussion platform with modern features

### Media & File Management
- **Immich** - Self-hosted photo and video backup solution

### Development Tools
- **Gitea** - Self-hosted Git service with web interface
- **Docker Registry** - Private container image registry

### Communication & Marketing
- **Keila** - Newsletter and email marketing platform
- **Listmonk** - High-performance newsletter and mailing list manager

### Databases & Caching
- **PostgreSQL** - Advanced open-source relational database
- **MySQL** - Popular relational database management system
- **Redis** - In-memory data structure store and cache
- **Memcached** - Distributed memory caching system

### AI & Machine Learning
- **vLLM** - High-performance LLM inference server with OpenAI-compatible API

## Core Commands Reference

### Setup & Initialization
```bash
wild-init                    # Initialize new cloud directory
wild-setup                   # Complete infrastructure deployment
wild-setup-cluster          # Deploy Kubernetes cluster only
wild-setup-services         # Deploy cluster services only
```

### Application Management
```bash
wild-apps-list              # List available applications
wild-app-add <app>          # Configure application
wild-app-deploy <app>       # Deploy to cluster
wild-app-delete <app>       # Remove application
wild-app-doctor <app>       # Run diagnostics
```

### Configuration Management
```bash
wild-config <key>           # Read configuration value
wild-config-set <key> <val> # Set configuration value
wild-secret <key>           # Read secret value
wild-secret-set <key> <val> # Set secret value
```

### Operations & Monitoring
```bash
wild-health                 # System health check
wild-dashboard-token        # Get dashboard access token
wild-backup                 # Backup system and apps
wild-app-backup <app>       # Backup specific application
```

## Best Practices

### Security
- Never commit `secrets.yaml` to version control
- Use strong, unique passwords for all services
- Regularly update system and application images
- Monitor certificate expiration and renewal
- Implement network policies for production workloads

### Configuration Management
- Store `config.yaml` in version control with proper .gitignore
- Document configuration changes in commit messages
- Use branches for experimental configurations
- Backup configuration files before major changes
- Test configuration changes in development environment

### Operations
- Monitor cluster health with `wild-health`
- Set up regular backup schedules with `wild-backup`
- Keep applications updated with latest security patches
- Monitor disk usage and expand storage as needed
- Document custom configurations and procedures

### Development
- Follow Wild Cloud app structure conventions
- Use proper Kubernetes security contexts
- Include comprehensive health checks and probes
- Test applications thoroughly before deployment
- Document application-specific configuration requirements

## Common Use Cases

### Personal Blog/Website
```bash
# Deploy Ghost blog with custom domain
wild-config-set apps.ghost.domain "blog.yourdomain.com"
wild-app-add ghost
wild-app-deploy ghost
```

### Photo Management
```bash
# Deploy Immich for photo backup and management
wild-app-add postgresql immich
wild-app-deploy postgresql immich
```

### Development Environment
```bash
# Set up Git hosting and container registry
wild-app-add gitea docker-registry
wild-app-deploy gitea docker-registry
```

### AI/ML Workloads
```bash
# Deploy vLLM for local AI inference
wild-config-set apps.vllm.model "Qwen/Qwen2.5-7B-Instruct"
wild-app-add vllm
wild-app-deploy vllm
```

## Troubleshooting

### Common Issues

**Cluster Not Responding**:
```bash
# Check node status
kubectl get nodes
talosctl health

# Verify network connectivity
ping <node-ip>
```

**Applications Not Starting**:
```bash
# Check pod status
kubectl get pods -n <app-namespace>

# View logs
kubectl logs deployment/<app-name> -n <app-namespace>

# Run diagnostics
wild-app-doctor <app-name>
```

**Certificate Issues**:
```bash
# Check certificate status
kubectl get certificates -A

# View cert-manager logs
kubectl logs -n cert-manager deployment/cert-manager
```

**DNS Problems**:
```bash
# Test DNS resolution
nslookup <app-domain>

# Check external-dns logs
kubectl logs -n external-dns deployment/external-dns
```

### Getting Help

**Documentation**:
- Check `docs/` directory for detailed guides
- Review application-specific README files
- Consult Kubernetes and Talos documentation

**Community Support**:
- Report issues on GitHub repository
- Join community forums and discussions
- Share configurations and troubleshooting tips

**Professional Support**:
- Consider professional services for production deployments
- Engage with cloud infrastructure consultants
- Participate in training and certification programs

## Advanced Topics

### Custom Applications

Create your own Wild Cloud applications:

1. **Create App Directory**: `apps/myapp/`
2. **Define Manifest**: Include metadata and configuration defaults
3. **Create Templates**: Kubernetes resources with gomplate variables
4. **Test Deployment**: Use standard Wild Cloud workflow
5. **Share**: Contribute back to the community

### Multi-Environment Deployments

Manage multiple Wild Cloud instances:

- **Development**: Single-node cluster for testing
- **Staging**: Multi-node cluster mirroring production
- **Production**: Full HA cluster with monitoring and backups

### Integration with External Services

Extend Wild Cloud capabilities:

- **External DNS Providers**: Cloudflare, Route53, Google DNS
- **Backup Storage**: S3, Google Cloud Storage, Azure Blob
- **Monitoring**: Prometheus, Grafana, AlertManager
- **CI/CD**: GitLab CI, GitHub Actions, Jenkins

### Performance Optimization

Optimize for your workloads:

- **Resource Allocation**: CPU and memory limits/requests
- **Storage Performance**: NVMe SSDs, storage classes
- **Network Optimization**: Network policies, service mesh
- **Scaling**: Horizontal pod autoscaling, cluster autoscaling

Wild Cloud provides a solid foundation for personal cloud infrastructure while maintaining the flexibility to grow and adapt to changing needs. Whether you're running a simple blog or a complex multi-service application, Wild Cloud's enterprise-grade technologies ensure your infrastructure is reliable, secure, and maintainable.