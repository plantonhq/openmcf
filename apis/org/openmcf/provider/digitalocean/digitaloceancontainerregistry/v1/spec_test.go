package digitaloceancontainerregistryv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/provider/digitalocean"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
)

func TestDigitalOceanContainerRegistrySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanContainerRegistrySpec Custom Validation Tests")
}

var _ = ginkgo.Describe("DigitalOceanContainerRegistrySpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("digitalocean_container_registry", func() {

			ginkgo.It("should not return a validation error for minimal valid fields", func() {
				input := &DigitalOceanContainerRegistry{
					ApiVersion: "digital-ocean.openmcf.org/v1",
					Kind:       "DigitalOceanContainerRegistry",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-container-registry",
					},
					Spec: &DigitalOceanContainerRegistrySpec{
						Name:             "test-registry",
						SubscriptionTier: DigitalOceanContainerRegistryTier_starter,
						Region:           digitalocean.DigitalOceanRegion_nyc3,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
