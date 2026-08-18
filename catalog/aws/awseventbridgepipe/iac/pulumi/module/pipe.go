package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/pipes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// pipe creates the EventBridge Pipe and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the pipe's name (metadata.name) and its SOURCE are fixed for
//     life (replace-on-change) - so are the per-family stream/topic
//     positions; the TARGET swaps in place;
//   - DesiredState RUNNING/STOPPED flips consumption without deleting
//     (creates and state flips can take minutes; the provider waits
//     up to 30);
//   - Kafka/MQ credentials are Secrets Manager secret ARNs -
//     references, never credential values;
//   - AssignPublicIp is a string enum (ENABLED/DISABLED) and
//     IncludeExecutionDatas a list (["ALL"]) at the provider - the
//     module maps the spec's bools.
func pipe(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &pipes.PipeArgs{
		Name:    pulumi.String(locals.Target.Metadata.Name),
		Source:  pulumi.String(spec.Source.GetValue()),
		Target:  pulumi.String(spec.Target.GetValue()),
		RoleArn: pulumi.String(spec.RoleArn.GetValue()),
		Tags:    pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.DesiredState != "" {
		args.DesiredState = pulumi.String(spec.DesiredState)
	}
	if spec.Enrichment.GetValue() != "" {
		args.Enrichment = pulumi.String(spec.Enrichment.GetValue())
	}
	if spec.KmsKeyIdentifier.GetValue() != "" {
		args.KmsKeyIdentifier = pulumi.String(spec.KmsKeyIdentifier.GetValue())
	}
	if spec.SourceParameters != nil {
		args.SourceParameters = buildSourceParameters(spec.SourceParameters)
	}
	if spec.EnrichmentParameters != nil {
		enrichmentParameters := &pipes.PipeEnrichmentParametersArgs{}
		if spec.EnrichmentParameters.InputTemplate != "" {
			enrichmentParameters.InputTemplate = pulumi.String(spec.EnrichmentParameters.InputTemplate)
		}
		if spec.EnrichmentParameters.HttpParameters != nil {
			httpParameters := &pipes.PipeEnrichmentParametersHttpParametersArgs{
				HeaderParameters:      pulumi.ToStringMap(spec.EnrichmentParameters.HttpParameters.HeaderParameters),
				QueryStringParameters: pulumi.ToStringMap(spec.EnrichmentParameters.HttpParameters.QueryStringParameters),
			}
			if spec.EnrichmentParameters.HttpParameters.PathParameterValue != "" {
				httpParameters.PathParameterValues = pulumi.String(spec.EnrichmentParameters.HttpParameters.PathParameterValue)
			}
			enrichmentParameters.HttpParameters = httpParameters
		}
		args.EnrichmentParameters = enrichmentParameters
	}
	if spec.TargetParameters != nil {
		args.TargetParameters = buildTargetParameters(spec.TargetParameters)
	}
	if spec.LogConfiguration != nil {
		logConfiguration := &pipes.PipeLogConfigurationArgs{
			Level: pulumi.String(spec.LogConfiguration.Level),
		}
		if spec.LogConfiguration.IncludeExecutionData {
			logConfiguration.IncludeExecutionDatas = pulumi.StringArray{pulumi.String("ALL")}
		}
		if spec.LogConfiguration.CloudwatchLogs != nil {
			logConfiguration.CloudwatchLogsLogDestination = &pipes.PipeLogConfigurationCloudwatchLogsLogDestinationArgs{
				LogGroupArn: pulumi.String(spec.LogConfiguration.CloudwatchLogs.LogGroupArn.GetValue()),
			}
		}
		if spec.LogConfiguration.Firehose != nil {
			logConfiguration.FirehoseLogDestination = &pipes.PipeLogConfigurationFirehoseLogDestinationArgs{
				DeliveryStreamArn: pulumi.String(spec.LogConfiguration.Firehose.DeliveryStreamArn.GetValue()),
			}
		}
		if spec.LogConfiguration.S3 != nil {
			s3Destination := &pipes.PipeLogConfigurationS3LogDestinationArgs{
				BucketName:  pulumi.String(spec.LogConfiguration.S3.BucketName.GetValue()),
				BucketOwner: pulumi.String(spec.LogConfiguration.S3.BucketOwner),
			}
			if spec.LogConfiguration.S3.OutputFormat != "" {
				s3Destination.OutputFormat = pulumi.String(spec.LogConfiguration.S3.OutputFormat)
			}
			if spec.LogConfiguration.S3.Prefix != "" {
				s3Destination.Prefix = pulumi.String(spec.LogConfiguration.S3.Prefix)
			}
			logConfiguration.S3LogDestination = s3Destination
		}
		args.LogConfiguration = logConfiguration
	}

	createdPipe, err := pipes.NewPipe(ctx, "pipe", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create pipe")
	}

	ctx.Export(OpPipeArn, createdPipe.Arn)
	ctx.Export(OpPipeName, createdPipe.Name)
	return nil
}
