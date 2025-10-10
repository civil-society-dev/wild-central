package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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

