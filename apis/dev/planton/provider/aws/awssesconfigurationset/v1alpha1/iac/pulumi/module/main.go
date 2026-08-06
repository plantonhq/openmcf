package module

import (
	"github.com/pkg/errors"
	awssesconfigurationsetv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awssesconfigurationset/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources creates the SES configuration set, its event destinations, and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awssesconfigurationsetv1alpha1.AwsSesConfigurationSetStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsSesConfigurationSet.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	result, err := configurationSet(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create ses configuration set")
	}

	ctx.Export(OpConfigurationSetArn, result.Arn)
	ctx.Export(OpConfigurationSetName, result.ConfigurationSetName)

	return nil
}
