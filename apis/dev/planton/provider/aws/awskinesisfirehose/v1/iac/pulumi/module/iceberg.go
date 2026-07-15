package module

import (
	awskinesisfirehose "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awskinesisfirehose/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kinesis"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildIcebergArgs constructs the Iceberg destination configuration —
// direct commits into Glue-cataloged Apache Iceberg tables, with optional
// multi-table routing and unique-key upserts.
func buildIcebergArgs(dest *awskinesisfirehose.AwsKinesisFirehoseIcebergDestination, locals *Locals) (*kinesis.FirehoseDeliveryStreamIcebergConfigurationArgs, error) {
	args := &kinesis.FirehoseDeliveryStreamIcebergConfigurationArgs{
		CatalogArn: pulumi.String(dest.CatalogArn.GetValue()),
		RoleArn:    pulumi.String(dest.RoleArn.GetValue()),
	}

	// Append-only mode (ForceNew) — only sent when true so AWS owns the
	// default when unset.
	if dest.AppendOnly {
		args.AppendOnly = pulumi.BoolPtr(true)
	}

	// Destination tables (ForceNew; unique keys enable upserts)
	if len(dest.DestinationTables) > 0 {
		tables := kinesis.FirehoseDeliveryStreamIcebergConfigurationDestinationTableConfigurationArray{}
		for _, tbl := range dest.DestinationTables {
			tblArgs := &kinesis.FirehoseDeliveryStreamIcebergConfigurationDestinationTableConfigurationArgs{
				DatabaseName: pulumi.String(tbl.DatabaseName),
				TableName:    pulumi.String(tbl.TableName),
			}
			if tbl.S3ErrorOutputPrefix != "" {
				tblArgs.S3ErrorOutputPrefix = pulumi.StringPtr(tbl.S3ErrorOutputPrefix)
			}
			if len(tbl.UniqueKeys) > 0 {
				keys := pulumi.StringArray{}
				for _, k := range tbl.UniqueKeys {
					keys = append(keys, pulumi.String(k))
				}
				tblArgs.UniqueKeys = keys
			}
			tables = append(tables, tblArgs)
		}
		args.DestinationTableConfigurations = tables
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

	// S3 backup mode
	if dest.S3BackupMode != "" {
		args.S3BackupMode = pulumi.StringPtr(dest.S3BackupMode)
	}

	// S3 config (required)
	if cfg := dest.S3Config; cfg != nil {
		args.S3Configuration = buildIcebergS3Config(cfg)
	}

	// Processing pipeline (normalized typed processors)
	if proc := buildIcebergProcessing(dest.Processing); proc != nil {
		args.ProcessingConfiguration = proc
	}

	// CloudWatch logging
	if logging := dest.Logging; logging != nil && logging.Enabled {
		args.CloudwatchLoggingOptions = &kinesis.FirehoseDeliveryStreamIcebergConfigurationCloudwatchLoggingOptionsArgs{
			Enabled:       pulumi.Bool(true),
			LogGroupName:  pulumi.StringPtr(logging.LogGroupName),
			LogStreamName: pulumi.StringPtr(logging.LogStreamName),
		}
	}

	return args, nil
}

// buildIcebergS3Config builds the S3 backup configuration for Iceberg.
func buildIcebergS3Config(cfg *awskinesisfirehose.AwsKinesisFirehoseS3Config) *kinesis.FirehoseDeliveryStreamIcebergConfigurationS3ConfigurationArgs {
	args := &kinesis.FirehoseDeliveryStreamIcebergConfigurationS3ConfigurationArgs{
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

// buildIcebergProcessing adapts the normalized processor pipeline to the
// Iceberg destination's SDK types.
func buildIcebergProcessing(processing *awskinesisfirehose.AwsKinesisFirehoseProcessing) *kinesis.FirehoseDeliveryStreamIcebergConfigurationProcessingConfigurationArgs {
	normalized := normalizeProcessors(processing)
	if normalized == nil {
		return nil
	}
	processors := kinesis.FirehoseDeliveryStreamIcebergConfigurationProcessingConfigurationProcessorArray{}
	for _, p := range normalized {
		params := kinesis.FirehoseDeliveryStreamIcebergConfigurationProcessingConfigurationProcessorParameterArray{}
		for _, param := range p.Parameters {
			params = append(params, &kinesis.FirehoseDeliveryStreamIcebergConfigurationProcessingConfigurationProcessorParameterArgs{
				ParameterName:  pulumi.String(param.Name),
				ParameterValue: pulumi.String(param.Value),
			})
		}
		processors = append(processors, &kinesis.FirehoseDeliveryStreamIcebergConfigurationProcessingConfigurationProcessorArgs{
			Type:       pulumi.String(p.Type),
			Parameters: params,
		})
	}
	return &kinesis.FirehoseDeliveryStreamIcebergConfigurationProcessingConfigurationArgs{
		Enabled:    pulumi.Bool(true),
		Processors: processors,
	}
}
