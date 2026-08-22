package awss3tablebucketv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsS3TableBucketSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsS3TableBucketSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func minimalTableBucket() *AwsS3TableBucketSpec {
	return &AwsS3TableBucketSpec{Region: "us-east-1"}
}

func analyticsNamespace() *AwsS3TablesNamespace {
	return &AwsS3TablesNamespace{
		Name: "analytics",
		Tables: []*AwsS3TablesTable{{
			Name: "events",
			IcebergSchema: &AwsS3TablesIcebergSchema{
				Fields: []*AwsS3TablesSchemaField{
					{Name: "event_id", Type: "string", Required: true},
					{Name: "occurred_at", Type: "timestamp"},
				},
			},
		}},
	}
}

var _ = ginkgo.Describe("AwsS3TableBucketSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal empty bucket", func() {
			gomega.Expect(protovalidate.Validate(minimalTableBucket())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace with a schema-bearing table", func() {
			spec := minimalTableBucket()
			spec.Namespaces = []*AwsS3TablesNamespace{analyticsNamespace()}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts KMS encryption with a key", func() {
			spec := minimalTableBucket()
			spec.Encryption = &AwsS3TablesEncryption{
				SseAlgorithm: "aws:kms",
				KmsKeyArn:    literal("arn:aws:kms:us-east-1:111122223333:key/abc"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts maintenance dials and replication", func() {
			spec := minimalTableBucket()
			spec.UnreferencedFileRemoval = &AwsS3TablesUnreferencedFileRemoval{
				NonCurrentDays:   20,
				UnreferencedDays: 5,
			}
			spec.Replication = &AwsS3TablesReplication{
				Role:                       literal("arn:aws:iam::111122223333:role/replication"),
				DestinationTableBucketArns: []string{"arn:aws:s3tables:us-west-2:111122223333:bucket/replica"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a table with per-table maintenance and policy", func() {
			namespace := analyticsNamespace()
			namespace.Tables[0].Compaction = &AwsS3TablesCompaction{TargetFileSizeMb: 256}
			namespace.Tables[0].SnapshotManagement = &AwsS3TablesSnapshotManagement{
				MaxSnapshotAgeHours: 240,
				MinSnapshotsToKeep:  3,
			}
			namespace.Tables[0].ResourcePolicy = `{"Version":"2012-10-17","Statement":[]}`
			spec := minimalTableBucket()
			spec.Namespaces = []*AwsS3TablesNamespace{namespace}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects duplicate namespace names", func() {
			spec := minimalTableBucket()
			spec.Namespaces = []*AwsS3TablesNamespace{analyticsNamespace(), analyticsNamespace()}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate table names within a namespace", func() {
			namespace := analyticsNamespace()
			namespace.Tables = append(namespace.Tables, &AwsS3TablesTable{Name: "events"})
			spec := minimalTableBucket()
			spec.Namespaces = []*AwsS3TablesNamespace{namespace}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a KMS key under AES256", func() {
			spec := minimalTableBucket()
			spec.Encryption = &AwsS3TablesEncryption{
				SseAlgorithm: "AES256",
				KmsKeyArn:    literal("arn:aws:kms:us-east-1:111122223333:key/abc"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase namespace name", func() {
			namespace := analyticsNamespace()
			namespace.Name = "Analytics"
			spec := minimalTableBucket()
			spec.Namespaces = []*AwsS3TablesNamespace{namespace}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a hyphenated table name", func() {
			namespace := analyticsNamespace()
			namespace.Tables[0].Name = "raw-events"
			spec := minimalTableBucket()
			spec.Namespaces = []*AwsS3TablesNamespace{namespace}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty schema field set", func() {
			namespace := analyticsNamespace()
			namespace.Tables[0].IcebergSchema = &AwsS3TablesIcebergSchema{}
			spec := minimalTableBucket()
			spec.Namespaces = []*AwsS3TablesNamespace{namespace}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate schema field names", func() {
			namespace := analyticsNamespace()
			namespace.Tables[0].IcebergSchema.Fields = append(
				namespace.Tables[0].IcebergSchema.Fields,
				&AwsS3TablesSchemaField{Name: "event_id", Type: "long"},
			)
			spec := minimalTableBucket()
			spec.Namespaces = []*AwsS3TablesNamespace{namespace}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects replication without destinations", func() {
			spec := minimalTableBucket()
			spec.Replication = &AwsS3TablesReplication{
				Role: literal("arn:aws:iam::111122223333:role/replication"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a non-s3tables destination arn", func() {
			spec := minimalTableBucket()
			spec.Replication = &AwsS3TablesReplication{
				Role:                       literal("arn:aws:iam::111122223333:role/replication"),
				DestinationTableBucketArns: []string{"arn:aws:s3:::plain-bucket"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects six replication destinations", func() {
			spec := minimalTableBucket()
			arns := []string{}
			for _, suffix := range []string{"a", "b", "c", "d", "e", "f"} {
				arns = append(arns, "arn:aws:s3tables:us-west-2:111122223333:bucket/"+suffix)
			}
			spec.Replication = &AwsS3TablesReplication{
				Role:                       literal("arn:aws:iam::111122223333:role/replication"),
				DestinationTableBucketArns: arns,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
