package module

import (
	testcloudresourcegenericv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/_test/testcloudresourcegeneric/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the module entry point. The test kind provisions nothing
// outside the engine's own state: outputs derive deterministically from the
// resolved inputs, which is exactly what makes the kind useful -- the full
// manifest -> stack-input -> module -> outputs pipeline runs with zero
// credentials, zero network, and zero cost, and any input change is visible
// as an output change.
func Resources(ctx *pulumi.Context, stackInput *testcloudresourcegenericv1alpha1.TestCloudResourceGenericStackInput) error {
	target := stackInput.Target

	name := target.Metadata.Name
	commands := target.Spec.GetCommands()

	tags := make(pulumi.StringArray, 0, len(commands))
	for _, c := range commands {
		tags = append(tags, pulumi.String(c))
	}

	// Output names and meanings mirror stack_outputs.proto exactly; the id
	// composes the kind's registry id_prefix, so output shapes look exactly
	// like a real kind's.
	ctx.Export("id", pulumi.String("tcrg-"+name))
	ctx.Export("name", pulumi.String(name))
	ctx.Export("endpoint", pulumi.String("test://"+name))
	ctx.Export("tags", tags)
	return nil
}
