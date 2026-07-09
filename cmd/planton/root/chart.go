package root

import (
	"github.com/plantonhq/planton/cmd/planton/root/chart"
	"github.com/spf13/cobra"
)

// Chart groups chart authoring tools. It is part of the standalone binary's
// command set (like e2e), NOT the embedded engine set in RegisterCommands:
// the Planton Platform CLI mounts its own `chart` command tree (build/publish
// against the control plane), and the two must not collide at the embedding
// seam. Offline validation needs no control plane, so it lives here.
var Chart = &cobra.Command{
	Use:   "chart",
	Short: "InfraChart authoring tools",
}

func init() {
	Chart.AddCommand(
		chart.Validate,
	)
}
