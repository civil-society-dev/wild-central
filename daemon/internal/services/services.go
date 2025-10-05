package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wild-cloud/wild-central/daemon/internal/storage"
	"github.com/wild-cloud/wild-central/daemon/internal/tools"
)

// Manager handles base service operations
type Manager struct {
	dataDir     string
	servicesDir string                      // Path to services directory
	manifests   map[string]*ServiceManifest // Cached service manifests
}

// NewManager creates a new services manager
func NewManager(dataDir, servicesDir string) *Manager {
	m := &Manager{
		dataDir:     dataDir,
		servicesDir: servicesDir,
	}

	// Load all service manifests
	manifests, err := LoadAllManifests(servicesDir)
	if err != nil {
		// Log error but continue - services without manifests will fall back to hardcoded map
		fmt.Printf("Warning: failed to load service manifests: %v\n", err)
		manifests = make(map[string]*ServiceManifest)
	}
	m.manifests = manifests

	return m
}

// Service represents a base service
type Service struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Version     string   `json:"version"`
	Namespace   string   `json:"namespace"`
	Dependencies []string `json:"dependencies,omitempty"`
}

// Base services in Wild Cloud (kept for reference/validation)
var BaseServices = []string{
	"metallb",      // Load balancer
	"traefik",      // Ingress controller
	"cert-manager", // Certificate management
	"longhorn",     // Storage
}

// serviceDeployments maps service directory names to their actual namespace and deployment name
var serviceDeployments = map[string]struct {
	namespace      string
	deploymentName string
}{
	"cert-manager":           {"cert-manager", "cert-manager"},
	"coredns":                {"kube-system", "coredns"},
	"docker-registry":        {"docker-registry", "docker-registry"},
	"externaldns":            {"externaldns", "external-dns"},
	"kubernetes-dashboard":   {"kubernetes-dashboard", "kubernetes-dashboard"},
	"longhorn":               {"longhorn-system", "longhorn-ui"},
	"metallb":                {"metallb-system", "controller"},
	"nfs":                    {"nfs-system", "nfs-server"},
	"node-feature-discovery": {"node-feature-discovery", "node-feature-discovery"},
	"nvidia-device-plugin":   {"nvidia-device-plugin", "nvidia-device-plugin-daemonset"},
	"smtp":                   {"smtp-system", "smtp"},
	"traefik":                {"traefik", "traefik"},
	"utils":                  {"utils-system", "utils"},
}

// getKubeconfigPath returns the path to the kubeconfig for an instance
func (m *Manager) getKubeconfigPath(instanceName string) string {
	return filepath.Join(m.dataDir, "instances", instanceName, ".kubeconfig")
}

// checkServiceStatus checks if a service is deployed
func (m *Manager) checkServiceStatus(instanceName, serviceName string) string {
	kubeconfigPath := m.getKubeconfigPath(instanceName)

	// If kubeconfig doesn't exist, cluster isn't bootstrapped
	if !storage.FileExists(kubeconfigPath) {
		return "not-deployed"
	}

	var namespace, deploymentName string

	// Try to get from manifest first
	if manifest, ok := m.manifests[serviceName]; ok {
		namespace = manifest.Namespace
		deploymentName = manifest.GetDeploymentName()
	} else {
		// Fall back to hardcoded map for services not yet migrated
		deployment, ok := serviceDeployments[serviceName]
		if !ok {
			// Service not found anywhere, assume not deployed
			return "not-deployed"
		}
		namespace = deployment.namespace
		deploymentName = deployment.deploymentName
	}

	kubectl := tools.NewKubectl(kubeconfigPath)
	if kubectl.DeploymentExists(deploymentName, namespace) {
		return "deployed"
	}

	return "not-deployed"
}

// List returns all base services and their status
func (m *Manager) List(instanceName string) ([]Service, error) {
	services := []Service{}

	// Discover services from the services directory
	entries, err := os.ReadDir(m.servicesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read services directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip non-directories like README.md
		}

		name := entry.Name()

		// Get service info from manifest if available
		var namespace, description, version string
		var dependencies []string

		if manifest, ok := m.manifests[name]; ok {
			namespace = manifest.Namespace
			description = manifest.Description
			version = manifest.Category // Using category as version for now
			dependencies = manifest.Dependencies
		} else {
			// Fall back to hardcoded map
			namespace = name + "-system" // default
			if deployment, ok := serviceDeployments[name]; ok {
				namespace = deployment.namespace
			}
		}

		service := Service{
			Name:         name,
			Status:       m.checkServiceStatus(instanceName, name),
			Namespace:    namespace,
			Description:  description,
			Version:      version,
			Dependencies: dependencies,
		}

		services = append(services, service)
	}

	return services, nil
}

// Get returns a specific service
func (m *Manager) Get(instanceName, serviceName string) (*Service, error) {
	// Get the correct namespace from the map
	namespace := serviceName + "-system" // default
	if deployment, ok := serviceDeployments[serviceName]; ok {
		namespace = deployment.namespace
	}

	service := &Service{
		Name:      serviceName,
		Status:    m.checkServiceStatus(instanceName, serviceName),
		Namespace: namespace,
	}

	return service, nil
}

// Install installs a base service
func (m *Manager) Install(instanceName, serviceName string, fetch, deploy bool) error {
	// Get service manifests
	serviceDir := filepath.Join(m.servicesDir, serviceName)
	if !storage.FileExists(serviceDir) {
		return fmt.Errorf("service %s not found", serviceName)
	}

	manifestsFile := filepath.Join(serviceDir, "manifests.yaml")
	if !storage.FileExists(manifestsFile) {
		return fmt.Errorf("service %s has no manifests", serviceName)
	}

	if !deploy {
		// Just configuration, don't deploy
		return nil
	}

	// Apply manifests
	cmd := exec.Command("kubectl", "apply", "-f", manifestsFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install service: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// InstallAll installs all base services
func (m *Manager) InstallAll(instanceName string, fetch, deploy bool) error {
	for _, serviceName := range BaseServices {
		if err := m.Install(instanceName, serviceName, fetch, deploy); err != nil {
			return fmt.Errorf("failed to install %s: %w", serviceName, err)
		}
	}

	return nil
}

// Delete removes a service
func (m *Manager) Delete(instanceName, serviceName string) error {
	serviceDir := filepath.Join(m.servicesDir, serviceName)
	manifestsFile := filepath.Join(serviceDir, "manifests.yaml")

	if !storage.FileExists(manifestsFile) {
		return fmt.Errorf("service %s not found", serviceName)
	}

	cmd := exec.Command("kubectl", "delete", "-f", manifestsFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete service: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetStatus returns detailed status for a service
func (m *Manager) GetStatus(instanceName, serviceName string) (*Service, error) {
	// Get the correct namespace from the map
	namespace := serviceName + "-system" // default
	if deployment, ok := serviceDeployments[serviceName]; ok {
		namespace = deployment.namespace
	}

	service := &Service{
		Name:      serviceName,
		Namespace: namespace,
		Status:    m.checkServiceStatus(instanceName, serviceName),
	}

	return service, nil
}

// GetManifest returns the manifest for a service
func (m *Manager) GetManifest(serviceName string) (*ServiceManifest, error) {
	if manifest, ok := m.manifests[serviceName]; ok {
		return manifest, nil
	}
	return nil, fmt.Errorf("service %s not found or has no manifest", serviceName)
}

// GetServiceConfig returns the service configuration fields from the manifest
func (m *Manager) GetServiceConfig(serviceName string) (map[string]ConfigDefinition, error) {
	manifest, err := m.GetManifest(serviceName)
	if err != nil {
		return nil, err
	}
	return manifest.ServiceConfig, nil
}

// GetConfigReferences returns the config references from the manifest
func (m *Manager) GetConfigReferences(serviceName string) ([]string, error) {
	manifest, err := m.GetManifest(serviceName)
	if err != nil {
		return nil, err
	}
	return manifest.ConfigReferences, nil
}
