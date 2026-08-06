package module

import (
	"github.com/pkg/errors"
	gcpserverlessvpcconnectorv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpserverlessvpcconnector/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the Pulumi program entry-point for the
// GcpServerlessVpcConnector component.
func Resources(ctx *pulumi.Context, stackInput *gcpserverlessvpcconnectorv1alpha1.GcpServerlessVpcConnectorStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	createdConnector, err := connector(ctx, locals, gcpProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create serverless vpc access connector")
	}

	ctx.Export(OpName, createdConnector.Name)
	ctx.Export(OpSelfLink, createdConnector.SelfLink)
	ctx.Export(OpState, createdConnector.State)
	// Emit the plain spec region, not the resource attribute: on released
	// 6.x provider lines regional attributes can surface as region
	// self-links rather than plain names, and the output proto documents a
	// plain region name.
	ctx.Export(OpRegion, pulumi.String(locals.GcpServerlessVpcConnector.Spec.Region))

	return nil
}
