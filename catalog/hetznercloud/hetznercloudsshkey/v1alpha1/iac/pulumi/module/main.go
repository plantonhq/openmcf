package module

import (
	"github.com/pkg/errors"
	hetznercloudsshkeyv1alpha1 "github.com/plantonhq/planton/catalog/hetznercloud/hetznercloudsshkey/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/hetznercloud/pulumihcloudprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the module entry point.
func Resources(
	ctx *pulumi.Context,
	stackInput *hetznercloudsshkeyv1alpha1.HetznerCloudSshKeyStackInput,
) error {
	locals := initializeLocals(ctx, stackInput)

	hcloudProvider, err := pulumihcloudprovider.Get(
		ctx,
		stackInput.ProviderConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to setup hetzner cloud provider")
	}

	if err := sshKey(ctx, locals, hcloudProvider); err != nil {
		return errors.Wrap(err, "failed to create ssh key")
	}

	return nil
}
