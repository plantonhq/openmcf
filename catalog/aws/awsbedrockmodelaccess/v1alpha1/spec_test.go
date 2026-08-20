package awsbedrockmodelaccessv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsBedrockModelAccessSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockModelAccessSpec Validation Suite")
}

// minimalModelAccess is the smallest valid manifest: region + model id.
func minimalModelAccess() *AwsBedrockModelAccessSpec {
	return &AwsBedrockModelAccessSpec{
		Region:  "us-west-2",
		ModelId: "amazon.nova-micro-v1:0",
	}
}

var _ = ginkgo.Describe("AwsBedrockModelAccessSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalModelAccess())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a use-case form", func() {
			ginkgo.It("should not return a validation error", func() {
				form, err := structpb.NewStruct(map[string]any{
					"companyName":    "Example Corp",
					"companyWebsite": "https://example.com",
					"intendedUsers":  "Internal employees",
					"industryOption": "Technology",
					"useCases":       "Customer support assistant",
				})
				gomega.Expect(err).To(gomega.BeNil())

				spec := minimalModelAccess()
				spec.ModelId = "anthropic.claude-3-5-haiku-20241022-v1:0"
				spec.UseCaseForm = form
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			})
		})
	})

	// -----------------------------------------------------------------
	// Invalid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("should reject a missing region", func() {
			spec := minimalModelAccess()
			spec.Region = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a missing model id", func() {
			spec := minimalModelAccess()
			spec.ModelId = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
