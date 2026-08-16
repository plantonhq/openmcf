package awsconfigconformancepackv1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsConfigConformancePackSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsConfigConformancePackSpec Validation Suite")
}

const sampleTemplate = `Resources:
  S3BucketVersioningEnabled:
    Type: AWS::Config::ConfigRule
    Properties:
      ConfigRuleName: s3-bucket-versioning-enabled
      Source:
        Owner: AWS
        SourceIdentifier: S3_BUCKET_VERSIONING_ENABLED
`

// minimalPack is the smallest valid instance: region plus an inline
// template, deployed at account scope.
func minimalPack() *AwsConfigConformancePackSpec {
	return &AwsConfigConformancePackSpec{
		Region:       "us-west-2",
		TemplateBody: sampleTemplate,
	}
}

var _ = ginkgo.Describe("AwsConfigConformancePackSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal account-scope pack", func() {
			gomega.Expect(protovalidate.Validate(minimalPack())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an S3-templated pack with parameters", func() {
			spec := minimalPack()
			spec.TemplateBody = ""
			spec.TemplateS3Uri = "s3://my-templates/operational-best-practices.yaml"
			spec.InputParameters = []*AwsConfigConformancePackInputParameter{{
				ParameterName:  "AccessKeysRotatedParamMaxAccessKeyAge",
				ParameterValue: "90",
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an account pack with BOTH template forms (AWS prefers the S3 one)", func() {
			spec := minimalPack()
			spec.TemplateS3Uri = "s3://my-templates/pack.yaml"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an organization pack with exclusions", func() {
			spec := minimalPack()
			spec.OrganizationScope = true
			spec.ExcludedAccounts = []string{"123456789012", "210987654321"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalPack()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a pack without any template", func() {
			spec := minimalPack()
			spec.TemplateBody = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an organization pack with both template forms", func() {
			spec := minimalPack()
			spec.OrganizationScope = true
			spec.TemplateS3Uri = "s3://my-templates/pack.yaml"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects exclusions on an account-scope pack", func() {
			spec := minimalPack()
			spec.ExcludedAccounts = []string{"123456789012"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a template URI outside s3://", func() {
			spec := minimalPack()
			spec.TemplateBody = ""
			spec.TemplateS3Uri = "https://example.com/pack.yaml"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a template body beyond 51200 characters", func() {
			spec := minimalPack()
			spec.TemplateBody = strings.Repeat("a", 51201)
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more than 60 input parameters", func() {
			spec := minimalPack()
			for i := 0; i < 61; i++ {
				spec.InputParameters = append(spec.InputParameters, &AwsConfigConformancePackInputParameter{
					ParameterName:  "p",
					ParameterValue: "v",
				})
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a parameter without a value", func() {
			spec := minimalPack()
			spec.InputParameters = []*AwsConfigConformancePackInputParameter{{ParameterName: "p"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed excluded account id", func() {
			spec := minimalPack()
			spec.OrganizationScope = true
			spec.ExcludedAccounts = []string{"12345"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
