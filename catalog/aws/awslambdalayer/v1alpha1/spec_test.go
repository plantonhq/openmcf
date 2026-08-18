package awslambdalayerv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsLambdaLayerSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsLambdaLayerSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func minimalConfig() *AwsLambdaLayerSpec {
	return &AwsLambdaLayerSpec{
		Region: "us-west-2",
		Code: &AwsLambdaLayerS3Code{
			Bucket: literal("layer-artifacts"),
			Key:    "layers/shared-utils.zip",
		},
	}
}

var _ = ginkgo.Describe("AwsLambdaLayerSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal configuration", func() {
			gomega.Expect(protovalidate.Validate(minimalConfig())).To(gomega.BeNil())
		})

		ginkgo.It("accepts runtimes, architectures, and license metadata", func() {
			spec := minimalConfig()
			spec.CompatibleRuntimes = []string{"python3.13", "python3.12"}
			spec.CompatibleArchitectures = []string{"x86_64", "arm64"}
			spec.LicenseInfo = "Apache-2.0"
			spec.Description = "Shared Python utilities"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an account-scoped share grant", func() {
			spec := minimalConfig()
			spec.Permissions = []*AwsLambdaLayerPermission{
				{StatementId: "share-tools", Principal: "222233334444"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an organization-scoped wildcard grant", func() {
			spec := minimalConfig()
			spec.Permissions = []*AwsLambdaLayerPermission{
				{StatementId: "org-share", Principal: "*", OrganizationId: "o-a1b2c3d4e5"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing code source", func() {
			spec := minimalConfig()
			spec.Code = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects code without a bucket", func() {
			spec := minimalConfig()
			spec.Code.Bucket = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects code without a key", func() {
			spec := minimalConfig()
			spec.Code.Key = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown architecture", func() {
			spec := minimalConfig()
			spec.CompatibleArchitectures = []string{"armv7"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate permission statement ids", func() {
			spec := minimalConfig()
			spec.Permissions = []*AwsLambdaLayerPermission{
				{StatementId: "dup", Principal: "222233334444"},
				{StatementId: "dup", Principal: "555566667777"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an organization grant on a specific account principal", func() {
			spec := minimalConfig()
			spec.Permissions = []*AwsLambdaLayerPermission{
				{StatementId: "bad-org", Principal: "222233334444", OrganizationId: "o-a1b2c3d4e5"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a principal that is neither an account id nor a wildcard", func() {
			spec := minimalConfig()
			spec.Permissions = []*AwsLambdaLayerPermission{
				{StatementId: "bad-principal", Principal: "arn:aws:iam::123:root"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed organization id", func() {
			spec := minimalConfig()
			spec.Permissions = []*AwsLambdaLayerPermission{
				{StatementId: "bad-org-id", Principal: "*", OrganizationId: "org-12345"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a statement id with illegal characters", func() {
			spec := minimalConfig()
			spec.Permissions = []*AwsLambdaLayerPermission{
				{StatementId: "bad id!", Principal: "*"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
