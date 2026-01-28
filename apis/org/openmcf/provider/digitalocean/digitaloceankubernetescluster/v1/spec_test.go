package digitaloceankubernetesclusterv1

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"

	"buf.build/go/protovalidate"
	"github.com/plantonhq/openmcf/apis/org/openmcf/provider/digitalocean"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
)

func TestDigitalOceanKubernetesClusterSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanKubernetesClusterSpec Custom Validation Tests")
}

var _ = ginkgo.Describe("DigitalOceanKubernetesClusterSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("digitalocean_kubernetes_cluster", func() {

			ginkgo.It("should not return a validation error for minimal valid fields", func() {
				input := &DigitalOceanKubernetesCluster{
					ApiVersion: "digital-ocean.openmcf.org/v1",
					Kind:       "DigitalOceanKubernetesCluster",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-k8s-cluster",
					},
					Spec: &DigitalOceanKubernetesClusterSpec{
						ClusterName:       "test-cluster",
						Region:            digitalocean.DigitalOceanRegion_nyc3,
						KubernetesVersion: "1.26.3",
						Vpc: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-vpc"},
						},
						DefaultNodePool: &DigitalOceanKubernetesClusterDefaultNodePool{
							Size:      "s-2vcpu-4gb",
							NodeCount: 3,
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
