package root

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"connectrpc.com/connect"
	credentialv1 "github.com/plantonhq/openmcf/apis/org/openmcf/app/credential/v1"
	credentialv1connect "github.com/plantonhq/openmcf/apis/org/openmcf/app/credential/v1/credentialv1connect"
	"github.com/spf13/cobra"
)

var CredentialListCmd = &cobra.Command{
	Use:   "credential:list",
	Short: "list all credentials from backend",
	Long:  "List all credentials from the backend service. Optionally filter by provider.",
	Run:   credentialListHandler,
}

func init() {
	CredentialListCmd.Flags().StringP("provider", "p", "", "filter credentials by provider (gcp, aws, azure)")
}

func credentialListHandler(cmd *cobra.Command, args []string) {
	// Get backend URL from configuration
	backendURL, err := GetBackendURL()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Get provider filter if provided
	providerStr, _ := cmd.Flags().GetString("provider")

	// Create Connect-RPC client
	client := credentialv1connect.NewCredentialQueryControllerClient(
		http.DefaultClient,
		backendURL,
	)

	// Prepare request
	req := &credentialv1.ListCredentialsRequest{}
	if providerStr != "" {
		var provider credentialv1.Credential_CredentialProvider
		switch providerStr {
		case "gcp":
			provider = credentialv1.Credential_GCP
		case "aws":
			provider = credentialv1.Credential_AWS
		case "azure":
			provider = credentialv1.Credential_AZURE
		default:
			fmt.Printf("Error: Invalid provider '%s'. Valid values: gcp, aws, azure\n", providerStr)
			os.Exit(1)
		}
		req.Provider = provider
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Make the API call
	resp, err := client.List(ctx, connect.NewRequest(req))
	if err != nil {
		if connect.CodeOf(err) == connect.CodeUnavailable {
			fmt.Printf("Error: Cannot connect to backend service at %s. Please check:\n", backendURL)
			fmt.Printf("  1. The backend service is running\n")
			fmt.Printf("  2. The backend URL is correct\n")
			fmt.Printf("  3. Network connectivity\n")
			os.Exit(1)
		}
		fmt.Printf("Error fetching credentials: %v\n", err)
		os.Exit(1)
	}

	credentials := resp.Msg.Credentials

	if len(credentials) == 0 {
		if providerStr != "" {
			fmt.Printf("No credentials found with provider '%s'\n", providerStr)
		} else {
			fmt.Println("No credentials found")
		}
		return
	}

	// Display results in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Table header
	fmt.Fprintln(w, "ID\tNAME\tPROVIDER\tCREATED")

	// Table rows
	for _, cred := range credentials {
		providerName := "UNKNOWN"
		switch cred.Provider {
		case credentialv1.Credential_GCP:
			providerName = "GCP"
		case credentialv1.Credential_AWS:
			providerName = "AWS"
		case credentialv1.Credential_AZURE:
			providerName = "Azure"
		}

		createdAt := "N/A"
		if cred.CreatedAt != nil {
			createdAt = cred.CreatedAt.AsTime().Format("2006-01-02 15:04:05")
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			cred.Id,
			cred.Name,
			providerName,
			createdAt,
		)
	}

	// Summary
	fmt.Printf("\nTotal: %d credential(s)", len(credentials))
	if providerStr != "" {
		fmt.Printf(" (filtered by provider: %s)", providerStr)
	}
	fmt.Println()
}
