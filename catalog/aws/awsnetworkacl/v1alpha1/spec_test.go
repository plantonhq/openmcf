package awsnetworkaclv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsNetworkAclSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsNetworkAclSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func httpsRule(no int32) *AwsNetworkAclRule {
	return &AwsNetworkAclRule{
		RuleNo:    no,
		Action:    "allow",
		Protocol:  "tcp",
		CidrBlock: "0.0.0.0/0",
		FromPort:  443,
		ToPort:    443,
	}
}

func minimalAcl() *AwsNetworkAclSpec {
	return &AwsNetworkAclSpec{
		Region:  "us-east-1",
		VpcId:   svr("vpc-0123456789abcdef0"),
		Ingress: []*AwsNetworkAclRule{httpsRule(100)},
		Egress:  []*AwsNetworkAclRule{httpsRule(100)},
	}
}

var _ = ginkgo.Describe("AwsNetworkAclSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal ACL", func() {
			gomega.Expect(protovalidate.Validate(minimalAcl())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a rule-less ACL (implicit deny-everything)", func() {
			spec := minimalAcl()
			spec.Ingress = nil
			spec.Egress = nil
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an all-protocols rule with unset ports", func() {
			spec := minimalAcl()
			spec.Ingress = []*AwsNetworkAclRule{{
				RuleNo:    100,
				Action:    "deny",
				Protocol:  "-1",
				CidrBlock: "10.0.0.0/8",
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an IPv6 rule", func() {
			spec := minimalAcl()
			spec.Ingress = []*AwsNetworkAclRule{{
				RuleNo:        100,
				Action:        "allow",
				Protocol:      "6",
				Ipv6CidrBlock: "::/0",
				FromPort:      443,
				ToPort:        443,
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an ICMP echo rule with type 0 (presence-typed boundary)", func() {
			spec := minimalAcl()
			spec.Ingress = []*AwsNetworkAclRule{{
				RuleNo:    100,
				Action:    "allow",
				Protocol:  "icmp",
				CidrBlock: "10.0.0.0/8",
				IcmpType:  proto.Int32(0),
				IcmpCode:  proto.Int32(0),
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the same rule_no across directions", func() {
			gomega.Expect(protovalidate.Validate(minimalAcl())).To(gomega.BeNil())
		})

		ginkgo.It("accepts the rule_no boundaries", func() {
			spec := minimalAcl()
			spec.Ingress = []*AwsNetworkAclRule{httpsRule(1), httpsRule(32766)}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing vpc reference", func() {
			spec := minimalAcl()
			spec.VpcId = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate rule numbers within a direction", func() {
			spec := minimalAcl()
			spec.Ingress = []*AwsNetworkAclRule{httpsRule(100), httpsRule(100)}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule with both address families", func() {
			spec := minimalAcl()
			spec.Ingress[0].Ipv6CidrBlock = "::/0"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule with neither address family", func() {
			spec := minimalAcl()
			spec.Ingress[0].CidrBlock = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects ports on an all-protocols rule", func() {
			spec := minimalAcl()
			spec.Ingress[0].Protocol = "-1"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects ICMP selectors on a TCP rule", func() {
			spec := minimalAcl()
			spec.Ingress[0].IcmpType = proto.Int32(8)
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects all-types with a specific code", func() {
			spec := minimalAcl()
			spec.Ingress = []*AwsNetworkAclRule{{
				RuleNo:    100,
				Action:    "allow",
				Protocol:  "icmp",
				CidrBlock: "10.0.0.0/8",
				IcmpType:  proto.Int32(-1),
				IcmpCode:  proto.Int32(3),
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects rule numbers in AWS's reserved catch-all range", func() {
			spec := minimalAcl()
			spec.Ingress = []*AwsNetworkAclRule{httpsRule(32767)}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown action", func() {
			spec := minimalAcl()
			spec.Ingress[0].Action = "ALLOW"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
