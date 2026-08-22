package digitaloceansshkeyv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestDigitalOceanSshKeySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanSshKeySpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanSshKeySpec validations", func() {

	makeValidSpec := func() *DigitalOceanSshKeySpec {
		return &DigitalOceanSshKeySpec{
			KeyName:   "deploy-key",
			PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIB5K1XlWQxr9nMytXvvFyzYZaFVNSTAmTUYbSGXPqIQd e2e@example.com",
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a minimal valid spec", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing key_name", func() {
			spec := makeValidSpec()
			spec.KeyName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing public_key", func() {
			spec := makeValidSpec()
			spec.PublicKey = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
