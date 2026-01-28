package azurevpcv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
)

func TestAzureVpcSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AzureVpcSpec Custom Validation Tests")
}

var _ = ginkgo.Describe("AzureVpcSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("azure_vpc", func() {

			ginkgo.It("should not return a validation error for minimal valid fields", func() {
				input := &AzureVpc{
					ApiVersion: "azure.openmcf.org/v1",
					Kind:       "AzureVpc",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-vpc",
					},
					Spec: &AzureVpcSpec{
						AddressSpaceCidr: "10.0.0.0/16",
						NodesSubnetCidr:  "10.0.0.0/18",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
