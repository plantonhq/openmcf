package awsorganizationpolicyv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsOrganizationPolicySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsOrganizationPolicySpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func ref(name string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
			ValueFrom: &foreignkeyv1.ValueFromRef{Name: name},
		},
	}
}

func denyLeaveOrganization() *structpb.Struct {
	content, err := structpb.NewStruct(map[string]any{
		"Version": "2012-10-17",
		"Statement": []any{map[string]any{
			"Effect":   "Deny",
			"Action":   "organizations:LeaveOrganization",
			"Resource": "*",
		}},
	})
	if err != nil {
		panic(err)
	}
	return content
}

// minimalPolicy is the smallest valid instance: an unattached SCP.
func minimalPolicy() *AwsOrganizationPolicySpec {
	return &AwsOrganizationPolicySpec{
		Region:     "us-west-2",
		PolicyName: "Deny Leaving The Organization",
		Content:    denyLeaveOrganization(),
	}
}

var _ = ginkgo.Describe("AwsOrganizationPolicySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal unattached SCP", func() {
			gomega.Expect(protovalidate.Validate(minimalPolicy())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an explicit type with attachments to every target class", func() {
			spec := minimalPolicy()
			spec.Type = "SERVICE_CONTROL_POLICY"
			spec.Description = "Guardrail for every workload account"
			spec.Attachments = []*AwsOrganizationPolicyAttachment{
				{TargetId: svr("r-abcd")},
				{TargetId: svr("ou-abcd-12345678")},
				{TargetId: svr("111111111111")},
				{TargetId: ref("workloads-ou")},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts each 2026 policy type", func() {
			for _, policyType := range []string{"RESOURCE_CONTROL_POLICY", "DECLARATIVE_POLICY_EC2", "BEDROCK_POLICY", "SECURITYHUB_POLICY", "UPGRADE_ROLLOUT_POLICY"} {
				spec := minimalPolicy()
				spec.Type = policyType
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			}
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalPolicy()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing policy name", func() {
			spec := minimalPolicy()
			spec.PolicyName = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing content document", func() {
			spec := minimalPolicy()
			spec.Content = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown policy type", func() {
			spec := minimalPolicy()
			spec.Type = "FIREWALL_POLICY"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a description above 512 characters", func() {
			spec := minimalPolicy()
			description := ""
			for i := 0; i < 513; i++ {
				description += "a"
			}
			spec.Description = description
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a literal attachment target of the wrong shape", func() {
			spec := minimalPolicy()
			spec.Attachments = []*AwsOrganizationPolicyAttachment{{TargetId: svr("vpc-12345678")}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate literal attachment targets", func() {
			spec := minimalPolicy()
			spec.Attachments = []*AwsOrganizationPolicyAttachment{
				{TargetId: svr("ou-abcd-12345678")},
				{TargetId: svr("ou-abcd-12345678")},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an attachment without a target", func() {
			spec := minimalPolicy()
			spec.Attachments = []*AwsOrganizationPolicyAttachment{{}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
