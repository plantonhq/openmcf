package webapp

import (
	"github.com/spf13/cobra"
)

// WebAppCmd is the main command for web app management
var WebAppCmd = &cobra.Command{
	Use:   "webapp",
	Short: "manage the OpenMCF web app",
	Long: `Manage the OpenMCF web app - a unified Docker container that provides
a web interface for managing cloud resources and deployments.

The web app includes:
  - MongoDB for data persistence
  - Backend API (port 50051)
  - Frontend web UI (port 3000)

Quick start:
  1. planton webapp init      # Initialize and pull Docker image
  2. planton webapp start     # Start the web app
  3. Open http://localhost:3000 in your browser`,
}

func init() {
	// Add all subcommands
	WebAppCmd.AddCommand(InitCmd)
	WebAppCmd.AddCommand(StartCmd)
	WebAppCmd.AddCommand(StopCmd)
	WebAppCmd.AddCommand(StatusCmd)
	WebAppCmd.AddCommand(LogsCmd)
	WebAppCmd.AddCommand(RestartCmd)
	WebAppCmd.AddCommand(UninstallCmd)
}
