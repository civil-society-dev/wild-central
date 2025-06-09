package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Cloud struct {
		Domain         string `yaml:"domain"`
		InternalDomain string `yaml:"internalDomain"`
		DNS            struct {
			IP string `yaml:"ip"`
		} `yaml:"dns"`
		Router struct {
			IP string `yaml:"ip"`
		} `yaml:"router"`
		DHCPRange string `yaml:"dhcpRange"`
		Dnsmasq   struct {
			Interface string `yaml:"interface"`
		} `yaml:"dnsmasq"`
	} `yaml:"cloud"`
	Cluster struct {
		EndpointIP string `yaml:"endpointIp"`
		Nodes      struct {
			Talos struct {
				Version string `yaml:"version"`
			} `yaml:"talos"`
		} `yaml:"nodes"`
	} `yaml:"cluster"`
}

type App struct {
	config *Config
}

func main() {
	app := &App{}
	
	if err := app.loadConfig(); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	router := mux.NewRouter()
	app.setupRoutes(router)

	addr := fmt.Sprintf("%s:%d", app.config.Server.Host, app.config.Server.Port)
	log.Printf("Starting wild-cloud-central server on %s", addr)
	
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func (app *App) loadConfig() error {
	configPath := "/etc/wild-cloud-central/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "./config.yaml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = "/build/config.yaml"
		}
	}

	log.Printf("Loading config from: %s", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config file %s: %w", configPath, err)
	}

	app.config = &Config{}
	if err := yaml.Unmarshal(data, app.config); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}

	// Set defaults
	if app.config.Server.Port == 0 {
		app.config.Server.Port = 8080
	}
	if app.config.Server.Host == "" {
		app.config.Server.Host = "0.0.0.0"
	}

	return nil
}

func (app *App) setupRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/health", app.healthHandler).Methods("GET")
	router.HandleFunc("/api/v1/config", app.getConfigHandler).Methods("GET")
	router.HandleFunc("/api/v1/config", app.updateConfigHandler).Methods("PUT")
	router.HandleFunc("/api/v1/dnsmasq/config", app.getDnsmasqConfigHandler).Methods("GET")
	router.HandleFunc("/api/v1/dnsmasq/restart", app.restartDnsmasqHandler).Methods("POST")
	router.HandleFunc("/api/v1/pxe/assets", app.downloadPXEAssetsHandler).Methods("POST")
	
	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
}

func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "wild-cloud-central",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *App) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.config)
}

func (app *App) updateConfigHandler(w http.ResponseWriter, r *http.Request) {
	var newConfig Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	app.config = &newConfig
	
	// Persist config to file
	if err := app.saveConfig(); err != nil {
		log.Printf("Failed to save config: %v", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}
	
	// Regenerate and apply dnsmasq config
	if err := app.updateDnsmasqConfig(); err != nil {
		log.Printf("Failed to update dnsmasq config: %v", err)
		http.Error(w, "Failed to update dnsmasq config", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (app *App) getDnsmasqConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := app.generateDnsmasqConfig()
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(config))
}

func (app *App) restartDnsmasqHandler(w http.ResponseWriter, r *http.Request) {
	// Update dnsmasq config first
	if err := app.updateDnsmasqConfig(); err != nil {
		log.Printf("Failed to update dnsmasq config: %v", err)
		http.Error(w, "Failed to update dnsmasq config", http.StatusInternalServerError)
		return
	}
	
	// Restart dnsmasq service
	cmd := exec.Command("systemctl", "restart", "dnsmasq")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to restart dnsmasq: %v", err)
		http.Error(w, "Failed to restart dnsmasq service", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restarted"})
}

func (app *App) downloadPXEAssetsHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.downloadTalosAssets(); err != nil {
		log.Printf("Failed to download PXE assets: %v", err)
		http.Error(w, "Failed to download PXE assets", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "downloaded"})
}

func (app *App) generateDnsmasqConfig() string {
	template := `# Configuration file for dnsmasq.

# Basic Settings
interface=%s
listen-address=%s
domain-needed
bogus-priv
no-resolv

# DNS Forwarding
server=/%s/%s
server=/%s/%s
server=1.1.1.1
server=8.8.8.8

# --- DHCP Settings ---
dhcp-range=%s,12h
dhcp-option=3,%s
dhcp-option=6,%s

# --- PXE Booting ---
enable-tftp
tftp-root=/var/ftpd

dhcp-match=set:efi-x86_64,option:client-arch,7
dhcp-boot=tag:efi-x86_64,ipxe.efi
dhcp-boot=tag:!efi-x86_64,undionly.kpxe

dhcp-match=set:efi-arm64,option:client-arch,11
dhcp-boot=tag:efi-arm64,ipxe-arm64.efi

dhcp-userclass=set:ipxe,iPXE
dhcp-boot=tag:ipxe,http://%s/boot.ipxe

log-queries
log-dhcp
`

	return fmt.Sprintf(template,
		app.config.Cloud.Dnsmasq.Interface,
		app.config.Cloud.DNS.IP,
		app.config.Cloud.Domain,
		app.config.Cluster.EndpointIP,
		app.config.Cloud.InternalDomain,
		app.config.Cluster.EndpointIP,
		app.config.Cloud.DHCPRange,
		app.config.Cloud.Router.IP,
		app.config.Cloud.DNS.IP,
		app.config.Cloud.DNS.IP,
	)
}

func (app *App) saveConfig() error {
	configPath := "/etc/wild-cloud-central/config.yaml"
	if _, err := os.Stat("/etc/wild-cloud-central"); os.IsNotExist(err) {
		configPath = "./config.yaml"
	}

	data, err := yaml.Marshal(app.config)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

func (app *App) updateDnsmasqConfig() error {
	config := app.generateDnsmasqConfig()
	
	// Write to dnsmasq config file
	if err := os.WriteFile("/etc/dnsmasq.conf", []byte(config), 0644); err != nil {
		return fmt.Errorf("writing dnsmasq config: %w", err)
	}

	return nil
}

func (app *App) downloadTalosAssets() error {
	// Create directories
	assetsDir := "/var/www/html/talos"
	if err := os.MkdirAll(filepath.Join(assetsDir, "amd64"), 0755); err != nil {
		return fmt.Errorf("creating assets directory: %w", err)
	}

	// Create Talos bare metal configuration (schematic format)
	bareMetalConfig := `customization:
  extraKernelArgs:
    - net.ifnames=0
  systemExtensions:
    officialExtensions:
      - siderolabs/gvisor
      - siderolabs/intel-ucode`

	// Create Talos schematic
	var buf bytes.Buffer
	buf.WriteString(bareMetalConfig)
	
	resp, err := http.Post("https://factory.talos.dev/schematics", "text/yaml", &buf)
	if err != nil {
		return fmt.Errorf("creating Talos schematic: %w", err)
	}
	defer resp.Body.Close()

	var schematic struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&schematic); err != nil {
		return fmt.Errorf("decoding schematic response: %w", err)
	}

	log.Printf("Created Talos schematic with ID: %s", schematic.ID)

	// Download kernel
	kernelURL := fmt.Sprintf("https://pxe.factory.talos.dev/image/%s/%s/kernel-amd64", 
		schematic.ID, app.config.Cluster.Nodes.Talos.Version)
	if err := app.downloadFile(kernelURL, filepath.Join(assetsDir, "amd64", "vmlinuz")); err != nil {
		return fmt.Errorf("downloading kernel: %w", err)
	}

	// Download initramfs
	initramfsURL := fmt.Sprintf("https://pxe.factory.talos.dev/image/%s/%s/initramfs-amd64.xz", 
		schematic.ID, app.config.Cluster.Nodes.Talos.Version)
	if err := app.downloadFile(initramfsURL, filepath.Join(assetsDir, "amd64", "initramfs.xz")); err != nil {
		return fmt.Errorf("downloading initramfs: %w", err)
	}

	// Create boot.ipxe file
	bootScript := fmt.Sprintf(`#!ipxe
imgfree
kernel http://%s/amd64/vmlinuz talos.platform=metal console=tty0 init_on_alloc=1 slab_nomerge pti=on consoleblank=0 nvme_core.io_timeout=4294967295 printk.devkmsg=on ima_template=ima-ng ima_appraise=fix ima_hash=sha512 selinux=1 net.ifnames=0
initrd http://%s/amd64/initramfs.xz
boot
`, app.config.Cloud.DNS.IP, app.config.Cloud.DNS.IP)

	if err := os.WriteFile(filepath.Join(assetsDir, "boot.ipxe"), []byte(bootScript), 0644); err != nil {
		return fmt.Errorf("writing boot script: %w", err)
	}

	// Download iPXE bootloaders
	if err := os.MkdirAll("/var/ftpd", 0755); err != nil {
		return fmt.Errorf("creating ftpd directory: %w", err)
	}

	bootloaders := map[string]string{
		"http://boot.ipxe.org/ipxe.efi":           "/var/ftpd/ipxe.efi",
		"http://boot.ipxe.org/undionly.kpxe":     "/var/ftpd/undionly.kpxe",
		"http://boot.ipxe.org/arm64-efi/ipxe.efi": "/var/ftpd/ipxe-arm64.efi",
	}

	for url, path := range bootloaders {
		if err := app.downloadFile(url, path); err != nil {
			return fmt.Errorf("downloading %s: %w", url, err)
		}
	}

	log.Printf("Successfully downloaded PXE assets")
	return nil
}

func (app *App) downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}