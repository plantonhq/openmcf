package awsssmdocumentv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsSsmDocumentSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSsmDocumentSpec Validation Suite")
}

// minimalDocument is the smallest valid instance: a JSON Command
// document.
func minimalDocument() *AwsSsmDocumentSpec {
	return &AwsSsmDocumentSpec{
		Region:       "us-west-2",
		Content:      `{"schemaVersion":"2.2","mainSteps":[]}`,
		DocumentType: "Command",
	}
}

var _ = ginkgo.Describe("AwsSsmDocumentSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal Command document", func() {
			gomega.Expect(protovalidate.Validate(minimalDocument())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a YAML Automation document with target type and version name", func() {
			spec := minimalDocument()
			spec.DocumentType = "Automation"
			spec.DocumentFormat = "YAML"
			spec.TargetType = "/AWS::EC2::Instance"
			spec.VersionName = "release-1.0.0"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts attachment sources and sharing", func() {
			spec := minimalDocument()
			spec.DocumentType = "Package"
			spec.AttachmentSources = []*AwsSsmDocumentAttachmentSource{{
				Key:    "S3FileUrl",
				Name:   "installer.zip",
				Values: []string{"https://s3.amazonaws.com/my-bucket/installer.zip"},
			}}
			spec.ShareWithAccountIds = []string{"123456789012", "All"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects missing content", func() {
			spec := minimalDocument()
			spec.Content = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown document type", func() {
			spec := minimalDocument()
			spec.DocumentType = "Playbook"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown document format", func() {
			spec := minimalDocument()
			spec.DocumentFormat = "XML"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target type without the leading slash", func() {
			spec := minimalDocument()
			spec.TargetType = "AWS::EC2::Instance"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a version name below 3 characters", func() {
			spec := minimalDocument()
			spec.VersionName = "v1"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an attachment source without values", func() {
			spec := minimalDocument()
			spec.AttachmentSources = []*AwsSsmDocumentAttachmentSource{{Key: "SourceUrl"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an attachment source with an unknown key", func() {
			spec := minimalDocument()
			spec.AttachmentSources = []*AwsSsmDocumentAttachmentSource{{
				Key:    "GitUrl",
				Values: []string{"https://example.com"},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects sharing with a malformed account id", func() {
			spec := minimalDocument()
			spec.ShareWithAccountIds = []string{"12345"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects sharing with lowercase 'all'", func() {
			spec := minimalDocument()
			spec.ShareWithAccountIds = []string{"all"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
