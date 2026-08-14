package module

import (
	"github.com/pkg/errors"
	awssagemakermodelv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssagemakermodel/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// model creates the SageMaker model and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - every argument is create-time only (the provider's update is
//     tags-only) - any spec change replaces the model, which is AWS's
//     own contract (roll a new model, repoint the endpoint);
//   - primary_container and container share one schema upstream (two
//     bridged types here); the spec's exactly-one rule decides which
//     form renders;
//   - the s3_data_source wrapper is single-valued upstream (the
//     expander reads index 0 only) - rendered as a one-element array.
func model(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &sagemaker.ModelArgs{
		// The component's name IS the model name.
		Name:             pulumi.String(locals.ModelName),
		ExecutionRoleArn: pulumi.String(spec.ExecutionRoleArn.GetValue()),
		Tags:             pulumi.ToStringMap(locals.AwsTags),
	}

	// Isolate the model container: no network calls in or out.
	if spec.EnableNetworkIsolation {
		args.EnableNetworkIsolation = pulumi.Bool(true)
	}

	// How pipeline containers are invoked - only meaningful with the
	// pipeline form (spec-validated).
	if spec.InferenceExecutionMode != "" {
		args.InferenceExecutionConfig = &sagemaker.ModelInferenceExecutionConfigArgs{
			Mode: pulumi.String(spec.InferenceExecutionMode),
		}
	}

	// The single-container form.
	if spec.PrimaryContainer != nil {
		args.PrimaryContainer = primaryContainerArgs(spec.PrimaryContainer)
	}

	// The inference-pipeline form (2-15 containers, same schema as the
	// primary container upstream).
	var containers sagemaker.ModelContainerArray
	for _, c := range spec.Containers {
		containers = append(containers, pipelineContainerArgs(c))
	}
	if len(containers) > 0 {
		args.Containers = containers
	}

	// Attach the containers to your VPC (private serving).
	if spec.VpcConfig != nil {
		var subnets pulumi.StringArray
		for _, s := range spec.VpcConfig.SubnetIds {
			subnets = append(subnets, pulumi.String(s.GetValue()))
		}
		var securityGroups pulumi.StringArray
		for _, s := range spec.VpcConfig.SecurityGroupIds {
			securityGroups = append(securityGroups, pulumi.String(s.GetValue()))
		}
		args.VpcConfig = &sagemaker.ModelVpcConfigArgs{
			Subnets:          subnets,
			SecurityGroupIds: securityGroups,
		}
	}

	createdModel, err := sagemaker.NewModel(ctx, locals.ModelName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create model")
	}

	ctx.Export(OpModelName, createdModel.Name)
	ctx.Export(OpModelArn, createdModel.Arn)

	return nil
}

// primaryContainerArgs maps the shared container message onto the
// primary-container bridged type.
func primaryContainerArgs(c *awssagemakermodelv1alpha1.AwsSagemakerModelContainer) *sagemaker.ModelPrimaryContainerArgs {
	args := &sagemaker.ModelPrimaryContainerArgs{}
	if c.Image != "" {
		args.Image = pulumi.String(c.Image)
	}
	if c.ModelPackageArn != "" {
		args.ModelPackageName = pulumi.String(c.ModelPackageArn)
	}
	if c.ContainerHostname != "" {
		args.ContainerHostname = pulumi.String(c.ContainerHostname)
	}
	if len(c.Environment) > 0 {
		args.Environment = pulumi.ToStringMap(c.Environment)
	}
	if c.Mode != "" {
		args.Mode = pulumi.String(c.Mode)
	}
	if c.ModelDataUrl != "" {
		args.ModelDataUrl = pulumi.String(c.ModelDataUrl)
	}
	if c.InferenceSpecificationName != "" {
		args.InferenceSpecificationName = pulumi.String(c.InferenceSpecificationName)
	}
	if c.ModelDataSource != nil {
		// The wrapper is single-valued upstream - a one-element array.
		args.ModelDataSource = &sagemaker.ModelPrimaryContainerModelDataSourceArgs{
			S3DataSources: sagemaker.ModelPrimaryContainerModelDataSourceS3DataSourceArray{
				primaryS3DataSourceArgs(c.ModelDataSource),
			},
		}
	}
	if len(c.AdditionalModelDataSources) > 0 {
		var additional sagemaker.ModelPrimaryContainerAdditionalModelDataSourceArray
		for _, a := range c.AdditionalModelDataSources {
			additional = append(additional, &sagemaker.ModelPrimaryContainerAdditionalModelDataSourceArgs{
				ChannelName: pulumi.String(a.ChannelName),
				S3DataSources: sagemaker.ModelPrimaryContainerAdditionalModelDataSourceS3DataSourceArray{
					primaryAdditionalS3DataSourceArgs(a.Source),
				},
			})
		}
		args.AdditionalModelDataSources = additional
	}
	// MultiModel caching - only meaningful in MultiModel mode
	// (spec-validated).
	if c.MultiModelCache != "" {
		args.MultiModelConfig = &sagemaker.ModelPrimaryContainerMultiModelConfigArgs{
			ModelCacheSetting: pulumi.String(c.MultiModelCache),
		}
	}
	// Pull from a private VPC-reachable registry instead of ECR.
	if c.ImageConfig != nil {
		imageConfig := &sagemaker.ModelPrimaryContainerImageConfigArgs{
			RepositoryAccessMode: pulumi.String(c.ImageConfig.RepositoryAccessMode),
		}
		if c.ImageConfig.RepositoryCredentialsProviderArn != "" {
			imageConfig.RepositoryAuthConfig = &sagemaker.ModelPrimaryContainerImageConfigRepositoryAuthConfigArgs{
				RepositoryCredentialsProviderArn: pulumi.String(c.ImageConfig.RepositoryCredentialsProviderArn),
			}
		}
		args.ImageConfig = imageConfig
	}
	return args
}

func primaryS3DataSourceArgs(s *awssagemakermodelv1alpha1.AwsSagemakerModelS3DataSource) *sagemaker.ModelPrimaryContainerModelDataSourceS3DataSourceArgs {
	args := &sagemaker.ModelPrimaryContainerModelDataSourceS3DataSourceArgs{
		S3Uri:           pulumi.String(s.S3Uri),
		S3DataType:      pulumi.String(s.S3DataType),
		CompressionType: pulumi.String(s.CompressionType),
	}
	if s.AcceptEula {
		args.ModelAccessConfig = &sagemaker.ModelPrimaryContainerModelDataSourceS3DataSourceModelAccessConfigArgs{
			AcceptEula: pulumi.Bool(true),
		}
	}
	return args
}

func primaryAdditionalS3DataSourceArgs(s *awssagemakermodelv1alpha1.AwsSagemakerModelS3DataSource) *sagemaker.ModelPrimaryContainerAdditionalModelDataSourceS3DataSourceArgs {
	args := &sagemaker.ModelPrimaryContainerAdditionalModelDataSourceS3DataSourceArgs{
		S3Uri:           pulumi.String(s.S3Uri),
		S3DataType:      pulumi.String(s.S3DataType),
		CompressionType: pulumi.String(s.CompressionType),
	}
	if s.AcceptEula {
		args.ModelAccessConfig = &sagemaker.ModelPrimaryContainerAdditionalModelDataSourceS3DataSourceModelAccessConfigArgs{
			AcceptEula: pulumi.Bool(true),
		}
	}
	return args
}

func containerS3DataSourceArgs(s *awssagemakermodelv1alpha1.AwsSagemakerModelS3DataSource) *sagemaker.ModelContainerModelDataSourceS3DataSourceArgs {
	args := &sagemaker.ModelContainerModelDataSourceS3DataSourceArgs{
		S3Uri:           pulumi.String(s.S3Uri),
		S3DataType:      pulumi.String(s.S3DataType),
		CompressionType: pulumi.String(s.CompressionType),
	}
	if s.AcceptEula {
		args.ModelAccessConfig = &sagemaker.ModelContainerModelDataSourceS3DataSourceModelAccessConfigArgs{
			AcceptEula: pulumi.Bool(true),
		}
	}
	return args
}

func containerAdditionalS3DataSourceArgs(s *awssagemakermodelv1alpha1.AwsSagemakerModelS3DataSource) *sagemaker.ModelContainerAdditionalModelDataSourceS3DataSourceArgs {
	args := &sagemaker.ModelContainerAdditionalModelDataSourceS3DataSourceArgs{
		S3Uri:           pulumi.String(s.S3Uri),
		S3DataType:      pulumi.String(s.S3DataType),
		CompressionType: pulumi.String(s.CompressionType),
	}
	if s.AcceptEula {
		args.ModelAccessConfig = &sagemaker.ModelContainerAdditionalModelDataSourceS3DataSourceModelAccessConfigArgs{
			AcceptEula: pulumi.Bool(true),
		}
	}
	return args
}

// pipelineContainerArgs maps the shared container message onto the
// pipeline-container bridged type (same schema, distinct type).
func pipelineContainerArgs(c *awssagemakermodelv1alpha1.AwsSagemakerModelContainer) *sagemaker.ModelContainerArgs {
	args := &sagemaker.ModelContainerArgs{}
	if c.Image != "" {
		args.Image = pulumi.String(c.Image)
	}
	if c.ModelPackageArn != "" {
		args.ModelPackageName = pulumi.String(c.ModelPackageArn)
	}
	if c.ContainerHostname != "" {
		args.ContainerHostname = pulumi.String(c.ContainerHostname)
	}
	if len(c.Environment) > 0 {
		args.Environment = pulumi.ToStringMap(c.Environment)
	}
	if c.Mode != "" {
		args.Mode = pulumi.String(c.Mode)
	}
	if c.ModelDataUrl != "" {
		args.ModelDataUrl = pulumi.String(c.ModelDataUrl)
	}
	if c.InferenceSpecificationName != "" {
		args.InferenceSpecificationName = pulumi.String(c.InferenceSpecificationName)
	}
	if c.ModelDataSource != nil {
		args.ModelDataSource = &sagemaker.ModelContainerModelDataSourceArgs{
			S3DataSources: sagemaker.ModelContainerModelDataSourceS3DataSourceArray{
				containerS3DataSourceArgs(c.ModelDataSource),
			},
		}
	}
	if len(c.AdditionalModelDataSources) > 0 {
		var additional sagemaker.ModelContainerAdditionalModelDataSourceArray
		for _, a := range c.AdditionalModelDataSources {
			additional = append(additional, &sagemaker.ModelContainerAdditionalModelDataSourceArgs{
				ChannelName: pulumi.String(a.ChannelName),
				S3DataSources: sagemaker.ModelContainerAdditionalModelDataSourceS3DataSourceArray{
					containerAdditionalS3DataSourceArgs(a.Source),
				},
			})
		}
		args.AdditionalModelDataSources = additional
	}
	if c.MultiModelCache != "" {
		args.MultiModelConfig = &sagemaker.ModelContainerMultiModelConfigArgs{
			ModelCacheSetting: pulumi.String(c.MultiModelCache),
		}
	}
	if c.ImageConfig != nil {
		imageConfig := &sagemaker.ModelContainerImageConfigArgs{
			RepositoryAccessMode: pulumi.String(c.ImageConfig.RepositoryAccessMode),
		}
		if c.ImageConfig.RepositoryCredentialsProviderArn != "" {
			imageConfig.RepositoryAuthConfig = &sagemaker.ModelContainerImageConfigRepositoryAuthConfigArgs{
				RepositoryCredentialsProviderArn: pulumi.String(c.ImageConfig.RepositoryCredentialsProviderArn),
			}
		}
		args.ImageConfig = imageConfig
	}
	return args
}
