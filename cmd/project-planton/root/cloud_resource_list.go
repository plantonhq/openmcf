package root

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"connectrpc.com/connect"
	cloudresourcev1 "github.com/plantonhq/openmcf/apis/org/openmcf/app/cloudresource/v1"
	cloudresourcev1connect "github.com/plantonhq/openmcf/apis/org/openmcf/app/cloudresource/v1/cloudresourcev1connect"
	commonsv1 "github.com/plantonhq/openmcf/apis/org/openmcf/app/commons"
	"github.com/spf13/cobra"
)

var CloudResourceListCmd = &cobra.Command{
	Use:   "cloud-resource:list",
	Short: "list all cloud resources from backend",
	Long:  "List all cloud resources from the backend service. Optionally filter by kind.",
	Run:   cloudResourceListHandler,
}

func init() {
	CloudResourceListCmd.Flags().StringP("kind", "k", "", "filter cloud resources by kind")
	CloudResourceListCmd.Flags().Int32P("page", "p", 0, "page number (0-indexed, default: 0)")
	CloudResourceListCmd.Flags().Int32P("size", "s", 20, "page size (default: 20)")
}

func cloudResourceListHandler(cmd *cobra.Command, args []string) {
	// Get backend URL from configuration
	backendURL, err := GetBackendURL()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Get kind filter if provided
	kind, _ := cmd.Flags().GetString("kind")

	// Get pagination parameters (defaults: page=0, size=20)
	pageNum, _ := cmd.Flags().GetInt32("page")
	pageSize, _ := cmd.Flags().GetInt32("size")

	// Create Connect-RPC client
	client := cloudresourcev1connect.NewCloudResourceQueryControllerClient(
		http.DefaultClient,
		backendURL,
	)

	// Prepare request
	req := &cloudresourcev1.ListCloudResourcesRequest{}
	if kind != "" {
		req.Kind = kind
	}

	// Always include pagination (defaults applied in service if not provided)
	req.PageInfo = &commonsv1.PageInfo{
		Num:  pageNum,
		Size: pageSize,
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
		fmt.Printf("Error fetching cloud resources: %v\n", err)
		os.Exit(1)
	}

	resources := resp.Msg.Resources

	if len(resources) == 0 {
		if kind != "" {
			fmt.Printf("No cloud resources found with kind '%s'\n", kind)
		} else {
			fmt.Println("No cloud resources found")
		}
		return
	}

	// Display results in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Table header
	fmt.Fprintln(w, "ID\tNAME\tKIND\tCREATED")

	// Table rows
	for _, resource := range resources {
		createdAt := "N/A"
		if resource.CreatedAt != nil {
			createdAt = resource.CreatedAt.AsTime().Format("2006-01-02 15:04:05")
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			resource.Id,
			resource.Name,
			resource.Kind,
			createdAt,
		)
	}

	// Summary
	fmt.Printf("\nTotal: %d cloud resource(s)", len(resources))
	if kind != "" {
		fmt.Printf(" (filtered by kind: %s)", kind)
	}
	if resp.Msg.TotalPages > 0 {
		fmt.Printf(" | Page %d of %d", pageNum+1, resp.Msg.TotalPages)
	} else {
		fmt.Printf(" | Page %d of 1", pageNum+1)
	}
	fmt.Println()
}
