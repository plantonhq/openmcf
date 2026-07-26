package module

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Output name constants — one per KubernetesTektonOperatorStackOutputs
// field.
const (
	OpNamespace = "namespace"
)

// exportOutputs publishes the release manifest's fixed handles (the
// final apply group is passed so the export carries its dependency).
func exportOutputs(ctx *pulumi.Context, _ *Locals, _ pulumi.Resource) error {
	ctx.Export(OpNamespace, pulumi.String(vars.Namespace))
	return nil
}
