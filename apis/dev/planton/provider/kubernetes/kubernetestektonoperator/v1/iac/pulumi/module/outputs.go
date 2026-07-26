package module

import (
	pulumiyaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Output name constants — one per KubernetesTektonOperatorStackOutputs
// field.
const (
	OpNamespace = "namespace"
)

// exportOutputs publishes the release manifest's fixed handles (the
// manifest resource is passed so the export carries its dependency).
func exportOutputs(ctx *pulumi.Context, _ *Locals, _ *pulumiyaml.ConfigFile) error {
	ctx.Export(OpNamespace, pulumi.String(vars.Namespace))
	return nil
}
