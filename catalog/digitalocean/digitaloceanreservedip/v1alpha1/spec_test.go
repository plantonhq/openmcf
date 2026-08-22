package digitaloceanreservedipv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	digitalocean "github.com/plantonhq/planton/catalog/digitalocean"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanReservedIpSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanReservedIpSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanReservedIpSpec validations", func() {

	makeValidSpec := func() *DigitalOceanReservedIpSpec {
		return &DigitalOceanReservedIpSpec{
			Region: digitalocean.DigitalOceanRegion_nyc3,
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a minimal valid spec (region only, ipv4 by default)", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing region", func() {
			spec := makeValidSpec()
			spec.Region = digitalocean.DigitalOceanRegion_digital_ocean_region_unspecified
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("IP version", func() {
		ginkgo.It("accepts both explicit versions", func() {
			for _, version := range []string{"ipv4", "ipv6"} {
				spec := makeValidSpec()
				spec.IpVersion = version
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			}
		})

		ginkgo.It("rejects an unknown ip_version", func() {
			spec := makeValidSpec()
			spec.IpVersion = "v4"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Droplet assignment", func() {
		ginkgo.It("accepts a literal numeric droplet id", func() {
			spec := makeValidSpec()
			spec.Droplet = &fk.StringValueOrRef{
				LiteralOrRef: &fk.StringValueOrRef_Value{Value: "123456789"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a droplet reference by name", func() {
			spec := makeValidSpec()
			spec.Droplet = &fk.StringValueOrRef{
				LiteralOrRef: &fk.StringValueOrRef_ValueFrom{
					ValueFrom: &fk.ValueFromRef{Name: "my-droplet"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})
})
