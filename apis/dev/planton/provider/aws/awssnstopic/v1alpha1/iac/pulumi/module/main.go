package module

import (
	"github.com/pkg/errors"
	awssnstopicv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awssnstopic/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates SNS topic creation and exports outputs. Subscriptions
// are first-class AwsSnsSubscription resources that reference this topic's
// exported topic_arn.
func Resources(ctx *pulumi.Context, stackInput *awssnstopicv1alpha1.AwsSnsTopicStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if _, err := topic(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "sns topic")
	}

	return nil
}
