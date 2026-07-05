package module

import (
	"github.com/pkg/errors"
	awscloudfrontv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awscloudfront/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *awscloudfrontv1.AwsCloudFrontStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsCloudFront.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	dist, err := createDistribution(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create cloudfront distribution")
	}

	// Names match the Terraform module's outputs.tf key-for-key so both
	// engines present one contract to consumers.
	ctx.Export(OpDistributionId, dist.ID())
	ctx.Export(OpDistributionArn, dist.Arn)
	ctx.Export(OpDomainName, dist.DomainName)
	ctx.Export(OpHostedZoneId, dist.HostedZoneId)
	ctx.Export(OpStatus, dist.Status)

	return nil
}
