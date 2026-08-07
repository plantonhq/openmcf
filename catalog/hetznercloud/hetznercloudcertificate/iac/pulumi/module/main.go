package module

import (
	"github.com/pkg/errors"
	hetznercloudcertificatev1alpha1 "github.com/plantonhq/planton/catalog/hetznercloud/hetznercloudcertificate/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/hetznercloud/pulumihcloudprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(
	ctx *pulumi.Context,
	stackInput *hetznercloudcertificatev1alpha1.HetznerCloudCertificateStackInput,
) error {
	locals := initializeLocals(ctx, stackInput)

	hcloudProvider, err := pulumihcloudprovider.Get(
		ctx,
		stackInput.ProviderConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to setup hetzner cloud provider")
	}

	if err := certificate(ctx, locals, hcloudProvider); err != nil {
		return errors.Wrap(err, "failed to create certificate")
	}

	return nil
}
