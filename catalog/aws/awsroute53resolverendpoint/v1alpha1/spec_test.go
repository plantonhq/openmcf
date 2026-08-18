package awsroute53resolverendpointv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsRoute53ResolverEndpointSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRoute53ResolverEndpointSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func minimalInbound() *AwsRoute53ResolverEndpointSpec {
	return &AwsRoute53ResolverEndpointSpec{
		Region:    "us-west-2",
		Direction: "INBOUND",
		IpAddresses: []*AwsRoute53ResolverEndpointIpAddress{
			{SubnetId: literal("subnet-0123456789abcdef0")},
			{SubnetId: literal("subnet-0123456789abcdef1")},
		},
		SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{literal("sg-0123456789abcdef0")},
	}
}

func outboundWithForwardRule() *AwsRoute53ResolverEndpointSpec {
	spec := minimalInbound()
	spec.Direction = "OUTBOUND"
	spec.Rules = []*AwsRoute53ResolverEndpointRule{
		{
			Name:       "corp-forward",
			DomainName: "corp.example.com",
			RuleType:   "FORWARD",
			TargetIps: []*AwsRoute53ResolverEndpointRuleTarget{
				{Ip: "10.0.0.53"},
			},
			VpcIds: []*foreignkeyv1.StringValueOrRef{literal("vpc-0123456789abcdef0")},
		},
	}
	return spec
}

var _ = ginkgo.Describe("AwsRoute53ResolverEndpointSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal inbound endpoint", func() {
			gomega.Expect(protovalidate.Validate(minimalInbound())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an outbound endpoint with a FORWARD rule and association", func() {
			gomega.Expect(protovalidate.Validate(outboundWithForwardRule())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a SYSTEM rule with no targets", func() {
			spec := outboundWithForwardRule()
			spec.Rules = append(spec.Rules, &AwsRoute53ResolverEndpointRule{
				Name:       "keep-recursive",
				DomainName: "internal.corp.example.com",
				RuleType:   "SYSTEM",
			})
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an IPv6 target with an explicit port and protocol", func() {
			spec := outboundWithForwardRule()
			spec.Rules[0].TargetIps = []*AwsRoute53ResolverEndpointRuleTarget{
				{Ipv6: "2001:db8::53", Port: 8053, Protocol: "DoH"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a dualstack endpoint speaking two protocols", func() {
			spec := minimalInbound()
			spec.EndpointType = "DUALSTACK"
			spec.Protocols = []string{"Do53", "DoH"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts explicit metrics toggles", func() {
			spec := minimalInbound()
			off := false
			spec.RniEnhancedMetricsEnabled = &off
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a FORWARD rule on an inbound endpoint", func() {
			spec := outboundWithForwardRule()
			spec.Direction = "INBOUND"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a single ip address entry", func() {
			spec := minimalInbound()
			spec.IpAddresses = spec.IpAddresses[:1]
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an endpoint without security groups", func() {
			spec := minimalInbound()
			spec.SecurityGroupIds = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a FORWARD rule without targets", func() {
			spec := outboundWithForwardRule()
			spec.Rules[0].TargetIps = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a SYSTEM rule carrying targets", func() {
			spec := outboundWithForwardRule()
			spec.Rules[0].RuleType = "SYSTEM"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a RECURSIVE rule type", func() {
			spec := outboundWithForwardRule()
			spec.Rules[0].RuleType = "RECURSIVE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate rule names", func() {
			spec := outboundWithForwardRule()
			spec.Rules = append(spec.Rules, &AwsRoute53ResolverEndpointRule{
				Name:       "corp-forward",
				DomainName: "other.example.com",
				RuleType:   "SYSTEM",
			})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target with both ip and ipv6", func() {
			spec := outboundWithForwardRule()
			spec.Rules[0].TargetIps = []*AwsRoute53ResolverEndpointRuleTarget{
				{Ip: "10.0.0.53", Ipv6: "2001:db8::53"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target with neither address family", func() {
			spec := outboundWithForwardRule()
			spec.Rules[0].TargetIps = []*AwsRoute53ResolverEndpointRuleTarget{
				{Port: 53},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown direction", func() {
			spec := minimalInbound()
			spec.Direction = "SIDEWAYS"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate protocols", func() {
			spec := minimalInbound()
			spec.Protocols = []string{"Do53", "Do53"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target port out of range", func() {
			spec := outboundWithForwardRule()
			spec.Rules[0].TargetIps[0].Port = 70000
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
