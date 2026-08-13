package module

import (
	awskinesisfirehose "github.com/plantonhq/planton/catalog/aws/awskinesisfirehose/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kinesis"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildOpenSearchServerlessArgs constructs the OpenSearch Serverless
// destination configuration. Unlike the domain-based OpenSearch destination
// there is no domain ARN, no index rotation, and no document-ID option — the
// collection endpoint plus a fixed index name is the whole target.
func buildOpenSearchServerlessArgs(dest *awskinesisfirehose.AwsKinesisFirehoseOpenSearchServerlessDestination, locals *Locals) (*kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationArgs, error) {
	args := &kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationArgs{
		CollectionEndpoint: pulumi.String(dest.CollectionEndpoint),
		IndexName:          pulumi.String(dest.IndexName),
		RoleArn:            pulumi.String(dest.RoleArn.GetValue()),
	}

	// Buffering (0-900s / 1-100 MiB — the 100 MiB cap is CEL-enforced)
	if b := dest.Buffering; b != nil {
		if b.IntervalInSeconds > 0 {
			args.BufferingInterval = pulumi.IntPtr(int(b.IntervalInSeconds))
		}
		if b.SizeInMbs > 0 {
			args.BufferingSize = pulumi.IntPtr(int(b.SizeInMbs))
		}
	}

	if dest.RetryDurationInSeconds > 0 {
		args.RetryDuration = pulumi.IntPtr(int(dest.RetryDurationInSeconds))
	}

	// S3 backup mode (ForceNew)
	if dest.S3BackupMode != "" {
		args.S3BackupMode = pulumi.StringPtr(dest.S3BackupMode)
	}

	// S3 config (required)
	if cfg := dest.S3Config; cfg != nil {
		args.S3Configuration = buildOpenSearchServerlessS3Config(cfg)
	}

	// Processing pipeline (normalized typed processors)
	if proc := buildOpenSearchServerlessProcessing(dest.Processing); proc != nil {
		args.ProcessingConfiguration = proc
	}

	// CloudWatch logging
	if logging := dest.Logging; logging != nil && logging.Enabled {
		args.CloudwatchLoggingOptions = &kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationCloudwatchLoggingOptionsArgs{
			Enabled:       pulumi.Bool(true),
			LogGroupName:  pulumi.StringPtr(logging.LogGroupName),
			LogStreamName: pulumi.StringPtr(logging.LogStreamName),
		}
	}

	// VPC config (ForceNew)
	if vpc := dest.VpcConfig; vpc != nil {
		vpcArgs := &kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationVpcConfigArgs{
			RoleArn: pulumi.String(vpc.RoleArn.GetValue()),
		}
		subnetIds := pulumi.StringArray{}
		for _, s := range vpc.SubnetIds {
			subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
		}
		vpcArgs.SubnetIds = subnetIds

		sgIds := pulumi.StringArray{}
		for _, s := range vpc.SecurityGroupIds {
			sgIds = append(sgIds, pulumi.String(s.GetValue()))
		}
		vpcArgs.SecurityGroupIds = sgIds
		args.VpcConfig = vpcArgs
	}

	return args, nil
}

// buildOpenSearchServerlessS3Config builds the S3 backup configuration for
// OpenSearch Serverless.
func buildOpenSearchServerlessS3Config(cfg *awskinesisfirehose.AwsKinesisFirehoseS3Config) *kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationS3ConfigurationArgs {
	args := &kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationS3ConfigurationArgs{
		BucketArn: pulumi.String(cfg.BucketArn.GetValue()),
		RoleArn:   pulumi.String(cfg.RoleArn.GetValue()),
	}
	if cfg.Prefix != "" {
		args.Prefix = pulumi.StringPtr(cfg.Prefix)
	}
	if cfg.ErrorOutputPrefix != "" {
		args.ErrorOutputPrefix = pulumi.StringPtr(cfg.ErrorOutputPrefix)
	}
	if cfg.CompressionFormat != "" {
		args.CompressionFormat = pulumi.StringPtr(cfg.CompressionFormat)
	}
	if cfg.KmsKeyArn != nil {
		args.KmsKeyArn = pulumi.StringPtr(cfg.KmsKeyArn.GetValue())
	}
	if b := cfg.Buffering; b != nil {
		if b.IntervalInSeconds > 0 {
			args.BufferingInterval = pulumi.IntPtr(int(b.IntervalInSeconds))
		}
		if b.SizeInMbs > 0 {
			args.BufferingSize = pulumi.IntPtr(int(b.SizeInMbs))
		}
	}
	// CloudWatch logging for the backup S3 leg -- distinct from the
	// destination-level logging options.
	if log := cfg.Logging; log != nil && log.Enabled {
		args.CloudwatchLoggingOptions = &kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationS3ConfigurationCloudwatchLoggingOptionsArgs{
			Enabled:       pulumi.BoolPtr(true),
			LogGroupName:  pulumi.StringPtr(log.LogGroupName),
			LogStreamName: pulumi.StringPtr(log.LogStreamName),
		}
	}
	return args
}

// buildOpenSearchServerlessProcessing adapts the normalized processor
// pipeline to the OpenSearch Serverless destination's SDK types.
func buildOpenSearchServerlessProcessing(processing *awskinesisfirehose.AwsKinesisFirehoseProcessing) *kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationProcessingConfigurationArgs {
	normalized := normalizeProcessors(processing)
	if normalized == nil {
		return nil
	}
	processors := kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationProcessingConfigurationProcessorArray{}
	for _, p := range normalized {
		params := kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationProcessingConfigurationProcessorParameterArray{}
		for _, param := range p.Parameters {
			params = append(params, &kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationProcessingConfigurationProcessorParameterArgs{
				ParameterName:  pulumi.String(param.Name),
				ParameterValue: pulumi.String(param.Value),
			})
		}
		processors = append(processors, &kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationProcessingConfigurationProcessorArgs{
			Type:       pulumi.String(p.Type),
			Parameters: params,
		})
	}
	return &kinesis.FirehoseDeliveryStreamOpensearchserverlessConfigurationProcessingConfigurationArgs{
		Enabled:    pulumi.Bool(true),
		Processors: processors,
	}
}
