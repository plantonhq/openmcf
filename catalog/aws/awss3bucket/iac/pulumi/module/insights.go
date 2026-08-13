package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// abac creates the attribute-based-access-control satellite when the spec
// states a posture: unset leaves the bucket at AWS's default (disabled)
// without managing the setting.
func abac(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec
	if spec.AbacStatus == "" {
		return nil
	}

	if _, err := s3.NewBucketAbac(ctx, "abac", &s3.BucketAbacArgs{
		Bucket: createdBucket.ID(),
		AbacStatus: &s3.BucketAbacAbacStatusArgs{
			Status: pulumi.String(spec.AbacStatus),
		},
	}, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "failed to configure abac")
	}
	return nil
}

// analyticsConfigurations creates one storage-class-analysis resource per
// named spec entry — a many-per-bucket satellite keyed by name, mirroring
// the Terraform module's for_each. AWS's export accepts only its V_1 schema
// and CSV format, so those provider arguments are left to their defaults and
// only the destination is user surface.
func analyticsConfigurations(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec

	for _, c := range spec.AnalyticsConfigurations {
		args := &s3.AnalyticsConfigurationArgs{
			Bucket: createdBucket.ID(),
			Name:   pulumi.String(c.Name),
		}

		if c.FilterPrefix != "" || len(c.FilterTags) > 0 {
			filter := &s3.AnalyticsConfigurationFilterArgs{}
			if c.FilterPrefix != "" {
				filter.Prefix = pulumi.StringPtr(c.FilterPrefix)
			}
			if len(c.FilterTags) > 0 {
				filter.Tags = pulumi.ToStringMap(c.FilterTags)
			}
			args.Filter = filter
		}

		if c.Export != nil {
			destination := &s3.AnalyticsConfigurationStorageClassAnalysisDataExportDestinationS3BucketDestinationArgs{
				BucketArn: pulumi.String(c.Export.BucketArn.GetValue()),
			}
			if c.Export.BucketAccountId != "" {
				destination.BucketAccountId = pulumi.StringPtr(c.Export.BucketAccountId)
			}
			if c.Export.Prefix != "" {
				destination.Prefix = pulumi.StringPtr(c.Export.Prefix)
			}
			args.StorageClassAnalysis = &s3.AnalyticsConfigurationStorageClassAnalysisArgs{
				DataExport: &s3.AnalyticsConfigurationStorageClassAnalysisDataExportArgs{
					Destination: &s3.AnalyticsConfigurationStorageClassAnalysisDataExportDestinationArgs{
						S3BucketDestination: destination,
					},
				},
			}
		}

		// The configuration name is the per-instance identity, mirroring the
		// Terraform module's for_each key.
		if _, err := s3.NewAnalyticsConfiguration(ctx, "analytics-"+c.Name, args, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "failed to configure analytics %q", c.Name)
		}
	}
	return nil
}

// inventoryConfigurations creates one inventory-report resource per named
// spec entry. `enabled` is always sent (AWS's InventoryConfiguration
// requires the IsEnabled member either way), derived from the spec's
// `disabled` so an unset spec matches AWS's active default.
func inventoryConfigurations(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec

	for _, c := range spec.InventoryConfigurations {
		args := &s3.InventoryArgs{
			Bucket:                 createdBucket.ID(),
			Name:                   pulumi.String(c.Name),
			Enabled:                pulumi.BoolPtr(!c.Disabled),
			IncludedObjectVersions: pulumi.String(c.IncludedObjectVersions),
			Schedule: &s3.InventoryScheduleArgs{
				Frequency: pulumi.String(c.Frequency),
			},
		}

		if len(c.OptionalFields) > 0 {
			args.OptionalFields = pulumi.ToStringArray(c.OptionalFields)
		}

		if c.FilterPrefix != "" {
			args.Filter = &s3.InventoryFilterArgs{
				Prefix: pulumi.StringPtr(c.FilterPrefix),
			}
		}

		destinationBucket := &s3.InventoryDestinationBucketArgs{
			BucketArn: pulumi.String(c.Destination.BucketArn.GetValue()),
			Format:    pulumi.String(c.Destination.Format),
		}
		if c.Destination.AccountId != "" {
			destinationBucket.AccountId = pulumi.StringPtr(c.Destination.AccountId)
		}
		if c.Destination.Prefix != "" {
			destinationBucket.Prefix = pulumi.StringPtr(c.Destination.Prefix)
		}
		// CEL guarantees at most one encryption arm; the block is emitted
		// only when an arm is actually chosen (an empty encryption block
		// would send an empty API struct).
		if c.Destination.SseKmsKeyId.GetValue() != "" {
			destinationBucket.Encryption = &s3.InventoryDestinationBucketEncryptionArgs{
				SseKms: &s3.InventoryDestinationBucketEncryptionSseKmsArgs{
					KeyId: pulumi.String(c.Destination.SseKmsKeyId.GetValue()),
				},
			}
		} else if c.Destination.SseS3 {
			destinationBucket.Encryption = &s3.InventoryDestinationBucketEncryptionArgs{
				SseS3: &s3.InventoryDestinationBucketEncryptionSseS3Args{},
			}
		}
		args.Destination = &s3.InventoryDestinationArgs{
			Bucket: destinationBucket,
		}

		if _, err := s3.NewInventory(ctx, "inventory-"+c.Name, args, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "failed to configure inventory %q", c.Name)
		}
	}
	return nil
}

// metricsConfigurations creates one CloudWatch request-metrics resource per
// named spec entry. No filter block means metrics for every request against
// the bucket.
func metricsConfigurations(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec

	for _, c := range spec.MetricsConfigurations {
		args := &s3.BucketMetricArgs{
			Bucket: createdBucket.ID(),
			Name:   pulumi.String(c.Name),
		}

		if c.FilterAccessPointArn != "" || c.FilterPrefix != "" || len(c.FilterTags) > 0 {
			filter := &s3.BucketMetricFilterArgs{}
			if c.FilterAccessPointArn != "" {
				filter.AccessPoint = pulumi.StringPtr(c.FilterAccessPointArn)
			}
			if c.FilterPrefix != "" {
				filter.Prefix = pulumi.StringPtr(c.FilterPrefix)
			}
			if len(c.FilterTags) > 0 {
				filter.Tags = pulumi.ToStringMap(c.FilterTags)
			}
			args.Filter = filter
		}

		if _, err := s3.NewBucketMetric(ctx, "metric-"+c.Name, args, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "failed to configure metrics %q", c.Name)
		}
	}
	return nil
}

// metadataConfiguration creates the S3 Metadata satellite when configured.
// The provider requires both table blocks: the journal table's expiration
// policy and the inventory table's state are stated explicitly either way.
// The destination table bucket/namespace are AWS-managed read-only state.
func metadataConfiguration(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec
	if spec.MetadataConfiguration == nil {
		return nil
	}

	inventoryState := "DISABLED"
	if spec.MetadataConfiguration.InventoryTableEnabled {
		inventoryState = "ENABLED"
	}
	inventoryTable := &s3.BucketMetadataConfigurationMetadataConfigurationInventoryTableConfigurationArgs{
		ConfigurationState: pulumi.String(inventoryState),
	}
	if enc := spec.MetadataConfiguration.InventoryTableEncryption; enc != nil {
		encArgs := &s3.BucketMetadataConfigurationMetadataConfigurationInventoryTableConfigurationEncryptionConfigurationArgs{
			SseAlgorithm: pulumi.String(enc.SseAlgorithm),
		}
		if enc.KmsKeyArn.GetValue() != "" {
			encArgs.KmsKeyArn = pulumi.StringPtr(enc.KmsKeyArn.GetValue())
		}
		inventoryTable.EncryptionConfiguration = encArgs
	}

	journalExpiration := "DISABLED"
	if spec.MetadataConfiguration.JournalRecordExpiration.Enabled {
		journalExpiration = "ENABLED"
	}
	recordExpiration := &s3.BucketMetadataConfigurationMetadataConfigurationJournalTableConfigurationRecordExpirationArgs{
		Expiration: pulumi.String(journalExpiration),
	}
	// days is only legal alongside ENABLED (CEL keeps it 0 when disabled).
	if spec.MetadataConfiguration.JournalRecordExpiration.Days > 0 {
		recordExpiration.Days = pulumi.IntPtr(int(spec.MetadataConfiguration.JournalRecordExpiration.Days))
	}
	journalTable := &s3.BucketMetadataConfigurationMetadataConfigurationJournalTableConfigurationArgs{
		RecordExpiration: recordExpiration,
	}
	if enc := spec.MetadataConfiguration.JournalTableEncryption; enc != nil {
		encArgs := &s3.BucketMetadataConfigurationMetadataConfigurationJournalTableConfigurationEncryptionConfigurationArgs{
			SseAlgorithm: pulumi.String(enc.SseAlgorithm),
		}
		if enc.KmsKeyArn.GetValue() != "" {
			encArgs.KmsKeyArn = pulumi.StringPtr(enc.KmsKeyArn.GetValue())
		}
		journalTable.EncryptionConfiguration = encArgs
	}

	if _, err := s3.NewBucketMetadataConfiguration(ctx, "metadata-configuration", &s3.BucketMetadataConfigurationArgs{
		Bucket: createdBucket.ID(),
		MetadataConfiguration: &s3.BucketMetadataConfigurationMetadataConfigurationArgs{
			InventoryTableConfiguration: inventoryTable,
			JournalTableConfiguration:   journalTable,
		},
	}, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "failed to configure metadata tables")
	}
	return nil
}
