package awsssmparameterv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsSsmParameterSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSsmParameterSpec Validation Suite")
}

// minimalParameter is the smallest valid instance: a flat String
// parameter with a plain value.
func minimalParameter() *AwsSsmParameterSpec {
	return &AwsSsmParameterSpec{
		Region:        "us-west-2",
		ParameterName: "app-log-level",
		Type:          "String",
		Value:         "info",
	}
}

var _ = ginkgo.Describe("AwsSsmParameterSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal flat parameter", func() {
			gomega.Expect(protovalidate.Validate(minimalParameter())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a fully qualified hierarchical name", func() {
			spec := minimalParameter()
			spec.ParameterName = "/prod/db/url"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a SecureString carrying its value in secure_value", func() {
			spec := minimalParameter()
			spec.Type = "SecureString"
			spec.Value = ""
			spec.SecureValue = "s3cret"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a String parameter sourced from secure_value", func() {
			spec := minimalParameter()
			spec.Value = ""
			spec.SecureValue = "not-for-plan-output"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts tier, data_type, allowed_pattern, and overwrite", func() {
			spec := minimalParameter()
			spec.Tier = "Advanced"
			spec.DataType = "text"
			spec.AllowedPattern = "^\\w+$"
			spec.Overwrite = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalParameter()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a hierarchical name without the leading slash", func() {
			spec := minimalParameter()
			spec.ParameterName = "prod/db/url"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a name with illegal characters", func() {
			spec := minimalParameter()
			spec.ParameterName = "app log level"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown type", func() {
			spec := minimalParameter()
			spec.Type = "SecretString"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects both value arms set", func() {
			spec := minimalParameter()
			spec.SecureValue = "also-set"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects neither value arm set", func() {
			spec := minimalParameter()
			spec.Value = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a SecureString carrying its value in the plain arm", func() {
			spec := minimalParameter()
			spec.Type = "SecureString"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown tier", func() {
			spec := minimalParameter()
			spec.Tier = "Premium"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown data type", func() {
			spec := minimalParameter()
			spec.DataType = "aws:rds:snapshot"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
