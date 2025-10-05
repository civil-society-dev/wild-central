package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/wild-cloud/wild-central/daemon/internal/cluster"
	"github.com/wild-cloud/wild-central/daemon/internal/operations"
)

// ClusterGenerateConfig generates cluster configuration
func (api *API) ClusterGenerateConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Parse request
	var config cluster.ClusterConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if config.ClusterName == "" {
		respondError(w, http.StatusBadRequest, "cluster_name is required")
		return
	}

	if config.VIP == "" {
		respondError(w, http.StatusBadRequest, "vip is required")
		return
	}

	// Generate configuration
	clusterMgr := cluster.NewManager(api.dataDir)
	if err := clusterMgr.GenerateConfig(instanceName, &config); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate config: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Cluster configuration generated successfully",
	})
}

// ClusterBootstrap bootstraps the cluster
func (api *API) ClusterBootstrap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Parse request
	var req struct {
		Node string `json:"node"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Node == "" {
		respondError(w, http.StatusBadRequest, "node is required")
		return
	}

	// Start bootstrap operation
	opsMgr := operations.NewManager(api.dataDir)
	opID, err := opsMgr.Start(instanceName, "bootstrap", req.Node)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start operation: %v", err))
		return
	}

	// Bootstrap in background
	go func() {
		clusterMgr := cluster.NewManager(api.dataDir)
		opsMgr.UpdateStatus(instanceName, opID, "running")

		if err := clusterMgr.Bootstrap(instanceName, req.Node); err != nil {
			opsMgr.Update(instanceName, opID, "failed", err.Error(), 0)
		} else {
			opsMgr.Update(instanceName, opID, "completed", "Bootstrap completed", 100)
		}
	}()

	respondJSON(w, http.StatusAccepted, map[string]string{
		"operation_id": opID,
		"message":      "Bootstrap initiated",
	})
}

// ClusterGetStatus returns cluster status
func (api *API) ClusterGetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Get status
	clusterMgr := cluster.NewManager(api.dataDir)
	status, err := clusterMgr.GetStatus(instanceName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get status: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, status)
}

// ClusterHealth returns cluster health checks
func (api *API) ClusterHealth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Get health checks
	clusterMgr := cluster.NewManager(api.dataDir)
	checks, err := clusterMgr.Health(instanceName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get health: %v", err))
		return
	}

	// Determine overall status
	overallStatus := "healthy"
	for _, check := range checks {
		if check.Status == "failing" {
			overallStatus = "unhealthy"
			break
		} else if check.Status == "warning" && overallStatus == "healthy" {
			overallStatus = "degraded"
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": overallStatus,
		"checks": checks,
	})
}

// ClusterGetKubeconfig returns the kubeconfig
func (api *API) ClusterGetKubeconfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Get kubeconfig
	clusterMgr := cluster.NewManager(api.dataDir)
	kubeconfig, err := clusterMgr.GetKubeconfig(instanceName)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Kubeconfig not found: %v", err))
		return
	}

	// Return as plain text
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(kubeconfig))
}

// ClusterGetTalosconfig returns the talosconfig
func (api *API) ClusterGetTalosconfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Get talosconfig
	clusterMgr := cluster.NewManager(api.dataDir)
	talosconfig, err := clusterMgr.GetTalosconfig(instanceName)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Talosconfig not found: %v", err))
		return
	}

	// Return as plain text
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(talosconfig))
}

// ClusterReset resets the cluster
func (api *API) ClusterReset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// Parse request
	var req struct {
		Confirm bool `json:"confirm"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !req.Confirm {
		respondError(w, http.StatusBadRequest, "Must confirm cluster reset")
		return
	}

	// Start reset operation
	opsMgr := operations.NewManager(api.dataDir)
	opID, err := opsMgr.Start(instanceName, "reset", instanceName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start operation: %v", err))
		return
	}

	// Reset in background
	go func() {
		clusterMgr := cluster.NewManager(api.dataDir)
		opsMgr.UpdateStatus(instanceName, opID, "running")

		if err := clusterMgr.Reset(instanceName, req.Confirm); err != nil {
			opsMgr.Update(instanceName, opID, "failed", err.Error(), 0)
		} else {
			opsMgr.Update(instanceName, opID, "completed", "Cluster reset completed", 100)
		}
	}()

	respondJSON(w, http.StatusAccepted, map[string]string{
		"operation_id": opID,
		"message":      "Cluster reset initiated",
	})
}
