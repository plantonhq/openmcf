package module

import (
	"github.com/pkg/errors"
	gcpalloydbuserv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpalloydbuser/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpalloydbuserv1alpha1.GcpAlloydbUserStackInput) error {
	locals := initializeLocals(stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := user(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create alloydb user")
	}

	return nil
}
