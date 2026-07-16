package awskmskeyv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsKmsKeySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsKmsKeySpec Validation Suite")
}

var _ = ginkgo.Describe("AwsKmsKeySpec validations", func() {
	var spec *AwsKmsKeySpec

	ginkgo.BeforeEach(func() {
		spec = &AwsKmsKeySpec{
			Region: "us-west-2",
		}
	})

	// -----------------------------------------------------------------
	// Happy paths
	// -----------------------------------------------------------------

	ginkgo.It("accepts a minimal symmetric key", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a fully configured symmetric encryption key", func() {
		spec.Description = "encrypts the orders database"
		spec.KeySpec = "SYMMETRIC_DEFAULT"
		spec.KeyUsage = "ENCRYPT_DECRYPT"
		spec.Policy = `{"Version":"2012-10-17","Statement":[{"Sid":"root","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:root"},"Action":"kms:*","Resource":"*"}]}`
		spec.EnableKeyRotation = true
		spec.RotationPeriodInDays = 180
		spec.MultiRegion = true
		spec.DeletionWindowDays = 7
		spec.Aliases = []string{"alias/orders-db", "alias/orders-db-legacy"}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts an RSA signing key", func() {
		spec.KeySpec = "RSA_4096"
		spec.KeyUsage = "SIGN_VERIFY"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts an RSA encryption key", func() {
		spec.KeySpec = "RSA_2048"
		spec.KeyUsage = "ENCRYPT_DECRYPT"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts an ECC signing key", func() {
		spec.KeySpec = "ECC_NIST_P256"
		spec.KeyUsage = "SIGN_VERIFY"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts an HMAC key", func() {
		spec.KeySpec = "HMAC_256"
		spec.KeyUsage = "GENERATE_VERIFY_MAC"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a disabled key", func() {
		spec.Disabled = true
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Shape validation
	// -----------------------------------------------------------------

	ginkgo.It("rejects an unknown key spec", func() {
		spec.KeySpec = "AES_128"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an unknown key usage", func() {
		spec.KeyUsage = "KEY_AGREEMENT"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a deletion window outside 7-30", func() {
		spec.DeletionWindowDays = 45
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Rotation couplings
	// -----------------------------------------------------------------

	ginkgo.It("rejects rotation on an asymmetric key", func() {
		spec.KeySpec = "RSA_2048"
		spec.KeyUsage = "ENCRYPT_DECRYPT"
		spec.EnableKeyRotation = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects rotation on an HMAC key", func() {
		spec.KeySpec = "HMAC_256"
		spec.KeyUsage = "GENERATE_VERIFY_MAC"
		spec.EnableKeyRotation = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a rotation period without rotation enabled", func() {
		spec.RotationPeriodInDays = 180
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a rotation period below the AWS floor", func() {
		spec.EnableKeyRotation = true
		spec.RotationPeriodInDays = 30
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Spec/usage couplings
	// -----------------------------------------------------------------

	ginkgo.It("rejects SIGN_VERIFY on a symmetric key", func() {
		spec.KeyUsage = "SIGN_VERIFY"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects GENERATE_VERIFY_MAC on a non-HMAC key", func() {
		spec.KeySpec = "RSA_2048"
		spec.KeyUsage = "GENERATE_VERIFY_MAC"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an HMAC key without GENERATE_VERIFY_MAC", func() {
		spec.KeySpec = "HMAC_384"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an ECC key without SIGN_VERIFY", func() {
		spec.KeySpec = "ECC_NIST_P384"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Aliases
	// -----------------------------------------------------------------

	ginkgo.It("rejects an alias without the alias/ prefix", func() {
		spec.Aliases = []string{"orders-db"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects the reserved alias/aws/ prefix", func() {
		spec.Aliases = []string{"alias/aws/orders"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects duplicate aliases", func() {
		spec.Aliases = []string{"alias/orders-db", "alias/orders-db"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})
})
