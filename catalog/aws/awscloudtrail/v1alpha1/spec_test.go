package awscloudtrailv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsCloudTrailSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCloudTrailSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalTrail is the smallest valid instance: region + delivery
// bucket.
func minimalTrail() *AwsCloudTrailSpec {
	return &AwsCloudTrailSpec{
		Region:       "us-west-2",
		S3BucketName: svr("my-trail-bucket"),
	}
}

var _ = ginkgo.Describe("AwsCloudTrailSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal trail", func() {
			gomega.Expect(protovalidate.Validate(minimalTrail())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a multi-region audit trail with validation and insights", func() {
			spec := minimalTrail()
			spec.IsMultiRegionTrail = true
			spec.EnableLogFileValidation = true
			spec.InsightTypes = []string{"ApiCallRateInsight", "ApiErrorRateInsight"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts classic event selectors with data resources", func() {
			spec := minimalTrail()
			spec.EventSelectors = []*AwsCloudTrailEventSelector{{
				ReadWriteType: "All",
				DataResources: []*AwsCloudTrailDataResource{{
					Type:   "AWS::S3::Object",
					Values: []string{"arn:aws:s3:::my-data-bucket/"},
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts advanced event selectors", func() {
			spec := minimalTrail()
			spec.AdvancedEventSelectors = []*AwsCloudTrailAdvancedEventSelector{{
				Name: "Management events only",
				FieldSelectors: []*AwsCloudTrailFieldSelector{{
					Field:  "eventCategory",
					Equals: []string{"Management"},
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts CloudWatch Logs mirroring with both group and role", func() {
			spec := minimalTrail()
			spec.CloudwatchLogs = &AwsCloudTrailCloudwatchLogs{
				LogGroupArn: svr("arn:aws:logs:us-west-2:123456789012:log-group:trail:*"),
				RoleArn:     svr("arn:aws:iam::123456789012:role/trail-to-cw"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an organization trail with a delegated admin", func() {
			spec := minimalTrail()
			spec.IsOrganizationTrail = true
			spec.OrganizationDelegatedAdminAccountId = "123456789012"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalTrail()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing delivery bucket", func() {
			spec := minimalTrail()
			spec.S3BucketName = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects mixing classic and advanced selectors", func() {
			spec := minimalTrail()
			spec.EventSelectors = []*AwsCloudTrailEventSelector{{}}
			spec.AdvancedEventSelectors = []*AwsCloudTrailAdvancedEventSelector{{
				FieldSelectors: []*AwsCloudTrailFieldSelector{{
					Field:  "eventCategory",
					Equals: []string{"Data"},
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more than five classic selectors", func() {
			spec := minimalTrail()
			for i := 0; i < 6; i++ {
				spec.EventSelectors = append(spec.EventSelectors, &AwsCloudTrailEventSelector{})
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a field selector without any condition", func() {
			spec := minimalTrail()
			spec.AdvancedEventSelectors = []*AwsCloudTrailAdvancedEventSelector{{
				FieldSelectors: []*AwsCloudTrailFieldSelector{{Field: "eventCategory"}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an advanced selector without field selectors", func() {
			spec := minimalTrail()
			spec.AdvancedEventSelectors = []*AwsCloudTrailAdvancedEventSelector{{Name: "empty"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown data-resource type", func() {
			spec := minimalTrail()
			spec.EventSelectors = []*AwsCloudTrailEventSelector{{
				DataResources: []*AwsCloudTrailDataResource{{
					Type:   "AWS::EC2::Instance",
					Values: []string{"arn:aws:ec2"},
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects CloudWatch Logs without the role", func() {
			spec := minimalTrail()
			spec.CloudwatchLogs = &AwsCloudTrailCloudwatchLogs{
				LogGroupArn: svr("arn:aws:logs:us-west-2:123456789012:log-group:trail:*"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a delegated admin on a non-organization trail", func() {
			spec := minimalTrail()
			spec.OrganizationDelegatedAdminAccountId = "123456789012"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed delegated-admin account id", func() {
			spec := minimalTrail()
			spec.IsOrganizationTrail = true
			spec.OrganizationDelegatedAdminAccountId = "12345"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown insight type", func() {
			spec := minimalTrail()
			spec.InsightTypes = []string{"ApiVolumeInsight"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an excluded management source outside the AWS vocabulary", func() {
			spec := minimalTrail()
			spec.EventSelectors = []*AwsCloudTrailEventSelector{{
				ExcludeManagementEventSources: []string{"s3.amazonaws.com"},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
