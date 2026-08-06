package module

import (
	awskinesisfirehose "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awskinesisfirehose/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kinesis"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildRedshiftArgs constructs the Redshift destination configuration.
func buildRedshiftArgs(dest *awskinesisfirehose.AwsKinesisFirehoseRedshiftDestination, locals *Locals) (*kinesis.FirehoseDeliveryStreamRedshiftConfigurationArgs, error) {
	args := &kinesis.FirehoseDeliveryStreamRedshiftConfigurationArgs{
		ClusterJdbcurl: pulumi.String(dest.ClusterJdbcurl),
		RoleArn:        pulumi.String(dest.RoleArn.GetValue()),
		DataTableName:  pulumi.String(dest.DataTableName),
	}

	if dest.DataTableColumns != "" {
		args.DataTableColumns = pulumi.StringPtr(dest.DataTableColumns)
	}
	if dest.CopyOptions != "" {
		args.CopyOptions = pulumi.StringPtr(dest.CopyOptions)
	}
	// Authentication — plaintext pair XOR Secrets Manager (CEL-enforced)
	if dest.Username != "" {
		args.Username = pulumi.StringPtr(dest.Username)
	}
	if dest.Password != "" {
		args.Password = pulumi.StringPtr(dest.Password)
	}
	if sm := dest.SecretsManager; sm != nil {
		smArgs := &kinesis.FirehoseDeliveryStreamRedshiftConfigurationSecretsManagerConfigurationArgs{
			Enabled:   pulumi.BoolPtr(true),
			SecretArn: pulumi.StringPtr(sm.SecretArn.GetValue()),
		}
		if sm.RoleArn != nil {
			smArgs.RoleArn = pulumi.StringPtr(sm.RoleArn.GetValue())
		}
		args.SecretsManagerConfiguration = smArgs
	}

	if dest.RetryDurationInSeconds > 0 {
		args.RetryDuration = pulumi.IntPtr(int(dest.RetryDurationInSeconds))
	}

	// S3 intermediate config (required for Redshift COPY)
	if cfg := dest.S3Config; cfg != nil {
		args.S3Configuration = buildRedshiftS3Config(cfg)
	}

	// S3 backup
	if dest.S3BackupMode != "" {
		args.S3BackupMode = pulumi.StringPtr(dest.S3BackupMode)
	}
	if dest.S3Backup != nil {
		args.S3BackupConfiguration = buildRedshiftS3BackupConfig(dest.S3Backup)
	}

	// Processing pipeline (normalized typed processors)
	if proc := buildRedshiftProcessing(dest.Processing); proc != nil {
		args.ProcessingConfiguration = proc
	}

	// CloudWatch logging
	if dest.Logging != nil && dest.Logging.Enabled {
		args.CloudwatchLoggingOptions = buildRedshiftCloudwatchLogging(dest.Logging)
	}

	return args, nil
}

// buildRedshiftS3Config builds the S3 intermediate staging configuration for Redshift.
func buildRedshiftS3Config(cfg *awskinesisfirehose.AwsKinesisFirehoseS3Config) *kinesis.FirehoseDeliveryStreamRedshiftConfigurationS3ConfigurationArgs {
	args := &kinesis.FirehoseDeliveryStreamRedshiftConfigurationS3ConfigurationArgs{
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

// buildRedshiftS3BackupConfig builds the S3 backup configuration for Redshift source records.
func buildRedshiftS3BackupConfig(cfg *awskinesisfirehose.AwsKinesisFirehoseS3Config) *kinesis.FirehoseDeliveryStreamRedshiftConfigurationS3BackupConfigurationArgs {
	args := &kinesis.FirehoseDeliveryStreamRedshiftConfigurationS3BackupConfigurationArgs{
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

// buildRedshiftProcessing adapts the normalized processor pipeline to the
// Redshift destination's SDK types.
func buildRedshiftProcessing(processing *awskinesisfirehose.AwsKinesisFirehoseProcessing) *kinesis.FirehoseDeliveryStreamRedshiftConfigurationProcessingConfigurationArgs {
	normalized := normalizeProcessors(processing)
	if normalized == nil {
		return nil
	}
	processors := kinesis.FirehoseDeliveryStreamRedshiftConfigurationProcessingConfigurationProcessorArray{}
	for _, p := range normalized {
		params := kinesis.FirehoseDeliveryStreamRedshiftConfigurationProcessingConfigurationProcessorParameterArray{}
		for _, param := range p.Parameters {
			params = append(params, &kinesis.FirehoseDeliveryStreamRedshiftConfigurationProcessingConfigurationProcessorParameterArgs{
				ParameterName:  pulumi.String(param.Name),
				ParameterValue: pulumi.String(param.Value),
			})
		}
		processors = append(processors, &kinesis.FirehoseDeliveryStreamRedshiftConfigurationProcessingConfigurationProcessorArgs{
			Type:       pulumi.String(p.Type),
			Parameters: params,
		})
	}
	return &kinesis.FirehoseDeliveryStreamRedshiftConfigurationProcessingConfigurationArgs{
		Enabled:    pulumi.Bool(true),
		Processors: processors,
	}
}

// buildRedshiftCloudwatchLogging constructs CloudWatch logging for Redshift.
func buildRedshiftCloudwatchLogging(logging *awskinesisfirehose.AwsKinesisFirehoseCloudwatchLogging) *kinesis.FirehoseDeliveryStreamRedshiftConfigurationCloudwatchLoggingOptionsArgs {
	if logging == nil || !logging.Enabled {
		return nil
	}
	return &kinesis.FirehoseDeliveryStreamRedshiftConfigurationCloudwatchLoggingOptionsArgs{
		Enabled:       pulumi.Bool(true),
		LogGroupName:  pulumi.StringPtr(logging.LogGroupName),
		LogStreamName: pulumi.StringPtr(logging.LogStreamName),
	}
}
