package module

import (
	"github.com/pkg/errors"
	ocipublicipv1alpha1 "github.com/plantonhq/planton/catalog/oci/ocipublicip/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/oci/pulumiociprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *ocipublicipv1alpha1.OciPublicIpStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	ociProvider, err := pulumiociprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup oci provider")
	}

	if err := publicIp(ctx, locals, ociProvider); err != nil {
		return errors.Wrap(err, "failed to create public ip")
	}

	return nil
}
