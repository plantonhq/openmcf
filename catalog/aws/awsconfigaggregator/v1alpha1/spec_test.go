package awsconfigaggregatorv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsConfigAggregatorSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsConfigAggregatorSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalAggregator is the smallest valid collector: one source
// account in one region.
func minimalAggregator() *AwsConfigAggregatorSpec {
	return &AwsConfigAggregatorSpec{
		Region: "us-west-2",
		Aggregation: &AwsConfigAggregatorAggregation{
			AccountSource: &AwsConfigAggregatorAccountSource{
				AccountIds: []string{"123456789012"},
				Regions:    []string{"us-west-2"},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsConfigAggregatorSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal account-source aggregator", func() {
			gomega.Expect(protovalidate.Validate(minimalAggregator())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an all-regions account source", func() {
			spec := minimalAggregator()
			spec.Aggregation.AccountSource.Regions = nil
			spec.Aggregation.AccountSource.AllRegions = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an organization source with a role", func() {
			spec := &AwsConfigAggregatorSpec{
				Region: "us-west-2",
				Aggregation: &AwsConfigAggregatorAggregation{
					OrganizationSource: &AwsConfigAggregatorOrganizationSource{
						RoleArn:    svr("arn:aws:iam::123456789012:role/config-org-aggregator"),
						AllRegions: true,
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a grants-only source-account deployment", func() {
			spec := &AwsConfigAggregatorSpec{
				Region: "us-west-2",
				Authorizations: []*AwsConfigAggregatorAuthorization{{
					AccountId:           "210987654321",
					AuthorizedAwsRegion: "us-west-2",
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an account playing both roles", func() {
			spec := minimalAggregator()
			spec.Authorizations = []*AwsConfigAggregatorAuthorization{{
				AccountId:           "210987654321",
				AuthorizedAwsRegion: "eu-west-1",
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalAggregator()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an instance managing neither arm", func() {
			spec := &AwsConfigAggregatorSpec{Region: "us-west-2"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an aggregation with both sources", func() {
			spec := minimalAggregator()
			spec.Aggregation.OrganizationSource = &AwsConfigAggregatorOrganizationSource{
				RoleArn:    svr("arn:aws:iam::123456789012:role/config-org-aggregator"),
				AllRegions: true,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an aggregation with neither source", func() {
			spec := minimalAggregator()
			spec.Aggregation.AccountSource = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an account source without accounts", func() {
			spec := minimalAggregator()
			spec.Aggregation.AccountSource.AccountIds = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed source account id", func() {
			spec := minimalAggregator()
			spec.Aggregation.AccountSource.AccountIds = []string{"12345"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an account source with neither regions nor all_regions", func() {
			spec := minimalAggregator()
			spec.Aggregation.AccountSource.Regions = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an organization source without the role", func() {
			spec := &AwsConfigAggregatorSpec{
				Region: "us-west-2",
				Aggregation: &AwsConfigAggregatorAggregation{
					OrganizationSource: &AwsConfigAggregatorOrganizationSource{AllRegions: true},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an organization source with neither regions nor all_regions", func() {
			spec := &AwsConfigAggregatorSpec{
				Region: "us-west-2",
				Aggregation: &AwsConfigAggregatorAggregation{
					OrganizationSource: &AwsConfigAggregatorOrganizationSource{
						RoleArn: svr("arn:aws:iam::123456789012:role/config-org-aggregator"),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a grant with a malformed aggregator account id", func() {
			spec := &AwsConfigAggregatorSpec{
				Region: "us-west-2",
				Authorizations: []*AwsConfigAggregatorAuthorization{{
					AccountId:           "not-an-account",
					AuthorizedAwsRegion: "us-west-2",
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a grant without the aggregator region", func() {
			spec := &AwsConfigAggregatorSpec{
				Region: "us-west-2",
				Authorizations: []*AwsConfigAggregatorAuthorization{{
					AccountId: "210987654321",
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
