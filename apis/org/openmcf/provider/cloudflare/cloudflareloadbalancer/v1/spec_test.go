package cloudflareloadbalancerv1

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"

	"buf.build/go/protovalidate"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
)

func TestCloudflareLoadBalancerSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "CloudflareLoadBalancerSpec Custom Validation Tests")
}

var _ = ginkgo.Describe("CloudflareLoadBalancerSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("cloudflare_load_balancer", func() {

			ginkgo.It("should not return a validation error for minimal valid fields", func() {
				input := &CloudflareLoadBalancer{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareLoadBalancer",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-load-balancer",
					},
					Spec: &CloudflareLoadBalancerSpec{
						Hostname: "lb.example.com",
						ZoneId: &foreignkeyv1.StringValueOrRef{
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "test-zone-123"},
						},
						Origins: []*CloudflareLoadBalancerOrigin{
							{
								Name:    "origin-1",
								Address: "192.168.1.1",
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
