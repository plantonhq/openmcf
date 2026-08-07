package module

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Output keys for stack outputs. They match the field names in
// KubernetesConfigMapStackOutputs so downstream resources can compose on them.
const (
	OutputConfigMapName = "configmap_name"
	OutputNamespace     = "namespace"
)

// exportOutputs exports all stack outputs
func exportOutputs(ctx *pulumi.Context, locals *Locals) error {
	ctx.Export(OutputConfigMapName, pulumi.String(locals.ConfigMapName))
	ctx.Export(OutputNamespace, pulumi.String(locals.Namespace))

	return nil
}
