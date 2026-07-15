package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apprunner"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// autoScalingConfiguration registers the App Runner auto scaling
// configuration version. AWS versions this resource: every settable value is
// create-time immutable, so any change destroys this revision and registers
// the NEXT revision under the same configuration name -- referencing services
// then pick up the new revision-carrying ARN through the resource graph.
func autoScalingConfiguration(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.AwsAppRunnerAutoScalingConfiguration.Spec

	args := &apprunner.AutoScalingConfigurationVersionArgs{
		// The cloud name is metadata.name -- the same basis the Terraform
		// module uses, so both engines register the same physical identity.
		AutoScalingConfigurationName: pulumi.String(locals.AwsAppRunnerAutoScalingConfiguration.Metadata.Name),
		Tags:                         pulumi.ToStringMap(locals.AwsTags),
	}

	// Unset optionals fall through to AWS's own defaults (100 concurrency /
	// 25 max / 1 min) -- the platform normally materializes the spec
	// defaults before the module runs, so these are explicit in practice.
	if spec.MaxConcurrency != nil {
		args.MaxConcurrency = pulumi.Int(int(*spec.MaxConcurrency))
	}
	if spec.MaxSize != nil {
		args.MaxSize = pulumi.Int(int(*spec.MaxSize))
	}
	if spec.MinSize != nil {
		args.MinSize = pulumi.Int(int(*spec.MinSize))
	}

	createdConfiguration, err := apprunner.NewAutoScalingConfigurationVersion(ctx,
		locals.AwsAppRunnerAutoScalingConfiguration.Metadata.Name, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create auto scaling configuration version")
	}

	// Export outputs matching AwsAppRunnerAutoScalingConfigurationStackOutputs.
	ctx.Export(OpConfigurationArn, createdConfiguration.Arn)
	ctx.Export(OpConfigurationRevision, createdConfiguration.AutoScalingConfigurationRevision)
	ctx.Export(OpLatest, createdConfiguration.Latest)

	return nil
}
