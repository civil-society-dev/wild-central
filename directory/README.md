# Wild Cloud Directory

The Wild Cloud Directory is a catalog of applications, services, and resources that can be installed on Wild Cloud instances. This directory is managed by the Wild Cloud organization and serves as a template repository.

## Purpose

The directory contains **static template data** that defines what can be installed, not runtime data about what has been installed. It is the source of truth for:

- Available applications and their configurations
- Cluster services and their manifests
- Default configurations and templates

## Directory Structure

```
directory/
├── apps/                    # Application catalog
│   ├── gitea/              # Each app has its own directory
│   │   ├── manifest.yaml   # App metadata and default config
│   │   └── deploy/         # Kubernetes manifests for deployment
│   ├── mysql/
│   ├── nextcloud/
│   └── ...
└── setup/                   # Cluster setup resources
    └── cluster-services/    # Core cluster services catalog
        ├── metallb/
        │   ├── wild-manifest.yaml  # Service metadata and config
        │   └── manifests.yaml      # Kubernetes manifests
        ├── traefik/
        ├── cert-manager/
        └── ...
```

## Configuration

The directory location **must be set** via the `WILD_DIRECTORY` environment variable:

```bash
export WILD_DIRECTORY=/path/to/wild-central/directory
```

**Required:** The daemon will not start without this environment variable set.

**Development usage:**

```bash
export WILD_DIRECTORY=~/repos/wild-central/directory
cd ~/repos/wild-central/daemon
./daemon
```

**Production usage:**

```bash
export WILD_DIRECTORY=/opt/wild-central/directory
export WILD_CENTRAL_DATA=/var/lib/wild-central
/usr/local/bin/wild-daemon
```

## Separation from Runtime Data

The Wild Cloud Directory is separate from Wild Central Data (`$WILD_CENTRAL_DATA`):

| Wild Cloud Directory | Wild Central Data |
|---------------------|-------------------|
| Template catalog (read-only) | Instance state (read-write) |
| Managed by Wild Cloud org | Managed by user |
| Apps and services definitions | Deployed instances |
| Shared across all instances | Instance-specific configuration |
| Version controlled | Not version controlled |

## Usage

When a user installs an application or service:

1. **Source:** The daemon reads the template from the Directory
2. **Copy:** Files are copied to the instance directory in Wild Central Data
3. **Configure:** User-specific configuration is applied
4. **Deploy:** The configured resources are deployed to the cluster

## Apps Directory

Contains application manifests that define:
- App metadata (name, description, version)
- Default configuration values
- Required secrets
- Dependencies on other apps
- Kubernetes deployment manifests

## Setup Directory

Contains cluster service manifests that define:
- Service metadata (name, namespace, description)
- Configuration references (values that must already be set)
- Service configuration (values to prompt user for)
- Dependencies on other services
- Kubernetes deployment resources

## Development vs Production

**Development:**
- Set to repository path: `WILD_DIRECTORY=~/repos/wild-central/directory`
- Edit templates directly in the repository
- Changes take effect immediately on daemon restart

**Production:**
- Use absolute path: `WILD_DIRECTORY=/opt/wild-central/directory`
- Can be mounted as read-only volume
- Update by replacing entire directory and restarting daemon

## Maintenance

The Wild Cloud Directory is:
- Version controlled with the wild-central repository
- Updated by the Wild Cloud organization
- Safe to update without affecting running instances
- Can be shared across multiple Wild Central installations

## See Also

- [Building Wild Central](../ai/BUILDING_WILD_CENTRAL.md) - Architecture overview
- [Wild Daemon](../daemon/README.md) - Daemon implementation
- [Wild CLI](../cli/README.md) - CLI implementation
