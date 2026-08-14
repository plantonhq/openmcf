package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// customModel starts the Bedrock model-customization job that produces
// the custom model, and exports outputs. Every argument is
// create-time-immutable -- a customization job cannot be altered once
// started, so any change replaces the job and model.
func customModel(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &bedrock.CustomModelArgs{
		CustomModelName: pulumi.String(locals.CustomModelName),
		JobName:         pulumi.String(locals.JobName),
		// The base foundation model being customized (foundation-model ARN).
		BaseModelIdentifier: pulumi.String(spec.BaseModelArn),
		// Per-base-model training knobs (epochCount, batchSize, learningRate, ...).
		Hyperparameters: pulumi.ToStringMap(spec.Hyperparameters),
		// The role Bedrock assumes to read training data and write
		// outputs. Must trust bedrock.amazonaws.com.
		RoleArn: pulumi.String(spec.RoleArn.GetValue()),
		TrainingDataConfig: &bedrock.CustomModelTrainingDataConfigArgs{
			S3Uri: pulumi.String(spec.TrainingDataS3Uri),
		},
		OutputDataConfig: &bedrock.CustomModelOutputDataConfigArgs{
			S3Uri: pulumi.String(spec.OutputDataS3Uri),
		},
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// Sent only when set: the provider attribute is Optional+Computed and
	// AWS defaults to FINE_TUNING.
	if spec.CustomizationType != "" {
		args.CustomizationType = pulumi.String(spec.CustomizationType)
	}

	// Customer-managed key for the resulting model when referenced;
	// Bedrock-managed key otherwise.
	if spec.CustomModelKmsKeyArn.GetValue() != "" {
		args.CustomModelKmsKeyId = pulumi.String(spec.CustomModelKmsKeyArn.GetValue())
	}

	// Up to 10 validation datasets -- Bedrock reports per-dataset metrics
	// on the finished job.
	if len(spec.ValidationDataS3Uris) > 0 {
		var validators bedrock.CustomModelValidationDataConfigValidatorArray
		for _, uri := range spec.ValidationDataS3Uris {
			validators = append(validators, &bedrock.CustomModelValidationDataConfigValidatorArgs{
				S3Uri: pulumi.String(uri),
			})
		}
		args.ValidationDataConfig = &bedrock.CustomModelValidationDataConfigArgs{
			Validators: validators,
		}
	}

	// Run the job's data access inside the caller's VPC (both members
	// required together -- CEL-enforced).
	if vpc := spec.VpcConfig; vpc != nil {
		subnetIds := make([]string, 0, len(vpc.SubnetIds))
		for _, s := range vpc.SubnetIds {
			subnetIds = append(subnetIds, s.GetValue())
		}
		securityGroupIds := make([]string, 0, len(vpc.SecurityGroupIds))
		for _, s := range vpc.SecurityGroupIds {
			securityGroupIds = append(securityGroupIds, s.GetValue())
		}
		args.VpcConfig = &bedrock.CustomModelVpcConfigArgs{
			SubnetIds:        pulumi.ToStringArray(subnetIds),
			SecurityGroupIds: pulumi.ToStringArray(securityGroupIds),
		}
	}

	createdModel, err := bedrock.NewCustomModel(ctx, locals.CustomModelName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create custom model customization job")
	}

	ctx.Export(OpCustomModelArn, createdModel.CustomModelArn)
	ctx.Export(OpCustomModelName, createdModel.CustomModelName)
	ctx.Export(OpJobArn, createdModel.JobArn)
	ctx.Export(OpJobStatus, createdModel.JobStatus)

	return nil
}
