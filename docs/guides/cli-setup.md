# Wild Central CLI Setup Guide

This guide shows how to set up a complete Wild Cloud cluster using the `wild` CLI.

## Prerequisites

```bash
# Ensure wild daemon is running
wild daemon status

# Create and select your instance
wild instance create test-cloud
wild instance use test-cloud
```

## 1. Cluster Configuration

Configure your cluster's basic settings:

```bash
# Set operator email
wild config set operator.email "your-email@example.com"

# Set cluster name
wild config set cluster.name "wild-cluster"

# Configure network settings
wild config set cloud.router.ip "192.168.8.1"
wild config set cloud.dns.ip "192.168.1.50"
wild config set cloud.dhcpRange "192.168.1.34,192.168.1.79"
wild config set cloud.dnsmasq.interface "eth0"

# Configure MetalLB IP pool
wild config set cluster.ipAddressPool "192.168.1.80-192.168.1.89"
wild config set cluster.loadBalancerIp "192.168.1.80"

# Configure control plane VIP
wild config set cluster.nodes.control.vip "192.168.1.30"

# Set Talos version and schematic
wild config set cluster.nodes.talos.version "v1.11.2"
wild config set cluster.nodes.talos.schematicId "434a0300db532066f1098e05ac068159371d00f0aba0a3103a0e826e83825c82"
```

## 2. Generate Cluster Configuration

```bash
wild cluster config generate
```

## 3. Configure talosctl Context

```bash
wild cluster talosconfig --persist
source <(wild instance env)
```

Set up three control plane nodes for high availability:

```bash
# Control plane node 1
wild node detect 192.168.8.122
wild node add test-control-1 controlplane --currentIp 192.168.8.122 --ip 192.168.8.31 --disk /dev/sdb --interface enp4s0
wild node apply test-control-1
```

After the first control plane node is up, bootstrap the cluster!

```bash
wild cluster bootstrap test-control-1
```

Now repeat the detect, add, apply for control nodes 2 and 3.

### Worker Nodes (Optional)

Add worker nodes as needed:

```bash
wild node detect 192.168.8.100
wild node add test-worker-1 worker --targetIp 192.168.8.100 --disk /dev/sda --interface eth0
wild node apply test-worker-1
```

## 5. Services Setup

Install cluster services in dependency order:

```bash
wild service install metallb
wild service install longhorn
wild service install traefik
wild service install coredns
wild service install cert-manager
wild service install externaldns
wild service install kubernetes-dashboard
wild service install nfs
wild service install docker-registry
```

**Or install all services at once:**

```bash
wild services install --all
```

## 6. Verification

Verify your cluster is healthy:

```bash
# Check cluster health
wild health

# Check nodes
kubectl get nodes

# Get dashboard token
wild dashboard token

# Access dashboard at:
# https://dashboard.<your-internal-domain>
```

## Key Differences from v.PoC

The new Wild Central CLI differs from the v.PoC bash scripts in several important ways:

1. **Context-based**: Use `wild context` instead of `WC_HOME` environment variable
2. **Unified CLI**: Single `wild` command instead of many `wild-*` scripts
3. **Simplified node setup**: `wild node setup` handles detection, config generation, and deployment
4. **Service management**: `wild service install` instead of `wild-service-setup`
5. **Daemon-backed**: All operations go through `wildd` API instead of direct file manipulation

## Troubleshooting

If you encounter issues:

```bash
# Check daemon logs
wild daemon logs

# Check cluster health
wild health

# View node status
wild node list

# View service status
wild service list
```

For more help, see the [Troubleshooting Guide](../troubleshooting.md).
