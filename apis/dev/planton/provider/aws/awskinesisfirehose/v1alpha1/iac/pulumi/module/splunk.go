package module

import (
	awskinesisfirehose "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awskinesisfirehose/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kinesis"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildSplunkArgs constructs the Splunk destination configuration. Note:
// Splunk has no destination-level role ARN — the S3 configuration's role
// carries the backup permissions, and HEC authorization is the token (or a
// Secrets Manager secret).
func buildSplunkArgs(dest *awskinesisfirehose.AwsKinesisFirehoseSplunkDestination, locals *Locals) (*kinesis.FirehoseDeliveryStreamSplunkConfigurationArgs, error) {
	args := &kinesis.FirehoseDeliveryStreamSplunkConfigurationArgs{
		HecEndpoint: pulumi.String(dest.HecEndpoint),
	}

	if dest.HecEndpointType != "" {
		args.HecEndpointType = pulumi.StringPtr(dest.HecEndpointType)
	}

	// Authentication — plaintext HEC token XOR Secrets Manager (CEL-enforced)
	if dest.HecToken != "" {
		args.HecToken = pulumi.StringPtr(dest.HecToken)
	}
	if sm := dest.SecretsManager; sm != nil {
		smArgs := &kinesis.FirehoseDeliveryStreamSplunkConfigurationSecretsManagerConfigurationArgs{
			Enabled:   pulumi.BoolPtr(true),
			SecretArn: pulumi.StringPtr(sm.SecretArn.GetValue()),
		}
		if sm.RoleArn != nil {
			smArgs.RoleArn = pulumi.StringPtr(sm.RoleArn.GetValue())
		}
		args.SecretsManagerConfiguration = smArgs
	}

	if dest.HecAcknowledgmentTimeoutInSeconds > 0 {
		args.HecAcknowledgmentTimeout = pulumi.IntPtr(int(dest.HecAcknowledgmentTimeoutInSeconds))
	}

	// Buffering (Splunk caps: 0-60s interval, 1-5 MiB size — CEL-enforced)
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

	// S3 backup mode
	if dest.S3BackupMode != "" {
		args.S3BackupMode = pulumi.StringPtr(dest.S3BackupMode)
	}

	// S3 config (required)
	if cfg := dest.S3Config; cfg != nil {
		args.S3Configuration = buildSplunkS3Config(cfg)
	}

	// Processing pipeline (normalized typed processors)
	if proc := buildSplunkProcessing(dest.Processing); proc != nil {
		args.ProcessingConfiguration = proc
	}

	// CloudWatch logging
	if logging := dest.Logging; logging != nil && logging.Enabled {
		args.CloudwatchLoggingOptions = &kinesis.FirehoseDeliveryStreamSplunkConfigurationCloudwatchLoggingOptionsArgs{
			Enabled:       pulumi.Bool(true),
			LogGroupName:  pulumi.StringPtr(logging.LogGroupName),
			LogStreamName: pulumi.StringPtr(logging.LogStreamName),
		}
	}

	return args, nil
}

// buildSplunkS3Config builds the S3 backup configuration for Splunk.
func buildSplunkS3Config(cfg *awskinesisfirehose.AwsKinesisFirehoseS3Config) *kinesis.FirehoseDeliveryStreamSplunkConfigurationS3ConfigurationArgs {
	args := &kinesis.FirehoseDeliveryStreamSplunkConfigurationS3ConfigurationArgs{
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

// buildSplunkProcessing adapts the normalized processor pipeline to the
// Splunk destination's SDK types.
func buildSplunkProcessing(processing *awskinesisfirehose.AwsKinesisFirehoseProcessing) *kinesis.FirehoseDeliveryStreamSplunkConfigurationProcessingConfigurationArgs {
	normalized := normalizeProcessors(processing)
	if normalized == nil {
		return nil
	}
	processors := kinesis.FirehoseDeliveryStreamSplunkConfigurationProcessingConfigurationProcessorArray{}
	for _, p := range normalized {
		params := kinesis.FirehoseDeliveryStreamSplunkConfigurationProcessingConfigurationProcessorParameterArray{}
		for _, param := range p.Parameters {
			params = append(params, &kinesis.FirehoseDeliveryStreamSplunkConfigurationProcessingConfigurationProcessorParameterArgs{
				ParameterName:  pulumi.String(param.Name),
				ParameterValue: pulumi.String(param.Value),
			})
		}
		processors = append(processors, &kinesis.FirehoseDeliveryStreamSplunkConfigurationProcessingConfigurationProcessorArgs{
			Type:       pulumi.String(p.Type),
			Parameters: params,
		})
	}
	return &kinesis.FirehoseDeliveryStreamSplunkConfigurationProcessingConfigurationArgs{
		Enabled:    pulumi.Bool(true),
		Processors: processors,
	}
}
