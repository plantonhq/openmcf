package awsmskserverlessclusterv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/shared"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsMskServerlessClusterSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsMskServerlessClusterSpec Validation Tests")
}

func validMinimalSpec() *AwsMskServerlessCluster {
	return &AwsMskServerlessCluster{
		ApiVersion: "aws.planton.dev/v1alpha1",
		Kind:       "AwsMskServerlessCluster",
		Metadata: &shared.CloudResourceMetadata{
			Name: "test-msk-serverless",
		},
		Spec: &AwsMskServerlessClusterSpec{
			Region: "us-west-2",
			SubnetIds: []*foreignkeyv1.StringValueOrRef{
				{LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "subnet-aaa"}},
				{LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "subnet-bbb"}},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsMskServerlessClusterSpec Validation Tests", func() {

	// ===== HAPPY PATH TESTS =====

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("should accept a minimal valid cluster (subnets only; AWS attaches the default SG)", func() {
			input := validMinimalSpec()
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept a cluster with security groups attached", func() {
			input := validMinimalSpec()
			input.Spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{
				{LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "sg-broker-123"}},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept a cluster with the maximum of 5 security groups", func() {
			input := validMinimalSpec()
			for _, id := range []string{"sg-1", "sg-2", "sg-3", "sg-4", "sg-5"} {
				input.Spec.SecurityGroupIds = append(input.Spec.SecurityGroupIds,
					&foreignkeyv1.StringValueOrRef{LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: id}})
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept a cluster with valueFrom references", func() {
			input := validMinimalSpec()
			input.Spec.SubnetIds = []*foreignkeyv1.StringValueOrRef{
				{LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
					ValueFrom: &foreignkeyv1.ValueFromRef{
						Kind: cloudresourcekind.CloudResourceKind_AwsSubnet,
						Name: "private-subnet-a",
					},
				}},
				{LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
					ValueFrom: &foreignkeyv1.ValueFromRef{
						Kind: cloudresourcekind.CloudResourceKind_AwsSubnet,
						Name: "private-subnet-b",
					},
				}},
			}
			input.Spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{
				{LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
					ValueFrom: &foreignkeyv1.ValueFromRef{
						Kind: cloudresourcekind.CloudResourceKind_AwsSecurityGroup,
						Name: "broker-sg",
					},
				}},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	// ===== FAILURE TESTS =====

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("required fields", func() {
			ginkgo.It("should fail when region is empty", func() {
				input := validMinimalSpec()
				input.Spec.Region = ""
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should fail when subnet_ids is empty", func() {
				input := validMinimalSpec()
				input.Spec.SubnetIds = nil
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("bounds", func() {
			ginkgo.It("should fail when more than 5 security groups are attached", func() {
				input := validMinimalSpec()
				for _, id := range []string{"sg-1", "sg-2", "sg-3", "sg-4", "sg-5", "sg-6"} {
					input.Spec.SecurityGroupIds = append(input.Spec.SecurityGroupIds,
						&foreignkeyv1.StringValueOrRef{LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: id}})
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("API envelope validation", func() {
			ginkgo.It("should fail for wrong apiVersion", func() {
				input := validMinimalSpec()
				input.ApiVersion = "gcp.planton.dev/v1alpha1"
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should fail for wrong kind", func() {
				input := validMinimalSpec()
				input.Kind = "AwsMskCluster"
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should fail when metadata is missing", func() {
				input := validMinimalSpec()
				input.Metadata = nil
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should fail when spec is missing", func() {
				input := validMinimalSpec()
				input.Spec = nil
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
