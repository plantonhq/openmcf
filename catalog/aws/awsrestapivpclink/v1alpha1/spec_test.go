package awsrestapivpclinkv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsRestApiVpcLinkSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRestApiVpcLinkSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalVpcLink is the smallest valid link: region + target NLB.
func minimalVpcLink() *AwsRestApiVpcLinkSpec {
	return &AwsRestApiVpcLinkSpec{
		Region:    "us-west-2",
		TargetArn: svr("arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/net/orders/abc"),
	}
}

var _ = ginkgo.Describe("AwsRestApiVpcLinkSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalVpcLink())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a description", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalVpcLink()
				spec.Description = "orders service NLB"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a link without its target NLB", func() {
			spec := minimalVpcLink()
			spec.TargetArn = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing region", func() {
			spec := minimalVpcLink()
			spec.Region = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an oversized description", func() {
			spec := minimalVpcLink()
			long := make([]byte, 1025)
			for i := range long {
				long[i] = 'x'
			}
			spec.Description = string(long)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
