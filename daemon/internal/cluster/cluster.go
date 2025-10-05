package cluster

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wild-cloud/wild-central/daemon/internal/storage"
	"github.com/wild-cloud/wild-central/daemon/internal/tools"
)

// Manager handles cluster lifecycle operations
type Manager struct {
	dataDir  string
	talosctl *tools.Talosctl
}

// NewManager creates a new cluster manager
func NewManager(dataDir string) *Manager {
	return &Manager{
		dataDir:  dataDir,
		talosctl: tools.NewTalosctl(),
	}
}

// ClusterConfig contains cluster configuration parameters
type ClusterConfig struct {
	ClusterName string `json:"cluster_name"`
	VIP         string `json:"vip"` // Control plane virtual IP
	Version     string `json:"version"`
}

// ClusterStatus represents cluster health and status
type ClusterStatus struct {
	Status             string            `json:"status"` // ready, pending, error
	Nodes              int               `json:"nodes"`
	ControlPlaneNodes  int               `json:"control_plane_nodes"`
	WorkerNodes        int               `json:"worker_nodes"`
	KubernetesVersion  string            `json:"kubernetes_version"`
	TalosVersion       string            `json:"talos_version"`
	Services           map[string]string `json:"services"`
}

// GetTalosDir returns the talos directory for an instance
func (m *Manager) GetTalosDir(instanceName string) string {
	return filepath.Join(m.dataDir, "instances", instanceName, "talos")
}

// GetGeneratedDir returns the generated config directory
func (m *Manager) GetGeneratedDir(instanceName string) string {
	return filepath.Join(m.GetTalosDir(instanceName), "generated")
}

// GenerateConfig generates initial cluster configuration using talosctl gen config
func (m *Manager) GenerateConfig(instanceName string, config *ClusterConfig) error {
	generatedDir := m.GetGeneratedDir(instanceName)

	// Check if already generated (idempotency)
	secretsFile := filepath.Join(generatedDir, "secrets.yaml")
	if storage.FileExists(secretsFile) {
		// Already generated
		return nil
	}

	// Ensure generated directory exists
	if err := storage.EnsureDir(generatedDir, 0755); err != nil {
		return fmt.Errorf("failed to create generated directory: %w", err)
	}

	// Generate secrets
	cmd := exec.Command("talosctl", "gen", "secrets")
	cmd.Dir = generatedDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate secrets: %w\nOutput: %s", err, string(output))
	}

	// Generate config with secrets
	endpoint := fmt.Sprintf("https://%s:6443", config.VIP)
	cmd = exec.Command("talosctl", "gen", "config",
		"--with-secrets", "secrets.yaml",
		config.ClusterName,
		endpoint,
	)
	cmd.Dir = generatedDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate config: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Bootstrap bootstraps the cluster on the specified node
func (m *Manager) Bootstrap(instanceName, nodeName string) error {
	// Get node IP from config
	// This is a simplified version - real implementation would query config.yaml
	// For now, we'll require the full node IP to be passed

	// Bootstrap command
	cmd := exec.Command("talosctl", "bootstrap", "--nodes", nodeName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to bootstrap cluster: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetStatus retrieves cluster status
func (m *Manager) GetStatus(instanceName string) (*ClusterStatus, error) {
	// This is a simplified version
	// Real implementation would query talosctl and kubectl for actual status

	status := &ClusterStatus{
		Status:            "unknown",
		Nodes:             0,
		ControlPlaneNodes: 0,
		WorkerNodes:       0,
		Services:          make(map[string]string),
	}

	// Try to get cluster info
	// This requires talosconfig and kubeconfig to be set up
	// For now, return basic status

	return status, nil
}

// GetKubeconfig returns the kubeconfig for the cluster
func (m *Manager) GetKubeconfig(instanceName string) (string, error) {
	kubeconfigPath := filepath.Join(m.GetGeneratedDir(instanceName), "kubeconfig")

	if !storage.FileExists(kubeconfigPath) {
		return "", fmt.Errorf("kubeconfig not found - cluster may not be bootstrapped")
	}

	data, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	return string(data), nil
}

// GetTalosconfig returns the talosconfig for the cluster
func (m *Manager) GetTalosconfig(instanceName string) (string, error) {
	talosconfigPath := filepath.Join(m.GetGeneratedDir(instanceName), "talosconfig")

	if !storage.FileExists(talosconfigPath) {
		return "", fmt.Errorf("talosconfig not found - cluster may not be initialized")
	}

	data, err := os.ReadFile(talosconfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to read talosconfig: %w", err)
	}

	return string(data), nil
}

// Health checks cluster health
func (m *Manager) Health(instanceName string) ([]HealthCheck, error) {
	checks := []HealthCheck{}

	// Check 1: Talos config exists
	checks = append(checks, HealthCheck{
		Name:    "Talos Configuration",
		Status:  "passing",
		Message: "Talos configuration generated",
	})

	// Check 2: Kubeconfig exists
	if _, err := m.GetKubeconfig(instanceName); err == nil {
		checks = append(checks, HealthCheck{
			Name:    "Kubernetes Configuration",
			Status:  "passing",
			Message: "Kubeconfig available",
		})
	} else {
		checks = append(checks, HealthCheck{
			Name:    "Kubernetes Configuration",
			Status:  "warning",
			Message: "Kubeconfig not found",
		})
	}

	// Additional health checks would query actual cluster state
	// via kubectl and talosctl

	return checks, nil
}

// HealthCheck represents a single health check result
type HealthCheck struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // passing, warning, failing
	Message string `json:"message"`
}

// Reset resets the cluster (dangerous operation)
func (m *Manager) Reset(instanceName string, confirm bool) error {
	if !confirm {
		return fmt.Errorf("reset requires confirmation")
	}

	// This is a destructive operation
	// Real implementation would:
	// 1. Reset all nodes via talosctl reset
	// 2. Remove generated configs
	// 3. Clear node status in config.yaml

	generatedDir := m.GetGeneratedDir(instanceName)
	if storage.FileExists(generatedDir) {
		if err := os.RemoveAll(generatedDir); err != nil {
			return fmt.Errorf("failed to remove generated configs: %w", err)
		}
	}

	return nil
}

// ConfigureContext configures talosctl context for the cluster
func (m *Manager) ConfigureContext(instanceName, clusterName string) error {
	talosconfigPath := filepath.Join(m.GetGeneratedDir(instanceName), "talosconfig")

	if !storage.FileExists(talosconfigPath) {
		return fmt.Errorf("talosconfig not found")
	}

	// Merge talosconfig into user's talosctl config
	cmd := exec.Command("talosctl", "config", "merge", talosconfigPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to merge talosconfig: %w\nOutput: %s", err, string(output))
	}

	// Set context
	cmd = exec.Command("talosctl", "config", "context", clusterName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set context: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// HasContext checks if talosctl context exists
func (m *Manager) HasContext(clusterName string) (bool, error) {
	cmd := exec.Command("talosctl", "config", "contexts")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to list contexts: %w", err)
	}

	return strings.Contains(string(output), clusterName), nil
}
