package awsbedrockagentcoretokenvaultv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBedrockAgentCoreTokenVaultSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockAgentCoreTokenVaultSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// customerManagedVault pairs the customer key type with its KMS key.
func customerManagedVault() *AwsBedrockAgentCoreTokenVaultSpec {
	return &AwsBedrockAgentCoreTokenVaultSpec{
		Region:    "us-west-2",
		KeyType:   "CustomerManagedKey",
		KmsKeyArn: svr("arn:aws:kms:us-west-2:123456789012:key/11111111-2222-3333-4444-555555555555"),
	}
}

// serviceManagedVault is the revert-to-AWS-owned posture.
func serviceManagedVault() *AwsBedrockAgentCoreTokenVaultSpec {
	return &AwsBedrockAgentCoreTokenVaultSpec{
		Region:  "us-west-2",
		KeyType: "ServiceManagedKey",
	}
}

var _ = ginkgo.Describe("AwsBedrockAgentCoreTokenVaultSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with a customer-managed key", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(customerManagedVault())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the service-managed posture", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(serviceManagedVault())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("targeting a named vault", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := customerManagedVault()
				spec.TokenVaultId = "default"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects CustomerManagedKey without a KMS key", func() {
			spec := customerManagedVault()
			spec.KmsKeyArn = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("kms_key_arn is required"))
		})

		ginkgo.It("rejects ServiceManagedKey with a KMS key", func() {
			spec := serviceManagedVault()
			spec.KmsKeyArn = svr("arn:aws:kms:us-west-2:123456789012:key/11111111-2222-3333-4444-555555555555")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must be unset"))
		})

		ginkgo.It("rejects an unknown key type", func() {
			spec := serviceManagedVault()
			spec.KeyType = "BringYourOwnKey"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing key type", func() {
			spec := serviceManagedVault()
			spec.KeyType = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing region", func() {
			spec := serviceManagedVault()
			spec.Region = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
