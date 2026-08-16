package awsiamsamlproviderv1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsIamSamlProviderSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsIamSamlProviderSpec Validation Suite")
}

// sampleMetadata builds a metadata document over AWS's 1000-character
// floor (real IdP metadata is far larger; only the length matters to
// the schema).
func sampleMetadata() string {
	return "<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" entityID=\"https://idp.example.com\">" +
		strings.Repeat("<!-- certificate padding -->", 40) +
		"</EntityDescriptor>"
}

func minimalProvider() *AwsIamSamlProviderSpec {
	return &AwsIamSamlProviderSpec{
		Region:               "us-east-1",
		SamlMetadataDocument: sampleMetadata(),
	}
}

var _ = ginkgo.Describe("AwsIamSamlProviderSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal provider", func() {
			gomega.Expect(protovalidate.Validate(minimalProvider())).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalProvider()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing metadata document", func() {
			spec := minimalProvider()
			spec.SamlMetadataDocument = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a metadata document under AWS's 1000-character floor", func() {
			spec := minimalProvider()
			spec.SamlMetadataDocument = "<EntityDescriptor/>"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
