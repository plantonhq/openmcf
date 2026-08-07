package module

import (
	"github.com/pkg/errors"
	gcprouternatv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcprouternat/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the Pulumi program entry-point for the GcpRouterNat component.
func Resources(ctx *pulumi.Context, stackInput *gcprouternatv1alpha1.GcpRouterNatStackInput) error {
	locals := initializeLocals(stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if _, err = routerNat(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create router nat resources")
	}

	return nil
}
