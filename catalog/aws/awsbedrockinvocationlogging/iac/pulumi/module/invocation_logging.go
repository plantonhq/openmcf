package module

import (
	"github.com/pkg/errors"
	awsbedrockinvocationloggingv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockinvocationlogging/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrockmodel"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// invocationLogging manages the region's ONE Bedrock invocation
// logging configuration and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the four data-type toggles default to TRUE upstream; the spec's
//     presence-typed optional bools render only when set (unset lets
//     the provider apply AWS's enabled default, explicit false is
//     sent);
//   - at least one delivery destination is guaranteed by the spec's
//     CEL, so the arm conditionals below never both skip;
//   - AWS validates the CloudWatch role's permission chain at apply
//     ("Failed to validate permissions for log group") and the
//     provider retries through IAM propagation lag;
//   - destroy DELETES the configuration -- the region reverts to no
//     invocation logging.
func invocationLogging(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	loggingConfig := &bedrockmodel.InvocationLoggingConfigurationLoggingConfigArgs{}
	if spec.TextDataDeliveryEnabled != nil {
		loggingConfig.TextDataDeliveryEnabled = pulumi.Bool(spec.GetTextDataDeliveryEnabled())
	}
	if spec.ImageDataDeliveryEnabled != nil {
		loggingConfig.ImageDataDeliveryEnabled = pulumi.Bool(spec.GetImageDataDeliveryEnabled())
	}
	if spec.EmbeddingDataDeliveryEnabled != nil {
		loggingConfig.EmbeddingDataDeliveryEnabled = pulumi.Bool(spec.GetEmbeddingDataDeliveryEnabled())
	}
	if spec.VideoDataDeliveryEnabled != nil {
		loggingConfig.VideoDataDeliveryEnabled = pulumi.Bool(spec.GetVideoDataDeliveryEnabled())
	}

	if spec.Cloudwatch != nil {
		cloudwatchConfig := &bedrockmodel.InvocationLoggingConfigurationLoggingConfigCloudwatchConfigArgs{
			LogGroupName: pulumi.String(spec.Cloudwatch.LogGroupName.GetValue()),
			RoleArn:      pulumi.String(spec.Cloudwatch.RoleArn.GetValue()),
		}
		// Where CloudWatch delivery spills payloads larger than a log
		// event (256 KB); without it oversized bodies are truncated.
		if spec.Cloudwatch.LargeDataDeliveryS3 != nil {
			cloudwatchConfig.LargeDataDeliveryS3Config = &bedrockmodel.InvocationLoggingConfigurationLoggingConfigCloudwatchConfigLargeDataDeliveryS3ConfigArgs{
				BucketName: pulumi.String(spec.Cloudwatch.LargeDataDeliveryS3.BucketName.GetValue()),
				KeyPrefix:  keyPrefixOrNil(spec.Cloudwatch.LargeDataDeliveryS3),
			}
		}
		loggingConfig.CloudwatchConfig = cloudwatchConfig
	}

	if spec.S3 != nil {
		loggingConfig.S3Config = &bedrockmodel.InvocationLoggingConfigurationLoggingConfigS3ConfigArgs{
			BucketName: pulumi.String(spec.S3.BucketName.GetValue()),
			KeyPrefix:  keyPrefixOrNil(spec.S3),
		}
	}

	created, err := bedrockmodel.NewInvocationLoggingConfiguration(ctx, "invocation-logging",
		&bedrockmodel.InvocationLoggingConfigurationArgs{
			LoggingConfig: loggingConfig,
		}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "put invocation logging configuration")
	}

	ctx.Export(OpConfiguredRegion, created.ID())
	return nil
}

// keyPrefixOrNil renders the optional key prefix only when set,
// matching the Terraform module's empty-to-null conversion.
func keyPrefixOrNil(s3 *awsbedrockinvocationloggingv1alpha1.AwsBedrockInvocationLoggingS3) pulumi.StringPtrInput {
	if s3.KeyPrefix == "" {
		return nil
	}
	return pulumi.String(s3.KeyPrefix)
}
