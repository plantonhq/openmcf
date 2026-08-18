package awscloudwatchlogaccountpolicyv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsCloudwatchLogAccountPolicySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCloudwatchLogAccountPolicySpec Validation Suite")
}

func samplePolicyDocument() *structpb.Struct {
	doc, err := structpb.NewStruct(map[string]any{
		"Fields": []any{"requestId", "sourceIp"},
	})
	if err != nil {
		panic(err)
	}
	return doc
}

func minimalPolicy() *AwsCloudwatchLogAccountPolicySpec {
	return &AwsCloudwatchLogAccountPolicySpec{
		Region:         "us-east-1",
		PolicyName:     "account-field-index",
		PolicyType:     "FIELD_INDEX_POLICY",
		PolicyDocument: samplePolicyDocument(),
	}
}

var _ = ginkgo.Describe("AwsCloudwatchLogAccountPolicySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal account policy", func() {
			gomega.Expect(protovalidate.Validate(minimalPolicy())).To(gomega.BeNil())
		})

		ginkgo.It("accepts every policy type", func() {
			for _, policyType := range []string{
				"DATA_PROTECTION_POLICY",
				"SUBSCRIPTION_FILTER_POLICY",
				"FIELD_INDEX_POLICY",
				"TRANSFORMER_POLICY",
				"METRIC_EXTRACTION_POLICY",
			} {
				spec := minimalPolicy()
				spec.PolicyType = policyType
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			}
		})

		ginkgo.It("accepts selection criteria", func() {
			spec := minimalPolicy()
			spec.SelectionCriteria = "LogGroupNamePrefix IN [\"my-service\"]"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing policy name", func() {
			spec := minimalPolicy()
			spec.PolicyName = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown policy type", func() {
			spec := minimalPolicy()
			spec.PolicyType = "RETENTION_POLICY"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing policy document", func() {
			spec := minimalPolicy()
			spec.PolicyDocument = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
