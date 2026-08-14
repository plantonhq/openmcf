package module

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	awsbedrockknowledgebasev1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockknowledgebase/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// knowledgeBase creates the Bedrock knowledge base with its derived type
// discriminators, the vector store backend, and the folded data sources,
// and exports outputs.
//
// Nearly the whole surface is create-time only upstream; the provider
// retries the IAM/data-access propagation classes at create -- the module
// adds none.
func knowledgeBase(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &bedrock.AgentKnowledgeBaseArgs{
		// Create-time naming basis; doubles as the Name tag. metadata.name
		// on both engines.
		Name:    pulumi.String(locals.KnowledgeBaseName),
		RoleArn: pulumi.String(spec.RoleArn.GetValue()),
		Tags:    pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// -------------------------------------------------------------------
	// Knowledge-base type (the discriminator is derived from which spec
	// arm is set -- exactly one, per the spec's CEL guards)
	// -------------------------------------------------------------------
	kbConfiguration := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationArgs{}

	switch {
	case spec.Vector != nil:
		kbConfiguration.Type = pulumi.String("VECTOR")
		vector := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationArgs{
			EmbeddingModelArn: pulumi.String(spec.Vector.EmbeddingModelArn),
		}
		if spec.Vector.EmbeddingModel != nil {
			m := spec.Vector.EmbeddingModel
			bedrockModel := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationArgs{}
			if m.Dimensions != 0 {
				bedrockModel.Dimensions = pulumi.Int(int(m.Dimensions))
			}
			if m.EmbeddingDataType != "" {
				bedrockModel.EmbeddingDataType = pulumi.String(m.EmbeddingDataType)
			}
			if m.AudioSegmentationSeconds != 0 {
				bedrockModel.Audio = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationAudioArgs{
					SegmentationConfiguration: &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationAudioSegmentationConfigurationArgs{
						FixedLengthDuration: pulumi.Int(int(m.AudioSegmentationSeconds)),
					},
				}
			}
			if m.VideoSegmentationSeconds != 0 {
				bedrockModel.Video = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationVideoArgs{
					SegmentationConfiguration: &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationVideoSegmentationConfigurationArgs{
						FixedLengthDuration: pulumi.Int(int(m.VideoSegmentationSeconds)),
					},
				}
			}
			vector.EmbeddingModelConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationEmbeddingModelConfigurationArgs{
				BedrockEmbeddingModelConfiguration: bedrockModel,
			}
		}
		// The supplemental S3 location rides a fixed two-level wrapper
		// upstream (storage_location type S3 + s3_location.uri); the spec
		// carries the URI leaf directly.
		if spec.Vector.SupplementalDataS3Uri != "" {
			vector.SupplementalDataStorageConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationSupplementalDataStorageConfigurationArgs{
				StorageLocations: bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationSupplementalDataStorageConfigurationStorageLocationArray{
					&bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationSupplementalDataStorageConfigurationStorageLocationArgs{
						Type: pulumi.String("S3"),
						S3Location: &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationVectorKnowledgeBaseConfigurationSupplementalDataStorageConfigurationStorageLocationS3LocationArgs{
							Uri: pulumi.String(spec.Vector.SupplementalDataS3Uri),
						},
					},
				},
			}
		}
		kbConfiguration.VectorKnowledgeBaseConfiguration = vector

	case spec.Managed != nil:
		kbConfiguration.Type = pulumi.String("MANAGED")
		managed := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationManagedKnowledgeBaseConfigurationArgs{}
		if spec.Managed.EmbeddingModelArn != "" {
			managed.EmbeddingModelArn = pulumi.String(spec.Managed.EmbeddingModelArn)
		}
		// Derived discriminator, ALWAYS sent: AWS's embeddingModelType is
		// CUSTOM exactly when an embedding-model ARN is brought, MANAGED
		// otherwise. The provider marks the attribute Optional+Computed
		// (UseStateForUnknown); leaving it unset makes the bridge fail the
		// apply with "unexpected unknown property value" AFTER AWS created
		// the knowledge base -- stranding it outside state (live-caught
		// 2026-08-13). Sending the derived value keeps it known at plan.
		if spec.Managed.EmbeddingModelArn != "" {
			managed.EmbeddingModelType = pulumi.String("CUSTOM")
		} else {
			managed.EmbeddingModelType = pulumi.String("MANAGED")
		}
		if spec.Managed.EmbeddingModel != nil {
			m := spec.Managed.EmbeddingModel
			bedrockModel := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationManagedKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationArgs{}
			if m.Dimensions != 0 {
				bedrockModel.Dimensions = pulumi.Int(int(m.Dimensions))
			}
			if m.EmbeddingDataType != "" {
				bedrockModel.EmbeddingDataType = pulumi.String(m.EmbeddingDataType)
			}
			if m.AudioSegmentationSeconds != 0 {
				bedrockModel.Audio = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationManagedKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationAudioArgs{
					SegmentationConfiguration: &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationManagedKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationAudioSegmentationConfigurationArgs{
						FixedLengthDuration: pulumi.Int(int(m.AudioSegmentationSeconds)),
					},
				}
			}
			if m.VideoSegmentationSeconds != 0 {
				bedrockModel.Video = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationManagedKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationVideoArgs{
					SegmentationConfiguration: &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationManagedKnowledgeBaseConfigurationEmbeddingModelConfigurationBedrockEmbeddingModelConfigurationVideoSegmentationConfigurationArgs{
						FixedLengthDuration: pulumi.Int(int(m.VideoSegmentationSeconds)),
					},
				}
			}
			managed.EmbeddingModelConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationManagedKnowledgeBaseConfigurationEmbeddingModelConfigurationArgs{
				BedrockEmbeddingModelConfiguration: bedrockModel,
			}
		}
		if spec.Managed.KmsKeyArn.GetValue() != "" {
			managed.ServerSideEncryptionConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationManagedKnowledgeBaseConfigurationServerSideEncryptionConfigurationArgs{
				KmsKeyArn: pulumi.String(spec.Managed.KmsKeyArn.GetValue()),
			}
		}
		kbConfiguration.ManagedKnowledgeBaseConfiguration = managed

	case spec.Kendra != nil:
		kbConfiguration.Type = pulumi.String("KENDRA")
		kbConfiguration.KendraKnowledgeBaseConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationKendraKnowledgeBaseConfigurationArgs{
			KendraIndexArn: pulumi.String(spec.Kendra.KendraIndexArn),
		}

	case spec.Sql != nil:
		kbConfiguration.Type = pulumi.String("SQL")
		// REDSHIFT is the only SQL engine AWS defines -- the module owns
		// the constant.
		redshift := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationArgs{}

		queryEngine := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryEngineConfigurationArgs{}
		if spec.Sql.Provisioned != nil {
			queryEngine.Type = pulumi.String("PROVISIONED")
			provisionedAuth := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryEngineConfigurationProvisionedConfigurationAuthConfigurationArgs{
				Type: pulumi.String(spec.Sql.Provisioned.Auth.Type),
			}
			if spec.Sql.Provisioned.Auth.DatabaseUser != "" {
				provisionedAuth.DatabaseUser = pulumi.String(spec.Sql.Provisioned.Auth.DatabaseUser)
			}
			if spec.Sql.Provisioned.Auth.UsernamePasswordSecretArn.GetValue() != "" {
				provisionedAuth.UsernamePasswordSecretArn = pulumi.String(spec.Sql.Provisioned.Auth.UsernamePasswordSecretArn.GetValue())
			}
			queryEngine.ProvisionedConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryEngineConfigurationProvisionedConfigurationArgs{
				ClusterIdentifier: pulumi.String(spec.Sql.Provisioned.ClusterIdentifier.GetValue()),
				AuthConfiguration: provisionedAuth,
			}
		}
		if spec.Sql.Serverless != nil {
			queryEngine.Type = pulumi.String("SERVERLESS")
			serverlessAuth := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryEngineConfigurationServerlessConfigurationAuthConfigurationArgs{
				Type: pulumi.String(spec.Sql.Serverless.Auth.Type),
			}
			if spec.Sql.Serverless.Auth.UsernamePasswordSecretArn.GetValue() != "" {
				serverlessAuth.UsernamePasswordSecretArn = pulumi.String(spec.Sql.Serverless.Auth.UsernamePasswordSecretArn.GetValue())
			}
			queryEngine.ServerlessConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryEngineConfigurationServerlessConfigurationArgs{
				WorkgroupArn:      pulumi.String(spec.Sql.Serverless.WorkgroupArn.GetValue()),
				AuthConfiguration: serverlessAuth,
			}
		}
		redshift.QueryEngineConfiguration = queryEngine

		warehouse := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationStorageConfigurationArgs{}
		if spec.Sql.Warehouse.DataCatalog != nil {
			warehouse.Type = pulumi.String("AWS_DATA_CATALOG")
			warehouse.AwsDataCatalogConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationStorageConfigurationAwsDataCatalogConfigurationArgs{
				TableNames: pulumi.ToStringArray(spec.Sql.Warehouse.DataCatalog.TableNames),
			}
		}
		if spec.Sql.Warehouse.Redshift != nil {
			warehouse.Type = pulumi.String("REDSHIFT")
			warehouse.RedshiftConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationStorageConfigurationRedshiftConfigurationArgs{
				DatabaseName: pulumi.String(spec.Sql.Warehouse.Redshift.DatabaseName),
			}
		}
		redshift.StorageConfiguration = warehouse

		if spec.Sql.QueryGeneration != nil {
			qg := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryGenerationConfigurationArgs{}
			if spec.Sql.QueryGeneration.ExecutionTimeoutSeconds != 0 {
				qg.ExecutionTimeoutSeconds = pulumi.Int(int(spec.Sql.QueryGeneration.ExecutionTimeoutSeconds))
			}
			if len(spec.Sql.QueryGeneration.CuratedQueries) > 0 || len(spec.Sql.QueryGeneration.Tables) > 0 {
				genContext := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryGenerationConfigurationGenerationContextArgs{}
				var curated bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryGenerationConfigurationGenerationContextCuratedQueryArray
				for _, q := range spec.Sql.QueryGeneration.CuratedQueries {
					curated = append(curated, &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryGenerationConfigurationGenerationContextCuratedQueryArgs{
						NaturalLanguage: pulumi.String(q.NaturalLanguage),
						Sql:             pulumi.String(q.Sql),
					})
				}
				genContext.CuratedQueries = curated
				var tables bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryGenerationConfigurationGenerationContextTableArray
				for _, t := range spec.Sql.QueryGeneration.Tables {
					table := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryGenerationConfigurationGenerationContextTableArgs{
						Name: pulumi.String(t.Name),
					}
					if t.Description != "" {
						table.Description = pulumi.String(t.Description)
					}
					if t.Inclusion != "" {
						table.Inclusion = pulumi.String(t.Inclusion)
					}
					var columns bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryGenerationConfigurationGenerationContextTableColumnArray
					for _, c := range t.Columns {
						column := &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationRedshiftConfigurationQueryGenerationConfigurationGenerationContextTableColumnArgs{
							Name: pulumi.String(c.Name),
						}
						if c.Description != "" {
							column.Description = pulumi.String(c.Description)
						}
						if c.Inclusion != "" {
							column.Inclusion = pulumi.String(c.Inclusion)
						}
						columns = append(columns, column)
					}
					table.Columns = columns
					tables = append(tables, table)
				}
				genContext.Tables = tables
				qg.GenerationContext = genContext
			}
			redshift.QueryGenerationConfiguration = qg
		}

		kbConfiguration.SqlKnowledgeBaseConfiguration = &bedrock.AgentKnowledgeBaseKnowledgeBaseConfigurationSqlKnowledgeBaseConfigurationArgs{
			Type:                  pulumi.String("REDSHIFT"),
			RedshiftConfiguration: redshift,
		}
	}
	args.KnowledgeBaseConfiguration = kbConfiguration

	// -------------------------------------------------------------------
	// Vector store (required with the vector type, absent otherwise --
	// the spec's CEL guards enforce the pairing)
	// -------------------------------------------------------------------
	if spec.Storage != nil {
		storage := &bedrock.AgentKnowledgeBaseStorageConfigurationArgs{}
		switch {
		case spec.Storage.OpensearchServerless != nil:
			s := spec.Storage.OpensearchServerless
			storage.Type = pulumi.String("OPENSEARCH_SERVERLESS")
			storage.OpensearchServerlessConfiguration = &bedrock.AgentKnowledgeBaseStorageConfigurationOpensearchServerlessConfigurationArgs{
				CollectionArn:   pulumi.String(s.CollectionArn.GetValue()),
				VectorIndexName: pulumi.String(s.VectorIndexName),
				FieldMapping: &bedrock.AgentKnowledgeBaseStorageConfigurationOpensearchServerlessConfigurationFieldMappingArgs{
					VectorField:   pulumi.String(s.FieldMapping.VectorField),
					TextField:     pulumi.String(s.FieldMapping.TextField),
					MetadataField: pulumi.String(s.FieldMapping.MetadataField),
				},
			}
		case spec.Storage.OpensearchManaged != nil:
			s := spec.Storage.OpensearchManaged
			storage.Type = pulumi.String("OPENSEARCH_MANAGED_CLUSTER")
			storage.OpensearchManagedClusterConfiguration = &bedrock.AgentKnowledgeBaseStorageConfigurationOpensearchManagedClusterConfigurationArgs{
				DomainArn:       pulumi.String(s.DomainArn.GetValue()),
				DomainEndpoint:  pulumi.String(s.DomainEndpoint),
				VectorIndexName: pulumi.String(s.VectorIndexName),
				FieldMapping: &bedrock.AgentKnowledgeBaseStorageConfigurationOpensearchManagedClusterConfigurationFieldMappingArgs{
					VectorField:   pulumi.String(s.FieldMapping.VectorField),
					TextField:     pulumi.String(s.FieldMapping.TextField),
					MetadataField: pulumi.String(s.FieldMapping.MetadataField),
				},
			}
		case spec.Storage.S3Vectors != nil:
			s := spec.Storage.S3Vectors
			storage.Type = pulumi.String("S3_VECTORS")
			s3Vectors := &bedrock.AgentKnowledgeBaseStorageConfigurationS3VectorsConfigurationArgs{}
			if s.IndexArn != "" {
				s3Vectors.IndexArn = pulumi.String(s.IndexArn)
			}
			if s.IndexName != "" {
				s3Vectors.IndexName = pulumi.String(s.IndexName)
			}
			if s.VectorBucketArn != "" {
				s3Vectors.VectorBucketArn = pulumi.String(s.VectorBucketArn)
			}
			storage.S3VectorsConfiguration = s3Vectors
		case spec.Storage.Rds != nil:
			s := spec.Storage.Rds
			storage.Type = pulumi.String("RDS")
			rdsMapping := &bedrock.AgentKnowledgeBaseStorageConfigurationRdsConfigurationFieldMappingArgs{
				VectorField:     pulumi.String(s.FieldMapping.VectorField),
				TextField:       pulumi.String(s.FieldMapping.TextField),
				MetadataField:   pulumi.String(s.FieldMapping.MetadataField),
				PrimaryKeyField: pulumi.String(s.FieldMapping.PrimaryKeyField),
			}
			if s.FieldMapping.CustomMetadataField != "" {
				rdsMapping.CustomMetadataField = pulumi.String(s.FieldMapping.CustomMetadataField)
			}
			storage.RdsConfiguration = &bedrock.AgentKnowledgeBaseStorageConfigurationRdsConfigurationArgs{
				ResourceArn:          pulumi.String(s.ResourceArn.GetValue()),
				CredentialsSecretArn: pulumi.String(s.CredentialsSecretArn.GetValue()),
				DatabaseName:         pulumi.String(s.DatabaseName),
				TableName:            pulumi.String(s.TableName),
				FieldMapping:         rdsMapping,
			}
		case spec.Storage.Pinecone != nil:
			s := spec.Storage.Pinecone
			storage.Type = pulumi.String("PINECONE")
			pinecone := &bedrock.AgentKnowledgeBaseStorageConfigurationPineconeConfigurationArgs{
				ConnectionString:     pulumi.String(s.ConnectionString),
				CredentialsSecretArn: pulumi.String(s.CredentialsSecretArn.GetValue()),
				FieldMapping: &bedrock.AgentKnowledgeBaseStorageConfigurationPineconeConfigurationFieldMappingArgs{
					TextField:     pulumi.String(s.FieldMapping.TextField),
					MetadataField: pulumi.String(s.FieldMapping.MetadataField),
				},
			}
			if s.Namespace != "" {
				pinecone.Namespace = pulumi.String(s.Namespace)
			}
			storage.PineconeConfiguration = pinecone
		case spec.Storage.MongodbAtlas != nil:
			s := spec.Storage.MongodbAtlas
			storage.Type = pulumi.String("MONGO_DB_ATLAS")
			mongo := &bedrock.AgentKnowledgeBaseStorageConfigurationMongoDbAtlasConfigurationArgs{
				Endpoint:             pulumi.String(s.Endpoint),
				DatabaseName:         pulumi.String(s.DatabaseName),
				CollectionName:       pulumi.String(s.CollectionName),
				VectorIndexName:      pulumi.String(s.VectorIndexName),
				CredentialsSecretArn: pulumi.String(s.CredentialsSecretArn.GetValue()),
				FieldMapping: &bedrock.AgentKnowledgeBaseStorageConfigurationMongoDbAtlasConfigurationFieldMappingArgs{
					VectorField:   pulumi.String(s.FieldMapping.VectorField),
					TextField:     pulumi.String(s.FieldMapping.TextField),
					MetadataField: pulumi.String(s.FieldMapping.MetadataField),
				},
			}
			if s.TextIndexName != "" {
				mongo.TextIndexName = pulumi.String(s.TextIndexName)
			}
			if s.EndpointServiceName != "" {
				mongo.EndpointServiceName = pulumi.String(s.EndpointServiceName)
			}
			storage.MongoDbAtlasConfiguration = mongo
		case spec.Storage.NeptuneAnalytics != nil:
			s := spec.Storage.NeptuneAnalytics
			storage.Type = pulumi.String("NEPTUNE_ANALYTICS")
			storage.NeptuneAnalyticsConfiguration = &bedrock.AgentKnowledgeBaseStorageConfigurationNeptuneAnalyticsConfigurationArgs{
				GraphArn: pulumi.String(s.GraphArn),
				FieldMapping: &bedrock.AgentKnowledgeBaseStorageConfigurationNeptuneAnalyticsConfigurationFieldMappingArgs{
					TextField:     pulumi.String(s.FieldMapping.TextField),
					MetadataField: pulumi.String(s.FieldMapping.MetadataField),
				},
			}
		case spec.Storage.RedisEnterpriseCloud != nil:
			s := spec.Storage.RedisEnterpriseCloud
			storage.Type = pulumi.String("REDIS_ENTERPRISE_CLOUD")
			redis := &bedrock.AgentKnowledgeBaseStorageConfigurationRedisEnterpriseCloudConfigurationArgs{
				Endpoint:             pulumi.String(s.Endpoint),
				VectorIndexName:      pulumi.String(s.VectorIndexName),
				CredentialsSecretArn: pulumi.String(s.CredentialsSecretArn.GetValue()),
			}
			if s.FieldMapping != nil {
				redisMapping := &bedrock.AgentKnowledgeBaseStorageConfigurationRedisEnterpriseCloudConfigurationFieldMappingArgs{}
				if s.FieldMapping.VectorField != "" {
					redisMapping.VectorField = pulumi.String(s.FieldMapping.VectorField)
				}
				if s.FieldMapping.TextField != "" {
					redisMapping.TextField = pulumi.String(s.FieldMapping.TextField)
				}
				if s.FieldMapping.MetadataField != "" {
					redisMapping.MetadataField = pulumi.String(s.FieldMapping.MetadataField)
				}
				redis.FieldMapping = redisMapping
			}
			storage.RedisEnterpriseCloudConfiguration = redis
		}
		args.StorageConfiguration = storage
	}

	createdKnowledgeBase, err := bedrock.NewAgentKnowledgeBase(ctx, locals.KnowledgeBaseName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create knowledge base")
	}

	ctx.Export(OpKnowledgeBaseId, createdKnowledgeBase.ID())
	ctx.Export(OpKnowledgeBaseArn, createdKnowledgeBase.Arn)

	// Document connectors keyed by their stable entry names. Iteration is
	// name-sorted for deterministic previews.
	dataSourceIds := pulumi.StringMap{}
	for _, d := range sortedDataSources(spec.DataSources) {
		dataSourceArgs, err := dataSourceArgs(createdKnowledgeBase, d)
		if err != nil {
			return errors.Wrapf(err, "render data source %q", d.Name)
		}
		createdDataSource, err := bedrock.NewAgentDataSource(ctx, "data-source-"+d.Name, dataSourceArgs,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdKnowledgeBase}))
		if err != nil {
			return errors.Wrapf(err, "create data source %q", d.Name)
		}
		dataSourceIds[d.Name] = createdDataSource.DataSourceId
	}
	ctx.Export(OpDataSourceIds, dataSourceIds)

	return nil
}

// dataSourceArgs renders one data source entry -- the connector arm plus
// the ingestion pipeline.
func dataSourceArgs(kb *bedrock.AgentKnowledgeBase, d *awsbedrockknowledgebasev1alpha1.AwsBedrockKnowledgeBaseDataSource) (*bedrock.AgentDataSourceArgs, error) {
	args := &bedrock.AgentDataSourceArgs{
		Name:            pulumi.String(d.Name),
		KnowledgeBaseId: kb.ID(),
	}
	if d.Description != "" {
		args.Description = pulumi.String(d.Description)
	}
	if d.DataDeletionPolicy != "" {
		args.DataDeletionPolicy = pulumi.String(d.DataDeletionPolicy)
	}
	if d.KmsKeyArn.GetValue() != "" {
		args.ServerSideEncryptionConfiguration = &bedrock.AgentDataSourceServerSideEncryptionConfigurationArgs{
			KmsKeyArn: pulumi.String(d.KmsKeyArn.GetValue()),
		}
	}

	configuration := &bedrock.AgentDataSourceDataSourceConfigurationArgs{}
	switch {
	case d.S3 != nil:
		configuration.Type = pulumi.String("S3")
		s3 := &bedrock.AgentDataSourceDataSourceConfigurationS3ConfigurationArgs{
			BucketArn: pulumi.String(d.S3.BucketArn.GetValue()),
		}
		if d.S3.InclusionPrefix != "" {
			s3.InclusionPrefixes = pulumi.StringArray{pulumi.String(d.S3.InclusionPrefix)}
		}
		if d.S3.BucketOwnerAccountId != "" {
			s3.BucketOwnerAccountId = pulumi.String(d.S3.BucketOwnerAccountId)
		}
		configuration.S3Configuration = s3

	case d.Web != nil:
		configuration.Type = pulumi.String("WEB")
		var seedUrls bedrock.AgentDataSourceDataSourceConfigurationWebConfigurationSourceConfigurationUrlConfigurationSeedUrlArray
		for _, u := range d.Web.SeedUrls {
			seedUrls = append(seedUrls, &bedrock.AgentDataSourceDataSourceConfigurationWebConfigurationSourceConfigurationUrlConfigurationSeedUrlArgs{
				Url: pulumi.String(u),
			})
		}
		web := &bedrock.AgentDataSourceDataSourceConfigurationWebConfigurationArgs{
			SourceConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationWebConfigurationSourceConfigurationArgs{
				UrlConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationWebConfigurationSourceConfigurationUrlConfigurationArgs{
					SeedUrls: seedUrls,
				},
			},
		}
		crawler := &bedrock.AgentDataSourceDataSourceConfigurationWebConfigurationCrawlerConfigurationArgs{}
		if d.Web.Scope != "" {
			crawler.Scope = pulumi.String(d.Web.Scope)
		}
		if len(d.Web.InclusionFilters) > 0 {
			crawler.InclusionFilters = pulumi.ToStringArray(d.Web.InclusionFilters)
		}
		if len(d.Web.ExclusionFilters) > 0 {
			crawler.ExclusionFilters = pulumi.ToStringArray(d.Web.ExclusionFilters)
		}
		if d.Web.UserAgent != "" {
			crawler.UserAgent = pulumi.String(d.Web.UserAgent)
		}
		if d.Web.MaxPages != 0 || d.Web.RateLimit != 0 {
			limits := &bedrock.AgentDataSourceDataSourceConfigurationWebConfigurationCrawlerConfigurationCrawlerLimitsArgs{}
			if d.Web.MaxPages != 0 {
				limits.MaxPages = pulumi.Int(int(d.Web.MaxPages))
			}
			if d.Web.RateLimit != 0 {
				limits.RateLimit = pulumi.Int(int(d.Web.RateLimit))
			}
			crawler.CrawlerLimits = limits
		}
		web.CrawlerConfiguration = crawler
		configuration.WebConfiguration = web

	case d.Confluence != nil:
		configuration.Type = pulumi.String("CONFLUENCE")
		confluence := &bedrock.AgentDataSourceDataSourceConfigurationConfluenceConfigurationArgs{
			SourceConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationConfluenceConfigurationSourceConfigurationArgs{
				// SAAS is the only Confluence host type AWS defines -- the
				// module owns the constant.
				HostType:             pulumi.String("SAAS"),
				HostUrl:              pulumi.String(d.Confluence.HostUrl),
				AuthType:             pulumi.String(d.Confluence.AuthType),
				CredentialsSecretArn: pulumi.String(d.Confluence.CredentialsSecretArn.GetValue()),
			},
		}
		if len(d.Confluence.Filters) > 0 {
			var filters bedrock.AgentDataSourceDataSourceConfigurationConfluenceConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterFilterArray
			for _, f := range d.Confluence.Filters {
				filter := &bedrock.AgentDataSourceDataSourceConfigurationConfluenceConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterFilterArgs{
					ObjectType: pulumi.String(f.ObjectType),
				}
				if len(f.InclusionFilters) > 0 {
					filter.InclusionFilters = pulumi.ToStringArray(f.InclusionFilters)
				}
				if len(f.ExclusionFilters) > 0 {
					filter.ExclusionFilters = pulumi.ToStringArray(f.ExclusionFilters)
				}
				filters = append(filters, filter)
			}
			confluence.CrawlerConfiguration = &bedrock.AgentDataSourceDataSourceConfigurationConfluenceConfigurationCrawlerConfigurationArgs{
				FilterConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationConfluenceConfigurationCrawlerConfigurationFilterConfigurationArgs{
					// PATTERN is the only filter type AWS defines -- the
					// module owns the constant.
					Type: pulumi.String("PATTERN"),
					PatternObjectFilters: bedrock.AgentDataSourceDataSourceConfigurationConfluenceConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterArray{
						&bedrock.AgentDataSourceDataSourceConfigurationConfluenceConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterArgs{
							Filters: filters,
						},
					},
				},
			}
		}
		configuration.ConfluenceConfiguration = confluence

	case d.Salesforce != nil:
		configuration.Type = pulumi.String("SALESFORCE")
		salesforce := &bedrock.AgentDataSourceDataSourceConfigurationSalesforceConfigurationArgs{
			SourceConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationSalesforceConfigurationSourceConfigurationArgs{
				// OAUTH2_CLIENT_CREDENTIALS is the only Salesforce auth
				// type AWS defines -- the module owns the constant.
				AuthType:             pulumi.String("OAUTH2_CLIENT_CREDENTIALS"),
				HostUrl:              pulumi.String(d.Salesforce.HostUrl),
				CredentialsSecretArn: pulumi.String(d.Salesforce.CredentialsSecretArn.GetValue()),
			},
		}
		if len(d.Salesforce.Filters) > 0 {
			var filters bedrock.AgentDataSourceDataSourceConfigurationSalesforceConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterFilterArray
			for _, f := range d.Salesforce.Filters {
				filter := &bedrock.AgentDataSourceDataSourceConfigurationSalesforceConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterFilterArgs{
					ObjectType: pulumi.String(f.ObjectType),
				}
				if len(f.InclusionFilters) > 0 {
					filter.InclusionFilters = pulumi.ToStringArray(f.InclusionFilters)
				}
				if len(f.ExclusionFilters) > 0 {
					filter.ExclusionFilters = pulumi.ToStringArray(f.ExclusionFilters)
				}
				filters = append(filters, filter)
			}
			salesforce.CrawlerConfiguration = &bedrock.AgentDataSourceDataSourceConfigurationSalesforceConfigurationCrawlerConfigurationArgs{
				FilterConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationSalesforceConfigurationCrawlerConfigurationFilterConfigurationArgs{
					Type: pulumi.String("PATTERN"),
					PatternObjectFilters: bedrock.AgentDataSourceDataSourceConfigurationSalesforceConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterArray{
						&bedrock.AgentDataSourceDataSourceConfigurationSalesforceConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterArgs{
							Filters: filters,
						},
					},
				},
			}
		}
		configuration.SalesforceConfiguration = salesforce

	case d.Sharepoint != nil:
		configuration.Type = pulumi.String("SHAREPOINT")
		sourceConfiguration := &bedrock.AgentDataSourceDataSourceConfigurationSharePointConfigurationSourceConfigurationArgs{
			// ONLINE is the only SharePoint host type AWS defines -- the
			// module owns the constant.
			HostType:             pulumi.String("ONLINE"),
			SiteUrls:             pulumi.ToStringArray(d.Sharepoint.SiteUrls),
			Domain:               pulumi.String(d.Sharepoint.Domain),
			AuthType:             pulumi.String(d.Sharepoint.AuthType),
			CredentialsSecretArn: pulumi.String(d.Sharepoint.CredentialsSecretArn.GetValue()),
		}
		if d.Sharepoint.TenantId != "" {
			sourceConfiguration.TenantId = pulumi.String(d.Sharepoint.TenantId)
		}
		sharepoint := &bedrock.AgentDataSourceDataSourceConfigurationSharePointConfigurationArgs{
			SourceConfiguration: sourceConfiguration,
		}
		if len(d.Sharepoint.Filters) > 0 {
			var filters bedrock.AgentDataSourceDataSourceConfigurationSharePointConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterFilterArray
			for _, f := range d.Sharepoint.Filters {
				filter := &bedrock.AgentDataSourceDataSourceConfigurationSharePointConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterFilterArgs{
					ObjectType: pulumi.String(f.ObjectType),
				}
				if len(f.InclusionFilters) > 0 {
					filter.InclusionFilters = pulumi.ToStringArray(f.InclusionFilters)
				}
				if len(f.ExclusionFilters) > 0 {
					filter.ExclusionFilters = pulumi.ToStringArray(f.ExclusionFilters)
				}
				filters = append(filters, filter)
			}
			sharepoint.CrawlerConfiguration = &bedrock.AgentDataSourceDataSourceConfigurationSharePointConfigurationCrawlerConfigurationArgs{
				FilterConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationSharePointConfigurationCrawlerConfigurationFilterConfigurationArgs{
					Type: pulumi.String("PATTERN"),
					PatternObjectFilters: bedrock.AgentDataSourceDataSourceConfigurationSharePointConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterArray{
						&bedrock.AgentDataSourceDataSourceConfigurationSharePointConfigurationCrawlerConfigurationFilterConfigurationPatternObjectFilterArgs{
							Filters: filters,
						},
					},
				},
			}
		}
		configuration.SharePointConfiguration = sharepoint

	case d.ManagedConnector != nil:
		configuration.Type = pulumi.String("MANAGED_KNOWLEDGE_BASE_CONNECTOR")
		managed := &bedrock.AgentDataSourceDataSourceConfigurationManagedKnowledgeBaseConnectorConfigurationArgs{}
		if d.ManagedConnector.ConnectorParameters != nil {
			parametersJson, err := json.Marshal(d.ManagedConnector.ConnectorParameters.AsMap())
			if err != nil {
				return nil, errors.Wrap(err, "marshal connector parameters")
			}
			managed.ConnectorParameters = pulumi.String(string(parametersJson))
		}
		if d.ManagedConnector.DeletionProtection != nil {
			protection := &bedrock.AgentDataSourceDataSourceConfigurationManagedKnowledgeBaseConnectorConfigurationDeletionProtectionConfigurationArgs{
				DeletionProtectionStatus: pulumi.String(enabledOrDisabled(d.ManagedConnector.DeletionProtection.Enabled)),
			}
			if d.ManagedConnector.DeletionProtection.ThresholdPercent != 0 {
				protection.DeletionProtectionThreshold = pulumi.Int(int(d.ManagedConnector.DeletionProtection.ThresholdPercent))
			}
			managed.DeletionProtectionConfiguration = protection
		}
		if d.ManagedConnector.MediaExtraction != nil {
			managed.MediaExtractionConfiguration = &bedrock.AgentDataSourceDataSourceConfigurationManagedKnowledgeBaseConnectorConfigurationMediaExtractionConfigurationArgs{
				AudioExtractionConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationManagedKnowledgeBaseConnectorConfigurationMediaExtractionConfigurationAudioExtractionConfigurationArgs{
					AudioExtractionStatus: pulumi.String(enabledOrDisabled(d.ManagedConnector.MediaExtraction.Audio)),
				},
				ImageExtractionConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationManagedKnowledgeBaseConnectorConfigurationMediaExtractionConfigurationImageExtractionConfigurationArgs{
					ImageExtractionStatus: pulumi.String(enabledOrDisabled(d.ManagedConnector.MediaExtraction.Image)),
				},
				VideoExtractionConfiguration: &bedrock.AgentDataSourceDataSourceConfigurationManagedKnowledgeBaseConnectorConfigurationMediaExtractionConfigurationVideoExtractionConfigurationArgs{
					VideoExtractionStatus: pulumi.String(enabledOrDisabled(d.ManagedConnector.MediaExtraction.Video)),
				},
			}
		}
		configuration.ManagedKnowledgeBaseConnectorConfiguration = managed
	}
	args.DataSourceConfiguration = configuration

	if d.VectorIngestion != nil {
		ingestion := &bedrock.AgentDataSourceVectorIngestionConfigurationArgs{}
		if d.VectorIngestion.Chunking != nil {
			chunking := &bedrock.AgentDataSourceVectorIngestionConfigurationChunkingConfigurationArgs{
				ChunkingStrategy: pulumi.String(d.VectorIngestion.Chunking.Strategy),
			}
			if d.VectorIngestion.Chunking.FixedSize != nil {
				chunking.FixedSizeChunkingConfiguration = &bedrock.AgentDataSourceVectorIngestionConfigurationChunkingConfigurationFixedSizeChunkingConfigurationArgs{
					MaxTokens:         pulumi.Int(int(d.VectorIngestion.Chunking.FixedSize.MaxTokens)),
					OverlapPercentage: pulumi.Int(int(d.VectorIngestion.Chunking.FixedSize.OverlapPercentage)),
				}
			}
			if d.VectorIngestion.Chunking.Hierarchical != nil {
				var levels bedrock.AgentDataSourceVectorIngestionConfigurationChunkingConfigurationHierarchicalChunkingConfigurationLevelConfigurationArray
				for _, l := range d.VectorIngestion.Chunking.Hierarchical.Levels {
					levels = append(levels, &bedrock.AgentDataSourceVectorIngestionConfigurationChunkingConfigurationHierarchicalChunkingConfigurationLevelConfigurationArgs{
						MaxTokens: pulumi.Int(int(l.MaxTokens)),
					})
				}
				chunking.HierarchicalChunkingConfiguration = &bedrock.AgentDataSourceVectorIngestionConfigurationChunkingConfigurationHierarchicalChunkingConfigurationArgs{
					OverlapTokens:       pulumi.Int(int(d.VectorIngestion.Chunking.Hierarchical.OverlapTokens)),
					LevelConfigurations: levels,
				}
			}
			if d.VectorIngestion.Chunking.Semantic != nil {
				chunking.SemanticChunkingConfiguration = &bedrock.AgentDataSourceVectorIngestionConfigurationChunkingConfigurationSemanticChunkingConfigurationArgs{
					BreakpointPercentileThreshold: pulumi.Int(int(d.VectorIngestion.Chunking.Semantic.BreakpointPercentileThreshold)),
					BufferSize:                    pulumi.Int(int(d.VectorIngestion.Chunking.Semantic.BufferSize)),
					// The provider spells this one singular.
					MaxToken: pulumi.Int(int(d.VectorIngestion.Chunking.Semantic.MaxTokens)),
				}
			}
			ingestion.ChunkingConfiguration = chunking
		}
		if d.VectorIngestion.Parsing != nil {
			parsing := &bedrock.AgentDataSourceVectorIngestionConfigurationParsingConfigurationArgs{
				ParsingStrategy: pulumi.String(d.VectorIngestion.Parsing.Strategy),
			}
			// MULTIMODAL is the only parsing modality AWS defines -- the
			// spec models it as a bool and the module owns the constant.
			if d.VectorIngestion.Parsing.Strategy == "BEDROCK_DATA_AUTOMATION" && d.VectorIngestion.Parsing.Multimodal {
				parsing.BedrockDataAutomationConfiguration = &bedrock.AgentDataSourceVectorIngestionConfigurationParsingConfigurationBedrockDataAutomationConfigurationArgs{
					ParsingModality: pulumi.String("MULTIMODAL"),
				}
			}
			if d.VectorIngestion.Parsing.FoundationModel != nil {
				foundationModel := &bedrock.AgentDataSourceVectorIngestionConfigurationParsingConfigurationBedrockFoundationModelConfigurationArgs{
					ModelArn: pulumi.String(d.VectorIngestion.Parsing.FoundationModel.ModelArn),
				}
				if d.VectorIngestion.Parsing.FoundationModel.Multimodal {
					foundationModel.ParsingModality = pulumi.String("MULTIMODAL")
				}
				if d.VectorIngestion.Parsing.FoundationModel.ParsingPrompt != "" {
					foundationModel.ParsingPrompt = &bedrock.AgentDataSourceVectorIngestionConfigurationParsingConfigurationBedrockFoundationModelConfigurationParsingPromptArgs{
						ParsingPromptString: pulumi.String(d.VectorIngestion.Parsing.FoundationModel.ParsingPrompt),
					}
				}
				parsing.BedrockFoundationModelConfiguration = foundationModel
			}
			ingestion.ParsingConfiguration = parsing
		}
		if d.VectorIngestion.CustomTransformation != nil {
			ingestion.CustomTransformationConfiguration = &bedrock.AgentDataSourceVectorIngestionConfigurationCustomTransformationConfigurationArgs{
				IntermediateStorage: &bedrock.AgentDataSourceVectorIngestionConfigurationCustomTransformationConfigurationIntermediateStorageArgs{
					S3Location: &bedrock.AgentDataSourceVectorIngestionConfigurationCustomTransformationConfigurationIntermediateStorageS3LocationArgs{
						Uri: pulumi.String(d.VectorIngestion.CustomTransformation.IntermediateS3Uri),
					},
				},
				Transformation: &bedrock.AgentDataSourceVectorIngestionConfigurationCustomTransformationConfigurationTransformationArgs{
					// POST_CHUNKING is the only transformation step AWS
					// defines -- the module owns the constant.
					StepToApply: pulumi.String("POST_CHUNKING"),
					TransformationFunction: &bedrock.AgentDataSourceVectorIngestionConfigurationCustomTransformationConfigurationTransformationTransformationFunctionArgs{
						TransformationLambdaConfiguration: &bedrock.AgentDataSourceVectorIngestionConfigurationCustomTransformationConfigurationTransformationTransformationFunctionTransformationLambdaConfigurationArgs{
							LambdaArn: pulumi.String(d.VectorIngestion.CustomTransformation.LambdaArn.GetValue()),
						},
					},
				},
			}
		}
		args.VectorIngestionConfiguration = ingestion
	}

	return args, nil
}

func enabledOrDisabled(enabled bool) string {
	if enabled {
		return "ENABLED"
	}
	return "DISABLED"
}

func sortedDataSources(in []*awsbedrockknowledgebasev1alpha1.AwsBedrockKnowledgeBaseDataSource) []*awsbedrockknowledgebasev1alpha1.AwsBedrockKnowledgeBaseDataSource {
	out := append([]*awsbedrockknowledgebasev1alpha1.AwsBedrockKnowledgeBaseDataSource{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
