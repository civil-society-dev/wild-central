package v1

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/wild-cloud/wild-central/daemon/internal/operations"
)

// OperationGet returns operation status
func (api *API) OperationGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	opID := vars["id"]

	// Extract instance name from query param or header
	instanceName := r.URL.Query().Get("instance")
	if instanceName == "" {
		respondError(w, http.StatusBadRequest, "instance parameter is required")
		return
	}

	// Get operation
	opsMgr := operations.NewManager(api.dataDir)
	op, err := opsMgr.GetByInstance(instanceName, opID)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Operation not found: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, op)
}

// OperationList returns all operations for an instance
func (api *API) OperationList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceName := vars["name"]

	// Validate instance exists
	if err := api.instance.ValidateInstance(instanceName); err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("Instance not found: %v", err))
		return
	}

	// List operations
	opsMgr := operations.NewManager(api.dataDir)
	ops, err := opsMgr.List(instanceName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list operations: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"operations": ops,
	})
}

// OperationCancel cancels an operation
func (api *API) OperationCancel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	opID := vars["id"]

	// Extract instance name from query param
	instanceName := r.URL.Query().Get("instance")
	if instanceName == "" {
		respondError(w, http.StatusBadRequest, "instance parameter is required")
		return
	}

	// Cancel operation
	opsMgr := operations.NewManager(api.dataDir)
	if err := opsMgr.Cancel(instanceName, opID); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to cancel operation: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Operation cancelled",
		"id":      opID,
	})
}
