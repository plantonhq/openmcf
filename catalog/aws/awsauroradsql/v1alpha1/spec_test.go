package awsauroradsqlv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsAuroraDsqlSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsAuroraDsqlSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func minimalConfig() *AwsAuroraDsqlSpec {
	return &AwsAuroraDsqlSpec{
		Region: "us-east-1",
	}
}

var _ = ginkgo.Describe("AwsAuroraDsqlSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal single-region cluster", func() {
			gomega.Expect(protovalidate.Validate(minimalConfig())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a customer-managed key and protection dials", func() {
			spec := minimalConfig()
			spec.DeletionProtectionEnabled = true
			spec.ForceDestroy = true
			spec.KmsEncryptionKey = literal("arn:aws:kms:us-east-1:111111111111:key/1234abcd-12ab-34cd-56ef-1234567890ab")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a multi-region pairing", func() {
			spec := minimalConfig()
			spec.MultiRegion = &AwsAuroraDsqlMultiRegion{
				WitnessRegion: "us-west-2",
				PeerClusterArns: []*foreignkeyv1.StringValueOrRef{
					literal("arn:aws:dsql:us-east-2:111111111111:cluster/peer1234"),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalConfig()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects multi_region without a witness region", func() {
			spec := minimalConfig()
			spec.MultiRegion = &AwsAuroraDsqlMultiRegion{
				PeerClusterArns: []*foreignkeyv1.StringValueOrRef{
					literal("arn:aws:dsql:us-east-2:111111111111:cluster/peer1234"),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects multi_region without peers", func() {
			spec := minimalConfig()
			spec.MultiRegion = &AwsAuroraDsqlMultiRegion{
				WitnessRegion: "us-west-2",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
