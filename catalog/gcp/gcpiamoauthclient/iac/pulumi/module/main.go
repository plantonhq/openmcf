package module

import (
	"github.com/pkg/errors"
	gcpiamoauthclientv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpiamoauthclient/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpiamoauthclientv1alpha1.GcpIamOauthClientStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := oauthClient(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create oauth client")
	}

	return nil
}
