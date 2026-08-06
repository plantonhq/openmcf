package module

import (
	"github.com/pkg/errors"
	awsdynamodbv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsdynamodb/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// table provisions the DynamoDB table itself: key schema and indexes,
// capacity, streams, global-table replication, encryption, and
// recovery. Create-only in AWS: the table name, the primary key schema,
// and every local secondary index (LSIs can never be added or removed
// later). Everything else -- billing mode, GSIs, streams, replicas,
// TTL, PITR, SSE key, table class -- edits in place.
func table(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*dynamodb.Table, error) {
	spec := locals.AwsDynamodb.Spec

	args := &dynamodb.TableArgs{
		// The table name is metadata.name on both engines, so a
		// manifest deploys the same table regardless of engine.
		Name: pulumi.String(locals.TableName),
		Tags: pulumi.ToStringMap(locals.AwsTags),

		DeletionProtectionEnabled: pulumi.Bool(spec.DeletionProtectionEnabled),
	}

	// Empty keeps the AWS default (PROVISIONED). Spec values are the
	// exact AWS API strings, so they pass through untranslated.
	if spec.BillingMode != "" {
		args.BillingMode = pulumi.String(spec.BillingMode)
	}

	// Key attributes only -- DynamoDB is schemaless beyond the keys.
	// Restore-created tables inherit schema from the source, so the
	// lists may be empty (CEL enforces the coupling).
	if len(spec.AttributeDefinitions) > 0 {
		attributes := dynamodb.TableAttributeArray{}
		for _, attribute := range spec.AttributeDefinitions {
			attributes = append(attributes, dynamodb.TableAttributeArgs{
				Name: pulumi.String(attribute.Name),
				Type: pulumi.String(attribute.Type),
			})
		}
		args.Attributes = attributes
	}

	// The provider models the table's primary key as hash_key/range_key
	// scalars; the spec's key_schema (one HASH, optional RANGE --
	// CEL-enforced) is lowered here.
	for _, element := range spec.KeySchema {
		switch element.KeyType {
		case "HASH":
			args.HashKey = pulumi.String(element.AttributeName)
		case "RANGE":
			args.RangeKey = pulumi.String(element.AttributeName)
		}
	}

	// Reserved capacity applies only when the effective billing mode is
	// PROVISIONED (CEL enforces the coupling either way).
	if spec.ProvisionedThroughput != nil {
		args.ReadCapacity = pulumi.Int(int(spec.ProvisionedThroughput.ReadCapacityUnits))
		args.WriteCapacity = pulumi.Int(int(spec.ProvisionedThroughput.WriteCapacityUnits))
	}

	// On-demand ceilings throttle instead of billing past the cap; -1
	// explicitly removes a previously-set ceiling.
	if spec.OnDemandThroughput != nil {
		args.OnDemandThroughput = onDemandThroughputArgs(spec.OnDemandThroughput)
	}

	// Warm throughput only ever increases -- AWS replaces the table on
	// a decrease (the provider marks it ForceNew for that reason).
	if spec.WarmThroughput != nil {
		warm := &dynamodb.TableWarmThroughputArgs{}
		if spec.WarmThroughput.ReadUnitsPerSecond != 0 {
			warm.ReadUnitsPerSecond = pulumi.Int(int(spec.WarmThroughput.ReadUnitsPerSecond))
		}
		if spec.WarmThroughput.WriteUnitsPerSecond != 0 {
			warm.WriteUnitsPerSecond = pulumi.Int(int(spec.WarmThroughput.WriteUnitsPerSecond))
		}
		args.WarmThroughput = warm
	}

	if len(spec.GlobalSecondaryIndexes) > 0 {
		globalSecondaryIndexes := dynamodb.TableGlobalSecondaryIndexArray{}
		for _, gsi := range spec.GlobalSecondaryIndexes {
			globalSecondaryIndexes = append(globalSecondaryIndexes, globalSecondaryIndexArgs(gsi))
		}
		args.GlobalSecondaryIndexes = globalSecondaryIndexes
	}

	// LSIs are create-only: they can never be added or removed after
	// the table exists, and their presence permanently caps each item
	// collection at 10 GB.
	if len(spec.LocalSecondaryIndexes) > 0 {
		localSecondaryIndexes := dynamodb.TableLocalSecondaryIndexArray{}
		for _, lsi := range spec.LocalSecondaryIndexes {
			lsiArgs := dynamodb.TableLocalSecondaryIndexArgs{
				Name:           pulumi.String(lsi.Name),
				RangeKey:       pulumi.String(lsi.RangeKey),
				ProjectionType: pulumi.String(lsi.Projection.Type),
			}
			if len(lsi.Projection.NonKeyAttributes) > 0 {
				lsiArgs.NonKeyAttributes = pulumi.ToStringArray(lsi.Projection.NonKeyAttributes)
			}
			localSecondaryIndexes = append(localSecondaryIndexes, lsiArgs)
		}
		args.LocalSecondaryIndexes = localSecondaryIndexes
	}

	if spec.Ttl != nil {
		ttlArgs := &dynamodb.TableTtlArgs{Enabled: pulumi.Bool(spec.Ttl.Enabled)}
		if spec.Ttl.AttributeName != "" {
			ttlArgs.AttributeName = pulumi.String(spec.Ttl.AttributeName)
		}
		args.Ttl = ttlArgs
	}

	// Streams carry item-level change data; global tables require the
	// NEW_AND_OLD_IMAGES view (CEL enforces it).
	args.StreamEnabled = pulumi.Bool(spec.StreamEnabled)
	if spec.StreamViewType != "" {
		args.StreamViewType = pulumi.String(spec.StreamViewType)
	}

	if spec.PointInTimeRecovery != nil {
		pitrArgs := &dynamodb.TablePointInTimeRecoveryArgs{
			Enabled: pulumi.Bool(spec.PointInTimeRecovery.Enabled),
		}
		if spec.PointInTimeRecovery.RecoveryPeriodInDays != 0 {
			pitrArgs.RecoveryPeriodInDays = pulumi.Int(int(spec.PointInTimeRecovery.RecoveryPeriodInDays))
		}
		args.PointInTimeRecovery = pitrArgs
	}

	// DynamoDB always encrypts; this block switches from the AWS-owned
	// key to the AWS-managed aws/dynamodb key (no KMS ARN) or a
	// customer-managed key (ARN set).
	if spec.ServerSideEncryption != nil {
		sseArgs := &dynamodb.TableServerSideEncryptionArgs{
			Enabled: pulumi.Bool(spec.ServerSideEncryption.Enabled),
		}
		if spec.ServerSideEncryption.KmsKeyArn.GetValue() != "" {
			sseArgs.KmsKeyArn = pulumi.String(spec.ServerSideEncryption.KmsKeyArn.GetValue())
		}
		args.ServerSideEncryption = sseArgs
	}

	// Empty keeps the AWS default (STANDARD).
	if spec.TableClass != "" {
		args.TableClass = pulumi.String(spec.TableClass)
	}

	// Global Tables v2: each replica is an active read/write copy in
	// another region; each region encrypts independently, so the KMS
	// key is per-replica.
	if len(spec.Replicas) > 0 {
		replicas := dynamodb.TableReplicaTypeArray{}
		for _, replica := range spec.Replicas {
			replicaArgs := dynamodb.TableReplicaTypeArgs{
				RegionName:                pulumi.String(replica.RegionName),
				PointInTimeRecovery:       pulumi.Bool(replica.PointInTimeRecovery),
				DeletionProtectionEnabled: pulumi.Bool(replica.DeletionProtectionEnabled),
				PropagateTags:             pulumi.Bool(replica.PropagateTags),
			}
			if replica.KmsKeyArn.GetValue() != "" {
				replicaArgs.KmsKeyArn = pulumi.String(replica.KmsKeyArn.GetValue())
			}
			if replica.ConsistencyMode != "" {
				replicaArgs.ConsistencyMode = pulumi.String(replica.ConsistencyMode)
			}
			replicas = append(replicas, replicaArgs)
		}
		args.Replicas = replicas
	}

	// The MRSC witness persists replicated writes for quorum but serves
	// no reads or writes; CEL pins the exact topology AWS accepts.
	if spec.GlobalTableWitness != nil {
		args.GlobalTableWitness = &dynamodb.TableGlobalTableWitnessArgs{
			RegionName: pulumi.String(spec.GlobalTableWitness.RegionName),
		}
	}

	// Create sources are mutually exclusive (CEL): a point-in-time
	// restore by name or ARN, a backup restore, or an S3 import.
	if spec.RestoreSourceName != "" {
		args.RestoreSourceName = pulumi.String(spec.RestoreSourceName)
	}
	if spec.RestoreSourceTableArn != "" {
		args.RestoreSourceTableArn = pulumi.String(spec.RestoreSourceTableArn)
	}
	if spec.RestoreDateTime != "" {
		args.RestoreDateTime = pulumi.String(spec.RestoreDateTime)
	}
	if spec.RestoreToLatestTime {
		args.RestoreToLatestTime = pulumi.Bool(true)
	}
	if spec.RestoreBackupArn != "" {
		args.RestoreBackupArn = pulumi.String(spec.RestoreBackupArn)
	}
	if spec.ImportTable != nil {
		args.ImportTable = importTableArgs(spec.ImportTable)
	}

	createdTable, err := dynamodb.NewTable(ctx, "table", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create DynamoDB table")
	}
	return createdTable, nil
}

// globalSecondaryIndexArgs lowers one GSI. The modern key_schema shape
// carries multi-attribute keys (1-4 HASH elements first, then 0-4
// RANGE elements); per-index capacity follows the table's billing mode
// (CEL enforces the coupling).
func globalSecondaryIndexArgs(gsi *awsdynamodbv1alpha1.AwsDynamodbGlobalSecondaryIndex) dynamodb.TableGlobalSecondaryIndexArgs {
	keySchemas := dynamodb.TableGlobalSecondaryIndexKeySchemaArray{}
	for _, element := range gsi.KeySchema {
		keySchemas = append(keySchemas, dynamodb.TableGlobalSecondaryIndexKeySchemaArgs{
			AttributeName: pulumi.String(element.AttributeName),
			KeyType:       pulumi.String(element.KeyType),
		})
	}

	gsiArgs := dynamodb.TableGlobalSecondaryIndexArgs{
		Name:           pulumi.String(gsi.Name),
		KeySchemas:     keySchemas,
		ProjectionType: pulumi.String(gsi.Projection.Type),
	}
	if len(gsi.Projection.NonKeyAttributes) > 0 {
		gsiArgs.NonKeyAttributes = pulumi.ToStringArray(gsi.Projection.NonKeyAttributes)
	}
	if gsi.ProvisionedThroughput != nil {
		gsiArgs.ReadCapacity = pulumi.Int(int(gsi.ProvisionedThroughput.ReadCapacityUnits))
		gsiArgs.WriteCapacity = pulumi.Int(int(gsi.ProvisionedThroughput.WriteCapacityUnits))
	}
	if gsi.OnDemandThroughput != nil {
		onDemand := &dynamodb.TableGlobalSecondaryIndexOnDemandThroughputArgs{}
		if gsi.OnDemandThroughput.MaxReadRequestUnits != 0 {
			onDemand.MaxReadRequestUnits = pulumi.Int(int(gsi.OnDemandThroughput.MaxReadRequestUnits))
		}
		if gsi.OnDemandThroughput.MaxWriteRequestUnits != 0 {
			onDemand.MaxWriteRequestUnits = pulumi.Int(int(gsi.OnDemandThroughput.MaxWriteRequestUnits))
		}
		gsiArgs.OnDemandThroughput = onDemand
	}
	if gsi.WarmThroughput != nil {
		warm := &dynamodb.TableGlobalSecondaryIndexWarmThroughputArgs{}
		if gsi.WarmThroughput.ReadUnitsPerSecond != 0 {
			warm.ReadUnitsPerSecond = pulumi.Int(int(gsi.WarmThroughput.ReadUnitsPerSecond))
		}
		if gsi.WarmThroughput.WriteUnitsPerSecond != 0 {
			warm.WriteUnitsPerSecond = pulumi.Int(int(gsi.WarmThroughput.WriteUnitsPerSecond))
		}
		gsiArgs.WarmThroughput = warm
	}
	return gsiArgs
}

// onDemandThroughputArgs lowers the table-level on-demand ceilings.
// 0 means "not configured" and is omitted; -1 passes through to
// explicitly remove a previously-set ceiling.
func onDemandThroughputArgs(spec *awsdynamodbv1alpha1.AwsDynamodbOnDemandThroughput) *dynamodb.TableOnDemandThroughputArgs {
	onDemand := &dynamodb.TableOnDemandThroughputArgs{}
	if spec.MaxReadRequestUnits != 0 {
		onDemand.MaxReadRequestUnits = pulumi.Int(int(spec.MaxReadRequestUnits))
	}
	if spec.MaxWriteRequestUnits != 0 {
		onDemand.MaxWriteRequestUnits = pulumi.Int(int(spec.MaxWriteRequestUnits))
	}
	return onDemand
}

// importTableArgs lowers the S3 import source that seeds a brand-new
// table -- billed as a one-time import instead of per-item writes.
func importTableArgs(spec *awsdynamodbv1alpha1.AwsDynamodbImportTable) *dynamodb.TableImportTableArgs {
	s3BucketSource := dynamodb.TableImportTableS3BucketSourceArgs{
		Bucket: pulumi.String(spec.S3Bucket.GetValue()),
	}
	if spec.S3BucketOwner != "" {
		s3BucketSource.BucketOwner = pulumi.String(spec.S3BucketOwner)
	}
	if spec.S3KeyPrefix != "" {
		s3BucketSource.KeyPrefix = pulumi.String(spec.S3KeyPrefix)
	}

	importArgs := &dynamodb.TableImportTableArgs{
		InputFormat:    pulumi.String(spec.InputFormat),
		S3BucketSource: s3BucketSource,
	}
	if spec.InputCompressionType != "" {
		importArgs.InputCompressionType = pulumi.String(spec.InputCompressionType)
	}
	if spec.Csv != nil {
		csvArgs := &dynamodb.TableImportTableInputFormatOptionsCsvArgs{}
		if spec.Csv.Delimiter != "" {
			csvArgs.Delimiter = pulumi.String(spec.Csv.Delimiter)
		}
		if len(spec.Csv.HeaderList) > 0 {
			csvArgs.HeaderLists = pulumi.ToStringArray(spec.Csv.HeaderList)
		}
		importArgs.InputFormatOptions = &dynamodb.TableImportTableInputFormatOptionsArgs{
			Csv: csvArgs,
		}
	}
	return importArgs
}
