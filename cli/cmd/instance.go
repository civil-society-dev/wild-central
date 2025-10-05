package cmd

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
)

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage instances",
	Long:  `Create, list, and manage Wild Cloud instances.`,
}

var instanceCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		resp, err := apiClient.Post("/api/v1/instances", map[string]string{
			"name": name,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Instance '%s' created successfully\n", name)
		if msg, ok := resp.Data["message"].(string); ok && msg != "" {
			fmt.Println(msg)
		}
		return nil
	},
}

var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.Get("/api/v1/instances")
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		if outputFormat == "yaml" {
			return printYAML(resp.Data)
		}

		instances := resp.GetArray("instances")
		if len(instances) == 0 {
			fmt.Println("No instances found")
			return nil
		}

		fmt.Println("INSTANCES:")
		for _, inst := range instances {
			if name, ok := inst.(string); ok {
				fmt.Printf("  - %s\n", name)
			}
		}
		return nil
	},
}

var instanceShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show instance details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		resp, err := apiClient.Get(fmt.Sprintf("/api/v1/instances/%s", name))
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(resp.Data)
		}

		if outputFormat == "yaml" {
			return printYAML(resp.Data)
		}

		// Pretty print instance info
		instanceName := resp.GetString("name")
		if instanceName != "" {
			fmt.Printf("Instance: %s\n", instanceName)
		}

		// Show key config values
		if config := resp.GetMap("config"); config != nil {
			// Check for cloud config (test-cloud structure)
			if cloud, ok := config["cloud"].(map[string]interface{}); ok {
				if domain, ok := cloud["domain"].(string); ok && domain != "" {
					fmt.Printf("Domain: %s\n", domain)
				}
				if baseDomain, ok := cloud["baseDomain"].(string); ok && baseDomain != "" {
					fmt.Printf("Base Domain: %s\n", baseDomain)
				}
			} else {
				// Check for direct config fields (test-cli structure)
				if domain, ok := config["domain"].(string); ok && domain != "" {
					fmt.Printf("Domain: %s\n", domain)
				}
				if baseDomain, ok := config["baseDomain"].(string); ok && baseDomain != "" {
					fmt.Printf("Base Domain: %s\n", baseDomain)
				}
			}

			// Show cluster info if available
			if cluster, ok := config["cluster"].(map[string]interface{}); ok {
				if clusterName, ok := cluster["name"].(string); ok && clusterName != "" {
					fmt.Printf("Cluster Name: %s\n", clusterName)
				}
			}
		}

		fmt.Println("\nUse -o json or -o yaml for full configuration")
		return nil
	},
}

var instanceDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete an instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Confirm deletion
		fmt.Printf("Are you sure you want to delete instance '%s'? (yes/no): ", name)
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Deletion cancelled")
			return nil
		}

		resp, err := apiClient.Delete(fmt.Sprintf("/api/v1/instances/%s", name))
		if err != nil {
			return err
		}

		fmt.Printf("Instance '%s' deleted successfully\n", name)
		if msg, ok := resp.Data["message"].(string); ok && msg != "" {
			fmt.Println(msg)
		}
		return nil
	},
}

var instanceCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current instance",
	Long:  `Display the instance that would be used by commands (from --instance flag or WILD_INSTANCE environment variable).`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check flag first (it takes precedence)
		inst := instanceName

		// If not set by flag, check environment
		if inst == "" {
			inst = os.Getenv("WILD_INSTANCE")
		}

		// If still not set, show helpful message
		if inst == "" {
			fmt.Println("No instance configured (use --instance flag or set WILD_INSTANCE environment variable)")
			return
		}

		fmt.Println(inst)
	},
}

func init() {
	instanceCmd.AddCommand(instanceCreateCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceShowCmd)
	instanceCmd.AddCommand(instanceDeleteCmd)
	instanceCmd.AddCommand(instanceCurrentCmd)
}
