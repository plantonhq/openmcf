package awscloudtraileventdatastorev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsCloudTrailEventDataStoreSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCloudTrailEventDataStoreSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalStore is the smallest valid instance: just the region (AWS
// defaults everything else, including the all-management-events
// selector).
func minimalStore() *AwsCloudTrailEventDataStoreSpec {
	return &AwsCloudTrailEventDataStoreSpec{
		Region: "us-west-2",
	}
}

func boolPtr(b bool) *bool { return &b }

var _ = ginkgo.Describe("AwsCloudTrailEventDataStoreSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal store", func() {
			gomega.Expect(protovalidate.Validate(minimalStore())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a fixed-pricing store with a bounded retention", func() {
			spec := minimalStore()
			spec.BillingMode = "FIXED_RETENTION_PRICING"
			spec.RetentionPeriodDays = 90
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a KMS-encrypted organization store", func() {
			spec := minimalStore()
			spec.KmsKeyId = svr("arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab")
			spec.OrganizationEnabled = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the deletable posture with ingestion paused", func() {
			spec := minimalStore()
			spec.TerminationProtectionEnabled = boolPtr(false)
			spec.Suspend = boolPtr(true)
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts advanced event selectors with an eventCategory condition", func() {
			spec := minimalStore()
			spec.AdvancedEventSelectors = []*AwsCloudTrailEventDataStoreAdvancedEventSelector{{
				Name: "S3 data events only",
				FieldSelectors: []*AwsCloudTrailEventDataStoreFieldSelector{
					{Field: "eventCategory", Equals: []string{"Data"}},
					{Field: "resources.type", Equals: []string{"AWS::S3::Object"}},
				},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the retention boundaries", func() {
			spec := minimalStore()
			spec.RetentionPeriodDays = 7
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			spec.RetentionPeriodDays = 2555
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalStore()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown billing mode", func() {
			spec := minimalStore()
			spec.BillingMode = "PAY_PER_QUERY"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a retention below 7 days", func() {
			spec := minimalStore()
			spec.RetentionPeriodDays = 6
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a retention above 2555 days", func() {
			spec := minimalStore()
			spec.RetentionPeriodDays = 2556
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an advanced selector without field selectors", func() {
			spec := minimalStore()
			spec.AdvancedEventSelectors = []*AwsCloudTrailEventDataStoreAdvancedEventSelector{{Name: "empty"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a field selector without any condition", func() {
			spec := minimalStore()
			spec.AdvancedEventSelectors = []*AwsCloudTrailEventDataStoreAdvancedEventSelector{{
				FieldSelectors: []*AwsCloudTrailEventDataStoreFieldSelector{{Field: "eventCategory"}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown selector field", func() {
			spec := minimalStore()
			spec.AdvancedEventSelectors = []*AwsCloudTrailEventDataStoreAdvancedEventSelector{{
				FieldSelectors: []*AwsCloudTrailEventDataStoreFieldSelector{{
					Field:  "requestParameters.bucketName",
					Equals: []string{"my-bucket"},
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a condition value beyond 2048 characters", func() {
			long := make([]byte, 2049)
			for i := range long {
				long[i] = 'a'
			}
			spec := minimalStore()
			spec.AdvancedEventSelectors = []*AwsCloudTrailEventDataStoreAdvancedEventSelector{{
				FieldSelectors: []*AwsCloudTrailEventDataStoreFieldSelector{{
					Field:  "eventCategory",
					Equals: []string{string(long)},
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
