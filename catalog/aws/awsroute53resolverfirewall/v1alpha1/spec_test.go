package awsroute53resolverfirewallv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsRoute53ResolverFirewallSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRoute53ResolverFirewallSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func ttl(v int64) *int64 { return &v }

func groupWithListAndRule() *AwsRoute53ResolverFirewallSpec {
	return &AwsRoute53ResolverFirewallSpec{
		Region: "us-west-2",
		DomainLists: []*AwsRoute53ResolverFirewallDomainList{
			{Name: "blocked", Domains: []string{"malware.example.", "phishing.example."}},
		},
		Rules: []*AwsRoute53ResolverFirewallRule{
			{
				Name:           "block-bad",
				Priority:       100,
				Action:         "BLOCK",
				DomainListName: "blocked",
			},
		},
		VpcAssociations: []*AwsRoute53ResolverFirewallVpcAssociation{
			{Name: "main", VpcId: literal("vpc-0123456789abcdef0"), Priority: 200},
		},
	}
}

var _ = ginkgo.Describe("AwsRoute53ResolverFirewallSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal empty rule group", func() {
			spec := &AwsRoute53ResolverFirewallSpec{Region: "us-west-2"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a group with an owned list, a BLOCK rule, and an association", func() {
			gomega.Expect(protovalidate.Validate(groupWithListAndRule())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an OVERRIDE block response with domain and ttl", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].BlockResponse = "OVERRIDE"
			spec.Rules[0].BlockOverrideDomain = "sinkhole.example.com"
			spec.Rules[0].BlockOverrideTtl = ttl(300)
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a threat-protection rule with a confidence threshold", func() {
			spec := groupWithListAndRule()
			spec.Rules = append(spec.Rules, &AwsRoute53ResolverFirewallRule{
				Name:                "block-tunneling",
				Priority:            200,
				Action:              "BLOCK",
				DnsThreatProtection: "DNS_TUNNELING",
				ConfidenceThreshold: "HIGH",
			})
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a rule matching an external managed list by id", func() {
			spec := groupWithListAndRule()
			spec.Rules = append(spec.Rules, &AwsRoute53ResolverFirewallRule{
				Name:         "aws-managed",
				Priority:     300,
				Action:       "ALERT",
				DomainListId: literal("rslvr-fdl-0123456789abcdef"),
				QType:        "TXT",
			})
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a rule with no match source", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].DomainListName = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule with two match sources", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].DnsThreatProtection = "DGA"
			spec.Rules[0].ConfidenceThreshold = "LOW"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects threat protection without a confidence threshold", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].DomainListName = ""
			spec.Rules[0].DnsThreatProtection = "DGA"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a confidence threshold without threat protection", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].ConfidenceThreshold = "LOW"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects OVERRIDE without the override record", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].BlockResponse = "OVERRIDE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects override fields without block_response OVERRIDE", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].BlockOverrideDomain = "sinkhole.example.com"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects block_response on an ALLOW rule", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].Action = "ALLOW"
			spec.Rules[0].BlockResponse = "NXDOMAIN"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate rule priorities", func() {
			spec := groupWithListAndRule()
			spec.Rules = append(spec.Rules, &AwsRoute53ResolverFirewallRule{
				Name:                "dup-priority",
				Priority:            100,
				Action:              "ALERT",
				DnsThreatProtection: "DGA",
				ConfidenceThreshold: "LOW",
			})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule naming a domain list that is not owned here", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].DomainListName = "missing"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate domain list names", func() {
			spec := groupWithListAndRule()
			spec.DomainLists = append(spec.DomainLists, &AwsRoute53ResolverFirewallDomainList{Name: "blocked"})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an association without a vpc", func() {
			spec := groupWithListAndRule()
			spec.VpcAssociations[0].VpcId = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an association without a priority", func() {
			spec := groupWithListAndRule()
			spec.VpcAssociations[0].Priority = 0
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range override ttl", func() {
			spec := groupWithListAndRule()
			spec.Rules[0].BlockResponse = "OVERRIDE"
			spec.Rules[0].BlockOverrideDomain = "sinkhole.example.com"
			spec.Rules[0].BlockOverrideTtl = ttl(700000)
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
