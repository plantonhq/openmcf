package awsiamgroupv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsIamGroupSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsIamGroupSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalGroup is the smallest valid instance: an empty group (the
// name comes from metadata.name).
func minimalGroup() *AwsIamGroupSpec {
	return &AwsIamGroupSpec{
		Region: "us-east-1",
	}
}

var _ = ginkgo.Describe("AwsIamGroupSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal empty group", func() {
			gomega.Expect(protovalidate.Validate(minimalGroup())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a path", func() {
			spec := minimalGroup()
			spec.Path = "/teams/"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the root path", func() {
			spec := minimalGroup()
			spec.Path = "/"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts users and managed policies", func() {
			spec := minimalGroup()
			spec.Users = []*foreignkeyv1.StringValueOrRef{svr("alice"), svr("bob")}
			spec.ManagedPolicyArns = []*foreignkeyv1.StringValueOrRef{svr("arn:aws:iam::aws:policy/ReadOnlyAccess")}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts inline policies", func() {
			document, err := structpb.NewStruct(map[string]any{
				"Version": "2012-10-17",
				"Statement": []any{map[string]any{
					"Effect":   "Allow",
					"Action":   "s3:ListBucket",
					"Resource": "*",
				}},
			})
			gomega.Expect(err).To(gomega.BeNil())
			spec := minimalGroup()
			spec.InlinePolicies = map[string]*structpb.Struct{"list-buckets": document}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalGroup()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a path without the trailing slash", func() {
			spec := minimalGroup()
			spec.Path = "/teams"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a path without the leading slash", func() {
			spec := minimalGroup()
			spec.Path = "teams/"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a path with empty segments", func() {
			spec := minimalGroup()
			spec.Path = "/teams//platform/"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
