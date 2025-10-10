package cmd

import (
	"fmt"
	"os"
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
				fmt.Printf("%-15s  %-12s  %-10s\n", "IP", "INTERFACE", "VERSION")
				fmt.Println("-----------------------------------------------")
				for _, node := range nodesFound {
					if m, ok := node.(map[string]interface{}); ok {
						fmt.Printf("%-15s  %-12s  %-10s\n",
							m["ip"], m["interface"], m["version"])
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

var nodeDetectCmd = &cobra.Command{
	Use:   "detect <ip>",
	Short: "Detect hardware on a single node",
	Long: `Detect hardware configuration on a single node in maintenance mode.

This queries the node for available network interfaces and disks, helping you
decide which hardware to use when adding the node to the cluster.

Example:
  wild node detect 192.168.1.31

Output shows:
  - Available network interfaces
  - Available disks with sizes
  - Recommended selections`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		nodeIP := args[0]

		// Call API to detect hardware
		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/nodes/detect", inst), map[string]string{
			"ip": nodeIP,
		})
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		if outputFormat == "yaml" {
			return printYAML(resp.Data)
		}

		// Text format - show hardware details
		fmt.Printf("Hardware detected for node at %s:\n\n", nodeIP)

		if iface := resp.GetString("interface"); iface != "" {
			fmt.Printf("Interface:        %s\n", iface)
		}

		if disks := resp.GetArray("disks"); len(disks) > 0 {
			fmt.Printf("\nAvailable Disks:\n")
			for _, diskData := range disks {
				diskMap, ok := diskData.(map[string]interface{})
				if !ok {
					continue
				}

				path, _ := diskMap["path"].(string)
				size, _ := diskMap["size"].(float64) // JSON numbers are float64

				// Format size in GB/TB
				sizeGB := size / (1024 * 1024 * 1024)
				var sizeStr string
				if sizeGB >= 1000 {
					sizeStr = fmt.Sprintf("%.1f TB", sizeGB/1024)
				} else {
					sizeStr = fmt.Sprintf("%.1f GB", sizeGB)
				}

				fmt.Printf("  - %s (%s)\n", path, sizeStr)
			}
		}

		if selected := resp.GetString("selected_disk"); selected != "" {
			fmt.Printf("\nRecommended Disk: %s\n", selected)
		}

		fmt.Printf("\nTo add this node:\n")
		fmt.Printf("  wild node add <hostname> <role> --current-ip %s --target-ip <target-ip> --disk <disk> --interface <interface>\n", nodeIP)

		return nil
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

		fmt.Printf("%-20s  %-12s  %-15s\n", "HOSTNAME", "ROLE", "TARGET IP")
		fmt.Println("-----------------------------------------------------")
		for _, node := range nodes {
			if m, ok := node.(map[string]interface{}); ok {
				fmt.Printf("%-20s  %-12s  %-15s\n",
					m["hostname"], m["role"], m["target_ip"])
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
	Use:   "add <hostname> <role>",
	Short: "Add a node to cluster configuration",
	Long: `Add a node to the cluster configuration with required hardware details.

Role must be either 'controlplane' or 'worker'.

The node configuration will be stored in the instance config and used during apply.

Examples:
  # Node in maintenance mode (PXE booted)
  wild node add control-1 controlplane --current-ip 192.168.1.100 --target-ip 192.168.1.31 --disk /dev/sda

  # Node already applied (unusual, only if config was removed manually)
  wild node add worker-1 worker --target-ip 192.168.1.32 --disk /dev/nvme0n1`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		// Get flags
		targetIP, _ := cmd.Flags().GetString("target-ip")
		currentIP, _ := cmd.Flags().GetString("current-ip")
		disk, _ := cmd.Flags().GetString("disk")
		iface, _ := cmd.Flags().GetString("interface")
		schematicID, _ := cmd.Flags().GetString("schematic-id")
		maintenance, _ := cmd.Flags().GetBool("maintenance")

		// Build request body
		body := map[string]interface{}{
			"hostname": args[0],
			"role":     args[1],
		}

		if targetIP != "" {
			body["target_ip"] = targetIP
		}
		if currentIP != "" {
			body["current_ip"] = currentIP
		}
		if disk != "" {
			body["disk"] = disk
		}
		if iface != "" {
			body["interface"] = iface
		}
		if schematicID != "" {
			body["schematic_id"] = schematicID
		}
		if maintenance {
			body["maintenance"] = true
		}

		_, err = apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/nodes", inst), body)
		if err != nil {
			return err
		}

		fmt.Printf("Node added: %s (%s)\n", args[0], args[1])
		if targetIP != "" {
			fmt.Printf("  Target IP: %s\n", targetIP)
		}
		if disk != "" {
			fmt.Printf("  Disk: %s\n", disk)
		}
		return nil
	},
}

var nodeApplyCmd = &cobra.Command{
	Use:   "apply <hostname>",
	Short: "Apply Talos configuration to node",
	Long: `Generate and apply Talos configuration to a node.

This command:
1. Auto-fetches patch templates if missing
2. Generates node-specific configuration from templates
3. Merges base config with node patch
4. Applies configuration to node (using --insecure if in maintenance mode)
5. Updates node state after successful application

Examples:
  # Apply to node in maintenance mode (PXE booted)
  wild node apply control-1

  # Re-apply to production node (updates configuration)
  wild node apply worker-1`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/nodes/%s/apply", inst, args[0]), nil)
		if err != nil {
			return err
		}

		fmt.Printf("Node configuration applied: %s\n", args[0])
		if msg := resp.GetString("message"); msg != "" {
			fmt.Printf("%s\n", msg)
		}
		return nil
	},
}

var nodeUpdateCmd = &cobra.Command{
	Use:   "update <hostname>",
	Short: "Update node configuration",
	Long: `Update existing node configuration with partial updates.

This command modifies node properties without requiring all fields.

Examples:
  # Update disk after hardware change
  wild node update worker-1 --disk /dev/sdb

  # Move node to maintenance mode
  wild node update control-1 --current-ip 192.168.1.100 --maintenance

  # Clear maintenance after successful apply
  wild node update control-1 --no-maintenance`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		// Get flags
		targetIP, _ := cmd.Flags().GetString("target-ip")
		currentIP, _ := cmd.Flags().GetString("current-ip")
		disk, _ := cmd.Flags().GetString("disk")
		iface, _ := cmd.Flags().GetString("interface")
		schematicID, _ := cmd.Flags().GetString("schematic-id")
		maintenance, _ := cmd.Flags().GetBool("maintenance")
		noMaintenance, _ := cmd.Flags().GetBool("no-maintenance")

		// Build request body with only provided fields
		body := map[string]interface{}{}

		if targetIP != "" {
			body["target_ip"] = targetIP
		}
		if currentIP != "" {
			body["current_ip"] = currentIP
		}
		if disk != "" {
			body["disk"] = disk
		}
		if iface != "" {
			body["interface"] = iface
		}
		if schematicID != "" {
			body["schematic_id"] = schematicID
		}
		if maintenance {
			body["maintenance"] = true
		}
		if noMaintenance {
			body["maintenance"] = false
		}

		if len(body) == 0 {
			return fmt.Errorf("no updates specified")
		}

		_, err = apiClient.Put(fmt.Sprintf("/api/v1/instances/%s/nodes/%s", inst, args[0]), body)
		if err != nil {
			return err
		}

		fmt.Printf("Node updated: %s\n", args[0])
		return nil
	},
}

var nodeFetchTemplatesCmd = &cobra.Command{
	Use:   "fetch-patch-templates",
	Short: "Fetch patch templates from directory",
	Long: `Copy latest patch templates from directory/setup/cluster-nodes/patch.templates to instance.

This command is automatically called by 'apply' if templates are missing.
You can use it manually to update templates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		_, err = apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/nodes/fetch-templates", inst), nil)
		if err != nil {
			return err
		}

		fmt.Println("Templates fetched successfully")
		return nil
	},
}

var nodeDeleteCmd = &cobra.Command{
	Use:   "delete <hostname>",
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
	nodeCmd.AddCommand(nodeDetectCmd)
	nodeCmd.AddCommand(nodeListCmd)
	nodeCmd.AddCommand(nodeShowCmd)
	nodeCmd.AddCommand(nodeAddCmd)
	nodeCmd.AddCommand(nodeApplyCmd)
	nodeCmd.AddCommand(nodeUpdateCmd)
	nodeCmd.AddCommand(nodeFetchTemplatesCmd)
	nodeCmd.AddCommand(nodeDeleteCmd)

	// Add flags to node add command
	nodeAddCmd.Flags().String("target-ip", "", "Target IP address for production")
	nodeAddCmd.Flags().String("current-ip", "", "Current IP address (for maintenance mode)")
	nodeAddCmd.Flags().String("disk", "", "Disk device (required, e.g., /dev/sda)")
	nodeAddCmd.Flags().String("interface", "", "Network interface (optional, e.g., eth0)")
	nodeAddCmd.Flags().String("schematic-id", "", "Talos schematic ID (optional, uses instance default)")
	nodeAddCmd.Flags().Bool("maintenance", false, "Mark node as in maintenance mode")

	// Add flags to node update command
	nodeUpdateCmd.Flags().String("target-ip", "", "Update target IP address")
	nodeUpdateCmd.Flags().String("current-ip", "", "Update current IP address")
	nodeUpdateCmd.Flags().String("disk", "", "Update disk device")
	nodeUpdateCmd.Flags().String("interface", "", "Update network interface")
	nodeUpdateCmd.Flags().String("schematic-id", "", "Update Talos schematic ID")
	nodeUpdateCmd.Flags().Bool("maintenance", false, "Set maintenance mode")
	nodeUpdateCmd.Flags().Bool("no-maintenance", false, "Clear maintenance mode")
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
	Use:   "bootstrap <node>",
	Short: "Bootstrap cluster on a control plane node",
	Long: `Bootstrap the Kubernetes cluster by initializing etcd on a control plane node.

This should be run once after the first control plane node is configured.

Example:
  wild cluster bootstrap test-control-1`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		nodeName := args[0]

		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/cluster/bootstrap", inst), map[string]string{
			"node": nodeName,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Cluster bootstrap started on node: %s\n", nodeName)
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

		persist, _ := cmd.Flags().GetBool("persist")

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/cluster/kubeconfig", inst))
		if err != nil {
			return err
		}

		kubeconfigContent := resp.GetString("kubeconfig")

		// If --persist flag is set, save to instance directory
		if persist {
			dataDir := config.GetWildCLIDataDir()
			instanceDir := fmt.Sprintf("%s/instances/%s", dataDir, inst)

			// Create instance directory if it doesn't exist
			if err := os.MkdirAll(instanceDir, 0755); err != nil {
				return fmt.Errorf("failed to create instance directory: %w", err)
			}

			kubeconfigPath := fmt.Sprintf("%s/kubeconfig", instanceDir)
			if err := os.WriteFile(kubeconfigPath, []byte(kubeconfigContent), 0600); err != nil {
				return fmt.Errorf("failed to write kubeconfig: %w", err)
			}

			fmt.Printf("Kubeconfig saved to %s\n", kubeconfigPath)
			return nil
		}

		// Default behavior: print to stdout
		fmt.Println(kubeconfigContent)
		return nil
	},
}

var clusterConfigGenerateCmd = &cobra.Command{
	Use:   "config generate",
	Short: "Generate cluster configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Post(fmt.Sprintf("/api/v1/instances/%s/cluster/config/generate", inst), nil)
		if err != nil {
			return err
		}

		fmt.Println("Cluster configuration generated successfully")
		if msg := resp.GetString("message"); msg != "" {
			fmt.Println(msg)
		}
		return nil
	},
}

var clusterTalosconfigCmd = &cobra.Command{
	Use:   "talosconfig",
	Short: "Get talosconfig",
	RunE: func(cmd *cobra.Command, args []string) error {
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		persist, _ := cmd.Flags().GetBool("persist")

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s/cluster/talosconfig", inst))
		if err != nil {
			return err
		}

		talosconfigContent := resp.GetString("talosconfig")

		// If --persist flag is set, save to instance directory
		if persist {
			dataDir := config.GetWildCLIDataDir()
			instanceDir := fmt.Sprintf("%s/instances/%s", dataDir, inst)

			// Create instance directory if it doesn't exist
			if err := os.MkdirAll(instanceDir, 0755); err != nil {
				return fmt.Errorf("failed to create instance directory: %w", err)
			}

			talosconfigPath := fmt.Sprintf("%s/talosconfig", instanceDir)
			if err := os.WriteFile(talosconfigPath, []byte(talosconfigContent), 0600); err != nil {
				return fmt.Errorf("failed to write talosconfig: %w", err)
			}

			fmt.Printf("Talosconfig saved to %s\n", talosconfigPath)
			return nil
		}

		// Default behavior: print to stdout
		fmt.Println(talosconfigContent)
		return nil
	},
}

func init() {
	clusterCmd.AddCommand(clusterBootstrapCmd)
	clusterCmd.AddCommand(clusterStatusCmd)
	clusterCmd.AddCommand(clusterHealthCmd)
	clusterCmd.AddCommand(clusterKubeconfigCmd)
	clusterCmd.AddCommand(clusterConfigGenerateCmd)
	clusterCmd.AddCommand(clusterTalosconfigCmd)

	// Add --persist flags to config commands
	clusterTalosconfigCmd.Flags().Bool("persist", false, "Save talosconfig to instance directory")
	clusterKubeconfigCmd.Flags().Bool("persist", false, "Save kubeconfig to instance directory")
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
		inst, err := getInstanceName()
		if err != nil {
			return err
		}

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/operations/%s?instance=%s", args[0], inst))
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
	// Get instance name
	inst, err := getInstanceName()
	if err != nil {
		return err
	}

	// Get base URL
	baseURL := daemonURL
	if baseURL == "" {
		baseURL = config.GetDaemonURL()
	}

	// Connect to SSE stream
	url := fmt.Sprintf("%s/api/v1/operations/%s/stream?instance=%s", baseURL, opID, inst)
	client := sse.NewClient(url)
	events := make(chan *sse.Event)

	err = client.SubscribeChan("messages", events)
	if err != nil {
		return fmt.Errorf("failed to subscribe to SSE: %w", err)
	}

	// Poll for completion in background
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			resp, err := apiClient.Get(fmt.Sprintf("/api/v1/operations/%s?instance=%s", opID, inst))
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
