package awsbedrockprovisionedthroughputv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBedrockProvisionedThroughputSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockProvisionedThroughputSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalThroughput is the smallest valid manifest: a no-commitment
// single-unit purchase.
func minimalThroughput() *AwsBedrockProvisionedThroughputSpec {
	return &AwsBedrockProvisionedThroughputSpec{
		Region:     "us-east-1",
		ModelArn:   svr("arn:aws:bedrock:us-east-1:123456789012:custom-model/amazon.titan-text-lite-v1/abc123"),
		ModelUnits: 1,
	}
}

var _ = ginkgo.Describe("AwsBedrockProvisionedThroughputSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with a no-commitment purchase", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalThroughput())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a committed purchase", func() {
			ginkgo.It("should accept OneMonth", func() {
				spec := minimalThroughput()
				spec.CommitmentDuration = "OneMonth"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should accept SixMonths", func() {
				spec := minimalThroughput()
				spec.CommitmentDuration = "SixMonths"
				spec.ModelUnits = 4
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
			spec := minimalThroughput()
			spec.Region = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a missing model reference", func() {
			spec := minimalThroughput()
			spec.ModelArn = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject zero model units", func() {
			spec := minimalThroughput()
			spec.ModelUnits = 0
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject negative model units", func() {
			spec := minimalThroughput()
			spec.ModelUnits = -1
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an unknown commitment duration", func() {
			spec := minimalThroughput()
			spec.CommitmentDuration = "TwelveMonths"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
