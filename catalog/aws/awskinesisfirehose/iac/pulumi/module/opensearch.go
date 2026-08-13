package module

import (
	awskinesisfirehose "github.com/plantonhq/planton/catalog/aws/awskinesisfirehose/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kinesis"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildOpenSearchArgs constructs the OpenSearch destination configuration.
func buildOpenSearchArgs(dest *awskinesisfirehose.AwsKinesisFirehoseOpenSearchDestination, locals *Locals) (*kinesis.FirehoseDeliveryStreamOpensearchConfigurationArgs, error) {
	args := &kinesis.FirehoseDeliveryStreamOpensearchConfigurationArgs{
		IndexName: pulumi.String(dest.IndexName),
		RoleArn:   pulumi.String(dest.RoleArn.GetValue()),
	}

	// Target: domain_arn XOR cluster_endpoint
	if dest.DomainArn != nil {
		args.DomainArn = pulumi.StringPtr(dest.DomainArn.GetValue())
	}
	if dest.ClusterEndpoint != "" {
		args.ClusterEndpoint = pulumi.StringPtr(dest.ClusterEndpoint)
	}

	if dest.IndexRotationPeriod != "" {
		args.IndexRotationPeriod = pulumi.StringPtr(dest.IndexRotationPeriod)
	}
	if dest.TypeName != "" {
		args.TypeName = pulumi.StringPtr(dest.TypeName)
	}

	// Document ID assignment (FIREHOSE_DEFAULT vs NO_DOCUMENT_ID)
	if dest.DefaultDocumentIdFormat != "" {
		args.DocumentIdOptions = &kinesis.FirehoseDeliveryStreamOpensearchConfigurationDocumentIdOptionsArgs{
			DefaultDocumentIdFormat: pulumi.String(dest.DefaultDocumentIdFormat),
		}
	}

	// Buffering
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

	// S3 backup
	if dest.S3BackupMode != "" {
		args.S3BackupMode = pulumi.StringPtr(dest.S3BackupMode)
	}

	// S3 config (required)
	if cfg := dest.S3Config; cfg != nil {
		args.S3Configuration = buildOpenSearchS3Config(cfg)
	}

	// Processing pipeline (normalized typed processors)
	if proc := buildOpenSearchProcessing(dest.Processing); proc != nil {
		args.ProcessingConfiguration = proc
	}

	// CloudWatch logging
	if dest.Logging != nil && dest.Logging.Enabled {
		args.CloudwatchLoggingOptions = buildOpenSearchCloudwatchLogging(dest.Logging)
	}

	// VPC config
	if vpc := dest.VpcConfig; vpc != nil {
		vpcArgs := &kinesis.FirehoseDeliveryStreamOpensearchConfigurationVpcConfigArgs{
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

// buildOpenSearchS3Config builds the S3 configuration for OpenSearch backup.
func buildOpenSearchS3Config(cfg *awskinesisfirehose.AwsKinesisFirehoseS3Config) *kinesis.FirehoseDeliveryStreamOpensearchConfigurationS3ConfigurationArgs {
	args := &kinesis.FirehoseDeliveryStreamOpensearchConfigurationS3ConfigurationArgs{
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
		args.CloudwatchLoggingOptions = &kinesis.FirehoseDeliveryStreamOpensearchConfigurationS3ConfigurationCloudwatchLoggingOptionsArgs{
			Enabled:       pulumi.BoolPtr(true),
			LogGroupName:  pulumi.StringPtr(log.LogGroupName),
			LogStreamName: pulumi.StringPtr(log.LogStreamName),
		}
	}
	return args
}

// buildOpenSearchProcessing adapts the normalized processor pipeline to the
// OpenSearch destination's SDK types.
func buildOpenSearchProcessing(processing *awskinesisfirehose.AwsKinesisFirehoseProcessing) *kinesis.FirehoseDeliveryStreamOpensearchConfigurationProcessingConfigurationArgs {
	normalized := normalizeProcessors(processing)
	if normalized == nil {
		return nil
	}
	processors := kinesis.FirehoseDeliveryStreamOpensearchConfigurationProcessingConfigurationProcessorArray{}
	for _, p := range normalized {
		params := kinesis.FirehoseDeliveryStreamOpensearchConfigurationProcessingConfigurationProcessorParameterArray{}
		for _, param := range p.Parameters {
			params = append(params, &kinesis.FirehoseDeliveryStreamOpensearchConfigurationProcessingConfigurationProcessorParameterArgs{
				ParameterName:  pulumi.String(param.Name),
				ParameterValue: pulumi.String(param.Value),
			})
		}
		processors = append(processors, &kinesis.FirehoseDeliveryStreamOpensearchConfigurationProcessingConfigurationProcessorArgs{
			Type:       pulumi.String(p.Type),
			Parameters: params,
		})
	}
	return &kinesis.FirehoseDeliveryStreamOpensearchConfigurationProcessingConfigurationArgs{
		Enabled:    pulumi.Bool(true),
		Processors: processors,
	}
}

// buildOpenSearchCloudwatchLogging constructs CloudWatch logging for OpenSearch.
func buildOpenSearchCloudwatchLogging(logging *awskinesisfirehose.AwsKinesisFirehoseCloudwatchLogging) *kinesis.FirehoseDeliveryStreamOpensearchConfigurationCloudwatchLoggingOptionsArgs {
	if logging == nil || !logging.Enabled {
		return nil
	}
	return &kinesis.FirehoseDeliveryStreamOpensearchConfigurationCloudwatchLoggingOptionsArgs{
		Enabled:       pulumi.Bool(true),
		LogGroupName:  pulumi.StringPtr(logging.LogGroupName),
		LogStreamName: pulumi.StringPtr(logging.LogStreamName),
	}
}
