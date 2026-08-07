package module

import (
	"github.com/pkg/errors"
	awslbtargetgroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awslbtargetgroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *awslbtargetgroupv1alpha1.AwsLbTargetGroupStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web identity,
	// or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsLbTargetGroup.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := targetGroup(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "failed to create target group")
	}

	return nil
}
