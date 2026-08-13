package awsrestapiusageplanv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsRestApiUsagePlanSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRestApiUsagePlanSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func boolPtr(v bool) *bool { return &v }

// minimalPlan is the smallest valid plan: region only (an unattached
// plan is legal in AWS).
func minimalPlan() *AwsRestApiUsagePlanSpec {
	return &AwsRestApiUsagePlanSpec{Region: "us-west-2"}
}

var _ = ginkgo.Describe("AwsRestApiUsagePlanSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalPlan())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalPlan()
				spec.Description = "Free tier"
				spec.ApiStages = []*AwsRestApiUsagePlanApiStage{
					{
						RestApiId: svr("abc123"),
						StageName: svr("prod"),
						MethodThrottles: []*AwsRestApiUsagePlanMethodThrottle{
							{Path: "/search/GET", BurstLimit: 10, RateLimit: 5},
						},
					},
				}
				spec.Quota = &AwsRestApiUsagePlanQuota{Limit: 1000, Period: "DAY"}
				spec.Throttle = &AwsRestApiUsagePlanThrottle{BurstLimit: 100, RateLimit: 50}
				spec.ApiKeys = []*AwsRestApiUsagePlanApiKey{
					{Name: "acme-mobile", Description: "acme mobile app", Enabled: boolPtr(true)},
					{Name: "acme-web", Value: "supersecretapikeyvalue123"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects duplicate API key names", func() {
			spec := minimalPlan()
			spec.ApiKeys = []*AwsRestApiUsagePlanApiKey{
				{Name: "k"},
				{Name: "k"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("api_keys entries must have unique names"))
		})

		ginkgo.It("rejects a quota offset outside its period", func() {
			spec := minimalPlan()
			spec.Quota = &AwsRestApiUsagePlanQuota{Limit: 100, Period: "DAY", Offset: 1}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("quota offset must be"))

			spec.Quota = &AwsRestApiUsagePlanQuota{Limit: 100, Period: "WEEK", Offset: 7}
			err = protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())

			spec.Quota = &AwsRestApiUsagePlanQuota{Limit: 100, Period: "MONTH", Offset: 28}
			err = protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts valid quota offsets", func() {
			spec := minimalPlan()
			spec.Quota = &AwsRestApiUsagePlanQuota{Limit: 100, Period: "WEEK", Offset: 6}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			spec.Quota = &AwsRestApiUsagePlanQuota{Limit: 100, Period: "MONTH", Offset: 27}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects an all-zero throttle", func() {
			spec := minimalPlan()
			spec.Throttle = &AwsRestApiUsagePlanThrottle{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("at least one of burst_limit or rate_limit"))
		})

		ginkgo.It("rejects duplicate method throttle paths in a stage", func() {
			spec := minimalPlan()
			spec.ApiStages = []*AwsRestApiUsagePlanApiStage{
				{
					RestApiId: svr("abc"),
					StageName: svr("prod"),
					MethodThrottles: []*AwsRestApiUsagePlanMethodThrottle{
						{Path: "/a/GET", RateLimit: 1},
						{Path: "/a/GET", RateLimit: 2},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("unique path values"))
		})

		ginkgo.It("rejects an api stage without its API", func() {
			spec := minimalPlan()
			spec.ApiStages = []*AwsRestApiUsagePlanApiStage{
				{StageName: svr("prod")},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a short explicit key value", func() {
			spec := minimalPlan()
			spec.ApiKeys = []*AwsRestApiUsagePlanApiKey{
				{Name: "k", Value: "tooshort"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid quota period", func() {
			spec := minimalPlan()
			spec.Quota = &AwsRestApiUsagePlanQuota{Limit: 100, Period: "YEAR"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
