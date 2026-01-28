package azurednszonev1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
)

func TestAzureDnsZoneSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AzureDnsZoneSpec Custom Validation Tests")
}

var _ = ginkgo.Describe("AzureDnsZoneSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("azure_dns_zone", func() {

			ginkgo.It("should not return a validation error for minimal valid fields", func() {
				input := &AzureDnsZone{
					ApiVersion: "azure.openmcf.org/v1",
					Kind:       "AzureDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &AzureDnsZoneSpec{
						ZoneName:      "example.com",
						ResourceGroup: "test-resource-group",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
