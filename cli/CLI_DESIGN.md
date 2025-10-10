# Wild CLI Design

## Overview

`wild-cli` is the command-line interface to the Wild Central daemon (wildd). It replaces the bash scripts from Wild Cloud v.PoC with a single Go binary that communicates with the daemon via REST API.

## Architecture

```
wild-cli (Go binary)
    ↓ HTTP/REST
wildd (daemon on local network)
    ↓ kubectl/talosctl
Kubernetes/Talos cluster
```

## Command Structure

### Context Management
```
wild context                      # Show current context
wild context list                 # List all contexts
wild context use <name>           # Switch context
wild context set <name> <url>     # Add/update context
wild context delete <name>        # Remove context
```

### Instance Management
```
wild instance create <name>       # Create new instance
wild instance list                # List all instances
wild instance delete <name>       # Delete instance
wild instance show <name>         # Show instance details
```

### Configuration
```
wild config get <key>             # Get config value
wild config set <key> <value>     # Set config value
wild config show                  # Show full config
wild config edit                  # Edit config in $EDITOR
```

### Secrets
```
wild secret get <key>             # Get secret (redacted)
wild secret set <key> <value>     # Set secret
wild secret generate <key>        # Generate random secret
```

### Node Management
```
wild node discover                # Discover nodes on network
wild node list                    # List configured nodes
wild node add <hostname> <role>   # Add node
wild node setup <hostname>        # Setup Talos on node
wild node delete <hostname>       # Remove node
wild node show <hostname>         # Show node details
```

### PXE Assets
```
wild pxe list                     # List PXE assets
wild pxe download <asset>         # Download PXE asset
wild pxe delete <asset>           # Delete PXE asset
```

### Cluster Operations
```
wild cluster bootstrap            # Bootstrap cluster
wild cluster status               # Get cluster status
wild cluster health               # Check cluster health
wild cluster config generate      # Generate cluster config
wild cluster kubeconfig           # Get kubeconfig
wild cluster talosconfig          # Get talosconfig
wild cluster reset                # Reset cluster
```

### Service Management
```
wild service list                 # List services
wild service install <service>    # Install service
wild service install-all          # Install all services
wild service delete <service>     # Delete service
wild service status <service>     # Get service status
```

### App Management
```
wild app list                     # List available apps
wild app list-deployed            # List deployed apps
wild app show <app>               # Show app details
wild app add <app>                # Add app to instance
wild app deploy <app>             # Deploy app
wild app delete <app>             # Delete app
wild app status <app>             # Get app status
```

### Backup & Restore
```
wild backup <app>                 # Backup app
wild backup list <app>            # List backups for app
wild restore <app>                # Restore app from latest backup
wild restore <app> --db-only      # Restore database only
wild restore <app> --pvc-only     # Restore PVCs only
```

### Utilities
```
wild health                       # Check cluster health
wild dashboard token              # Get dashboard token
wild node-ip                      # Get control plane IP
wild node-ips                     # Get all node IPs
wild version                      # Show versions (CLI + K8s + Talos)
```

### Operations
```
wild operation get <id>           # Get operation status
wild operation list               # List operations
wild operation cancel <id>        # Cancel operation
```

## Configuration

### Context File
```yaml
# ~/.wild/config.yaml
current-context: homelab

contexts:
  homelab:
    daemon-url: http://192.168.1.100:8080
    instance: production

  testlab:
    daemon-url: http://192.168.1.200:8080
    instance: test
```

### Environment Variables
```
WILD_DAEMON_URL    # Override daemon URL
WILD_CLI_DATA      # Override data directory (default: ~/.wildcloud)
```

## Command Mapping from v.PoC

| v.PoC Script | wild-cli Command |
|--------------|------------------|
| wild-init | wild instance create |
| wild-config | wild config get |
| wild-config-set | wild config set |
| wild-secret | wild secret get |
| wild-secret-set | wild secret set |
| wild-node-detect | wild node discover |
| wild-node-setup | wild node setup |
| wild-cluster-config-generate | wild cluster config generate |
| wild-setup-cluster | wild cluster bootstrap |
| wild-service-setup | wild service install |
| wild-setup-services | wild service install-all |
| wild-app-add | wild app add |
| wild-app-deploy | wild app deploy |
| wild-app-delete | wild app delete |
| wild-app-list | wild app list-deployed |
| wild-apps-list | wild app list |
| wild-app-backup | wild backup |
| wild-app-restore | wild restore |
| wild-health | wild health |
| wild-dashboard-token | wild dashboard token |
| wild-cluster-node-ip | wild node-ip |

## Implementation

### Technology Stack
- **Cobra**: CLI framework
- **Viper**: Configuration management
- **net/http**: HTTP client for API calls
- **encoding/json**: JSON parsing
- **Direct, simple code**: No unnecessary abstractions

### API Client
Simple HTTP client that:
- Reads context from config file
- Makes HTTP requests to daemon
- Handles JSON responses
- Pretty-prints output
- Returns appropriate exit codes

### Error Handling
- Clear error messages
- Appropriate exit codes (0 = success, 1 = error)
- API errors displayed to user
- Connection errors handled gracefully

### Output Formats
- Default: Human-readable tables
- `--json`: JSON output for scripting
- `--yaml`: YAML output (where appropriate)
- `--quiet`: Minimal output

### Progress Display
- Long operations: Show operation ID and status polling
- `--wait`: Block until operation completes
- `--no-wait`: Return immediately (default)

## Philosophy

Following the Wild Central philosophy:
- **Ruthless Simplicity**: Direct API calls, no complex state management
- **Thin Wrapper**: CLI is just a friendly interface to the API
- **Clear Errors**: Meaningful error messages
- **Unix Philosophy**: Do one thing well, composable commands
- **Consistency**: Similar commands have similar flags and output
