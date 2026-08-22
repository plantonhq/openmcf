package module

import (
	"github.com/pkg/errors"
	cloudflarelogpushjobv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarelogpushjob/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/cloudflare/pulumicloudflareprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the module entry point.
func Resources(
	ctx *pulumi.Context,
	stackInput *cloudflarelogpushjobv1alpha1.CloudflareLogpushJobStackInput,
) error {
	// 1. Prepare locals (metadata, credentials).
	locals := initializeLocals(ctx, stackInput)

	// 2. Create a Cloudflare provider from the supplied credential.
	cloudflareProvider, err := pulumicloudflareprovider.Get(
		ctx,
		stackInput.ProviderConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to setup cloudflare provider")
	}

	// 3. Create the logpush job (and the ownership challenge when asked).
	if err := logpushJob(ctx, locals, cloudflareProvider); err != nil {
		return errors.Wrap(err, "failed to create logpush job")
	}

	return nil
}
