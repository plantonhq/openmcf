package awsdynamodbv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsDynamodbSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsDynamodbSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

var _ = ginkgo.Describe("AwsDynamodbSpec validations", func() {
	var spec *AwsDynamodbSpec

	ginkgo.BeforeEach(func() {
		spec = &AwsDynamodbSpec{
			Region:      "us-west-2",
			BillingMode: "PAY_PER_REQUEST",
			AttributeDefinitions: []*AwsDynamodbAttribute{
				{Name: "pk", Type: "S"},
			},
			KeySchema: []*AwsDynamodbKeySchemaElement{
				{AttributeName: "pk", KeyType: "HASH"},
			},
		}
	})

	// -----------------------------------------------------------------
	// Happy paths
	// -----------------------------------------------------------------

	ginkgo.It("accepts a minimal valid PAY_PER_REQUEST table", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a PROVISIONED table with throughput", func() {
		spec.BillingMode = "PROVISIONED"
		spec.ProvisionedThroughput = &AwsDynamodbProvisionedThroughput{ReadCapacityUnits: 5, WriteCapacityUnits: 5}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts an empty billing_mode with throughput (the AWS PROVISIONED default)", func() {
		spec.BillingMode = ""
		spec.ProvisionedThroughput = &AwsDynamodbProvisionedThroughput{ReadCapacityUnits: 5, WriteCapacityUnits: 5}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a full-featured on-demand table", func() {
		spec.AttributeDefinitions = append(spec.AttributeDefinitions,
			&AwsDynamodbAttribute{Name: "sk", Type: "S"},
			&AwsDynamodbAttribute{Name: "gsi1pk", Type: "S"},
		)
		spec.KeySchema = append(spec.KeySchema, &AwsDynamodbKeySchemaElement{AttributeName: "sk", KeyType: "RANGE"})
		spec.OnDemandThroughput = &AwsDynamodbOnDemandThroughput{MaxReadRequestUnits: 1000, MaxWriteRequestUnits: 500}
		spec.WarmThroughput = &AwsDynamodbWarmThroughput{ReadUnitsPerSecond: 15000, WriteUnitsPerSecond: 5000}
		spec.GlobalSecondaryIndexes = []*AwsDynamodbGlobalSecondaryIndex{{
			Name:       "gsi1",
			KeySchema:  []*AwsDynamodbKeySchemaElement{{AttributeName: "gsi1pk", KeyType: "HASH"}},
			Projection: &AwsDynamodbProjection{Type: "ALL"},
		}}
		spec.LocalSecondaryIndexes = []*AwsDynamodbLocalSecondaryIndex{{
			Name:       "lsi1",
			RangeKey:   "gsi1pk",
			Projection: &AwsDynamodbProjection{Type: "KEYS_ONLY"},
		}}
		spec.Ttl = &AwsDynamodbTtl{Enabled: true, AttributeName: "expiresAt"}
		spec.StreamEnabled = true
		spec.StreamViewType = "NEW_AND_OLD_IMAGES"
		spec.PointInTimeRecovery = &AwsDynamodbPointInTimeRecovery{Enabled: true, RecoveryPeriodInDays: 14}
		spec.ServerSideEncryption = &AwsDynamodbServerSideEncryption{Enabled: true, KmsKeyArn: literal("arn:aws:kms:us-west-2:111122223333:key/abc")}
		spec.TableClass = "STANDARD_INFREQUENT_ACCESS"
		spec.DeletionProtectionEnabled = true
		spec.ContributorInsights = &AwsDynamodbContributorInsights{Enabled: true, Mode: "THROTTLED_KEYS", GsiIndexNames: []string{"gsi1"}}
		spec.ResourcePolicy = `{"Version":"2012-10-17","Statement":[]}`
		spec.KinesisStreamingDestination = &AwsDynamodbKinesisStreamingDestination{
			StreamArn: literal("arn:aws:kinesis:us-west-2:111122223333:stream/cdc"),
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a global table with two EVENTUAL replicas", func() {
		spec.StreamEnabled = true
		spec.StreamViewType = "NEW_AND_OLD_IMAGES"
		spec.Replicas = []*AwsDynamodbReplica{
			{RegionName: "us-east-1"},
			{RegionName: "eu-west-1", PointInTimeRecovery: true, PropagateTags: true},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts MRSC with two STRONG replicas", func() {
		spec.StreamEnabled = true
		spec.StreamViewType = "NEW_AND_OLD_IMAGES"
		spec.Replicas = []*AwsDynamodbReplica{
			{RegionName: "us-east-1", ConsistencyMode: "STRONG"},
			{RegionName: "us-east-2", ConsistencyMode: "STRONG"},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts MRSC with one STRONG replica plus a witness", func() {
		spec.StreamEnabled = true
		spec.StreamViewType = "NEW_AND_OLD_IMAGES"
		spec.Replicas = []*AwsDynamodbReplica{{RegionName: "us-east-1", ConsistencyMode: "STRONG"}}
		spec.GlobalTableWitness = &AwsDynamodbGlobalTableWitness{RegionName: "us-east-2"}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a restore-created table without key schema", func() {
		spec.AttributeDefinitions = nil
		spec.KeySchema = nil
		spec.RestoreSourceName = "source-table"
		spec.RestoreToLatestTime = true
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts an S3 import with CSV options", func() {
		spec.ImportTable = &AwsDynamodbImportTable{
			S3Bucket:    literal("seed-data-bucket"),
			InputFormat: "CSV",
			Csv:         &AwsDynamodbImportTableCsv{Delimiter: ";", HeaderList: []string{"pk", "payload"}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Field-level validations
	// -----------------------------------------------------------------

	ginkgo.It("fails when region is empty", func() {
		spec.Region = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails on an unknown billing_mode", func() {
		spec.BillingMode = "ON_DEMAND"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails on an unknown attribute type", func() {
		spec.AttributeDefinitions[0].Type = "BOOL"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when an attribute name is empty", func() {
		spec.AttributeDefinitions[0].Name = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails on an unknown stream_view_type", func() {
		spec.StreamEnabled = true
		spec.StreamViewType = "EVERYTHING"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails on an unknown table_class", func() {
		spec.TableClass = "GLACIER"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails on warm throughput below the AWS minimums", func() {
		spec.WarmThroughput = &AwsDynamodbWarmThroughput{ReadUnitsPerSecond: 100}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// CEL: key schema and attributes
	// -----------------------------------------------------------------

	ginkgo.It("fails when key_schema is empty without a restore source", func() {
		spec.KeySchema = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when attribute_definitions is empty without a restore source", func() {
		spec.AttributeDefinitions = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when key_schema has two HASH elements", func() {
		spec.AttributeDefinitions = append(spec.AttributeDefinitions, &AwsDynamodbAttribute{Name: "pk2", Type: "S"})
		spec.KeySchema = []*AwsDynamodbKeySchemaElement{
			{AttributeName: "pk", KeyType: "HASH"},
			{AttributeName: "pk2", KeyType: "HASH"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when key_schema starts with a RANGE element", func() {
		spec.AttributeDefinitions = append(spec.AttributeDefinitions, &AwsDynamodbAttribute{Name: "sk", Type: "S"})
		spec.KeySchema = []*AwsDynamodbKeySchemaElement{
			{AttributeName: "sk", KeyType: "RANGE"},
			{AttributeName: "pk", KeyType: "HASH"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a key element references an undeclared attribute", func() {
		spec.KeySchema = []*AwsDynamodbKeySchemaElement{{AttributeName: "ghost", KeyType: "HASH"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when an LSI range key references an undeclared attribute", func() {
		spec.LocalSecondaryIndexes = []*AwsDynamodbLocalSecondaryIndex{{
			Name:       "lsi1",
			RangeKey:   "ghost",
			Projection: &AwsDynamodbProjection{Type: "ALL"},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// CEL: capacity coupling
	// -----------------------------------------------------------------

	ginkgo.It("fails when PROVISIONED lacks throughput", func() {
		spec.BillingMode = "PROVISIONED"
		spec.ProvisionedThroughput = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when empty billing_mode (PROVISIONED default) lacks throughput", func() {
		spec.BillingMode = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when PAY_PER_REQUEST carries provisioned_throughput", func() {
		spec.ProvisionedThroughput = &AwsDynamodbProvisionedThroughput{ReadCapacityUnits: 1, WriteCapacityUnits: 1}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a PROVISIONED table carries on_demand_throughput", func() {
		spec.BillingMode = "PROVISIONED"
		spec.ProvisionedThroughput = &AwsDynamodbProvisionedThroughput{ReadCapacityUnits: 5, WriteCapacityUnits: 5}
		spec.OnDemandThroughput = &AwsDynamodbOnDemandThroughput{MaxReadRequestUnits: 100}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a PROVISIONED table has a GSI without throughput", func() {
		spec.BillingMode = "PROVISIONED"
		spec.ProvisionedThroughput = &AwsDynamodbProvisionedThroughput{ReadCapacityUnits: 5, WriteCapacityUnits: 5}
		spec.AttributeDefinitions = append(spec.AttributeDefinitions, &AwsDynamodbAttribute{Name: "gsi1pk", Type: "S"})
		spec.GlobalSecondaryIndexes = []*AwsDynamodbGlobalSecondaryIndex{{
			Name:       "gsi1",
			KeySchema:  []*AwsDynamodbKeySchemaElement{{AttributeName: "gsi1pk", KeyType: "HASH"}},
			Projection: &AwsDynamodbProjection{Type: "ALL"},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when an on-demand table's GSI carries provisioned throughput", func() {
		spec.AttributeDefinitions = append(spec.AttributeDefinitions, &AwsDynamodbAttribute{Name: "gsi1pk", Type: "S"})
		spec.GlobalSecondaryIndexes = []*AwsDynamodbGlobalSecondaryIndex{{
			Name:                  "gsi1",
			KeySchema:             []*AwsDynamodbKeySchemaElement{{AttributeName: "gsi1pk", KeyType: "HASH"}},
			Projection:            &AwsDynamodbProjection{Type: "ALL"},
			ProvisionedThroughput: &AwsDynamodbProvisionedThroughput{ReadCapacityUnits: 1, WriteCapacityUnits: 1},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// CEL: index shapes and projections
	// -----------------------------------------------------------------

	ginkgo.It("fails when a GSI key_schema starts with RANGE", func() {
		spec.AttributeDefinitions = append(spec.AttributeDefinitions, &AwsDynamodbAttribute{Name: "gsi1sk", Type: "S"})
		spec.GlobalSecondaryIndexes = []*AwsDynamodbGlobalSecondaryIndex{{
			Name:       "gsi1",
			KeySchema:  []*AwsDynamodbKeySchemaElement{{AttributeName: "gsi1sk", KeyType: "RANGE"}},
			Projection: &AwsDynamodbProjection{Type: "ALL"},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a GSI has five HASH elements", func() {
		keySchema := make([]*AwsDynamodbKeySchemaElement, 0, 5)
		for _, name := range []string{"h1", "h2", "h3", "h4", "h5"} {
			spec.AttributeDefinitions = append(spec.AttributeDefinitions, &AwsDynamodbAttribute{Name: name, Type: "S"})
			keySchema = append(keySchema, &AwsDynamodbKeySchemaElement{AttributeName: name, KeyType: "HASH"})
		}
		spec.GlobalSecondaryIndexes = []*AwsDynamodbGlobalSecondaryIndex{{
			Name:       "gsi1",
			KeySchema:  keySchema,
			Projection: &AwsDynamodbProjection{Type: "ALL"},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when INCLUDE projection has no non_key_attributes", func() {
		spec.AttributeDefinitions = append(spec.AttributeDefinitions, &AwsDynamodbAttribute{Name: "gsi1pk", Type: "S"})
		spec.GlobalSecondaryIndexes = []*AwsDynamodbGlobalSecondaryIndex{{
			Name:       "gsi1",
			KeySchema:  []*AwsDynamodbKeySchemaElement{{AttributeName: "gsi1pk", KeyType: "HASH"}},
			Projection: &AwsDynamodbProjection{Type: "INCLUDE"},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a non-INCLUDE projection carries non_key_attributes", func() {
		spec.AttributeDefinitions = append(spec.AttributeDefinitions, &AwsDynamodbAttribute{Name: "gsi1pk", Type: "S"})
		spec.GlobalSecondaryIndexes = []*AwsDynamodbGlobalSecondaryIndex{{
			Name:       "gsi1",
			KeySchema:  []*AwsDynamodbKeySchemaElement{{AttributeName: "gsi1pk", KeyType: "HASH"}},
			Projection: &AwsDynamodbProjection{Type: "ALL", NonKeyAttributes: []string{"extra"}},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// CEL: streams, TTL, PITR, SSE
	// -----------------------------------------------------------------

	ginkgo.It("fails when streams are enabled without a view type", func() {
		spec.StreamEnabled = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a view type is set with streams disabled", func() {
		spec.StreamViewType = "KEYS_ONLY"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when TTL is enabled without attribute_name", func() {
		spec.Ttl = &AwsDynamodbTtl{Enabled: true}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when TTL is disabled but attribute_name is set", func() {
		spec.Ttl = &AwsDynamodbTtl{Enabled: false, AttributeName: "expiresAt"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a recovery period is set with PITR disabled", func() {
		spec.PointInTimeRecovery = &AwsDynamodbPointInTimeRecovery{Enabled: false, RecoveryPeriodInDays: 7}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when the recovery period exceeds 35 days", func() {
		spec.PointInTimeRecovery = &AwsDynamodbPointInTimeRecovery{Enabled: true, RecoveryPeriodInDays: 36}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a KMS key is set with SSE disabled", func() {
		spec.ServerSideEncryption = &AwsDynamodbServerSideEncryption{
			Enabled:   false,
			KmsKeyArn: literal("arn:aws:kms:us-west-2:111122223333:key/abc"),
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// CEL: replicas and MRSC topology
	// -----------------------------------------------------------------

	ginkgo.It("fails when replicas are set without streams", func() {
		spec.Replicas = []*AwsDynamodbReplica{{RegionName: "us-east-1"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when replicas are set with the wrong stream view type", func() {
		spec.StreamEnabled = true
		spec.StreamViewType = "KEYS_ONLY"
		spec.Replicas = []*AwsDynamodbReplica{{RegionName: "us-east-1"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when STRONG and EVENTUAL replicas are mixed", func() {
		spec.StreamEnabled = true
		spec.StreamViewType = "NEW_AND_OLD_IMAGES"
		spec.Replicas = []*AwsDynamodbReplica{
			{RegionName: "us-east-1", ConsistencyMode: "STRONG"},
			{RegionName: "eu-west-1", ConsistencyMode: "EVENTUAL"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails MRSC with one STRONG replica and no witness", func() {
		spec.StreamEnabled = true
		spec.StreamViewType = "NEW_AND_OLD_IMAGES"
		spec.Replicas = []*AwsDynamodbReplica{{RegionName: "us-east-1", ConsistencyMode: "STRONG"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails MRSC with two STRONG replicas plus a witness", func() {
		spec.StreamEnabled = true
		spec.StreamViewType = "NEW_AND_OLD_IMAGES"
		spec.Replicas = []*AwsDynamodbReplica{
			{RegionName: "us-east-1", ConsistencyMode: "STRONG"},
			{RegionName: "us-east-2", ConsistencyMode: "STRONG"},
		}
		spec.GlobalTableWitness = &AwsDynamodbGlobalTableWitness{RegionName: "eu-west-1"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a witness accompanies EVENTUAL replicas", func() {
		spec.StreamEnabled = true
		spec.StreamViewType = "NEW_AND_OLD_IMAGES"
		spec.Replicas = []*AwsDynamodbReplica{{RegionName: "us-east-1"}}
		spec.GlobalTableWitness = &AwsDynamodbGlobalTableWitness{RegionName: "us-east-2"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// CEL: create sources (restore / import)
	// -----------------------------------------------------------------

	ginkgo.It("fails when two restore sources are set", func() {
		spec.RestoreSourceName = "source-table"
		spec.RestoreToLatestTime = true
		spec.RestoreBackupArn = "arn:aws:dynamodb:us-west-2:111122223333:table/source/backup/b1"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a restore source is set with import_table", func() {
		spec.RestoreSourceName = "source-table"
		spec.RestoreToLatestTime = true
		spec.ImportTable = &AwsDynamodbImportTable{S3Bucket: literal("bucket"), InputFormat: "ION"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails a PITR restore without a restore point", func() {
		spec.RestoreSourceName = "source-table"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails a PITR restore with both a timestamp and latest", func() {
		spec.RestoreSourceName = "source-table"
		spec.RestoreDateTime = "2026-07-04T06:00:00Z"
		spec.RestoreToLatestTime = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a restore point is set without a PITR source", func() {
		spec.RestoreToLatestTime = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when CSV options accompany a non-CSV import", func() {
		spec.ImportTable = &AwsDynamodbImportTable{
			S3Bucket:    literal("bucket"),
			InputFormat: "DYNAMODB_JSON",
			Csv:         &AwsDynamodbImportTableCsv{Delimiter: ";"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails on an unknown import input_format", func() {
		spec.ImportTable = &AwsDynamodbImportTable{S3Bucket: literal("bucket"), InputFormat: "PARQUET"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// CEL: contributor insights
	// -----------------------------------------------------------------

	ginkgo.It("fails when insights name a GSI that does not exist", func() {
		spec.ContributorInsights = &AwsDynamodbContributorInsights{Enabled: true, GsiIndexNames: []string{"ghost"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when insights fields are set while disabled", func() {
		spec.ContributorInsights = &AwsDynamodbContributorInsights{Enabled: false, Mode: "THROTTLED_KEYS"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails on an unknown kinesis timestamp precision", func() {
		spec.KinesisStreamingDestination = &AwsDynamodbKinesisStreamingDestination{
			StreamArn:                            literal("arn:aws:kinesis:us-west-2:111122223333:stream/cdc"),
			ApproximateCreationDateTimePrecision: "NANOSECOND",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})
})
