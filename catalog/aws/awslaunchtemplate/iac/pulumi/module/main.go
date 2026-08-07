package module

import (
	"github.com/pkg/errors"
	awslaunchtemplatev1alpha1 "github.com/plantonhq/planton/catalog/aws/awslaunchtemplate/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *awslaunchtemplatev1alpha1.AwsLaunchTemplateStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web identity,
	// or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsLaunchTemplate.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := launchTemplate(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "failed to create launch template")
	}

	return nil
}
