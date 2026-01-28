package root

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"connectrpc.com/connect"
	stackupdatev1 "github.com/plantonhq/openmcf/apis/org/openmcf/app/stackupdate/v1"
	stackupdatev1connect "github.com/plantonhq/openmcf/apis/org/openmcf/app/stackupdate/v1/stackupdatev1connect"
	"github.com/spf13/cobra"
)

var StackUpdateStreamOutputCmd = &cobra.Command{
	Use:   "stack-update:stream-output",
	Short: "stream real-time output from a stack-update",
	Long:  "Stream real-time deployment logs from a stack-update. Shows stdout and stderr output as it's generated during deployment.",
	Run:   stackUpdateStreamOutputHandler,
}

func init() {
	StackUpdateStreamOutputCmd.Flags().StringP("id", "i", "", "unique identifier of the stack-update (required)")
	StackUpdateStreamOutputCmd.Flags().Int32P("last-sequence", "s", 0, "last sequence number received (for resuming stream from a specific point)")
	StackUpdateStreamOutputCmd.MarkFlagRequired("id")
}

func stackUpdateStreamOutputHandler(cmd *cobra.Command, args []string) {
	// Get stack-update ID
	stackUpdateID, _ := cmd.Flags().GetString("id")
	if stackUpdateID == "" {
		fmt.Println("Error: --id flag is required. Provide the stack-update ID")
		fmt.Println("Usage: openmcf stack-update:stream-output --id=<stack-update-id>")
		os.Exit(1)
	}

	// Get last sequence number (optional)
	lastSequenceNum, _ := cmd.Flags().GetInt32("last-sequence")

	// Get backend URL from configuration
	backendURL, err := GetBackendURL()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create Connect-RPC client
	client := stackupdatev1connect.NewStackUpdateQueryControllerClient(
		http.DefaultClient,
		backendURL,
	)

	// Prepare request
	req := &stackupdatev1.StreamStackUpdateOutputRequest{
		JobId: stackUpdateID,
	}
	if lastSequenceNum > 0 {
		req.LastSequenceNum = lastSequenceNum
	}

	// Create context with cancellation support
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\n⚠️  Interrupt received, stopping stream...")
		cancel()
	}()

	// Start streaming
	fmt.Printf("🚀 Streaming output for stack-update: %s\n", stackUpdateID)
	if lastSequenceNum > 0 {
		fmt.Printf("   Resuming from sequence: %d\n", lastSequenceNum)
	}
	fmt.Println()

	stream, err := client.StreamStackUpdateOutput(ctx, connect.NewRequest(req))
	if err != nil {
		if connect.CodeOf(err) == connect.CodeUnavailable {
			fmt.Printf("❌ Error: Cannot connect to backend service at %s. Please check:\n", backendURL)
			fmt.Printf("  1. The backend service is running\n")
			fmt.Printf("  2. The backend URL is correct\n")
			fmt.Printf("  3. Network connectivity\n")
			os.Exit(1)
		}
		if connect.CodeOf(err) == connect.CodeNotFound {
			fmt.Printf("❌ Error: Stack update with ID '%s' not found\n", stackUpdateID)
			os.Exit(1)
		}
		fmt.Printf("❌ Error starting stream: %v\n", err)
		os.Exit(1)
	}

	// Track statistics
	messageCount := 0
	lastSequence := int32(0)

	// Iterate over the stream
	for stream.Receive() {
		response := stream.Msg()
		messageCount++
		lastSequence = response.SequenceNum

		// Format timestamp
		timestamp := "N/A"
		if response.Timestamp != nil {
			timestamp = response.Timestamp.AsTime().Format("15:04:05.000")
		}

		// Determine stream type prefix and color
		streamPrefix := "[stdout]"
		if response.StreamType == "stderr" {
			streamPrefix = "[stderr]"
		}

		// Display the message
		fmt.Printf("[%s] %s [Seq: %d] %s\n",
			timestamp,
			streamPrefix,
			response.SequenceNum,
			response.Content,
		)

		// Check if stream is completed
		if response.Status == "completed" {
			fmt.Printf("\n✅ Stream completed successfully\n")
			break
		}

		// Check if stream failed
		if response.Status == "failed" {
			fmt.Printf("\n❌ Stream failed\n")
			break
		}

		// Check if context was cancelled
		if ctx.Err() != nil {
			fmt.Printf("\n⚠️  Stream cancelled\n")
			break
		}
	}

	// Check for stream errors
	if err := stream.Err(); err != nil {
		if err == io.EOF {
			// EOF is normal when stream completes
			fmt.Printf("\n📊 Total messages received: %d (last sequence: %d)\n", messageCount, lastSequence)
			return
		}
		if connect.CodeOf(err) == connect.CodeCanceled {
			fmt.Printf("\n⚠️  Stream was cancelled\n")
			fmt.Printf("📊 Total messages received: %d (last sequence: %d)\n", messageCount, lastSequence)
			return
		}
		fmt.Printf("\n❌ Stream error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n📊 Total messages received: %d (last sequence: %d)\n", messageCount, lastSequence)
}
