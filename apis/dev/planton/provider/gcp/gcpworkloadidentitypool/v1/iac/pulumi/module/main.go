package module

import (
	"github.com/pkg/errors"
	gcpworkloadidentitypoolv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpworkloadidentitypool/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpworkloadidentitypoolv1.GcpWorkloadIdentityPoolStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := workloadIdentityPool(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create workload identity pool")
	}

	return nil
}
