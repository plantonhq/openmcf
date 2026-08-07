package module

import (
	"github.com/pkg/errors"
	hetznercloudnetworkv1alpha1 "github.com/plantonhq/planton/catalog/hetznercloud/hetznercloudnetwork/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/hetznercloud/pulumihcloudprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(
	ctx *pulumi.Context,
	stackInput *hetznercloudnetworkv1alpha1.HetznerCloudNetworkStackInput,
) error {
	locals := initializeLocals(ctx, stackInput)

	hcloudProvider, err := pulumihcloudprovider.Get(
		ctx,
		stackInput.ProviderConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to setup hetzner cloud provider")
	}

	if err := network(ctx, locals, hcloudProvider); err != nil {
		return errors.Wrap(err, "failed to create network")
	}

	return nil
}
