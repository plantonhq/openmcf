package module

import (
	"fmt"

	"github.com/pkg/errors"
	awss3tablebucketv1alpha1 "github.com/plantonhq/planton/catalog/aws/awss3tablebucket/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3tables"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// tableBucket creates the table bucket with its full contents -
// namespaces, Iceberg tables, resource policies, replication - and
// exports outputs.
//
// Lifecycle facts the render below depends on:
//   - namespaces key by name; tables key by "namespace.table" - both
//     are create-only at AWS (renames replace);
//   - a table's iceberg schema/properties are CREATE-ONLY input: the
//     provider never reads them back (schema evolution happens
//     through query engines), so they cannot drift and never
//     round-trip on import;
//   - the table Format argument is module-pinned to ICEBERG - the
//     provider's enum holds exactly that one value;
//   - policies are JSON-normalized by AWS; replication carries a
//     version_token optimistic-concurrency handshake the provider
//     manages.
func tableBucket(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	bucketArgs := &s3tables.TableBucketArgs{
		Name:         pulumi.String(locals.Target.Metadata.Name),
		ForceDestroy: pulumi.Bool(spec.ForceDestroy),
		Tags:         pulumi.ToStringMap(locals.AwsTags),
	}
	if encryption := spec.Encryption; encryption != nil {
		encryptionArgs := &s3tables.TableBucketEncryptionConfigurationArgs{
			SseAlgorithm: pulumi.String(encryption.SseAlgorithm),
		}
		if encryption.KmsKeyArn.GetValue() != "" {
			encryptionArgs.KmsKeyArn = pulumi.String(encryption.KmsKeyArn.GetValue())
		}
		bucketArgs.EncryptionConfiguration = encryptionArgs
	}
	if removal := spec.UnreferencedFileRemoval; removal != nil {
		status := "enabled"
		if removal.Disabled {
			status = "disabled"
		}
		settings := &s3tables.TableBucketMaintenanceConfigurationIcebergUnreferencedFileRemovalSettingsArgs{}
		if removal.NonCurrentDays > 0 {
			settings.NonCurrentDays = pulumi.Int(int(removal.NonCurrentDays))
		}
		if removal.UnreferencedDays > 0 {
			settings.UnreferencedDays = pulumi.Int(int(removal.UnreferencedDays))
		}
		bucketArgs.MaintenanceConfiguration = &s3tables.TableBucketMaintenanceConfigurationArgs{
			IcebergUnreferencedFileRemoval: &s3tables.TableBucketMaintenanceConfigurationIcebergUnreferencedFileRemovalArgs{
				Status:   pulumi.String(status),
				Settings: settings,
			},
		}
	}

	createdBucket, err := s3tables.NewTableBucket(ctx, "table_bucket", bucketArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create table bucket")
	}

	if spec.ResourcePolicy != "" {
		if _, err := s3tables.NewTableBucketPolicy(ctx, "bucket_policy", &s3tables.TableBucketPolicyArgs{
			TableBucketArn: createdBucket.Arn,
			ResourcePolicy: pulumi.String(spec.ResourcePolicy),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "bucket policy")
		}
	}

	if replication := spec.Replication; replication != nil {
		if _, err := s3tables.NewTableBucketReplication(ctx, "bucket_replication", &s3tables.TableBucketReplicationArgs{
			TableBucketArn: createdBucket.Arn,
			Role:           pulumi.String(replication.Role.GetValue()),
			Rule:           buildBucketReplicationRule(replication),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "bucket replication")
		}
	}

	tableArns := map[string]pulumi.StringOutput{}
	tableWarehouseLocations := map[string]pulumi.StringOutput{}

	for _, namespace := range spec.Namespaces {
		createdNamespace, err := s3tables.NewNamespace(ctx, "namespace-"+namespace.Name, &s3tables.NamespaceArgs{
			Namespace:      pulumi.String(namespace.Name),
			TableBucketArn: createdBucket.Arn,
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "namespace %s", namespace.Name)
		}

		for _, table := range namespace.Tables {
			// "namespace//table" - the import bridge's
			// address-key-segment convention; the output maps reuse
			// the key.
			tableKey := fmt.Sprintf("%s//%s", namespace.Name, table.Name)

			tableArgs := &s3tables.TableArgs{
				Name:           pulumi.String(table.Name),
				Namespace:      createdNamespace.Namespace,
				TableBucketArn: createdBucket.Arn,
				// The provider's enum holds exactly this one value.
				Format: pulumi.String("ICEBERG"),
				Tags:   pulumi.ToStringMap(locals.AwsTags),
			}

			if encryption := table.Encryption; encryption != nil {
				encryptionArgs := &s3tables.TableEncryptionConfigurationArgs{
					SseAlgorithm: pulumi.String(encryption.SseAlgorithm),
				}
				if encryption.KmsKeyArn.GetValue() != "" {
					encryptionArgs.KmsKeyArn = pulumi.String(encryption.KmsKeyArn.GetValue())
				}
				tableArgs.EncryptionConfiguration = encryptionArgs
			}

			if table.Compaction != nil || table.SnapshotManagement != nil {
				maintenance := &s3tables.TableMaintenanceConfigurationArgs{}
				if compaction := table.Compaction; compaction != nil {
					status := "enabled"
					if compaction.Disabled {
						status = "disabled"
					}
					settings := &s3tables.TableMaintenanceConfigurationIcebergCompactionSettingsArgs{}
					if compaction.TargetFileSizeMb > 0 {
						settings.TargetFileSizeMb = pulumi.Int(int(compaction.TargetFileSizeMb))
					}
					maintenance.IcebergCompaction = &s3tables.TableMaintenanceConfigurationIcebergCompactionArgs{
						Status:   pulumi.String(status),
						Settings: settings,
					}
				}
				if snapshotManagement := table.SnapshotManagement; snapshotManagement != nil {
					status := "enabled"
					if snapshotManagement.Disabled {
						status = "disabled"
					}
					settings := &s3tables.TableMaintenanceConfigurationIcebergSnapshotManagementSettingsArgs{}
					if snapshotManagement.MaxSnapshotAgeHours > 0 {
						settings.MaxSnapshotAgeHours = pulumi.Int(int(snapshotManagement.MaxSnapshotAgeHours))
					}
					if snapshotManagement.MinSnapshotsToKeep > 0 {
						settings.MinSnapshotsToKeep = pulumi.Int(int(snapshotManagement.MinSnapshotsToKeep))
					}
					maintenance.IcebergSnapshotManagement = &s3tables.TableMaintenanceConfigurationIcebergSnapshotManagementArgs{
						Status:   pulumi.String(status),
						Settings: settings,
					}
				}
				tableArgs.MaintenanceConfiguration = maintenance
			}

			if schema := table.IcebergSchema; schema != nil {
				fields := s3tables.TableMetadataIcebergSchemaFieldArray{}
				for _, field := range schema.Fields {
					fields = append(fields, &s3tables.TableMetadataIcebergSchemaFieldArgs{
						Name:     pulumi.String(field.Name),
						Type:     pulumi.String(field.Type),
						Required: pulumi.Bool(field.Required),
					})
				}
				iceberg := &s3tables.TableMetadataIcebergArgs{
					Schema: &s3tables.TableMetadataIcebergSchemaArgs{Fields: fields},
				}
				if len(table.Properties) > 0 {
					iceberg.Properties = pulumi.ToStringMap(table.Properties)
				}
				tableArgs.Metadata = &s3tables.TableMetadataArgs{Iceberg: iceberg}
			}

			createdTable, err := s3tables.NewTable(ctx, "table-"+tableKey, tableArgs, pulumi.Provider(provider))
			if err != nil {
				return errors.Wrapf(err, "table %s", tableKey)
			}
			tableArns[tableKey] = createdTable.Arn
			tableWarehouseLocations[tableKey] = createdTable.WarehouseLocation

			if table.ResourcePolicy != "" {
				if _, err := s3tables.NewTablePolicy(ctx, "table_policy-"+tableKey, &s3tables.TablePolicyArgs{
					Name:           createdTable.Name,
					Namespace:      createdNamespace.Namespace,
					TableBucketArn: createdBucket.Arn,
					ResourcePolicy: pulumi.String(table.ResourcePolicy),
				}, pulumi.Provider(provider)); err != nil {
					return errors.Wrapf(err, "table policy %s", tableKey)
				}
			}

			if replication := table.Replication; replication != nil {
				if _, err := s3tables.NewTableReplication(ctx, "table_replication-"+tableKey, &s3tables.TableReplicationArgs{
					TableArn: createdTable.Arn,
					Role:     pulumi.String(replication.Role.GetValue()),
					Rule:     buildTableReplicationRule(replication),
				}, pulumi.Provider(provider)); err != nil {
					return errors.Wrapf(err, "table replication %s", tableKey)
				}
			}
		}
	}

	ctx.Export(OpTableBucketArn, createdBucket.Arn)
	ctx.Export(OpOwnerAccountId, createdBucket.OwnerAccountId)
	ctx.Export(OpTableArns, stringOutputMap(tableArns))
	ctx.Export(OpTableWarehouseLocations, stringOutputMap(tableWarehouseLocations))
	return nil
}

func buildBucketReplicationRule(replication *awss3tablebucketv1alpha1.AwsS3TablesReplication) *s3tables.TableBucketReplicationRuleArgs {
	destinations := s3tables.TableBucketReplicationRuleDestinationArray{}
	for _, arn := range replication.DestinationTableBucketArns {
		destinations = append(destinations, &s3tables.TableBucketReplicationRuleDestinationArgs{
			DestinationTableBucketArn: pulumi.String(arn),
		})
	}
	return &s3tables.TableBucketReplicationRuleArgs{Destinations: destinations}
}

func buildTableReplicationRule(replication *awss3tablebucketv1alpha1.AwsS3TablesReplication) *s3tables.TableReplicationRuleArgs {
	destinations := s3tables.TableReplicationRuleDestinationArray{}
	for _, arn := range replication.DestinationTableBucketArns {
		destinations = append(destinations, &s3tables.TableReplicationRuleDestinationArgs{
			DestinationTableBucketArn: pulumi.String(arn),
		})
	}
	return &s3tables.TableReplicationRuleArgs{Destinations: destinations}
}

// stringOutputMap converts a Go map of outputs into an exportable
// pulumi map.
func stringOutputMap(in map[string]pulumi.StringOutput) pulumi.StringMap {
	out := pulumi.StringMap{}
	for key, value := range in {
		out[key] = value
	}
	return out
}
