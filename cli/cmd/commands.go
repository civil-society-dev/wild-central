package cmd

import (
	"fmt"
	"time"

	"github.com/r3labs/sse/v2"
	"github.com/spf13/cobra"

	"github.com/wild-cloud/wild-central/wild/internal/config"
	"github.com/wild-cloud/wild-central/wild/internal/prompt"
)

// Config commands
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage instance configuration",
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}
		key := args[0]

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/config", inst))
		if err != nil {
			return err
		}

		// Config is returned directly at top level
		if val, ok := resp.Data[key]; ok {
			fmt.Println(val)
		} else {
			return fmt.Errorf("key '%s' not found", key)
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}
		key, value := args[0], args[1]

		_, err = apiClient.Put(fmt.Sprintf("/api/v1/instances/%s/config", inst), map[string]string{
			key: value,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Configuration updated: %s = %s\n", key, value)
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show full configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/config", inst))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		return printYAML(resp.Data)
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
}

// Secret commands
var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secrets",
}

var secretGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get secret value (redacted)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}
		key := args[0]

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/secrets", inst))
		if err != nil {
			return err
		}

		if secrets := resp.GetMap("secrets"); secrets != nil {
			if val, ok := secrets[key]; ok {
				fmt.Println(val)
			} else {
				return fmt.Errorf("secret '%s' not found", key)
			}
		}
		return nil
	},
}

var secretSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set secret value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}
		key, value := args[0], args[1]

		_, err = apiClient.Put(fmt.Sprintf("/api/v1/instances/%s/secrets", inst), map[string]string{
			key: value,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Secret updated: %s\n", key)
		return nil
	},
}

func init() {
	secretCmd.AddCommand(secretGetCmd)
	secretCmd.AddCommand(secretSetCmd)
}

// Node commands
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage nodes",
}

var nodeDiscoverCmd = &cobra.Command{
	Use:   "discover <ip>...",
	Short: "Discover nodes on network",
	Long:  "Discover nodes on the network by scanning the provided IP addresses or ranges",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		fmt.Printf("Starting discovery for %d IP(s)...\n", len(args))
		_, err = apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/nodes/discover", inst), map[string]interface{}{
			"ip_list": args,
		})
		if err != nil {
			return err
		}

		// Poll for completion
		fmt.Println("Scanning nodes...")
		for {
			resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/discovery", inst))
			if err != nil {
				return err
			}

			active, _ := resp.Data["active"].(bool)
			if !active {
				// Discovery complete
				nodesFound := resp.GetArray("nodes_found")
				if len(nodesFound) == 0 {
					fmt.Println("\nNo Talos nodes found")
					return nil
				}

				fmt.Printf("\nFound %d node(s):\n\n", len(nodesFound))
				fmt.Printf("%-15s  %-17s  %-12s  %-10s\n", "IP", "MAC", "INTERFACE", "VERSION")
				fmt.Println("---------------------------------------------------------------")
				for _, node := range nodesFound {
					if m, ok := node.(map[string]interface{}); ok {
						fmt.Printf("%-15s  %-17s  %-12s  %-10s\n",
							m["ip"], m["mac"], m["interface"], m["version"])
					}
				}
				return nil
			}

			// Still running, wait a bit
			fmt.Print(".")
			time.Sleep(500 * time.Millisecond)
		}
	},
}

var nodeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/nodes", inst))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		if outputFormat == "yaml" {
			return printYAML(resp.Data)
		}

		nodes := resp.GetArray("nodes")
		if len(nodes) == 0 {
			fmt.Println("No nodes found")
			return nil
		}

		fmt.Printf("%-17s  %-20s  %-12s  %-15s\n", "MAC", "HOSTNAME", "ROLE", "TARGET IP")
		fmt.Println("-------------------------------------------------------------------------")
		for _, node := range nodes {
			if m, ok := node.(map[string]interface{}); ok {
				fmt.Printf("%-17s  %-20s  %-12s  %-15s\n",
					m["mac"], m["hostname"], m["role"], m["target_ip"])
			}
		}
		return nil
	},
}

var nodeShowCmd = &cobra.Command{
	Use:   "show <hostname>",
	Short: "Show node details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/nodes/%s", inst, args[0]))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		if outputFormat == "yaml" {
			return printYAML(resp.Data)
		}

		// Text format - show node details
		fmt.Printf("Hostname:     %s\n", resp.GetString("hostname"))
		fmt.Printf("Role:         %s\n", resp.GetString("role"))
		fmt.Printf("Target IP:    %s\n", resp.GetString("target_ip"))
		fmt.Printf("MAC:          %s\n", resp.GetString("mac"))
		fmt.Printf("Disk:         %s\n", resp.GetString("disk"))
		fmt.Printf("Interface:    %s\n", resp.GetString("interface"))
		fmt.Printf("Version:      %s\n", resp.GetString("version"))
		fmt.Printf("Schematic ID: %s\n", resp.GetString("schematic_id"))
		fmt.Printf("Configured:   %v\n", resp.Data["configured"])
		fmt.Printf("Deployed:     %v\n", resp.Data["deployed"])

		return nil
	},
}

var nodeAddCmd = &cobra.Command{
	Use:   "add <mac> <hostname> <role>",
	Short: "Add a node",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		_, err = apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/nodes", inst), map[string]string{
			"mac":      args[0],
			"hostname": args[1],
			"role":     args[2],
		})
		if err != nil {
			return err
		}

		fmt.Printf("Node added: %s (%s)\n", args[1], args[0])
		return nil
	},
}

var nodeSetupCmd = &cobra.Command{
	Use:   "setup <mac>",
	Short: "Setup Talos on node",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/nodes/%s/setup", inst, args[0]), nil)
		if err != nil {
			return err
		}

		fmt.Printf("Node setup started\n")
		if opID := resp.GetString("operation_id"); opID != "" {
			fmt.Printf("Operation ID: %s\n", opID)
		}
		return nil
	},
}

var nodeDeleteCmd = &cobra.Command{
	Use:   "delete <mac>",
	Short: "Delete a node",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		_, err = apiClient.Delete(fmt.Sprintf("/api/v1/instances/%s/nodes/%s", inst, args[0]))
		if err != nil {
			return err
		}

		fmt.Printf("Node deleted: %s\n", args[0])
		return nil
	},
}

func init() {
	nodeCmd.AddCommand(nodeDiscoverCmd)
	nodeCmd.AddCommand(nodeListCmd)
	nodeCmd.AddCommand(nodeShowCmd)
	nodeCmd.AddCommand(nodeAddCmd)
	nodeCmd.AddCommand(nodeSetupCmd)
	nodeCmd.AddCommand(nodeDeleteCmd)
}

// PXE commands
var pxeCmd = &cobra.Command{
	Use:   "pxe",
	Short: "Manage PXE assets",
}

var pxeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List PXE assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/pxe/assets", inst))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		assets := resp.GetArray("assets")
		if len(assets) == 0 {
			fmt.Println("No PXE assets found")
			return nil
		}

		fmt.Printf("%-20s  %-15s  %-12s\n", "NAME", "VERSION", "STATUS")
		fmt.Println("--------------------------------------------------")
		for _, asset := range assets {
			if m, ok := asset.(map[string]interface{}); ok {
				fmt.Printf("%-20s  %-15s  %-12s\n",
					m["name"], m["version"], m["status"])
			}
		}
		return nil
	},
}

var pxeDownloadCmd = &cobra.Command{
	Use:   "download <asset>",
	Short: "Download PXE asset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/pxe/assets/%s/download", inst, args[0]), nil)
		if err != nil {
			return err
		}

		fmt.Printf("Download started for: %s\n", args[0])
		if opID := resp.GetString("operation_id"); opID != "" {
			fmt.Printf("Operation ID: %s\n", opID)
		}
		return nil
	},
}

func init() {
	pxeCmd.AddCommand(pxeListCmd)
	pxeCmd.AddCommand(pxeDownloadCmd)
}

// Cluster commands
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage cluster",
}

var clusterBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/cluster/bootstrap", inst), nil)
		if err != nil {
			return err
		}

		fmt.Println("Cluster bootstrap started")
		if opID := resp.GetString("operation_id"); opID != "" {
			fmt.Printf("Operation ID: %s\n", opID)
		}
		return nil
	},
}

var clusterStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get cluster status",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/cluster/status", inst))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		return printYAML(resp.GetMap("status"))
	},
}

var clusterHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check cluster health",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/cluster/health", inst))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		return printYAML(resp.GetMap("health"))
	},
}

var clusterKubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Get kubeconfig",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/cluster/kubeconfig", inst))
		if err != nil {
			return err
		}

		fmt.Println(resp.GetString("kubeconfig"))
		return nil
	},
}

func init() {
	clusterCmd.AddCommand(clusterBootstrapCmd)
	clusterCmd.AddCommand(clusterStatusCmd)
	clusterCmd.AddCommand(clusterHealthCmd)
	clusterCmd.AddCommand(clusterKubeconfigCmd)
}

// Service commands
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/services", inst))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		services := resp.GetArray("services")
		if len(services) == 0 {
			fmt.Println("No services found")
			return nil
		}

		fmt.Printf("%-20s  %-12s\n", "NAME", "STATUS")
		fmt.Println("----------------------------------")
		for _, svc := range services {
			if m, ok := svc.(map[string]interface{}); ok {
				fmt.Printf("%-20s  %-12s\n", m["name"], m["status"])
			}
		}
		return nil
	},
}

// ServiceManifest matches the daemon's ServiceManifest structure
type ServiceManifest struct {
	Name             string                      `json:"name"`
	Description      string                      `json:"description"`
	Namespace        string                      `json:"namespace"`
	ConfigReferences []string                    `json:"configReferences"`
	ServiceConfig    map[string]ConfigDefinition `json:"serviceConfig"`
}

// ConfigDefinition defines config that should be prompted during service setup
type ConfigDefinition struct {
	Path    string `json:"path"`
	Prompt  string `json:"prompt"`
	Default string `json:"default"`
	Type    string `json:"type"`
}

// ConfigUpdate represents a single configuration update
type ConfigUpdate struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

var (
	fetchFlag    bool
	noDeployFlag bool
)

var serviceInstallCmd = &cobra.Command{
	Use:   "install <service>",
	Short: "Install a service with interactive configuration",
	Long: `Install and configure a cluster service.

This command orchestrates the complete service installation lifecycle:
  1. Fetch service files from Wild Cloud Directory (if needed or --fetch)
  2. Validate configuration requirements
  3. Prompt for any missing service configuration
  4. Update instance configuration
  5. Compile templates using gomplate
  6. Deploy service to cluster (unless --no-deploy)

Examples:
  # Configure and deploy (most common)
  wild service install metallb

  # Fetch fresh templates and deploy
  wild service install metallb --fetch

  # Configure only, skip deployment
  wild service install metallb --no-deploy

  # Fetch fresh templates, configure only
  wild service install metallb --fetch --no-deploy

  # Use cached templates (default if files exist)
  wild service install traefik
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceName := args[0]

		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		fmt.Printf("Installing service: %s\n", serviceName)

		// Step 1: Fetch service manifest
		fmt.Println("\nFetching service manifest...")
		manifestResp, err := apiClient.Get(fmt.Sprintf("/api/v1/services/%s/manifest", serviceName))
		if err != nil {
			return fmt.Errorf("failed to fetch manifest: %w", err)
		}

		// Parse manifest
		var manifest ServiceManifest
		manifestData := manifestResp.Data
		// API returns camelCase field names
		if name, ok := manifestData["name"].(string); ok {
			manifest.Name = name
		}
		if desc, ok := manifestData["description"].(string); ok {
			manifest.Description = desc
		}
		if namespace, ok := manifestData["namespace"].(string); ok {
			manifest.Namespace = namespace
		}
		if refs, ok := manifestData["configReferences"].([]interface{}); ok {
			manifest.ConfigReferences = make([]string, len(refs))
			for i, ref := range refs {
				if s, ok := ref.(string); ok {
					manifest.ConfigReferences[i] = s
				}
			}
		}
		if svcConfig, ok := manifestData["serviceConfig"].(map[string]interface{}); ok {
			manifest.ServiceConfig = make(map[string]ConfigDefinition)
			for key, val := range svcConfig {
				if cfgMap, ok := val.(map[string]interface{}); ok {
					cfg := ConfigDefinition{}
					if path, ok := cfgMap["path"].(string); ok {
						cfg.Path = path
					}
					if prompt, ok := cfgMap["prompt"].(string); ok {
						cfg.Prompt = prompt
					}
					if def, ok := cfgMap["default"].(string); ok {
						cfg.Default = def
					}
					if typ, ok := cfgMap["type"].(string); ok {
						cfg.Type = typ
					}
					manifest.ServiceConfig[key] = cfg
				}
			}
		}

		fmt.Printf("Service: %s - %s\n", manifest.Name, manifest.Description)

		// Step 2: Fetch current config
		fmt.Println("\nFetching instance configuration...")
		configResp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/config", inst))
		if err != nil {
			return fmt.Errorf("failed to fetch config: %w", err)
		}

		instanceConfig := configResp.Data

		// Step 3: Validate configReferences
		if len(manifest.ConfigReferences) > 0 {
			fmt.Println("\nValidating configuration references...")
			missingPaths := config.ValidatePaths(instanceConfig, manifest.ConfigReferences)
			if len(missingPaths) > 0 {
				fmt.Println("\nERROR: Missing required configuration values:")
				for _, path := range missingPaths {
					fmt.Printf("  - %s\n", path)
				}
				fmt.Println("\nPlease set these configuration values before installing this service.")
				fmt.Printf("Use: wild config set <key> <value>\n")
				return fmt.Errorf("missing required configuration")
			}
			fmt.Println("All required configuration references are present.")
		}

		// Step 4: Process serviceConfig - prompt for missing values
		var updates []ConfigUpdate

		if len(manifest.ServiceConfig) > 0 {
			fmt.Println("\nConfiguring service parameters...")

			for key, cfg := range manifest.ServiceConfig {
				// Check if path already set
				existingValue := config.GetValue(instanceConfig, cfg.Path)
				if existingValue != nil && existingValue != "" && existingValue != "null" {
					fmt.Printf("  %s: %v (already set)\n", cfg.Path, existingValue)
					continue
				}

				// Expand default template
				defaultValue := cfg.Default
				if defaultValue != "" {
					expanded, err := config.ExpandTemplate(defaultValue, instanceConfig)
					if err != nil {
						return fmt.Errorf("failed to expand template for %s: %w", key, err)
					}
					defaultValue = expanded
				}

				// Prompt user
				var value string
				switch cfg.Type {
				case "int":
					intVal, err := prompt.Int(cfg.Prompt, 0)
					if err != nil {
						return fmt.Errorf("failed to read input for %s: %w", key, err)
					}
					value = fmt.Sprintf("%d", intVal)
				case "bool":
					boolVal, err := prompt.Bool(cfg.Prompt, false)
					if err != nil {
						return fmt.Errorf("failed to read input for %s: %w", key, err)
					}
					if boolVal {
						value = "true"
					} else {
						value = "false"
					}
				default: // string
					var err error
					value, err = prompt.String(cfg.Prompt, defaultValue)
					if err != nil {
						return fmt.Errorf("failed to read input for %s: %w", key, err)
					}
				}

				// Add to updates
				updates = append(updates, ConfigUpdate{
					Path:  cfg.Path,
					Value: value,
				})

				fmt.Printf("  %s: %s\n", cfg.Path, value)
			}
		}

		// Step 5: Update configuration if needed
		if len(updates) > 0 {
			fmt.Println("\nUpdating instance configuration...")
			_, err = apiClient.Patch(
				fmt.Sprintf("/api/v1/instances/%s/config", inst),
				map[string]interface{}{
					"updates": updates,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to update configuration: %w", err)
			}
			fmt.Printf("Configuration updated (%d values)\n", len(updates))
		}

		// Step 6: Install service with lifecycle control
		if noDeployFlag {
			fmt.Println("\nConfiguring service...")
		} else {
			fmt.Println("\nInstalling service...")
		}

		installResp, err := apiClient.Post(
			fmt.Sprintf("/api/v1/instances/%s/services", inst),
			map[string]interface{}{
				"name":   serviceName,
				"fetch":  fetchFlag,
				"deploy": !noDeployFlag,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to install service: %w", err)
		}

		// Show appropriate success message
		if noDeployFlag {
			fmt.Printf("\n✓ Service configured: %s\n", serviceName)
			fmt.Printf("  Templates compiled and ready to deploy\n")
			fmt.Printf("  To deploy later, run: wild service install %s\n", serviceName)
		} else {
			// Stream installation output
			opID := installResp.GetString("operation_id")
			if opID != "" {
				fmt.Printf("Installing service: %s\n\n", serviceName)
				if err := streamOperationOutput(opID); err != nil {
					// If streaming fails, show operation ID for manual monitoring
					fmt.Printf("\nCouldn't stream output: %v\n", err)
					fmt.Printf("Operation ID: %s\n", opID)
					fmt.Printf("Monitor with: wild operation get %s\n", opID)
				} else {
					fmt.Printf("\n✓ Service installed successfully: %s\n", serviceName)
				}
			}
		}

		return nil
	},
}

func init() {
	serviceInstallCmd.Flags().BoolVar(&fetchFlag, "fetch", false, "Fetch fresh templates from directory before installing")
	serviceInstallCmd.Flags().BoolVar(&noDeployFlag, "no-deploy", false, "Configure and compile only, skip deployment")

	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceInstallCmd)
}

// App commands
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage applications",
}

var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.Get("/api/v1/apps")
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		apps := resp.GetArray("apps")
		if len(apps) == 0 {
			fmt.Println("No apps found")
			return nil
		}

		fmt.Printf("%-20s  %-30s\n", "NAME", "DESCRIPTION")
		fmt.Println("-----------------------------------------------------")
		for _, app := range apps {
			if m, ok := app.(map[string]interface{}); ok {
				fmt.Printf("%-20s  %-30s\n", m["name"], m["description"])
			}
		}
		return nil
	},
}

var appListDeployedCmd = &cobra.Command{
	Use:   "list-deployed",
	Short: "List deployed apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/apps", inst))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		apps := resp.GetArray("apps")
		if len(apps) == 0 {
			fmt.Println("No deployed apps found")
			return nil
		}

		fmt.Printf("%-20s  %-12s\n", "NAME", "STATUS")
		fmt.Println("----------------------------------")
		for _, app := range apps {
			if m, ok := app.(map[string]interface{}); ok {
				fmt.Printf("%-20s  %-12s\n", m["name"], m["status"])
			}
		}
		return nil
	},
}

var appAddCmd = &cobra.Command{
	Use:   "add <app>",
	Short: "Add app to instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		_, err = apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/apps", inst), map[string]string{
			"app": args[0],
		})
		if err != nil {
			return err
		}

		fmt.Printf("App added: %s\n", args[0])
		return nil
	},
}

var appDeployCmd = &cobra.Command{
	Use:   "deploy <app>",
	Short: "Deploy an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/apps/%s/deploy", inst, args[0]), nil)
		if err != nil {
			return err
		}

		fmt.Printf("App deployment started: %s\n", args[0])
		if opID := resp.GetString("operation_id"); opID != "" {
			fmt.Printf("Operation ID: %s\n", opID)
		}
		return nil
	},
}

var appDeleteCmd = &cobra.Command{
	Use:   "delete <app>",
	Short: "Delete an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		_, err = apiClient.Delete(fmt.Sprintf("/api/v1/instances/%s/apps/%s", inst, args[0]))
		if err != nil {
			return err
		}

		fmt.Printf("App deleted: %s\n", args[0])
		return nil
	},
}

func init() {
	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appListDeployedCmd)
	appCmd.AddCommand(appAddCmd)
	appCmd.AddCommand(appDeployCmd)
	appCmd.AddCommand(appDeleteCmd)
}

// Backup/Restore commands
var backupCmd = &cobra.Command{
	Use:   "backup <app>",
	Short: "Backup an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/apps/%s/backup", inst, args[0]), nil)
		if err != nil {
			return err
		}

		fmt.Printf("Backup started: %s\n", args[0])
		if opID := resp.GetString("operation_id"); opID != "" {
			fmt.Printf("Operation ID: %s\n", opID)
		}
		return nil
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore <app>",
	Short: "Restore an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/apps/%s/restore", inst, args[0]), nil)
		if err != nil {
			return err
		}

		fmt.Printf("Restore started: %s\n", args[0])
		if opID := resp.GetString("operation_id"); opID != "" {
			fmt.Printf("Operation ID: %s\n", opID)
		}
		return nil
	},
}

// Utility commands
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check cluster health",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.Get("/api/v1/utilities/health")
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		return printYAML(resp.Data)
	},
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Dashboard operations",
}

var dashboardTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get dashboard token",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.Get("/api/v1/utilities/dashboard/token")
		if err != nil {
			return err
		}

		fmt.Println(resp.GetString("token"))
		return nil
	},
}

func init() {
	dashboardCmd.AddCommand(dashboardTokenCmd)
}

var nodeIPCmd = &cobra.Command{
	Use:   "node-ip",
	Short: "Get control plane IP",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.Get("/api/v1/utilities/controlplane/ip")
		if err != nil {
			return err
		}

		fmt.Println(resp.GetString("ip"))
		return nil
	},
}

// Operation commands
var operationCmd = &cobra.Command{
	Use:   "operation",
	Short: "Manage operations",
}

var operationGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get operation status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/operations/%s?instance=%s", args[0], instanceName))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		return printYAML(resp.Data)
	},
}

var operationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.Get("/api/v1/operations")
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		ops := resp.GetArray("operations")
		if len(ops) == 0 {
			fmt.Println("No operations found")
			return nil
		}

		fmt.Printf("%-10s  %-12s  %-12s  %-10s\n", "ID", "TYPE", "STATUS", "PROGRESS")
		fmt.Println("--------------------------------------------------------")
		for _, op := range ops {
			if m, ok := op.(map[string]interface{}); ok {
				fmt.Printf("%-10s  %-12s  %-12s  %d%%\n",
					m["id"], m["type"], m["status"], int(m["progress"].(float64)))
			}
		}
		return nil
	},
}

func init() {
	operationCmd.AddCommand(operationGetCmd)
	operationCmd.AddCommand(operationListCmd)
}

// streamOperationOutput streams operation output via SSE
func streamOperationOutput(opID string) error {
	// Get base URL
	baseURL := daemonURL
	if baseURL == "" {
		baseURL = config.GetDaemonURL()
	}

	// Connect to SSE stream
	url := fmt.Sprintf("%s/api/v1/operations/%s/stream?instance=%s", baseURL, opID, instanceName)
	client := sse.NewClient(url)
	events := make(chan *sse.Event)

	err := client.SubscribeChan("messages", events)
	if err != nil {
		return fmt.Errorf("failed to subscribe to SSE: %w", err)
	}

	// Poll for completion in background
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			resp, err := apiClient.Get(fmt.Sprintf("/api/v1/operations/%s?instance=%s", opID, instanceName))
			if err == nil {
				status := resp.GetString("status")
				if status == "completed" || status == "failed" {
					time.Sleep(500 * time.Millisecond) // Give SSE time to flush
					done <- true
					return
				}
			}
		}
	}()

	// Stream events
	for {
		select {
		case msg, ok := <-events:
			if !ok {
				// Channel closed
				return nil
			}
			if msg != nil {
				// Check for completion event
				if string(msg.Event) == "complete" {
					return nil
				}
				// Print data with newline
				if msg.Data != nil {
					fmt.Println(string(msg.Data))
				}
			}
		case <-done:
			return nil
		}
	}
}
