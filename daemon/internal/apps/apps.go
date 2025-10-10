package apps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/wild-cloud/wild-central/daemon/internal/storage"
	"github.com/wild-cloud/wild-central/daemon/internal/tools"
)

// Manager handles application lifecycle operations
type Manager struct {
	dataDir string
	appsDir string // Path to apps directory in repo
}

// NewManager creates a new apps manager
func NewManager(dataDir, appsDir string) *Manager {
	return &Manager{
		dataDir: dataDir,
		appsDir: appsDir,
	}
}

// App represents an application
type App struct {
	Name         string            `json:"name" yaml:"name"`
	Description  string            `json:"description" yaml:"description"`
	Version      string            `json:"version" yaml:"version"`
	Category     string            `json:"category" yaml:"category"`
	Dependencies []string          `json:"dependencies" yaml:"dependencies"`
	Config       map[string]string `json:"config,omitempty" yaml:"config,omitempty"`
}

// DeployedApp represents a deployed application instance
type DeployedApp struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Version   string `json:"version"`
	Namespace string `json:"namespace"`
	URL       string `json:"url,omitempty"`
}

// ListAvailable lists all available apps from the apps directory
func (m *Manager) ListAvailable() ([]App, error) {
	if m.appsDir == "" {
		return []App{}, fmt.Errorf("apps directory not configured")
	}

	// Read apps directory
	entries, err := os.ReadDir(m.appsDir)
	if err != nil {
		return []App{}, fmt.Errorf("failed to read apps directory: %w", err)
	}

	apps := []App{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check for manifest.yaml
		appFile := filepath.Join(m.appsDir, entry.Name(), "manifest.yaml")
		if !storage.FileExists(appFile) {
			continue
		}

		// Parse manifest.yaml
		data, err := os.ReadFile(appFile)
		if err != nil {
			continue
		}

		var app App
		if err := yaml.Unmarshal(data, &app); err != nil {
			continue
		}

		app.Name = entry.Name() // Use directory name as app name
		apps = append(apps, app)
	}

	return apps, nil
}

// Get returns details for a specific available app
func (m *Manager) Get(appName string) (*App, error) {
	appFile := filepath.Join(m.appsDir, appName, "manifest.yaml")

	if !storage.FileExists(appFile) {
		return nil, fmt.Errorf("app %s not found", appName)
	}

	data, err := os.ReadFile(appFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read app file: %w", err)
	}

	var app App
	if err := yaml.Unmarshal(data, &app); err != nil {
		return nil, fmt.Errorf("failed to parse app file: %w", err)
	}

	app.Name = appName
	return &app, nil
}

// ListDeployed lists deployed apps for an instance
func (m *Manager) ListDeployed(instanceName string) ([]DeployedApp, error) {
	// This would query kubectl for deployed apps
	// For now, return placeholder
	apps := []DeployedApp{}

	// TODO: Query kubectl get namespaces with app labels
	// TODO: Query deployment status for each app

	return apps, nil
}

// Add adds an app to the instance configuration
func (m *Manager) Add(instanceName, appName string, config map[string]string) error {
	// Verify app exists
	app, err := m.Get(appName)
	if err != nil {
		return err
	}

	// Add app to config.yaml
	// Path: apps.deployed.{appName}
	configFile := filepath.Join(m.dataDir, "instances", instanceName, "config.yaml")

	// Store app configuration
	basePath := fmt.Sprintf("apps.deployed.%s", appName)

	// Use yq to set configuration
	// This is simplified - real implementation would use yq wrapper
	_ = app // Use app for validation

	// TODO: Set each config value via yq
	_ = configFile
	_ = basePath
	_ = config

	return nil
}

// Deploy deploys an app to the cluster
func (m *Manager) Deploy(instanceName, appName string) error {
	kubeconfigPath := tools.GetKubeconfigPath(m.dataDir, instanceName)

	// Get app manifests
	manifestsDir := filepath.Join(m.appsDir, appName, "manifests")
	if !storage.FileExists(manifestsDir) {
		return fmt.Errorf("app %s has no manifests", appName)
	}

	// Apply manifests with kubectl
	cmd := exec.Command("kubectl", "apply", "-f", manifestsDir)
	tools.WithKubeconfig(cmd, kubeconfigPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to deploy app: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Delete removes an app from the cluster
func (m *Manager) Delete(instanceName, appName string) error {
	kubeconfigPath := tools.GetKubeconfigPath(m.dataDir, instanceName)

	// Delete manifests with kubectl
	manifestsDir := filepath.Join(m.appsDir, appName, "manifests")
	if !storage.FileExists(manifestsDir) {
		return fmt.Errorf("app %s has no manifests", appName)
	}

	cmd := exec.Command("kubectl", "delete", "-f", manifestsDir)
	tools.WithKubeconfig(cmd, kubeconfigPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete app: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetStatus returns the status of a deployed app
func (m *Manager) GetStatus(instanceName, appName string) (*DeployedApp, error) {
	// Query kubectl for app status
	// This is simplified - real implementation would parse kubectl output

	app := &DeployedApp{
		Name:      appName,
		Status:    "unknown",
		Namespace: appName, // Assume namespace matches app name
	}

	// TODO: Query kubectl get deployment -n {namespace}
	// TODO: Parse status and pod counts

	return app, nil
}
