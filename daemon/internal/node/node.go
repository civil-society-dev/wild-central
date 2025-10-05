package node

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wild-cloud/wild-central/daemon/internal/config"
	"github.com/wild-cloud/wild-central/daemon/internal/tools"
)

// Manager handles node configuration and state management
type Manager struct {
	dataDir    string
	configMgr  *config.Manager
	talosctl   *tools.Talosctl
}

// NewManager creates a new node manager
func NewManager(dataDir string) *Manager {
	return &Manager{
		dataDir:   dataDir,
		configMgr: config.NewManager(),
		talosctl:  tools.NewTalosctl(),
	}
}

// Node represents a cluster node configuration
type Node struct {
	MAC        string `yaml:"mac" json:"mac"`
	Hostname   string `yaml:"hostname" json:"hostname"`
	Role       string `yaml:"role" json:"role"` // controlplane or worker
	TargetIP   string `yaml:"targetIp" json:"target_ip"`
	Interface  string `yaml:"interface,omitempty" json:"interface,omitempty"`
	Disk       string `yaml:"disk" json:"disk"`
	Version    string `yaml:"version,omitempty" json:"version,omitempty"`
	SchematicID string `yaml:"schematicId,omitempty" json:"schematic_id,omitempty"`
	Configured bool   `yaml:"configured,omitempty" json:"configured"`
	Deployed   bool   `yaml:"deployed,omitempty" json:"deployed"`
}

// HardwareInfo contains discovered hardware information
type HardwareInfo struct {
	MAC             string   `json:"mac"`
	IP              string   `json:"ip"`
	Interface       string   `json:"interface"`
	Disks           []string `json:"disks"`
	SelectedDisk    string   `json:"selected_disk"`
	MaintenanceMode bool     `json:"maintenance_mode"`
}

// SetupOptions contains options for node setup
type SetupOptions struct {
	Reconfigure bool `json:"reconfigure"`
	NoDeploy    bool `json:"no_deploy"`
}

// GetInstancePath returns the path to an instance's nodes directory
func (m *Manager) GetInstancePath(instanceName string) string {
	return filepath.Join(m.dataDir, "instances", instanceName)
}

// List returns all nodes for an instance
func (m *Manager) List(instanceName string) ([]Node, error) {
	instancePath := m.GetInstancePath(instanceName)
	configPath := filepath.Join(instancePath, "config.yaml")

	yq := tools.NewYQ()

	// Get all node hostnames from cluster.nodes.active
	output, err := yq.Exec("eval", ".cluster.nodes.active | keys", configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nodes: %w", err)
	}

	// Parse hostnames (yq returns YAML array)
	hostnamesYAML := string(output)
	if hostnamesYAML == "" || hostnamesYAML == "null\n" {
		return []Node{}, nil
	}

	// Get hostnames line by line
	var hostnames []string
	for _, line := range strings.Split(hostnamesYAML, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && line != "null" && line != "-" {
			// Remove leading "- " from YAML array
			hostname := line
			if len(hostname) > 2 && hostname[0:2] == "- " {
				hostname = hostname[2:]
			}
			if hostname != "" {
				hostnames = append(hostnames, hostname)
			}
		}
	}

	// Get details for each node
	var nodes []Node
	for _, hostname := range hostnames {
		basePath := fmt.Sprintf(".cluster.nodes.active.%s", hostname)

		// Get node fields
		role, _ := yq.Exec("eval", basePath+".role", configPath)
		targetIP, _ := yq.Exec("eval", basePath+".targetIp", configPath)
		disk, _ := yq.Exec("eval", basePath+".disk", configPath)
		iface, _ := yq.Exec("eval", basePath+".interface", configPath)
		version, _ := yq.Exec("eval", basePath+".version", configPath)
		schematicID, _ := yq.Exec("eval", basePath+".schematicId", configPath)
		mac, _ := yq.Exec("eval", basePath+".mac", configPath)
		configured, _ := yq.Exec("eval", basePath+".configured", configPath)
		deployed, _ := yq.Exec("eval", basePath+".deployed", configPath)

		node := Node{
			Hostname:    hostname,
			Role:        tools.CleanYQOutput(string(role)),
			TargetIP:    tools.CleanYQOutput(string(targetIP)),
			Disk:        tools.CleanYQOutput(string(disk)),
			Interface:   tools.CleanYQOutput(string(iface)),
			Version:     tools.CleanYQOutput(string(version)),
			SchematicID: tools.CleanYQOutput(string(schematicID)),
			MAC:         tools.CleanYQOutput(string(mac)),
			Configured:  tools.CleanYQOutput(string(configured)) == "true",
			Deployed:    tools.CleanYQOutput(string(deployed)) == "true",
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// Get returns a specific node by hostname or MAC
func (m *Manager) Get(instanceName, nodeIdentifier string) (*Node, error) {
	// Get all nodes
	nodes, err := m.List(instanceName)
	if err != nil {
		return nil, err
	}

	// Find node by hostname or MAC
	for _, node := range nodes {
		if node.Hostname == nodeIdentifier || node.MAC == nodeIdentifier {
			return &node, nil
		}
	}

	return nil, fmt.Errorf("node %s not found", nodeIdentifier)
}

// Add registers a new node in config.yaml
func (m *Manager) Add(instanceName string, node *Node) error {
	instancePath := m.GetInstancePath(instanceName)

	// Validate node data
	if node.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}
	if node.MAC == "" {
		return fmt.Errorf("MAC address is required")
	}
	if node.Role != "controlplane" && node.Role != "worker" {
		return fmt.Errorf("role must be 'controlplane' or 'worker'")
	}
	if node.Disk == "" {
		return fmt.Errorf("disk is required")
	}

	// Check if node already exists (idempotency)
	existing, err := m.Get(instanceName, node.Hostname)
	if err == nil && existing != nil {
		// Node already exists, this is idempotent
		return nil
	}

	// Add node to config.yaml
	// Path: cluster.nodes.active.{hostname}
	basePath := fmt.Sprintf("cluster.nodes.active.%s", node.Hostname)

	configPath := filepath.Join(instancePath, "config.yaml")
	yq := tools.NewYQ()

	// Set each field
	if err := yq.Set(configPath, basePath+".mac", node.MAC); err != nil {
		return fmt.Errorf("failed to set MAC: %w", err)
	}
	if err := yq.Set(configPath, basePath+".role", node.Role); err != nil {
		return fmt.Errorf("failed to set role: %w", err)
	}
	if err := yq.Set(configPath, basePath+".disk", node.Disk); err != nil {
		return fmt.Errorf("failed to set disk: %w", err)
	}
	if node.TargetIP != "" {
		if err := yq.Set(configPath, basePath+".targetIp", node.TargetIP); err != nil {
			return fmt.Errorf("failed to set targetIP: %w", err)
		}
	}
	if node.Interface != "" {
		if err := yq.Set(configPath, basePath+".interface", node.Interface); err != nil {
			return fmt.Errorf("failed to set interface: %w", err)
		}
	}
	if node.Version != "" {
		if err := yq.Set(configPath, basePath+".version", node.Version); err != nil {
			return fmt.Errorf("failed to set version: %w", err)
		}
	}

	return nil
}

// Delete removes a node from config.yaml
func (m *Manager) Delete(instanceName, nodeIdentifier string) error {
	// Get node to find hostname
	node, err := m.Get(instanceName, nodeIdentifier)
	if err != nil {
		return err
	}

	instancePath := m.GetInstancePath(instanceName)
	configPath := filepath.Join(instancePath, "config.yaml")

	// Delete node from config.yaml
	// Path: cluster.nodes.active.{hostname}
	nodePath := fmt.Sprintf("cluster.nodes.active.%s", node.Hostname)

	yq := tools.NewYQ()
	// Use yq to delete the node
	_, err = yq.Exec("eval", "-i", fmt.Sprintf("del(%s)", nodePath), configPath)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	return nil
}

// DetectHardware queries node hardware information via talosctl
func (m *Manager) DetectHardware(nodeIP string) (*HardwareInfo, error) {
	// Query node with insecure flag (maintenance mode)
	insecure := true

	// Try to get default interface (with default route)
	iface, err := m.talosctl.GetDefaultInterface(nodeIP, insecure)
	if err != nil {
		// Fall back to physical interface
		iface, err = m.talosctl.GetPhysicalInterface(nodeIP, insecure)
		if err != nil {
			return nil, fmt.Errorf("failed to detect interface: %w", err)
		}
	}

	// Get disks
	disks, err := m.talosctl.GetDisks(nodeIP, insecure)
	if err != nil {
		return nil, fmt.Errorf("failed to detect disks: %w", err)
	}

	// Select first disk as default
	var selectedDisk string
	if len(disks) > 0 {
		selectedDisk = disks[0]
	}

	return &HardwareInfo{
		IP:              nodeIP,
		Interface:       iface,
		Disks:           disks,
		SelectedDisk:    selectedDisk,
		MaintenanceMode: true,
	}, nil
}

// Setup performs node setup (configure + optional deploy)
func (m *Manager) Setup(instanceName, nodeIdentifier string, opts SetupOptions) error {
	// Get node configuration
	node, err := m.Get(instanceName, nodeIdentifier)
	if err != nil {
		return err
	}

	// Check if already configured (unless force reconfigure)
	if node.Configured && !opts.Reconfigure {
		if node.Deployed || opts.NoDeploy {
			// Already done
			return nil
		}
	}

	instancePath := m.GetInstancePath(instanceName)
	talosDir := filepath.Join(instancePath, "talos")

	// Generate Talos machine config if not exists
	controlplaneConfig := filepath.Join(talosDir, "controlplane.yaml")
	workerConfig := filepath.Join(talosDir, "worker.yaml")

	// TODO: Check if configs exist, generate if needed
	// This would call talosctl gen config

	// Apply configuration to node
	if !opts.NoDeploy {
		var configFile string
		if node.Role == "controlplane" {
			configFile = controlplaneConfig
		} else {
			configFile = workerConfig
		}

		// Apply config with insecure flag (maintenance mode)
		if err := m.talosctl.ApplyConfig(node.TargetIP, configFile, true); err != nil {
			return fmt.Errorf("failed to apply config: %w", err)
		}

		// Mark as deployed
		node.Deployed = true
		if err := m.updateNodeStatus(instanceName, node); err != nil {
			return fmt.Errorf("failed to update node status: %w", err)
		}
	}

	// Mark as configured
	node.Configured = true
	if err := m.updateNodeStatus(instanceName, node); err != nil {
		return fmt.Errorf("failed to update node status: %w", err)
	}

	return nil
}

// updateNodeStatus updates node status flags in config.yaml
func (m *Manager) updateNodeStatus(instanceName string, node *Node) error {
	instancePath := m.GetInstancePath(instanceName)
	configPath := filepath.Join(instancePath, "config.yaml")
	basePath := fmt.Sprintf("cluster.nodes.active.%s", node.Hostname)

	yq := tools.NewYQ()

	if node.Configured {
		if err := yq.Set(configPath, basePath+".configured", "true"); err != nil {
			return err
		}
	}
	if node.Deployed {
		if err := yq.Set(configPath, basePath+".deployed", "true"); err != nil {
			return err
		}
	}

	return nil
}
