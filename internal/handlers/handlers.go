package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"wild-cloud-central/internal/config"
	"wild-cloud-central/internal/data"
	"wild-cloud-central/internal/dnsmasq"
)

// App represents the application with its dependencies
type App struct {
	Config         *config.Config
	StartTime      time.Time
	DataManager    *data.Manager
	DnsmasqManager *dnsmasq.ConfigGenerator
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{
		StartTime:      time.Now(),
		DataManager:    data.NewManager(),
		DnsmasqManager: dnsmasq.NewConfigGenerator(),
	}
}

// HealthHandler handles health check requests
func (app *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "wild-cloud-central",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StatusHandler handles status requests for the UI
func (app *App) StatusHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(app.StartTime)
	
	response := map[string]interface{}{
		"status":    "running",
		"version":   "1.0.0",
		"uptime":    uptime.String(),
		"timestamp": time.Now().UnixMilli(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetConfigHandler handles configuration retrieval requests
func (app *App) GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	log.Printf("Config handler called - config nil? %v", app.Config == nil)
	if app.Config != nil {
		log.Printf("Config values - Domain: '%s', DNS IP: '%s', Talos Version: '%s'", 
			app.Config.Cloud.Domain, 
			app.Config.Cloud.DNS.IP, 
			app.Config.Cluster.Nodes.Talos.Version)
		log.Printf("isConfigEmpty() returns: %v", app.Config.IsEmpty())
	}
	
	// Check if config is empty/uninitialized
	if app.Config == nil || app.Config.IsEmpty() {
		response := map[string]interface{}{
			"configured": false,
			"message":    "No configuration found. Please POST a configuration to /api/v1/config to get started.",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	response := map[string]interface{}{
		"configured": true,
		"config":     app.Config,
	}
	json.NewEncoder(w).Encode(response)
}

// CreateConfigHandler handles configuration creation requests
func (app *App) CreateConfigHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow config creation if no config exists
	if app.Config != nil && !app.Config.IsEmpty() {
		http.Error(w, "Configuration already exists. Use PUT to update.", http.StatusConflict)
		return
	}

	var newConfig config.Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set defaults
	if newConfig.Server.Port == 0 {
		newConfig.Server.Port = 5055
	}
	if newConfig.Server.Host == "" {
		newConfig.Server.Host = "0.0.0.0"
	}

	app.Config = &newConfig

	// Persist config to file
	paths := app.DataManager.GetPaths()
	if err := config.Save(app.Config, paths.ConfigFile); err != nil {
		log.Printf("Failed to save config: %v", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

// UpdateConfigHandler handles configuration update requests
func (app *App) UpdateConfigHandler(w http.ResponseWriter, r *http.Request) {
	// Check if config exists
	if app.Config == nil || app.Config.IsEmpty() {
		http.Error(w, "No configuration exists. Use POST to create initial configuration.", http.StatusNotFound)
		return
	}

	var newConfig config.Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	app.Config = &newConfig

	// Persist config to file
	paths := app.DataManager.GetPaths()
	if err := config.Save(app.Config, paths.ConfigFile); err != nil {
		log.Printf("Failed to save config: %v", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	// Regenerate and apply dnsmasq config
	if err := app.DnsmasqManager.WriteConfig(app.Config, paths.DnsmasqConf); err != nil {
		log.Printf("Failed to update dnsmasq config: %v", err)
		http.Error(w, "Failed to update dnsmasq config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// CORSMiddleware adds CORS headers to responses
func (app *App) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}