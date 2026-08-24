package aws

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
)

func TestAwsProviderConfig(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsProviderConfig Suite")
}

// baseConfig returns a minimal valid static-credential provider config; the
// provider-block cases below layer their fields on top of it.
func baseConfig() *AwsProviderConfig {
	return &AwsProviderConfig{
		AccountId:       "123456789012",
		AccessKeyId:     "AKIAIOSFODNN7EXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMIK7MDENGbPxRfiCYEXAMPLEKEY00",
		Region:          "us-east-1",
	}
}

var _ = ginkgo.Describe("AwsProviderConfig provider-block surface", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("full provider-block surface: chained assume roles, default tags, endpoints, retries", func() {
			ginkgo.It("should not return a validation error", func() {
				input := baseConfig()
				input.AssumeRoleChain = []*AwsAssumeRole{
					{
						RoleArn:     "arn:aws:iam::111111111111:role/intermediate",
						SessionName: "planton-hop-1",
					},
					{
						RoleArn:           "arn:aws:iam::222222222222:role/deploy",
						ExternalId:        "expected-external-id",
						Duration:          "1h",
						Tags:              map[string]string{"Team": "platform"},
						TransitiveTagKeys: []string{"Team"},
						SourceIdentity:    "platform-engineer",
					},
				}
				input.DefaultTags = &AwsDefaultTags{
					Tags: map[string]string{"CostCenter": "eng", "ManagedBy": "planton"},
				}
				input.Endpoints = map[string]string{"sts": "https://sts.internal.example.com"}
				input.MaxRetries = proto.Int32(0)
				input.RetryMode = "adaptive"
				gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
			})
		})

		ginkgo.Context("retry_mode standard", func() {
			ginkgo.It("should not return a validation error", func() {
				input := baseConfig()
				input.RetryMode = "standard"
				gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
			})
		})

		ginkgo.Context("retry_mode empty (provider default)", func() {
			ginkgo.It("should not return a validation error", func() {
				gomega.Expect(protovalidate.Validate(baseConfig())).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("retry_mode outside the provider's vocabulary", func() {
			ginkgo.It("should return a validation error", func() {
				input := baseConfig()
				input.RetryMode = "exponential"
				gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("assume-role hop missing role_arn", func() {
			ginkgo.It("should return a validation error", func() {
				input := baseConfig()
				input.AssumeRoleChain = []*AwsAssumeRole{{SessionName: "no-role-arn"}}
				gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
			})
		})
	})
})
