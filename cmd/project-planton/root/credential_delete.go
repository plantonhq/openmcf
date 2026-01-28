package root

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	credentialv1 "github.com/plantonhq/openmcf/apis/org/openmcf/app/credential/v1"
	credentialv1connect "github.com/plantonhq/openmcf/apis/org/openmcf/app/credential/v1/credentialv1connect"
	"github.com/spf13/cobra"
)

var CredentialDeleteCmd = &cobra.Command{
	Use:   "credential:delete",
	Short: "delete a credential by ID",
	Long:  "Delete a credential by providing its unique ID. This action is irreversible.",
	Run:   credentialDeleteHandler,
}

func init() {
	CredentialDeleteCmd.Flags().StringP("id", "i", "", "unique identifier of the credential (required)")
	CredentialDeleteCmd.MarkFlagRequired("id")
}

func credentialDeleteHandler(cmd *cobra.Command, args []string) {
	// Get credential ID
	id, _ := cmd.Flags().GetString("id")
	if id == "" {
		fmt.Println("Error: --id flag is required. Provide the credential ID")
		fmt.Println("Usage: openmcf credential:delete --id=<credential-id>")
		os.Exit(1)
	}

	// Get backend URL from configuration
	backendURL, err := GetBackendURL()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create Connect-RPC client
	client := credentialv1connect.NewCredentialCommandControllerClient(
		http.DefaultClient,
		backendURL,
	)

	// Prepare request
	req := &credentialv1.DeleteCredentialRequest{
		Id: id,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Make the API call
	resp, err := client.Delete(ctx, connect.NewRequest(req))
	if err != nil {
		if connect.CodeOf(err) == connect.CodeUnavailable {
			fmt.Printf("Error: Cannot connect to backend service at %s. Please check:\n", backendURL)
			fmt.Printf("  1. The backend service is running\n")
			fmt.Printf("  2. The backend URL is correct\n")
			fmt.Printf("  3. Network connectivity\n")
			os.Exit(1)
		}
		if connect.CodeOf(err) == connect.CodeNotFound {
			fmt.Printf("Error: Credential with ID '%s' not found\n", id)
			os.Exit(1)
		}
		fmt.Printf("Error deleting credential: %v\n", err)
		os.Exit(1)
	}

	// Display success message
	fmt.Println("✅ " + resp.Msg.Message)
}
