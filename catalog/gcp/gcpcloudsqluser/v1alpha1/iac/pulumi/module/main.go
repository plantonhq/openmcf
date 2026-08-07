package module

import (
	"github.com/pkg/errors"
	gcpcloudsqluserv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpcloudsqluser/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the Pulumi program entry-point for the GcpCloudSqlUser
// component.
func Resources(ctx *pulumi.Context, stackInput *gcpcloudsqluserv1alpha1.GcpCloudSqlUserStackInput) error {
	locals := initializeLocals(stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := user(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create cloud sql user")
	}

	return nil
}
