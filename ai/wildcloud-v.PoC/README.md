# Wild Cloud Agent Context Documentation

This directory contains comprehensive documentation about the Wild Cloud project, designed to provide AI agents (like Claude Code) with the context needed to effectively help users with Wild Cloud development, deployment, and operations.

## Documentation Overview

### 📚 Core Documentation Files

1. **[overview.md](./overview.md)** - Complete project introduction and getting started guide
   - What Wild Cloud is and why it exists
   - Technology stack and architecture overview
   - Quick start guide and common use cases
   - Best practices and troubleshooting

2. **[bin-scripts.md](./bin-scripts.md)** - Complete CLI reference
   - All 34+ `wild-*` commands with usage examples
   - Command categories (setup, apps, config, operations)
   - Script dependencies and execution order
   - Common usage patterns

3. **[setup-process.md](./setup-process.md)** - Infrastructure deployment deep dive
   - Complete setup phases and dependencies
   - Talos Linux and Kubernetes cluster deployment
   - Core services installation (MetalLB, Traefik, cert-manager, etc.)
   - Network configuration and DNS management

4. **[apps-system.md](./apps-system.md)** - Application management system
   - App structure and lifecycle management
   - Template system and configuration
   - Available applications and their features
   - Creating custom applications

5. **[configuration-system.md](./configuration-system.md)** - Configuration and secrets management
   - `config.yaml` and `secrets.yaml` structure
   - Template processing with gomplate
   - Environment setup and validation
   - Security best practices

6. **[project-architecture.md](./project-architecture.md)** - Project structure and organization
   - Wild Cloud repository structure
   - User cloud directory layout
   - File permissions and security model
   - Development and deployment patterns

## Quick Reference Guide

### Essential Commands
```bash
# Setup & Initialization
wild-init                    # Initialize new cloud
wild-setup                   # Complete deployment
wild-health                  # System health check

# Application Management
wild-apps-list              # List available apps
wild-app-add <app>          # Configure app
wild-app-deploy <app>       # Deploy app

# Configuration
wild-config <key>           # Read config
wild-config-set <key> <val> # Set config
wild-secret <key>           # Read secret
```

### Key File Locations

**Wild Cloud Repository** (`WC_ROOT`):
- `bin/` - All CLI commands
- `apps/` - Application templates
- `setup/` - Infrastructure templates
- `docs/` - Documentation

**User Cloud Directory** (`WC_HOME`):
- `config.yaml` - Main configuration
- `secrets.yaml` - Sensitive data
- `apps/` - Deployed app configs
- `.wildcloud/` - Project marker

### Application Categories

- **Content**: Ghost (blog), Discourse (forum)
- **Media**: Immich (photos)
- **Development**: Gitea (Git), Docker Registry
- **Databases**: PostgreSQL, MySQL, Redis
- **AI/ML**: vLLM (LLM inference)

## Technology Stack Summary

### Core Infrastructure
- **Talos Linux** - Immutable Kubernetes OS
- **Kubernetes** - Container orchestration
- **MetalLB** - Load balancing
- **Traefik** - Ingress/reverse proxy
- **Longhorn** - Distributed storage
- **cert-manager** - TLS certificates

### Management Tools
- **gomplate** - Template processing
- **Kustomize** - Configuration management
- **restic** - Backup system
- **kubectl/talosctl** - Cluster management

## Common Agent Tasks

### When Users Ask About...

**"How do I deploy X?"**
- Check apps-system.md for application management
- Look for X in available applications list
- Reference app deployment lifecycle

**"Setup isn't working"**
- Review setup-process.md for troubleshooting
- Check bin-scripts.md for command options
- Verify prerequisites and dependencies

**"How do I configure Y?"**
- Check configuration-system.md for config management
- Look at project-architecture.md for file locations
- Review template processing documentation

**"What does wild-X command do?"**
- Reference bin-scripts.md for complete command documentation
- Check command categories and usage patterns
- Look at dependencies between commands

### Development Tasks

**Creating New Apps**:
1. Review apps-system.md "Creating Custom Apps" section
2. Follow Wild Cloud app structure conventions
3. Use project-architecture.md for file organization
4. Test with standard app deployment workflow

**Modifying Infrastructure**:
1. Check setup-process.md for infrastructure components
2. Review configuration-system.md for template processing
3. Understand project-architecture.md file relationships
4. Test changes carefully in development environment

**Troubleshooting Issues**:
1. Use bin-scripts.md for diagnostic commands
2. Check setup-process.md for component validation
3. Review configuration-system.md for config problems
4. Reference apps-system.md for application issues

## Best Practices for Agents

### Understanding User Context
- Always check if they're in a Wild Cloud directory (look for `.wildcloud/`)
- Determine if they need setup help vs operational help
- Consider their experience level (beginner vs advanced)
- Check what applications they're trying to deploy

### Providing Help
- Reference specific documentation sections for detailed info
- Provide exact command syntax from bin-scripts.md
- Explain prerequisites and dependencies
- Offer validation steps to verify success

### Safety Considerations
- Always recommend testing in development first
- Warn about destructive operations (delete, reset)
- Emphasize backup importance before major changes
- Explain security implications of configuration changes

### Common Gotchas
- `secrets.yaml` has restricted permissions (600)
- Templates need processing before deployment
- Dependencies between applications must be satisfied
- Node hardware detection requires maintenance mode boot

## Documentation Maintenance

This documentation should be updated when:
- New commands are added to `bin/`
- New applications are added to `apps/`
- Infrastructure components change
- Configuration schema evolves
- Best practices are updated

Each documentation file includes:
- Complete coverage of its topic area
- Practical examples and use cases
- Troubleshooting guidance
- References to related documentation

This comprehensive context should enable AI agents to provide expert-level assistance with Wild Cloud projects across all aspects of the system.