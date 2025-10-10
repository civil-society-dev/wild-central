# NODE_SETUP.md - Node Management System Design

## Overview

This document specifies the design for wild-central's node management system. The design follows Zen-like minimalism principles, prioritizing simplicity, clarity, and user experience.

## System Architecture

### Node Structure

```go
type Node struct {
    // Core Identity
    Hostname   string `yaml:"hostname" json:"hostname"`
    Role       string `yaml:"role" json:"role"` // controlplane or worker

    // Network Configuration
    TargetIP   string `yaml:"targetIp" json:"target_ip"`      // Final production IP
    CurrentIP  string `yaml:"currentIp,omitempty" json:"current_ip,omitempty"` // Temporary maintenance IP
    Interface  string `yaml:"interface,omitempty" json:"interface,omitempty"`

    // Hardware Configuration
    Disk       string `yaml:"disk" json:"disk"`

    // Talos Configuration
    Version    string `yaml:"version,omitempty" json:"version,omitempty"`
    SchematicID string `yaml:"schematicId,omitempty" json:"schematic_id,omitempty"`

    // State Tracking
    Maintenance bool   `yaml:"maintenance,omitempty" json:"maintenance"` // NEW: explicit maintenance flag
    Configured  bool   `yaml:"configured,omitempty" json:"configured"`
    Applied    bool   `yaml:"applied,omitempty" json:"applied"`
}
```

### State Model

Nodes exist in one of three states:

1. **Initial/Maintenance** - Node booted from PXE, accessible via CurrentIP with --insecure
2. **Configured** - Configuration generated but not yet applied
3. **Applied** - Configuration applied, node running on TargetIP

## Command Specifications

### 1. Fetch Patch Templates

```bash
wild node fetch-patch-templates
```

**Purpose**: Pull latest node patch templates from directory/ to instance.

**Operation**:
1. Copies `directory/setup/cluster-nodes/patch.templates/*.yaml` to instance
2. Overwrites existing templates in instance directory
3. Used automatically by apply if templates missing

**Error Handling**:
- ERROR if source templates not found
- ERROR if unable to write to destination

### 2. Add Node

```bash
wild node add <hostname> <role> [flags]

Flags:
  --target-ip string      Target IP address for production
  --current-ip string     Current IP if in maintenance mode
  --disk string          Disk device (required, e.g. /dev/sda)
  --interface string     Network interface (optional)
  --schematic-id string  Talos schematic ID (optional, uses instance default)
  --maintenance          Mark node as in maintenance mode
```

**Purpose**: Register a new node in cluster configuration.

**Validation**:
- ERROR if node with hostname already exists
- ERROR if role not "controlplane" or "worker"
- ERROR if disk not provided
- WARN if no IP addresses provided

**Operation**:
1. Check if node configuration exists - ERROR if yes
2. Validate required fields (hostname, role, disk)
3. If schematic-id not provided, use instance default
4. Set maintenance=true if --maintenance flag OR currentIP provided
5. Write to `config.yaml` under `cluster.nodes.active.<hostname>`

**Example**:
```bash
# Node in maintenance mode
wild node add control-1 controlplane --current-ip 192.168.1.100 --target-ip 192.168.1.31 --disk /dev/sda

# Node already applied (unusual to run this, only if the config was removed manually, otherwise this errors if the config already exists)
wild node add worker-1 worker --target-ip 192.168.1.32 --disk /dev/nvme0n1
```

### 3. Update Node

```bash
wild node update <hostname> [flags]

Flags:
  --target-ip string      Update target IP address
  --current-ip string     Update current IP address
  --disk string          Update disk device
  --interface string     Update network interface
  --schematic-id string  Update Talos schematic ID
  --maintenance          Set maintenance mode
  --no-maintenance       Clear maintenance mode
```

**Purpose**: Modify existing node configuration.

**Operation**:
1. Check if node exists - ERROR if not
2. Update only specified fields (partial update)
3. Handle maintenance flag logic:
   - If --maintenance: set maintenance=true
   - If --no-maintenance: set maintenance=false
   - If currentIP changed: auto-set maintenance=true
4. Update config.yaml

**Example**:
```bash
# Update disk after hardware change
wild node update worker-1 --disk /dev/sdb

# Move node to maintenance
wild node update control-1 --current-ip 192.168.1.100 --maintenance

# Clear maintenance after successful apply
wild node update control-1 --no-maintenance
```

### 4. Apply Node

```bash
wild node apply <hostname>
```

**Purpose**: Generate configuration and apply to node.

**Operation**:

1. **Pre-checks**:
   - Verify node exists in config
   - Auto-fetch templates if missing (silent)
   - Verify base configs exist (controlplane.yaml/worker.yaml)

2. **Generate Configuration**:
   - Read template based on role
   - Substitute variables: {{NODE_NAME}}, {{NODE_IP}}, {{SCHEMATIC_ID}}, {{VERSION}}
   - Create patch file in `patch/<hostname>.yaml`
   - Merge base + patch using `talosctl machineconfig patch`
   - Output to `final/<hostname>.yaml`

3. **Apply Configuration**:
   - Determine target IP and mode:
     - If maintenance=true OR currentIP set: use currentIP with --insecure
     - Else: use targetIP without --insecure
   - Apply config: `talosctl apply-config [--insecure] -n <ip> -f final/<hostname>.yaml`

4. **Post-application Updates**:
   - On success:
     - Set applied=true
     - Set configured=true
     - Set currentIP = targetIP (node now on production IP)
     - Set maintenance=false (node out of maintenance)
   - On failure:
     - Keep existing state
     - Return error with details

**Example**:
```bash
# Apply to node in maintenance
wild node apply control-1
# Uses currentIP with --insecure, then updates state

# Re-apply to production node
wild node apply worker-1
# Uses targetIP without --insecure
```

## Maintenance Mode Logic

### Determining Maintenance Mode

Node is in maintenance when:
- `maintenance` flag is true OR
- `currentIP` is set and differs from `targetIP`

### Using --insecure Flag

Use `--insecure` with talosctl when:
- Node is in maintenance mode (as determined above)
- During initial application from PXE boot
- When explicitly accessing via currentIP

### Automatic State Transitions

**On successful apply**:
- currentIP → targetIP (node moves to production IP)
- maintenance → false (exits maintenance mode)
- applied → true

**On node update with currentIP**:
- maintenance → true (enters maintenance mode)

## Error Handling Strategy

### Hard Errors (Stop Execution)

- Node already exists on `add`
- Node not found on `update` or `apply`
- Invalid role specification
- Missing required fields (hostname, role, disk)
- Base configuration files missing
- talosctl command failures

### Warnings (Continue with Notice)

- No IP addresses provided on `add`
- Templates missing (auto-fetch and continue)
- Node already applied (on re-apply)

### Silent Operations

- Template auto-fetch on apply
- State updates after successful operations
- Config.yaml field updates

## Implementation Flow Examples

### Complete Node Addition Flow

```bash
# 1. Detect hardware on PXE-booted node
wild node detect 192.168.1.100
# Output: Interface: eth0, Disk: /dev/sda

# 2. Add node in maintenance mode
wild node add control-1 controlplane \
  --current-ip 192.168.1.100 \
  --target-ip 192.168.1.31 \
  --disk /dev/sda \
  --interface eth0

# 3. Apply configuration
wild node apply control-1
# - Generates config from templates
# - Applies to 192.168.1.100 with --insecure
# - On success: currentIP→targetIP, maintenance→false

# 4. Node now running on production IP
```

### Update Applied Node Flow

```bash
# 1. Move node to maintenance for hardware change
wild node update worker-1 --current-ip 192.168.1.100 --maintenance

# 2. Update hardware config
wild node update worker-1 --disk /dev/nvme0n1

# 3. Re-apply with new configuration
wild node apply worker-1
# Uses 192.168.1.100 with --insecure
# On success: clears maintenance, sets currentIP=targetIP
```

## Design Rationale

### Why Explicit Maintenance Flag?

- **Clarity**: Makes intent explicit rather than inferring from IP state
- **Control**: Allows manual maintenance mode management
- **Safety**: Prevents accidental production changes

### Why Auto-fetch Templates?

- **Simplicity**: Reduces manual steps
- **Reliability**: Ensures templates always available
- **Convenience**: Just-in-time fetching when needed

### Why Standard Flag Syntax?

- **Consistency**: Matches kubectl/talosctl patterns
- **Usability**: Familiar to operators
- **Clarity**: Clear separation of flag and value

### Why Always Regenerate on Apply?

- **Correctness**: Ensures config matches current state
- **Simplicity**: No complex caching logic
- **Predictability**: Same command = same operation

## Current State & Changes Needed

### Current Implementation

**Existing Commands:**
- `wild node detect <ip>` - Hardware detection (keep as-is)
- `wild node add <hostname> <role> --ip --current-ip --disk --interface --schematic-id`
- `wild node setup <hostname> --reconfigure --no-apply --update-templates`

**Current Node Struct:**
```go
type Node struct {
    Hostname    string
    Role        string
    TargetIP    string
    CurrentIP   string  // Used to infer maintenance mode
    Interface   string
    Disk        string
    Version     string
    SchematicID string
    Configured  bool
    Applied    bool
    // Missing: explicit maintenance flag
}
```

**Current Behavior:**
- Maintenance mode inferred from currentIP presence
- `node setup` does both config generation and application
- `--no-apply` flag skips application
- `--reconfigure` forces regeneration
- `--update-templates` refreshes templates

### Required Changes

#### 1. Node Struct
- **ADD**: `Maintenance bool` field
- **KEEP**: All existing fields

#### 2. CLI Commands

**Rename:**
- `wild node setup` → `wild node apply`

**Add New:**
- `wild node update <hostname> [flags]` - Update existing node config
- `wild node fetch-patch-templates` - Manual template refresh

**Modify `wild node add`:**
- Change `--ip` to `--target-ip` and `--current-ip` (breaking change)
- Add `--maintenance` flag
- Error (not skip) if node exists
- Auto-set maintenance=true if currentIP provided

**Modify `wild node apply` (was setup):**
- Remove `--no-apply`, `--reconfigure`, `--update-templates` flags
- Always regenerate configs
- Auto-fetch templates if missing
- Post-application: set currentIP=targetIP, maintenance=false

#### 3. Backend Logic Changes

**In `Add()` method:**
```go
// Before adding, check if exists
existing, _ := m.Get(instanceName, node.Hostname)
if existing != nil {
    return fmt.Errorf("node %s already exists", node.Hostname)
}

// Set maintenance flag
if node.CurrentIP != "" || opts.Maintenance {
    node.Maintenance = true
}
```

**In `Apply()` method (was Setup):**
```go
// Always fetch templates if missing
if !templatesExist(templatesDir) {
    copyTemplatesFromDirectory(templatesDir)
}

// Always regenerate configs (no skip logic)

// Determine application mode
var deployIP string
var useInsecure bool
if node.Maintenance || (node.CurrentIP != "" && node.CurrentIP != node.TargetIP) {
    deployIP = node.CurrentIP
    useInsecure = true
} else {
    deployIP = node.TargetIP
    useInsecure = false
}

// After successful apply
if err := applyConfig(deployIP, config, useInsecure); err == nil {
    node.CurrentIP = node.TargetIP  // Move to production IP
    node.Maintenance = false         // Exit maintenance mode
    node.Applied = true
    node.Configured = true
    updateNodeStatus(instanceName, node)
}
```

**New `Update()` method:**
```go
func (m *Manager) Update(instanceName string, hostname string, updates map[string]string) error {
    node, err := m.Get(instanceName, hostname)
    if err != nil {
        return fmt.Errorf("node %s not found", hostname)
    }

    // Apply partial updates
    for key, value := range updates {
        switch key {
        case "target_ip":
            node.TargetIP = value
        case "current_ip":
            node.CurrentIP = value
            node.Maintenance = true  // Auto-set on currentIP change
        case "disk":
            node.Disk = value
        // ... etc
        }
    }

    // Write back to config
    return m.saveNode(instanceName, node)
}
```

**New `FetchTemplates()` method:**
```go
func (m *Manager) FetchTemplates(instanceName string) error {
    destDir := filepath.Join(m.GetInstancePath(instanceName), "setup", "cluster-nodes", "patch.templates")
    return m.copyTemplatesFromDirectory(destDir)
}
```

#### 4. API Endpoints

**Current:**
- POST `/api/v1/instances/{name}/nodes` - Add node
- POST `/api/v1/instances/{name}/nodes/{node}/setup` - Setup node
- POST `/api/v1/instances/{name}/nodes/detect` - Detect node

**New:**
- POST `/api/v1/instances/{name}/nodes` - Add node (modify)
- PUT `/api/v1/instances/{name}/nodes/{node}` - Update node (NEW)
- POST `/api/v1/instances/{name}/nodes/{node}/apply` - Apply node (rename from setup)
- POST `/api/v1/instances/{name}/nodes/fetch-templates` - Fetch templates (NEW)
- POST `/api/v1/instances/{name}/nodes/detect` - Detect node (keep)

#### 5. Config.yaml Schema

**Add to node entries:**
```yaml
cluster:
  nodes:
    active:
      control-1:
        hostname: control-1
        role: controlplane
        targetIp: 192.168.1.31
        currentIp: 192.168.1.100  # Optional
        disk: /dev/sda
        interface: eth0
        schematicId: "..."
        maintenance: true  # NEW FIELD
        configured: true
        applied: false
```

### Migration Strategy

For existing applications:

1. **Add maintenance field to all nodes:**
   ```go
   for each node in config:
       if node.CurrentIP != "" && node.CurrentIP != node.TargetIP:
           node.Maintenance = true
       else:
           node.Maintenance = false
   ```

2. **Update CLI binary** - Breaking changes require version bump

3. **Update documentation** - Reflect new command structure

### Testing Checklist

- [ ] Add node - errors if exists
- [ ] Add node - sets maintenance=true if currentIP provided
- [ ] Add node - sets maintenance=true if --maintenance flag
- [ ] Update node - partial field updates work
- [ ] Update node - sets maintenance=true on currentIP change
- [ ] Apply - auto-fetches templates if missing
- [ ] Apply - uses currentIP with --insecure in maintenance mode
- [ ] Apply - uses targetIP without --insecure when not in maintenance
- [ ] Apply - sets currentIP=targetIP on success
- [ ] Apply - sets maintenance=false on success
- [ ] Fetch templates - copies from directory/ correctly

## Summary

This design provides a clear, explicit node management system that eliminates ambiguity and follows standard CLI patterns. The explicit maintenance flag, automatic state transitions, and simplified command structure create a system that's both powerful and intuitive. All changes maintain the project's philosophy of ruthless simplicity while improving usability and reliability.
