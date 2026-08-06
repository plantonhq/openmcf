package module

import (
	"github.com/pkg/errors"
	gcpserviceconnectionpolicyv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpserviceconnectionpolicy/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpserviceconnectionpolicyv1alpha1.GcpServiceConnectionPolicyStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := serviceConnectionPolicy(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create service connection policy")
	}

	return nil
}
