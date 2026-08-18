package awscloudwatchlogresourcepolicyv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsCloudwatchLogResourcePolicySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCloudwatchLogResourcePolicySpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func sampleDocument() *structpb.Struct {
	doc, err := structpb.NewStruct(map[string]any{
		"Version": "2012-10-17",
		"Statement": []any{
			map[string]any{
				"Effect":    "Allow",
				"Principal": map[string]any{"Service": "route53.amazonaws.com"},
				"Action":    []any{"logs:CreateLogStream", "logs:PutLogEvents"},
				"Resource":  "arn:aws:logs:us-east-1:123456789012:log-group:/aws/route53/*",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return doc
}

func accountScopePolicy() *AwsCloudwatchLogResourcePolicySpec {
	return &AwsCloudwatchLogResourcePolicySpec{
		Region:         "us-east-1",
		PolicyName:     "route53-query-logging",
		PolicyDocument: sampleDocument(),
	}
}

var _ = ginkgo.Describe("AwsCloudwatchLogResourcePolicySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the account scope", func() {
			gomega.Expect(protovalidate.Validate(accountScopePolicy())).To(gomega.BeNil())
		})

		ginkgo.It("accepts the resource scope", func() {
			spec := accountScopePolicy()
			spec.PolicyName = ""
			spec.ResourceArn = svr("arn:aws:logs:us-east-1:123456789012:log-group:/app/api")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects both scopes at once", func() {
			spec := accountScopePolicy()
			spec.ResourceArn = svr("arn:aws:logs:us-east-1:123456789012:log-group:/app/api")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects neither scope", func() {
			spec := accountScopePolicy()
			spec.PolicyName = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a policy name with colons or asterisks", func() {
			spec := accountScopePolicy()
			spec.PolicyName = "route53:query*"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing policy document", func() {
			spec := accountScopePolicy()
			spec.PolicyDocument = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
