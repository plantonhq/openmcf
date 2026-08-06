package module

import (
	"github.com/pkg/errors"
	gcptargethttpproxyv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcptargethttpproxy/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcptargethttpproxyv1alpha1.GcpTargetHttpProxyStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := targetHttpProxy(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create target http proxy")
	}

	return nil
}
