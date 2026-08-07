package module

import (
	"github.com/pkg/errors"
	gcpcloudrunv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpcloudrun/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the Pulumi program entry-point for the GcpCloudRun component.
func Resources(ctx *pulumi.Context, stackInput *gcpcloudrunv1alpha1.GcpCloudRunStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Set up the GCP provider from the supplied credential.
	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	// Create the Cloud Run service (and its public-invoker grant when the
	// spec asks for one).
	createdService, err := service(ctx, locals, gcpProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud-run service")
	}

	ctx.Export(OpUrl, createdService.Uri)
	ctx.Export(OpServiceName, createdService.Name)
	ctx.Export(OpRevision, createdService.LatestReadyRevision)
	ctx.Export(OpLocation, createdService.Location)
	ctx.Export(OpUid, createdService.Uid)
	ctx.Export(OpUrls, createdService.Urls)

	return nil
}
