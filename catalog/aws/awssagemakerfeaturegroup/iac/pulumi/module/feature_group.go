package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// featureGroup creates the feature group and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - almost everything is create-time only (schema, stores, role) -
//     the ONLY in-place updates are the online store's TTL and the
//     throughput settings;
//   - provisioned capacity units pair with Provisioned mode
//     (spec-validated), mirroring the provider's create behavior;
//   - a Vector collection pairs exactly with its dimension
//     (spec-validated; AWS requires InMemory online storage for
//     collection types - server-enforced).
func featureGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	var featureDefinitions sagemaker.FeatureGroupFeatureDefinitionArray
	for _, f := range spec.FeatureDefinitions {
		definition := &sagemaker.FeatureGroupFeatureDefinitionArgs{
			FeatureName: pulumi.String(f.Name),
			FeatureType: pulumi.String(f.Type),
		}
		if f.CollectionType != "" {
			definition.CollectionType = pulumi.String(f.CollectionType)
		}
		// The provider's collection_config has exactly one member
		// (vector_config) - rendered exactly when the dimension is set
		// (spec pairs it with the Vector type).
		if f.VectorDimension != nil {
			definition.CollectionConfig = &sagemaker.FeatureGroupFeatureDefinitionCollectionConfigArgs{
				VectorConfig: &sagemaker.FeatureGroupFeatureDefinitionCollectionConfigVectorConfigArgs{
					Dimension: pulumi.Int(int(*f.VectorDimension)),
				},
			}
		}
		featureDefinitions = append(featureDefinitions, definition)
	}

	args := &sagemaker.FeatureGroupArgs{
		// The component's name IS the feature group name.
		FeatureGroupName:            pulumi.String(locals.FeatureGroupName),
		RecordIdentifierFeatureName: pulumi.String(spec.RecordIdentifierFeatureName),
		EventTimeFeatureName:        pulumi.String(spec.EventTimeFeatureName),
		RoleArn:                     pulumi.String(spec.RoleArn.GetValue()),
		FeatureDefinitions:          featureDefinitions,
		Tags:                        pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	if spec.OnlineStore != nil {
		online := &sagemaker.FeatureGroupOnlineStoreConfigArgs{
			EnableOnlineStore: pulumi.Bool(spec.OnlineStore.Enabled),
		}
		if spec.OnlineStore.StorageType != "" {
			online.StorageType = pulumi.String(spec.OnlineStore.StorageType)
		}
		if spec.OnlineStore.KmsKeyArn.GetValue() != "" {
			online.SecurityConfig = &sagemaker.FeatureGroupOnlineStoreConfigSecurityConfigArgs{
				KmsKeyId: pulumi.String(spec.OnlineStore.KmsKeyArn.GetValue()),
			}
		}
		// The one online-store surface that updates in place.
		if spec.OnlineStore.Ttl != nil {
			online.TtlDuration = &sagemaker.FeatureGroupOnlineStoreConfigTtlDurationArgs{
				Unit:  pulumi.String(spec.OnlineStore.Ttl.Unit),
				Value: pulumi.Int(int(spec.OnlineStore.Ttl.Value)),
			}
		}
		args.OnlineStoreConfig = online
	}

	if spec.OfflineStore != nil {
		s3Storage := &sagemaker.FeatureGroupOfflineStoreConfigS3StorageConfigArgs{
			S3Uri: pulumi.String(spec.OfflineStore.S3Uri),
		}
		if spec.OfflineStore.KmsKeyArn.GetValue() != "" {
			s3Storage.KmsKeyId = pulumi.String(spec.OfflineStore.KmsKeyArn.GetValue())
		}
		offline := &sagemaker.FeatureGroupOfflineStoreConfigArgs{
			S3StorageConfig: s3Storage,
		}
		if spec.OfflineStore.DisableGlueTableCreation {
			offline.DisableGlueTableCreation = pulumi.Bool(true)
		}
		if spec.OfflineStore.TableFormat != "" {
			offline.TableFormat = pulumi.String(spec.OfflineStore.TableFormat)
		}
		if spec.OfflineStore.DataCatalog != nil {
			offline.DataCatalogConfig = &sagemaker.FeatureGroupOfflineStoreConfigDataCatalogConfigArgs{
				Catalog:   pulumi.String(spec.OfflineStore.DataCatalog.Catalog),
				Database:  pulumi.String(spec.OfflineStore.DataCatalog.Database),
				TableName: pulumi.String(spec.OfflineStore.DataCatalog.TableName),
			}
		}
		args.OfflineStoreConfig = offline
	}

	if spec.Throughput != nil {
		throughput := &sagemaker.FeatureGroupThroughputConfigArgs{
			ThroughputMode: pulumi.String(spec.Throughput.Mode),
		}
		if spec.Throughput.ProvisionedReadCapacityUnits != nil {
			throughput.ProvisionedReadCapacityUnits = pulumi.Int(int(*spec.Throughput.ProvisionedReadCapacityUnits))
		}
		if spec.Throughput.ProvisionedWriteCapacityUnits != nil {
			throughput.ProvisionedWriteCapacityUnits = pulumi.Int(int(*spec.Throughput.ProvisionedWriteCapacityUnits))
		}
		args.ThroughputConfig = throughput
	}

	createdFeatureGroup, err := sagemaker.NewFeatureGroup(ctx, locals.FeatureGroupName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create feature group")
	}

	ctx.Export(OpFeatureGroupName, createdFeatureGroup.FeatureGroupName)
	ctx.Export(OpFeatureGroupArn, createdFeatureGroup.Arn)

	return nil
}
