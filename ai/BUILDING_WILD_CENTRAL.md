# Building Wild Cloud Central

The first version of Wild Cloud, the Proof of Concept version (v.PoC), was built as a collection of shell scripts that users would run from their local machines. This works well for early adopters who are comfortable with the command line, Talos, and Kubernetes.

To make Wild Cloud more accessible to a broader audience, we are developing Wild Central. Central is a single-purpose machine run on a LAN that will deliver:

- Wild Daemon: A lightweight service that runs on a local machine (e.g., a Raspberry Pi) to manage Wild Cloud instances on the local network.
- Wild App: A web-based interface (to Wild Daemon) for managing Wild Cloud instances.
- Wild CLI: A command-line interface (to Wild Daemon) for advanced users who prefer to manage Wild Cloud from the terminal.

## Background info

### Info about Wild Cloud v.PoC

- @docs/agent-context/wildcloud/README.md
- @docs/agent-context/wildcloud/overview.md
- @docs/agent-context/wildcloud/project-architecture.md
- @docs/agent-context/wildcloud/bin-scripts.md
- @docs/agent-context/wildcloud/configuration-system.md
- @docs/agent-context/wildcloud/setup-process.md
- @docs/agent-context/wildcloud/apps-system.md

### Info about Talos

- @docs/agent-context/talos-v1.11/README.md
- @docs/agent-context/talos-v1.11/architecture-and-components.md
- @docs/agent-context/talos-v1.11/cli-essentials.md
- @docs/agent-context/talos-v1.11/cluster-operations.md
- @docs/agent-context/talos-v1.11/discovery-and-networking.md
- @docs/agent-context/talos-v1.11/etcd-management.md
- @docs/agent-context/talos-v1.11/bare-metal-administration.md
- @docs/agent-context/talos-v1.11/troubleshooting-guide.md

## Architecture

### Old v.PoC Architecture

- WC_ROOT: The scripts used to set up and manage the Wild Cloud cluster. Currently, this is a set of shell scripts in $WC_ROOT/bin.
- WC_HOME: During setup, the user creates a Wild Cloud project directory (WC_HOME) on their local machine. This directory holds all configuration, secrets, and k8s manifests for their specific Wild Cloud deployment.
- Wild Cloud Apps Directory: The Wild Cloud apps are stored in the `apps/` directory within the WC_ROOT repository. Users can deploy these apps to their cluster using the scripts in WC_ROOT/bin.
- dnsmasq server: Scripts help the operator set up a dnsmasq server on a separate machine to provide LAN DNS services during node bootstrapping.

### New Wild Central Architecture

#### wildd: The Wild Cloud Daemon

wildd is a long-running service that provides an API and web interface for managing one or more Wild Cloud clusters. It runs on a dedicated device within the user's network.

wildd replaces functionality from the v.PoC scripts and the dnsmasq server. It is one API for managing multiple wild cloud instances on the LAN.

Both wild-app and wild-cli communicate with wildd to perform actions.

See: @daemon/BUILDING_WILD_DAEMON.md

#### wild-app

The web application that provides the user interface for Wild Cloud on Wild Central. It communicates with wildd to perform actions and display information.

See: @/app/BUILDING_WILD_APP.md

#### wild-cli

A command-line interface for advanced users who prefer to manage Wild Cloud from the terminal. It communicates with wildd to perform actions.

Mirrors all of the wild-* scripts from v.PoC, but adapted for the new architecture:

- One golang client (wild-cli) replaces many bash scripts (wild-*).
- Wrapper around wildd API instead of direct file manipulation.
- Multi-cloud: v.PoC scripts set the instance context with WC_HOME environment variable. In Central, wild-cli follows the "context" pattern like kubectl and talosctl, using `--context` or `WILD_CONTEXT` to select which wild cloud instance to manage, or defaulting to the "current" context.

See: @cli/BUILDING_WILD_CLI.md

#### Wild Central Data

Configured with $WILD_CENTRAL_DATA environment variable (default: /var/lib/wild-central).

Replaces multiple WC_HOMEs. All wild clouds managed on the LAN are configured here. These are still in easy to read YAML format and can be edited directly or through the webapp.

Wild Central data also holds the local app directory, logs, and artifacts, and overall state data.

#### Wild Cloud Apps Directory

The Wild Cloud apps are stored in the `apps/` directory within the WC_ROOT repository. Users can deploy these apps to their cluster using the webapp or wild-cli.

#### dnsmasq server

The Wild Cloud Central Daemon (wildd) includes functionality to manage a dnsmasq server on the same device, providing LAN DNS services during node bootstrapping.

## Philosophy

- Use talosctl and kubectl wherever possible to leverage existing tools and avoid reinventing the wheel. 
- A wild cloud instance is primarily data (YAML files for config, secrets, and manifests).
- Because a wild cloud instance is primarily data, a wild cloud instance can be managed by non-technical users through the webapp or by technical users by SSHing into the device (e.g. VSCode Remote SSH).


## Dev Environment Requirements

- Go 1.21+
- Node.js 20+
- Docker (for building container images)
- GNU Make (for build automation)

## Packaging and Installation

See @experimental/ai/WILD_CENTRAL_PACKAGING.md