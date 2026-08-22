package digitaloceanvpcpeeringv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanVpcPeeringSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanVpcPeeringSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanVpcPeeringSpec validations", func() {

	newVpcRef := func(vpcId string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: vpcId},
		}
	}

	makeValidSpec := func() *DigitalOceanVpcPeeringSpec {
		return &DigitalOceanVpcPeeringSpec{
			PeeringName: "app-to-data",
			Vpc_1:       newVpcRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
			Vpc_2:       newVpcRef("ffffffff-1111-2222-3333-444444444444"),
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a minimal valid spec", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts VPC references by name", func() {
			spec := makeValidSpec()
			spec.Vpc_1 = &fk.StringValueOrRef{
				LiteralOrRef: &fk.StringValueOrRef_ValueFrom{
					ValueFrom: &fk.ValueFromRef{Name: "vpc-app"},
				},
			}
			spec.Vpc_2 = &fk.StringValueOrRef{
				LiteralOrRef: &fk.StringValueOrRef_ValueFrom{
					ValueFrom: &fk.ValueFromRef{Name: "vpc-data"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing peering_name", func() {
			spec := makeValidSpec()
			spec.PeeringName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing vpc_1", func() {
			spec := makeValidSpec()
			spec.Vpc_1 = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing vpc_2", func() {
			spec := makeValidSpec()
			spec.Vpc_2 = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
