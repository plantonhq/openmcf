package webapp

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "check the status of OpenMCF web app",
	Long:  `Display the current status of the web app container and services`,
	Run:   statusHandler,
}

func statusHandler(cmd *cobra.Command, args []string) {
	fmt.Println("========================================")
	fmt.Println("📊 OpenMCF Web App Status")
	fmt.Println("========================================")
	fmt.Println()

	// Check if container exists
	if !containerExists() {
		fmt.Printf("❌ Container '%s' not found.\n", ContainerName)
		fmt.Println()
		fmt.Println("To initialize the web app, run:")
		fmt.Println("  planton webapp init")
		fmt.Println()
		os.Exit(1)
	}

	// Check if running
	running := isContainerRunning()

	fmt.Println("Container Information:")
	fmt.Printf("  Name:       %s\n", ContainerName)
	fmt.Printf("  Status:     %s\n", getContainerStatus())
	fmt.Printf("  Image:      %s:%s\n", DockerImageName, DockerImageTag)
	fmt.Println()

	if running {
		fmt.Println("Service Status:")
		printServiceStatus("MongoDB", "27017")
		printServiceStatus("Backend", BackendPort)
		printServiceStatus("Frontend", FrontendPort)
		fmt.Println()

		fmt.Println("Access URLs:")
		fmt.Printf("  🌐 Frontend:  http://localhost:%s\n", FrontendPort)
		fmt.Printf("  🔌 Backend:   http://localhost:%s\n", BackendPort)
		fmt.Println()

		fmt.Println("Data Volumes:")
		fmt.Printf("  MongoDB:     %s\n", MongoDBVolume)
		fmt.Printf("  Pulumi:      %s\n", PulumiVolume)
		fmt.Printf("  Go Cache:    %s\n", GoCacheVolume)
	} else {
		fmt.Println("The web app is not running.")
		fmt.Println()
		fmt.Println("To start the web app, run:")
		fmt.Println("  planton webapp start")
	}
	fmt.Println()
}

func getContainerStatus() string {
	cmd := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", ContainerName)
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	status := strings.TrimSpace(string(output))

	// Add emoji based on status
	switch status {
	case "running":
		return "🟢 running"
	case "exited":
		return "🔴 stopped"
	case "paused":
		return "🟡 paused"
	case "restarting":
		return "🔄 restarting"
	default:
		return status
	}
}

func printServiceStatus(serviceName, port string) {
	// Try to check if port is listening inside container
	cmd := exec.Command("docker", "exec", ContainerName, "sh", "-c",
		fmt.Sprintf("netstat -tuln 2>/dev/null | grep ':%s ' || ss -tuln 2>/dev/null | grep ':%s ' || echo 'unknown'", port, port))
	output, err := cmd.Output()

	status := "🟢 running"
	if err != nil || strings.Contains(string(output), "unknown") {
		status = "⚠️  status unknown"
	}

	fmt.Printf("  %-12s %s (port %s)\n", serviceName+":", status, port)
}
