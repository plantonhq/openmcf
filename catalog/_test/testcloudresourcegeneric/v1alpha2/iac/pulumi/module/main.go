package module

import (
	testcloudresourcegenericv1alpha2 "github.com/plantonhq/planton/catalog/_test/testcloudresourcegeneric/v1alpha2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the module entry point. The test kind provisions nothing
// outside the engine's own state: outputs derive deterministically from the
// resolved inputs, which is exactly what makes the kind useful -- the full
// manifest -> stack-input -> module -> outputs pipeline runs with zero
// credentials, zero network, and zero cost, and any input change is visible
// as an output change.
func Resources(ctx *pulumi.Context, stackInput *testcloudresourcegenericv1alpha2.TestCloudResourceGenericStackInput) error {
	target := stackInput.Target

	name := target.Metadata.Name
	steps := target.Spec.GetSteps()

	tags := make(pulumi.StringArray, 0, len(steps))
	for _, s := range steps {
		tags = append(tags, pulumi.String(s.GetCommand()))
	}

	// Output names and meanings mirror stack_outputs.proto exactly; the id
	// composes the kind's registry id_prefix, so output shapes look exactly
	// like a real kind's.
	ctx.Export("id", pulumi.String("tcrg-"+name))
	ctx.Export("name", pulumi.String(name))
	ctx.Export("url", pulumi.String("test://"+name))
	ctx.Export("tags", tags)
	return nil
}
