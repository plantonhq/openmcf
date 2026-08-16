package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// pipeline creates the SageMaker pipeline and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - everything except the name updates in place; creating a pipeline
//     is free (executions bill);
//   - the definition comes from exactly one place (spec-validated) -
//     inline JSON or an S3 object;
//   - AWS's describe API returns only the RESOLVED definition, never
//     the S3 location - the location is config-only on import and
//     S3-object drift is invisible to refresh (taught on the spec
//     field).
func pipeline(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &sagemaker.PipelineArgs{
		// The component's name IS the pipeline name.
		PipelineName: pulumi.String(locals.PipelineName),
		// Required by the provider - defaults to the pipeline name.
		PipelineDisplayName: pulumi.String(locals.DisplayName),
		RoleArn:             pulumi.String(spec.RoleArn.GetValue()),
		Tags:                pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.PipelineDescription = pulumi.String(spec.Description)
	}

	// The inline definition arm.
	if spec.Definition != nil {
		definitionBytes, err := json.Marshal(spec.Definition.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal pipeline definition")
		}
		args.PipelineDefinition = pulumi.String(string(definitionBytes))
	}

	// The S3-location arm.
	if spec.DefinitionS3Location != nil {
		s3Location := &sagemaker.PipelinePipelineDefinitionS3LocationArgs{
			Bucket:    pulumi.String(spec.DefinitionS3Location.Bucket.GetValue()),
			ObjectKey: pulumi.String(spec.DefinitionS3Location.ObjectKey),
		}
		if spec.DefinitionS3Location.VersionId != "" {
			s3Location.VersionId = pulumi.String(spec.DefinitionS3Location.VersionId)
		}
		args.PipelineDefinitionS3Location = s3Location
	}

	if spec.ParallelismMaxSteps != nil {
		args.ParallelismConfiguration = &sagemaker.PipelineParallelismConfigurationArgs{
			MaxParallelExecutionSteps: pulumi.Int(int(*spec.ParallelismMaxSteps)),
		}
	}

	createdPipeline, err := sagemaker.NewPipeline(ctx, locals.PipelineName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create pipeline")
	}

	ctx.Export(OpPipelineName, createdPipeline.PipelineName)
	ctx.Export(OpPipelineArn, createdPipeline.Arn)

	return nil
}
