package module

import (
	"github.com/pkg/errors"
	gcpcloudcomposeruserworkloadssecretv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcloudcomposeruserworkloadssecret/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpcloudcomposeruserworkloadssecretv1alpha1.GcpCloudComposerUserWorkloadsSecretStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := userWorkloadsSecret(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create composer user workloads secret")
	}

	return nil
}
