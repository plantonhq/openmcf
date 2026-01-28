package webapp

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	logsFollow  bool
	logsTail    string
	logsService string
)

var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "view logs from the OpenMCF web app",
	Long:  `View logs from the web app container. Use --follow to stream logs in real-time.`,
	Run:   logsHandler,
}

func init() {
	LogsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "follow log output (stream in real-time)")
	LogsCmd.Flags().StringVarP(&logsTail, "tail", "n", "100", "number of lines to show from the end of the logs")
	LogsCmd.Flags().StringVar(&logsService, "service", "", "filter logs by service (mongodb, backend, frontend)")
}

func logsHandler(cmd *cobra.Command, args []string) {
	// Check if container exists
	if !containerExists() {
		fmt.Printf("❌ Container '%s' not found.\n", ContainerName)
		fmt.Println("   Please run: planton webapp init")
		os.Exit(1)
	}

	// Build docker logs command
	dockerArgs := []string{"logs"}

	if logsFollow {
		dockerArgs = append(dockerArgs, "-f")
	}

	dockerArgs = append(dockerArgs, "--tail", logsTail)
	dockerArgs = append(dockerArgs, ContainerName)

	// Create command
	logsCmd := exec.Command("docker", dockerArgs...)
	logsCmd.Stdout = os.Stdout
	logsCmd.Stderr = os.Stderr

	// Handle Ctrl+C gracefully when following logs
	if logsFollow {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-sigChan
			if logsCmd.Process != nil {
				logsCmd.Process.Kill()
			}
			fmt.Println("\n✅ Stopped following logs")
			os.Exit(0)
		}()
	}

	// Print header if filtering by service
	if logsService != "" {
		fmt.Printf("📋 Showing logs for service: %s\n", logsService)
		fmt.Println("   Note: Service filtering is approximate (shows all container logs)")
		fmt.Println()
	}

	// Run the command
	if err := logsCmd.Run(); err != nil {
		if !logsFollow {
			fmt.Printf("❌ Failed to get logs: %v\n", err)
			os.Exit(1)
		}
	}
}
