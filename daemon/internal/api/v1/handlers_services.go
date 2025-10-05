package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"

	"github.com/wild-cloud/wild-central/daemon/internal/operations"
	"github.com/wild-cloud/wild-central/daemon/internal/services"
)

// ServicesList lists all base services
func (api *API) ServicesList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// List services
	servicesMgr := services.NewManager(api.dataDir, filepath.Join(api.directoryPath, "setup", "cluster-services"))
	svcList, err := servicesMgr.List(instanceName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list services: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"services": svcList,
	})
}

// ServicesGet returns a specific service
func (api *API) ServicesGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]
	serviceName := vars["service"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Get service
	servicesMgr := services.NewManager(api.dataDir, filepath.Join(api.directoryPath, "setup", "cluster-services"))
	service, err := servicesMgr.Get(instanceName, serviceName)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Service not found: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, service)
}

// ServicesInstall installs a service
func (api *API) ServicesInstall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Parse request
	var req struct {
		Name   string `json:"name"`
		Fetch  bool   `json:"fetch"`
		Deploy bool   `json:"deploy"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "service name is required")
		return
	}

	// Start install operation
	opsMgr := operations.NewManager(api.dataDir)
	opID, err := opsMgr.Start(instanceName, "install_service", req.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start operation: %v", err))
		return
	}

	// Install in background
	go func() {
		servicesMgr := services.NewManager(api.dataDir, filepath.Join(api.directoryPath, "setup", "cluster-services"))
		opsMgr.UpdateStatus(instanceName, opID, "running")

		if err := servicesMgr.Install(instanceName, req.Name, req.Fetch, req.Deploy); err != nil {
			opsMgr.Update(instanceName, opID, "failed", err.Error(), 0)
		} else {
			opsMgr.Update(instanceName, opID, "completed", "Service installed", 100)
		}
	}()

	respondJSON(w, http.StatusAccepted, map[string]string{
		"operation_id": opID,
		"message":      "Service installation initiated",
	})
}

// ServicesInstallAll installs all base services
func (api *API) ServicesInstallAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Parse request
	var req struct {
		Fetch  bool `json:"fetch"`
		Deploy bool `json:"deploy"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Use defaults if no body
		req.Deploy = true
	}

	// Start install operation
	opsMgr := operations.NewManager(api.dataDir)
	opID, err := opsMgr.Start(instanceName, "install_all_services", "all")
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start operation: %v", err))
		return
	}

	// Install in background
	go func() {
		servicesMgr := services.NewManager(api.dataDir, filepath.Join(api.directoryPath, "setup", "cluster-services"))
		opsMgr.UpdateStatus(instanceName, opID, "running")

		if err := servicesMgr.InstallAll(instanceName, req.Fetch, req.Deploy); err != nil {
			opsMgr.Update(instanceName, opID, "failed", err.Error(), 0)
		} else {
			opsMgr.Update(instanceName, opID, "completed", "All services installed", 100)
		}
	}()

	respondJSON(w, http.StatusAccepted, map[string]string{
		"operation_id": opID,
		"message":      "Services installation initiated",
	})
}

// ServicesDelete deletes a service
func (api *API) ServicesDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]
	serviceName := vars["service"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Start delete operation
	opsMgr := operations.NewManager(api.dataDir)
	opID, err := opsMgr.Start(instanceName, "delete_service", serviceName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start operation: %v", err))
		return
	}

	// Delete in background
	go func() {
		servicesMgr := services.NewManager(api.dataDir, filepath.Join(api.directoryPath, "setup", "cluster-services"))
		opsMgr.UpdateStatus(instanceName, opID, "running")

		if err := servicesMgr.Delete(instanceName, serviceName); err != nil {
			opsMgr.Update(instanceName, opID, "failed", err.Error(), 0)
		} else {
			opsMgr.Update(instanceName, opID, "completed", "Service deleted", 100)
		}
	}()

	respondJSON(w, http.StatusAccepted, map[string]string{
		"operation_id": opID,
		"message":      "Service deletion initiated",
	})
}

// ServicesGetStatus returns detailed service status
func (api *API) ServicesGetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]
	serviceName := vars["service"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Get status
	servicesMgr := services.NewManager(api.dataDir, filepath.Join(api.directoryPath, "setup", "cluster-services"))
	status, err := servicesMgr.GetStatus(instanceName, serviceName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get status: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, status)
}

// ServicesGetManifest returns the manifest for a service
func (api *API) ServicesGetManifest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]

	servicesMgr := services.NewManager(api.dataDir, filepath.Join(api.directoryPath, "setup", "cluster-services"))
	manifest, err := servicesMgr.GetManifest(serviceName)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Service not found: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, manifest)
}

// ServicesGetConfig returns the service configuration schema
func (api *API) ServicesGetConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]

	servicesMgr := services.NewManager(api.dataDir, filepath.Join(api.directoryPath, "setup", "cluster-services"))

	// Get manifest
	manifest, err := servicesMgr.GetManifest(serviceName)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Service not found: %v", err))
		return
	}

	// Return config schema
	response := map[string]interface{}{
		"configReferences": manifest.ConfigReferences,
		"serviceConfig":    manifest.ServiceConfig,
	}

	respondJSON(w, http.StatusOK, response)
}

// ServicesGetInstanceConfig returns current config values for a service instance
func (api *API) ServicesGetInstanceConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]
	serviceName := vars["service"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	servicesMgr := services.NewManager(api.dataDir, filepath.Join(api.directoryPath, "setup", "cluster-services"))

	// Get manifest to know which config paths to read
	manifest, err := servicesMgr.GetManifest(serviceName)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Service not found: %v", err))
		return
	}

	// Load instance config as map for dynamic path extraction
	configPath := filepath.Join(api.dataDir, "instances", instanceName, "config.yaml")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read instance config: %v", err))
		return
	}

	var instanceConfig map[string]interface{}
	if err := yaml.Unmarshal(configData, &instanceConfig); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse instance config: %v", err))
		return
	}

	// Extract values for all config paths
	configValues := make(map[string]interface{})

	// Add config references
	for _, path := range manifest.ConfigReferences {
		if value := getNestedValue(instanceConfig, path); value != nil {
			configValues[path] = value
		}
	}

	// Add service config
	for _, cfg := range manifest.ServiceConfig {
		if value := getNestedValue(instanceConfig, cfg.Path); value != nil {
			configValues[cfg.Path] = value
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"config": configValues,
	})
}

// getNestedValue retrieves a value from nested map using dot notation path
func getNestedValue(data map[string]interface{}, path string) interface{} {
	keys := strings.Split(path, ".")
	current := data

	for i, key := range keys {
		if i == len(keys)-1 {
			return current[key]
		}

		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}
