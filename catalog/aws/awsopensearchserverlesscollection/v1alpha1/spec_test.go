package awsopensearchserverlesscollectionv1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsOpenSearchServerlessCollectionSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsOpenSearchServerlessCollectionSpec Validation Suite")
}

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalCollection is the smallest valid manifest: a region alone (type and
// standby replicas carry defaults; encryption/network default module-side).
func minimalCollection() *AwsOpenSearchServerlessCollectionSpec {
	return &AwsOpenSearchServerlessCollectionSpec{
		Region: "us-west-2",
	}
}

// readWriteRule grants one principal full index access - the typical
// application rule.
func readWriteRule() *AwsOpenSearchServerlessCollectionAccessRule {
	return &AwsOpenSearchServerlessCollectionAccessRule{
		Principals: []*foreignkeyv1.StringValueOrRef{
			svr("arn:aws:iam::123456789012:role/app-role"),
		},
		IndexPermissions: []string{"aoss:ReadDocument", "aoss:WriteDocument"},
		IndexPatterns:    []string{"*"},
	}
}

var _ = ginkgo.Describe("AwsOpenSearchServerlessCollectionSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalCollection())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with full production configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsOpenSearchServerlessCollectionSpec{
					Region:              "us-west-2",
					Description:         "Application search collection",
					Type:                strPtr("SEARCH"),
					StandbyReplicas:     strPtr("ENABLED"),
					CollectionGroupName: "shared-capacity",
					Encryption: &AwsOpenSearchServerlessCollectionEncryption{
						KmsKeyArn: svr("arn:aws:kms:us-west-2:123456789012:key/abc-123"),
					},
					Network: &AwsOpenSearchServerlessCollectionNetwork{
						AllowFromPublic: boolPtr(false),
						VpcEndpointIds:  []string{"vpce-0123456789abcdef0"},
					},
					DataAccess: []*AwsOpenSearchServerlessCollectionAccessRule{
						readWriteRule(),
						{
							Principals:            []*foreignkeyv1.StringValueOrRef{svr("arn:aws:iam::123456789012:role/admin-role")},
							CollectionPermissions: []string{"aoss:*"},
							IndexPermissions:      []string{"aoss:*"},
						},
					},
					RetentionRules: []*AwsOpenSearchServerlessCollectionRetentionRule{
						{IndexPatterns: []string{"logs-*"}, MinIndexRetention: "30d"},
						{IndexPatterns: []string{"audit-*"}, Unlimited: true},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a vector search collection", func() {
			ginkgo.It("should accept vector acceleration on VECTORSEARCH", func() {
				spec := minimalCollection()
				spec.Type = strPtr("VECTORSEARCH")
				spec.ServerlessVectorAcceleration = "ENABLED"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with each collection type", func() {
			ginkgo.It("should accept SEARCH, TIMESERIES, and VECTORSEARCH", func() {
				for _, t := range []string{"SEARCH", "TIMESERIES", "VECTORSEARCH"} {
					spec := minimalCollection()
					spec.Type = strPtr(t)
					gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
				}
			})
		})
	})

	// -----------------------------------------------------------------
	// Invalid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("region and description", func() {
			ginkgo.It("should reject an empty region", func() {
				spec := minimalCollection()
				spec.Region = ""
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a description over 1000 characters", func() {
				spec := minimalCollection()
				spec.Description = strings.Repeat("x", 1001)
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("enumerated fields", func() {
			ginkgo.It("should reject an unknown collection type", func() {
				spec := minimalCollection()
				spec.Type = strPtr("ANALYTICS")
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown standby_replicas value", func() {
				spec := minimalCollection()
				spec.StandbyReplicas = strPtr("MAYBE")
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown vector acceleration value", func() {
				spec := minimalCollection()
				spec.Type = strPtr("VECTORSEARCH")
				spec.ServerlessVectorAcceleration = "TURBO"
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("collection group name", func() {
			ginkgo.It("should reject uppercase and short names", func() {
				spec := minimalCollection()
				spec.CollectionGroupName = "SharedCapacity"
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
				spec.CollectionGroupName = "ab"
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("vector acceleration coupling", func() {
			ginkgo.It("should reject acceleration on a non-vector collection", func() {
				spec := minimalCollection()
				spec.Type = strPtr("TIMESERIES")
				spec.ServerlessVectorAcceleration = "ENABLED"
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("network", func() {
			ginkgo.It("should reject private access without VPC endpoints", func() {
				spec := minimalCollection()
				spec.Network = &AwsOpenSearchServerlessCollectionNetwork{
					AllowFromPublic: boolPtr(false),
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject VPC endpoints with public access", func() {
				spec := minimalCollection()
				spec.Network = &AwsOpenSearchServerlessCollectionNetwork{
					AllowFromPublic: boolPtr(true),
					VpcEndpointIds:  []string{"vpce-0123456789abcdef0"},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a malformed VPC endpoint id", func() {
				spec := minimalCollection()
				spec.Network = &AwsOpenSearchServerlessCollectionNetwork{
					AllowFromPublic: boolPtr(false),
					VpcEndpointIds:  []string{"vpc-0123456789abcdef0"},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("data access rules", func() {
			ginkgo.It("should reject a rule with no principals", func() {
				spec := minimalCollection()
				spec.DataAccess = []*AwsOpenSearchServerlessCollectionAccessRule{
					{IndexPermissions: []string{"aoss:ReadDocument"}},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a rule that grants nothing", func() {
				spec := minimalCollection()
				spec.DataAccess = []*AwsOpenSearchServerlessCollectionAccessRule{
					{Principals: []*foreignkeyv1.StringValueOrRef{svr("arn:aws:iam::123456789012:role/app-role")}},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown permission value", func() {
				spec := minimalCollection()
				rule := readWriteRule()
				rule.IndexPermissions = []string{"aoss:DropIndex"}
				spec.DataAccess = []*AwsOpenSearchServerlessCollectionAccessRule{rule}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject index patterns without index permissions", func() {
				spec := minimalCollection()
				spec.DataAccess = []*AwsOpenSearchServerlessCollectionAccessRule{
					{
						Principals:            []*foreignkeyv1.StringValueOrRef{svr("arn:aws:iam::123456789012:role/app-role")},
						CollectionPermissions: []string{"aoss:DescribeCollectionItems"},
						IndexPatterns:         []string{"logs-*"},
					},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("retention rules", func() {
			ginkgo.It("should reject a rule without index patterns", func() {
				spec := minimalCollection()
				spec.RetentionRules = []*AwsOpenSearchServerlessCollectionRetentionRule{
					{MinIndexRetention: "24h"},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject both retention arms", func() {
				spec := minimalCollection()
				spec.RetentionRules = []*AwsOpenSearchServerlessCollectionRetentionRule{
					{IndexPatterns: []string{"*"}, MinIndexRetention: "24h", Unlimited: true},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject neither retention arm", func() {
				spec := minimalCollection()
				spec.RetentionRules = []*AwsOpenSearchServerlessCollectionRetentionRule{
					{IndexPatterns: []string{"*"}},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a malformed retention period", func() {
				spec := minimalCollection()
				spec.RetentionRules = []*AwsOpenSearchServerlessCollectionRetentionRule{
					{IndexPatterns: []string{"*"}, MinIndexRetention: "24x"},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})
	})
})
