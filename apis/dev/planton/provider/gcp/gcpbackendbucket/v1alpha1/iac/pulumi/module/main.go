package module

import (
	"github.com/pkg/errors"
	gcpbackendbucketv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpbackendbucket/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpbackendbucketv1alpha1.GcpBackendBucketStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := backendBucket(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create backend bucket")
	}

	return nil
}
