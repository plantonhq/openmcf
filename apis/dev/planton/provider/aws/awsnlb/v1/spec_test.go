package awsnlbv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsNlbSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsNlbSpec Validation Tests")
}

func literalRef(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

// minimalValidNlb is the common case: a region and one subnet mapping.
func minimalValidNlb() *AwsNlb {
	return &AwsNlb{
		ApiVersion: "aws.planton.dev/v1",
		Kind:       "AwsNlb",
		Metadata: &shared.CloudResourceMetadata{
			Name: "demo-nlb",
		},
		Spec: &AwsNlbSpec{
			Region: "us-west-2",
			SubnetMappings: []*AwsNlbSubnetMapping{
				{SubnetId: literalRef("subnet-12345678")},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsNlbSpec Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("aws_nlb", func() {

			ginkgo.It("should not return a validation error for a minimal NLB", func() {
				err := protovalidate.Validate(minimalValidNlb())
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for a static-IP internet-facing NLB", func() {
				input := minimalValidNlb()
				input.Spec.SubnetMappings = []*AwsNlbSubnetMapping{
					{SubnetId: literalRef("subnet-12345678"), AllocationId: literalRef("eipalloc-0123456789abcdef0")},
					{SubnetId: literalRef("subnet-12345679"), AllocationId: literalRef("eipalloc-0fedcba9876543210")},
				}
				input.Spec.CrossZoneLoadBalancingEnabled = true
				input.Spec.IpAddressType = "ipv4"
				input.Spec.DnsRecordClientRoutingPolicy = "availability_zone_affinity"
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for a fully tuned NLB", func() {
				input := minimalValidNlb()
				input.Spec.SecurityGroups = []*foreignkeyv1.StringValueOrRef{literalRef("sg-0123456789abcdef0")}
				input.Spec.Internal = true
				input.Spec.DeleteProtectionEnabled = true
				input.Spec.ZonalShiftEnabled = true
				input.Spec.EnforceSecurityGroupInboundRulesOnPrivateLinkTraffic = "on"
				input.Spec.AccessLogs = &AwsNlbAccessLogs{
					Bucket: literalRef("demo-nlb-logs"),
					Prefix: "nlb/demo",
				}
				input.Spec.Dns = &AwsNlbDns{
					Enabled:       true,
					Route53ZoneId: literalRef("Z0123456789ABCDEF"),
					Hostnames:     []string{"tcp.example.com"},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("aws_nlb", func() {

			ginkgo.It("should return a validation error when kind is wrong", func() {
				input := minimalValidNlb()
				input.Kind = "WrongKind"
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when region is empty", func() {
				input := minimalValidNlb()
				input.Spec.Region = ""
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when subnet mappings are empty", func() {
				input := minimalValidNlb()
				input.Spec.SubnetMappings = nil
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for a mapping without a subnet", func() {
				input := minimalValidNlb()
				input.Spec.SubnetMappings = []*AwsNlbSubnetMapping{
					{AllocationId: literalRef("eipalloc-0123456789abcdef0")},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an invalid ip_address_type", func() {
				input := minimalValidNlb()
				input.Spec.IpAddressType = "dualstack-without-public-ipv4"
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an invalid dns_record_client_routing_policy", func() {
				input := minimalValidNlb()
				input.Spec.DnsRecordClientRoutingPolicy = "sticky"
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an invalid private-link enforcement value", func() {
				input := minimalValidNlb()
				input.Spec.EnforceSecurityGroupInboundRulesOnPrivateLinkTraffic = "enabled"
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for access logs without a bucket", func() {
				input := minimalValidNlb()
				input.Spec.AccessLogs = &AwsNlbAccessLogs{Prefix: "nlb/demo"}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for duplicate DNS hostnames", func() {
				input := minimalValidNlb()
				input.Spec.Dns = &AwsNlbDns{
					Enabled:       true,
					Route53ZoneId: literalRef("Z0123456789ABCDEF"),
					Hostnames:     []string{"tcp.example.com", "tcp.example.com"},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
