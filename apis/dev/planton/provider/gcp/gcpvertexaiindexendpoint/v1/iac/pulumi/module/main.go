package module

import (
	"github.com/pkg/errors"
	gcpvertexaiindexendpointv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpvertexaiindexendpoint/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpvertexaiindexendpointv1.GcpVertexAiIndexEndpointStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := indexEndpoint(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create vertex ai index endpoint")
	}

	return nil
}
