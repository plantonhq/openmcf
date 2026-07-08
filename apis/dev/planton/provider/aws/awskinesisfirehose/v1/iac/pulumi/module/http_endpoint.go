package module

import (
	awskinesisfirehose "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awskinesisfirehose/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kinesis"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildHttpEndpointArgs constructs the HTTP endpoint destination configuration.
func buildHttpEndpointArgs(dest *awskinesisfirehose.AwsKinesisFirehoseHttpEndpointDestination, locals *Locals) (*kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationArgs, error) {
	args := &kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationArgs{
		Url: pulumi.String(dest.Url),
	}

	if dest.Name != "" {
		args.Name = pulumi.StringPtr(dest.Name)
	}
	if dest.AccessKey != "" {
		args.AccessKey = pulumi.StringPtr(dest.AccessKey)
	}
	if dest.RoleArn != nil {
		args.RoleArn = pulumi.StringPtr(dest.RoleArn.GetValue())
	}

	// Secrets Manager credentials (XOR with access_key — CEL-enforced;
	// setting the block IS the enable switch, ForceNew)
	if sm := dest.SecretsManager; sm != nil {
		smArgs := &kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationSecretsManagerConfigurationArgs{
			Enabled:   pulumi.BoolPtr(true),
			SecretArn: pulumi.StringPtr(sm.SecretArn.GetValue()),
		}
		if sm.RoleArn != nil {
			smArgs.RoleArn = pulumi.StringPtr(sm.RoleArn.GetValue())
		}
		args.SecretsManagerConfiguration = smArgs
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

	if dest.S3BackupMode != "" {
		args.S3BackupMode = pulumi.StringPtr(dest.S3BackupMode)
	}

	// S3 config (required)
	if cfg := dest.S3Config; cfg != nil {
		args.S3Configuration = buildHttpS3Config(cfg)
	}

	// Processing pipeline (normalized typed processors)
	if proc := buildHttpProcessing(dest.Processing); proc != nil {
		args.ProcessingConfiguration = proc
	}

	// CloudWatch logging
	if dest.Logging != nil && dest.Logging.Enabled {
		args.CloudwatchLoggingOptions = buildHttpCloudwatchLogging(dest.Logging)
	}

	// Request configuration
	if rc := dest.RequestConfig; rc != nil {
		reqArgs := &kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationRequestConfigurationArgs{}
		if rc.ContentEncoding != "" {
			reqArgs.ContentEncoding = pulumi.StringPtr(rc.ContentEncoding)
		}
		if len(rc.CommonAttributes) > 0 {
			attrs := kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationRequestConfigurationCommonAttributeArray{}
			for _, attr := range rc.CommonAttributes {
				attrs = append(attrs, &kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationRequestConfigurationCommonAttributeArgs{
					Name:  pulumi.String(attr.Name),
					Value: pulumi.String(attr.Value),
				})
			}
			reqArgs.CommonAttributes = attrs
		}
		args.RequestConfiguration = reqArgs
	}

	return args, nil
}

// buildHttpS3Config builds the S3 configuration for HTTP endpoint backup.
func buildHttpS3Config(cfg *awskinesisfirehose.AwsKinesisFirehoseS3Config) *kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationS3ConfigurationArgs {
	args := &kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationS3ConfigurationArgs{
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
	return args
}

// buildHttpProcessing adapts the normalized processor pipeline to the HTTP
// endpoint destination's SDK types.
func buildHttpProcessing(processing *awskinesisfirehose.AwsKinesisFirehoseProcessing) *kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationProcessingConfigurationArgs {
	normalized := normalizeProcessors(processing)
	if normalized == nil {
		return nil
	}
	processors := kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationProcessingConfigurationProcessorArray{}
	for _, p := range normalized {
		params := kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationProcessingConfigurationProcessorParameterArray{}
		for _, param := range p.Parameters {
			params = append(params, &kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationProcessingConfigurationProcessorParameterArgs{
				ParameterName:  pulumi.String(param.Name),
				ParameterValue: pulumi.String(param.Value),
			})
		}
		processors = append(processors, &kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationProcessingConfigurationProcessorArgs{
			Type:       pulumi.String(p.Type),
			Parameters: params,
		})
	}
	return &kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationProcessingConfigurationArgs{
		Enabled:    pulumi.Bool(true),
		Processors: processors,
	}
}

// buildHttpCloudwatchLogging constructs CloudWatch logging for HTTP endpoint.
func buildHttpCloudwatchLogging(logging *awskinesisfirehose.AwsKinesisFirehoseCloudwatchLogging) *kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationCloudwatchLoggingOptionsArgs {
	if logging == nil || !logging.Enabled {
		return nil
	}
	return &kinesis.FirehoseDeliveryStreamHttpEndpointConfigurationCloudwatchLoggingOptionsArgs{
		Enabled:       pulumi.Bool(true),
		LogGroupName:  pulumi.StringPtr(logging.LogGroupName),
		LogStreamName: pulumi.StringPtr(logging.LogStreamName),
	}
}
