package awsbedrockinferenceprofilev1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsBedrockInferenceProfileSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockInferenceProfileSpec Validation Suite")
}

// minimalProfile is the smallest valid manifest: region + model source.
func minimalProfile() *AwsBedrockInferenceProfileSpec {
	return &AwsBedrockInferenceProfileSpec{
		Region:    "us-west-2",
		SourceArn: "arn:aws:bedrock:us-west-2::foundation-model/amazon.nova-micro-v1:0",
	}
}

var _ = ginkgo.Describe("AwsBedrockInferenceProfileSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with a foundation-model source", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalProfile())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a system-defined cross-region profile source", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalProfile()
				spec.SourceArn = "arn:aws:bedrock:us-west-2:123456789012:inference-profile/us.amazon.nova-micro-v1:0"
				spec.Description = "cost tracking for the checkout service"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// -----------------------------------------------------------------
	// Invalid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("should reject a missing region", func() {
			spec := minimalProfile()
			spec.Region = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a missing source ARN", func() {
			spec := minimalProfile()
			spec.SourceArn = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a bare model id as the source", func() {
			spec := minimalProfile()
			spec.SourceArn = "amazon.nova-micro-v1:0"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a non-model ARN as the source", func() {
			spec := minimalProfile()
			spec.SourceArn = "arn:aws:bedrock:us-west-2:123456789012:guardrail/gr-abc"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a description over 200 characters", func() {
			spec := minimalProfile()
			spec.Description = strings.Repeat("d", 201)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
