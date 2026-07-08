package module

import (
	awskinesisfirehose "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awskinesisfirehose/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kinesis"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildSnowflakeArgs constructs the Snowflake destination configuration
// (Snowpipe Streaming — direct table inserts, no intermediate S3 staging).
func buildSnowflakeArgs(dest *awskinesisfirehose.AwsKinesisFirehoseSnowflakeDestination, locals *Locals) (*kinesis.FirehoseDeliveryStreamSnowflakeConfigurationArgs, error) {
	args := &kinesis.FirehoseDeliveryStreamSnowflakeConfigurationArgs{
		AccountUrl: pulumi.String(dest.AccountUrl),
		Database:   pulumi.String(dest.Database),
		Schema:     pulumi.String(dest.Schema),
		Table:      pulumi.String(dest.Table),
		RoleArn:    pulumi.String(dest.RoleArn.GetValue()),
	}

	// Authentication — inline key pair XOR Secrets Manager (CEL-enforced)
	if dest.User != "" {
		args.User = pulumi.StringPtr(dest.User)
	}
	if dest.PrivateKey != "" {
		args.PrivateKey = pulumi.StringPtr(dest.PrivateKey)
	}
	if dest.KeyPassphrase != "" {
		args.KeyPassphrase = pulumi.StringPtr(dest.KeyPassphrase)
	}
	if sm := dest.SecretsManager; sm != nil {
		smArgs := &kinesis.FirehoseDeliveryStreamSnowflakeConfigurationSecretsManagerConfigurationArgs{
			Enabled:   pulumi.BoolPtr(true),
			SecretArn: pulumi.StringPtr(sm.SecretArn.GetValue()),
		}
		if sm.RoleArn != nil {
			smArgs.RoleArn = pulumi.StringPtr(sm.RoleArn.GetValue())
		}
		args.SecretsManagerConfiguration = smArgs
	}

	// Data loading
	if dest.DataLoadingOption != "" {
		args.DataLoadingOption = pulumi.StringPtr(dest.DataLoadingOption)
	}
	if dest.ContentColumnName != "" {
		args.ContentColumnName = pulumi.StringPtr(dest.ContentColumnName)
	}
	if dest.MetadataColumnName != "" {
		args.MetadataColumnName = pulumi.StringPtr(dest.MetadataColumnName)
	}

	// Snowflake role (least-privilege ingestion role)
	if dest.SnowflakeRole != "" {
		args.SnowflakeRoleConfiguration = &kinesis.FirehoseDeliveryStreamSnowflakeConfigurationSnowflakeRoleConfigurationArgs{
			Enabled:       pulumi.BoolPtr(true),
			SnowflakeRole: pulumi.StringPtr(dest.SnowflakeRole),
		}
	}

	// PrivateLink connectivity
	if dest.PrivateLinkVpceId != "" {
		args.SnowflakeVpcConfiguration = &kinesis.FirehoseDeliveryStreamSnowflakeConfigurationSnowflakeVpcConfigurationArgs{
			PrivateLinkVpceId: pulumi.String(dest.PrivateLinkVpceId),
		}
	}

	// Buffering (Snowpipe Streaming defaults: 0s / 1 MiB)
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
		args.S3Configuration = buildSnowflakeS3Config(cfg)
	}

	// Processing pipeline (normalized typed processors)
	if proc := buildSnowflakeProcessing(dest.Processing); proc != nil {
		args.ProcessingConfiguration = proc
	}

	// CloudWatch logging
	if logging := dest.Logging; logging != nil && logging.Enabled {
		args.CloudwatchLoggingOptions = &kinesis.FirehoseDeliveryStreamSnowflakeConfigurationCloudwatchLoggingOptionsArgs{
			Enabled:       pulumi.Bool(true),
			LogGroupName:  pulumi.StringPtr(logging.LogGroupName),
			LogStreamName: pulumi.StringPtr(logging.LogStreamName),
		}
	}

	return args, nil
}

// buildSnowflakeS3Config builds the S3 backup configuration for Snowflake.
func buildSnowflakeS3Config(cfg *awskinesisfirehose.AwsKinesisFirehoseS3Config) *kinesis.FirehoseDeliveryStreamSnowflakeConfigurationS3ConfigurationArgs {
	args := &kinesis.FirehoseDeliveryStreamSnowflakeConfigurationS3ConfigurationArgs{
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

// buildSnowflakeProcessing adapts the normalized processor pipeline to the
// Snowflake destination's SDK types.
func buildSnowflakeProcessing(processing *awskinesisfirehose.AwsKinesisFirehoseProcessing) *kinesis.FirehoseDeliveryStreamSnowflakeConfigurationProcessingConfigurationArgs {
	normalized := normalizeProcessors(processing)
	if normalized == nil {
		return nil
	}
	processors := kinesis.FirehoseDeliveryStreamSnowflakeConfigurationProcessingConfigurationProcessorArray{}
	for _, p := range normalized {
		params := kinesis.FirehoseDeliveryStreamSnowflakeConfigurationProcessingConfigurationProcessorParameterArray{}
		for _, param := range p.Parameters {
			params = append(params, &kinesis.FirehoseDeliveryStreamSnowflakeConfigurationProcessingConfigurationProcessorParameterArgs{
				ParameterName:  pulumi.String(param.Name),
				ParameterValue: pulumi.String(param.Value),
			})
		}
		processors = append(processors, &kinesis.FirehoseDeliveryStreamSnowflakeConfigurationProcessingConfigurationProcessorArgs{
			Type:       pulumi.String(p.Type),
			Parameters: params,
		})
	}
	return &kinesis.FirehoseDeliveryStreamSnowflakeConfigurationProcessingConfigurationArgs{
		Enabled:    pulumi.Bool(true),
		Processors: processors,
	}
}
