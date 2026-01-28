package webapp

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart the OpenMCF web app",
	Long:  `Restart the OpenMCF web app container and all services`,
	Run:   restartHandler,
}

func restartHandler(cmd *cobra.Command, args []string) {
	fmt.Println("========================================")
	fmt.Println("🔄 Restarting OpenMCF Web App")
	fmt.Println("========================================")
	fmt.Println()

	// Check if container exists
	if !containerExists() {
		fmt.Printf("❌ Container '%s' not found.\n", ContainerName)
		fmt.Println("   Please run: planton webapp init")
		os.Exit(1)
	}

	// Restart the container
	fmt.Println("🔄 Restarting container...")
	restartCmd := exec.Command("docker", "restart", ContainerName)
	if err := restartCmd.Run(); err != nil {
		fmt.Printf("❌ Failed to restart container: %v\n", err)
		os.Exit(1)
	}

	// Wait for services to be healthy
	fmt.Println("⏳ Waiting for services to start (this may take 30-60 seconds)...")
	if err := waitForHealthy(60); err != nil {
		fmt.Printf("⚠️  Warning: %v\n", err)
		fmt.Println("   Check logs with: planton webapp logs")
	} else {
		fmt.Println("✅ All services are healthy")
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("✨ Web App Restarted Successfully!")
	fmt.Println("========================================")
	fmt.Println()
	printAccessInfo()
}
