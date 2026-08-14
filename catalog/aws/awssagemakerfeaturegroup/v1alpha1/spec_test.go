package awssagemakerfeaturegroupv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsSagemakerFeatureGroupSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSagemakerFeatureGroupSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalFeatureGroup is the smallest valid manifest: an online-only
// group with the two special features declared.
func minimalFeatureGroup() *AwsSagemakerFeatureGroupSpec {
	return &AwsSagemakerFeatureGroupSpec{
		Region:                      "us-west-2",
		RecordIdentifierFeatureName: "customer_id",
		EventTimeFeatureName:        "event_time",
		RoleArn:                     svr("arn:aws:iam::123456789012:role/sagemaker-execution"),
		FeatureDefinitions: []*AwsSagemakerFeatureGroupFeature{
			{Name: "customer_id", Type: "String"},
			{Name: "event_time", Type: "Fractional"},
			{Name: "lifetime_value", Type: "Fractional"},
		},
		OnlineStore: &AwsSagemakerFeatureGroupOnlineStore{Enabled: true},
	}
}

var _ = ginkgo.Describe("AwsSagemakerFeatureGroupSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalFeatureGroup())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full dual-store surface", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.Description = "customer features"
				spec.OnlineStore = &AwsSagemakerFeatureGroupOnlineStore{
					Enabled:     true,
					KmsKeyArn:   svr("arn:aws:kms:us-west-2:123456789012:key/abc"),
					StorageType: "InMemory",
					Ttl:         &AwsSagemakerFeatureGroupTtl{Unit: "Days", Value: 30},
				}
				spec.OfflineStore = &AwsSagemakerFeatureGroupOfflineStore{
					S3Uri:                    "s3://my-features/",
					KmsKeyArn:                svr("arn:aws:kms:us-west-2:123456789012:key/abc"),
					DisableGlueTableCreation: false,
					TableFormat:              "Iceberg",
					DataCatalog: &AwsSagemakerFeatureGroupDataCatalog{
						Catalog:   "AwsDataCatalog",
						Database:  "features",
						TableName: "customers",
					},
				}
				spec.Throughput = &AwsSagemakerFeatureGroupThroughput{
					Mode:                          "Provisioned",
					ProvisionedReadCapacityUnits:  proto.Int32(100),
					ProvisionedWriteCapacityUnits: proto.Int32(50),
				}
				spec.FeatureDefinitions = append(spec.FeatureDefinitions, &AwsSagemakerFeatureGroupFeature{
					Name:            "embedding",
					Type:            "Fractional",
					CollectionType:  "Vector",
					VectorDimension: proto.Int32(256),
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an offline-only store", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.OnlineStore = nil
				spec.OfflineStore = &AwsSagemakerFeatureGroupOfflineStore{S3Uri: "s3://my-features/"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with no store at all", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.OnlineStore = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an online store block that is not enabled and no offline store", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.OnlineStore = &AwsSagemakerFeatureGroupOnlineStore{Enabled: false}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate feature names", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.FeatureDefinitions = append(spec.FeatureDefinitions, &AwsSagemakerFeatureGroupFeature{
					Name: "customer_id",
					Type: "String",
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a record identifier that is not a feature", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.RecordIdentifierFeatureName = "missing"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an event time that is not a feature", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.EventTimeFeatureName = "missing"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a reserved feature name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.FeatureDefinitions = append(spec.FeatureDefinitions, &AwsSagemakerFeatureGroupFeature{
					Name: "is_deleted",
					Type: "String",
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a Vector feature missing its dimension", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.FeatureDefinitions = append(spec.FeatureDefinitions, &AwsSagemakerFeatureGroupFeature{
					Name:           "embedding",
					Type:           "Fractional",
					CollectionType: "Vector",
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a dimension on a non-Vector feature", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.FeatureDefinitions = append(spec.FeatureDefinitions, &AwsSagemakerFeatureGroupFeature{
					Name:            "scores",
					Type:            "Fractional",
					CollectionType:  "List",
					VectorDimension: proto.Int32(10),
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with capacity units in OnDemand mode", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.Throughput = &AwsSagemakerFeatureGroupThroughput{
					Mode:                         "OnDemand",
					ProvisionedReadCapacityUnits: proto.Int32(100),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a TTL missing its unit", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFeatureGroup()
				spec.OnlineStore.Ttl = &AwsSagemakerFeatureGroupTtl{Value: 30}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
