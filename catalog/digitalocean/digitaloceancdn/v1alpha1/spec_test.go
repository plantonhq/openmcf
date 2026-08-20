package digitaloceancdnv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestDigitalOceanCdnSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanCdnSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanCdnSpec validations", func() {

	newRef := func(value string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: value},
		}
	}

	makeValidSpec := func() *DigitalOceanCdnSpec {
		return &DigitalOceanCdnSpec{
			Origin: newRef("app-assets.nyc3.digitaloceanspaces.com"),
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a minimal valid spec", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the origin by reference", func() {
			spec := makeValidSpec()
			spec.Origin = &fk.StringValueOrRef{
				LiteralOrRef: &fk.StringValueOrRef_ValueFrom{
					ValueFrom: &fk.ValueFromRef{Name: "app-assets"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing origin", func() {
			spec := makeValidSpec()
			spec.Origin = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("TTL floor", func() {
		ginkgo.It("accepts ttl of 1 second", func() {
			spec := makeValidSpec()
			spec.Ttl = proto.Int32(1)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the provider default of 3600", func() {
			spec := makeValidSpec()
			spec.Ttl = proto.Int32(3600)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects an explicit zero ttl (unsendable upstream)", func() {
			spec := makeValidSpec()
			spec.Ttl = proto.Int32(0)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a negative ttl", func() {
			spec := makeValidSpec()
			spec.Ttl = proto.Int32(-60)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Custom domain and certificate pairing", func() {
		ginkgo.It("accepts a certificate without a custom domain", func() {
			spec := makeValidSpec()
			spec.Certificate = newRef("assets-cert")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a custom domain paired with a certificate", func() {
			spec := makeValidSpec()
			spec.Certificate = newRef("assets-cert")
			spec.CustomDomain = "assets.example.com"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the needs-cloudflare-cert sentinel", func() {
			spec := makeValidSpec()
			spec.Certificate = newRef("needs-cloudflare-cert")
			spec.CustomDomain = "assets.example.com"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects a custom domain without a certificate", func() {
			spec := makeValidSpec()
			spec.CustomDomain = "assets.example.com"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a custom domain that is not an FQDN", func() {
			spec := makeValidSpec()
			spec.Certificate = newRef("assets-cert")
			spec.CustomDomain = "not a domain"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a bare hostname without a TLD", func() {
			spec := makeValidSpec()
			spec.Certificate = newRef("assets-cert")
			spec.CustomDomain = "assets"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
