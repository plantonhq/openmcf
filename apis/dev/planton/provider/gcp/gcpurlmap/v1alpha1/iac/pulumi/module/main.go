package module

import (
	"github.com/pkg/errors"
	gcpurlmapv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpurlmap/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpurlmapv1alpha1.GcpUrlMapStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := urlMap(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create url map")
	}

	return nil
}
